package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/kaneshin/gate"
	"github.com/kaneshin/gate/cmd/internal"
)

var errNotSupportedPlatform = errors.New("not supported platform")
var errNotDefinedTarget = errors.New("not defined target")

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

func postToPixela(target, token, text string) error {
	el := strings.Split(target, "/")
	if len(el) != 2 {
		return errNotDefinedTarget
	}

	config := gate.NewConfig().WithHTTPClient(http.DefaultClient)
	config.WithAccessToken(token)
	config.WithID(el[0])
	svc := gate.NewPixelaService(config)

	payload := gate.GraphPayload{
		ID: el[1],
	}
	var err error
	switch text {
	case "i", "inc", "increment":
		_, err = svc.Increment(payload)
	case "d", "dec", "decrement":
		_, err = svc.Decrement(payload)
	default:
		payload.Date = time.Now().Format("20060102")
		payload.Quantity = text
		_, err = svc.PostGraphPayload(payload)
	}
	if err != nil {
		return err
	}
	return nil
}

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
	if len(el) != 2 {
		return errNotDefinedTarget
	}
	platform, ok := v[el[0]]
	if !ok {
		return errNotDefinedTarget
	}
	target := el[1]
	token, ok := platform[target]
	if !ok {
		return errNotDefinedTarget
	}

	switch el[0] {
	case "slack":
		err = postToSlackIncoming(token, text)
	case "line":
		err = postToLINENotify(token, text)
	case "pixela":
		err = postToPixela(target, token, text)
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

var postReg = regexp.MustCompile(`/post/(|.+\..+)$`)

func postHandler(w http.ResponseWriter, r *http.Request) {
	defer recovery()

	if !postReg.MatchString(r.URL.Path) {
		notFoundHandler(w, r)
		return
	}
	if r.Method != http.MethodPost {
		methodNotAllowdHandler(w, r)
		return
	}

	var text string
	var target string
	matches := postReg.FindStringSubmatch(r.URL.Path)
	if matches[1] == "" {
		// POST /notify/
		err := r.ParseForm()
		if err != nil {
			internalServerErrorHandler(w, r)
			log.Printf("Error: %v", err)
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
	log.Printf("[Post] %s", msg)
	w.Write([]byte(msg))
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	defer recovery()

	switch r.URL.Path {
	case internal.ConfigCLIJSON:
		if r.Method != http.MethodGet {
			methodNotAllowdHandler(w, r)
			return
		}

		gate := config.Gate
		gate.Host = privateIP()

		c := map[string]interface{}{
			"gate": gate,
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(c)
		if err != nil {
			internalServerErrorHandler(w, r)
			log.Printf("Error: %v", err)
			return
		}
		w.Write(buf.Bytes())
	}
}

var privateIPBlocks []*net.IPNet

func init() {
	for _, cidr := range []string{
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			log.Panic(fmt.Errorf("parse error on %q: %v", cidr, err))
		}
		privateIPBlocks = append(privateIPBlocks, block)
	}
}

func isPrivateIP(ip net.IP) bool {
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func privateIP() string {
	net.InterfaceAddrs()
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("Error: %v", err)
		return "0.0.0.0"
	}

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ok && isPrivateIP(ipnet.IP) {
			return ipnet.IP.String()
		}
	}
	return "0.0.0.0"
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
		Slack  map[string]string `json:"slack"`
		Line   map[string]string `json:"line"`
		Pixela map[string]string `json:"pixela"`
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

	http.HandleFunc("/post/", postHandler)
	http.HandleFunc("/config/", configHandler)

	scheme := config.Gate.Scheme
	if scheme == "" {
		scheme = "http" // default value
	}
	port := fmt.Sprintf(":%d", config.Gate.Port)
	if port == ":" {
		scheme = ":5731" // default value
	}
	fmt.Fprintf(os.Stdout, "Listening for HTTP on %s://%s%s\n\n", scheme, config.Gate.Host, port)
	fmt.Fprintf(os.Stdout,
		"Please run the command to fetch cli.json\n\n  curl -sL %s://%s%s%s > ~/.config/gate/cli.json\n\n",
		scheme, privateIP(), port, internal.ConfigCLIJSON)
	log.Fatal(http.ListenAndServe(port, http.DefaultServeMux))
}
