package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kaneshin/gate/cmd/internal"
)

var (
	target = flag.String("target", "", "comma-separeted available")
	code   = flag.Bool("code", false, "")
)

func run() error {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	text := strings.TrimSpace(string(b))
	if text != "" {
		if *target == "" {
			*target = internal.Config.DefaultTarget
		}

		if *code {
			text = fmt.Sprintf("```\n%s\n```", text)
		}

		val := url.Values{
			"target": strings.Split(*target, ","),
			"text":   []string{text},
		}

		url := fmt.Sprintf("%s:%d", internal.Config.Env.Host, internal.Config.Env.Port)
		resp, err := http.PostForm(url, val)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	}
	return nil
}

func main() {
	err := internal.Load()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	err = run()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
