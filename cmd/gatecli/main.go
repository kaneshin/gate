package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func run() error {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	text := strings.TrimSpace(string(b))
	if text == "" {
		return nil
	}

	if *code {
		text = fmt.Sprintf("```\n%s\n```", text)
	} else if *quote {
		reader := strings.NewReader(text)
		scanner := bufio.NewScanner(reader)
		scanner.Split(bufio.ScanLines)
		text = ""
		for scanner.Scan() {
			text += fmt.Sprintf("> %s\n", scanner.Text())
		}
	}

	if *target == "" {
		*target = config.Gate.Client.Default
	}

	val := url.Values{
		"target": []string{*target},
		"text":   []string{text},
	}

	scheme := config.Gate.Scheme
	if scheme == "" {
		scheme = "http" // default value
	}
	port := fmt.Sprintf(":%d", config.Gate.Port)
	if port == ":" {
		scheme = ":5731" // default value
	}

	url := fmt.Sprintf("%s://%s%s/post/", scheme, config.Gate.Host, port)
	resp, err := http.PostForm(url, val)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Print(string(b))
	return nil
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
}

var target = flag.String("target", "", "Post to the specified target")
var code = flag.Bool("code", false, "Be inline-code")
var quote = flag.Bool("quote", false, "Be quote-text")
var configPath = flag.String("config", "$HOME/.config/gate/cli.json", "")
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

	err = run()
	if err != nil {
		log.Fatal(err)
	}
}
