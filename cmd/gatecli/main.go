package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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

func main() {
	internal.Load()

	// Execute: echo "foo" | go run main.go
	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	str := strings.TrimSpace(string(body))

	if *target == "" {
		*target = internal.Config.DefaultTarget
	}
	val := url.Values{
		"target": strings.Split(*target, ","),
	}

	if body != nil {
		if *code {
			val.Set("text", "```"+str+"```")
		} else {
			val.Set("text", str)
		}

		url := fmt.Sprintf("%s:%d", internal.Config.Env.Host, internal.Config.Env.Port)
		resp, err := http.PostForm(url, val)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))
	}
}
