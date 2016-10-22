package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/kaneshin/gate/cmd/internal"
)

var (
	color = flag.String("color", "", "")
	image = flag.String("image", "", "")

	channel  = flag.String("channel", "", "")
	username = flag.String("username", "", "")
	emoji    = flag.String("emoji", "", "")

	code = flag.Bool("code", false, "")
)

var re = regexp.MustCompile("^https?.*\\.(png|jpg|jpeg|gif)($|\\?)")

func main() {
	flag.Parse()

	// Execute: echo "foo" | go run main.go
	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	str := strings.TrimSpace(string(body))

	if list := re.FindAllString(str, -1); len(list) > 0 {
		if *image == "" {
			*image = str
			body = nil
		}
	}

	val := url.Values{
		"color":          {*color},
		"image":          {*image},
		"slack.channel":  {*channel},
		"slack.username": {*username},
		"slack.emoji":    {*emoji},
	}

	if body != nil {
		if *code {
			val.Set("message", "```"+str+"```")
		} else {
			val.Set("message", str)
		}
	}

	if _, err := http.PostForm(
		fmt.Sprintf("%s:%d", internal.Config.Gate.Host, internal.Config.Gate.Port),
		val,
	); err != nil {
		log.Fatal(err)
	}
}
