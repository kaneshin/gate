package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	host    = flag.String("host", "http://localhost:8080", "")
	channel = flag.String("c", "", "")
)

func main() {
	flag.Parse()

	// Execute: echo "foo" | go run main.go
	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := http.PostForm(
		*host,
		url.Values{
			"message":       {string(body)},
			"slack.channel": {*channel},
		},
	); err != nil {
		log.Fatal(err)
	}
}
