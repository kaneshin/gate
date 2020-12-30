package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"

	"github.com/kaneshin/gate"
)

func postToSlackIncoming(url, text string) error {
	config := gate.NewConfig().WithHTTPClient(http.DefaultClient)
	svc := gate.NewSlackIncomingService(config).WithBaseURL(url)
	_, err := svc.PostTextPayload(gate.TextPayload{
		Text: text,
	})
	if err != nil {
		return err
	}
	return nil
}

func postToLINENotify(token, text string) error {
	config := gate.NewConfig().WithHTTPClient(http.DefaultClient)
	config.WithAccessToken(token)
	svc := gate.NewLINENotifyService(config)
	_, err := svc.PostMessagePayload(gate.MessagePayload{
		Message: text,
	})
	if err != nil {
		return err
	}
	return nil
}

var errNotSupportedPlatform = errors.New("not supported platform")
var errNotDefinedTarget = errors.New("not defined target")

func post(name, text string) error {
	var err error
	var buf bytes.Buffer
	var v map[string]map[string]string

	err = json.NewEncoder(&buf).Encode(config.Platforms)
	if err != nil {
		return err
	}
	err = json.NewDecoder(&buf).Decode(&v)
	if err != nil {
		return err
	}

	el := strings.Split(name, ".")
	platform, ok := v[el[0]]
	if !ok {
		return errNotDefinedTarget
	}
	target, ok := platform[el[1]]
	if !ok {
		return errNotDefinedTarget
	}

	switch el[0] {
	case "slack":
		err = postToSlackIncoming(target, text)
	case "line":
		err = postToLINENotify(target, text)
	default:
		err = errNotSupportedPlatform
	}
	if err != nil {
		return err
	}
	return nil
}

func recovery() {
	if r := recover(); r != nil {
		log.Printf("[Fatal] Recover: %v", r)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	b, err := httputil.DumpRequest(r, false)
	if err != nil {
		log.Printf("[Request] 404 Not Found: %s", err)
	} else {
		log.Printf("[Request] 404 Not Found: %s", string(b))
	}
	http.Error(w, "404 Not Found", http.StatusNotFound)
}

func methodNotAllowdHandler(w http.ResponseWriter, r *http.Request) {
	b, err := httputil.DumpRequest(r, false)
	if err != nil {
		log.Printf("[Request] 405 Method Not Allowed: %s", err)
	} else {
		log.Printf("[Request] 405 Method Not Allowed: %s", string(b))
	}
	http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
}

func internalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	b, err := httputil.DumpRequest(r, false)
	if err != nil {
		log.Printf("[Request] 500 Internal Server Error: %s", err)
	} else {
		log.Printf("[Request] 500 Internal Server Error: %s", string(b))
	}
	http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
}

var notifyReg = regexp.MustCompile(`/notify/(|.+\..+)$`)

func notifyHandler(w http.ResponseWriter, r *http.Request) {
	defer recovery()

	if !notifyReg.MatchString(r.URL.Path) {
		notFoundHandler(w, r)
		return
	}
	if r.Method != http.MethodPost {
		methodNotAllowdHandler(w, r)
		return
	}

	var text string
	var target string
	matches := notifyReg.FindStringSubmatch(r.URL.Path)
	if matches[1] == "" {
		// POST /notify/
		err := r.ParseForm()
		if err != nil {
			internalServerErrorHandler(w, r)
			return
		}
		target = r.FormValue("target")
		text = r.FormValue("text")
	} else {
		// POST /notify/[platform].[target]?t=[text]
		target = matches[1]
		text = r.URL.Query().Get("t")
	}

	var msg string
	err := post(target, text)
	if err != nil {
		msg = fmt.Sprintf("✘ %s: failed %s\n", target, err)
	} else {
		msg = fmt.Sprintf("✔ %s: success\n", target)
	}
	log.Printf("[Notify] %s", msg)
	w.Write([]byte(msg))
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	defer recovery()

	switch r.URL.Path {
	case "/config/cli.json":
		if r.Method != http.MethodGet {
			methodNotAllowdHandler(w, r)
			return
		}
		c := map[string]interface{}{
			"gate": config.Gate,
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(c)
		if err != nil {
			internalServerErrorHandler(w, r)
			return
		}
		w.Write(buf.Bytes())
	}
}

type Config struct {
	Gate struct {
		Scheme string `json:"scheme"`
		Host   string `json:"host"`
		Port   int    `json:"port"`
		Client struct {
			Default string `json:"default"`
		} `json:"client"`
	} `json:"gate"`
	Platforms struct {
		Slack map[string]string `json:"slack"`
		Line  map[string]string `json:"line"`
	} `json:"platforms"`
}

var configPath = flag.String("config", "$HOME/.config/gate/config.json", "")
var config Config

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	f, err := os.Open(os.ExpandEnv(*configPath))
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/notify/", notifyHandler)
	http.HandleFunc("/config/", configHandler)

	scheme := config.Gate.Scheme
	if scheme == "" {
		scheme = "http" // default value
	}
	port := fmt.Sprintf(":%d", config.Gate.Port)
	if port == ":" {
		scheme = ":5731" // default value
	}
	fmt.Fprintf(os.Stdout, "Listening for HTTP on %s://%s%s\n", scheme, config.Gate.Host, port)
	log.Fatal(http.ListenAndServe(port, http.DefaultServeMux))
}
