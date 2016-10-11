package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"

	"github.com/kaneshin/gate/gate"
	"github.com/kaneshin/gate/gate/facebook"
)

type (
	config struct {
		Slack struct {
			IncomingURL string `toml:"incoming_url"`
			Channel     string `toml:"channel"`
			Username    string `toml:"username"`
			IconEmoji   string `toml:"icon_emoji"`
		} `toml:"slack"`
		LINE struct {
			AccessToken string `toml:"access_token"`
		} `toml:"line"`
		Facebook struct {
			Messenger struct {
				ID          string `toml:"id"`
				AccessToken string `toml:"access_token"`
			} `toml:"messenger"`
		} `toml:"facebook"`
	}
)

var (
	configPath = flag.String("config", "$HOME/.gate.toml", "")
	port       = flag.Int("port", 8080, "")
)

func main() {
	flag.Parse()

	var c config
	if _, err := toml.DecodeFile(os.ExpandEnv(*configPath), &c); err != nil {
		log.Fatal(err)
		return
	}

	slackSvc := gate.NewSlackIncomingService(
		gate.NewConfig().
			WithHTTPClient(http.DefaultClient),
	).WithBaseURL(c.Slack.IncomingURL)

	lineSvc := gate.NewLINEService(
		gate.NewConfig().
			WithHTTPClient(http.DefaultClient).
			WithAccessToken(c.LINE.AccessToken),
	)

	facebookSvc := gate.NewFacebookService(
		gate.NewConfig().
			WithHTTPClient(http.DefaultClient).
			WithAccessToken(c.Facebook.Messenger.AccessToken),
	)

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		message := r.FormValue("message")

		{
			var channel string
			if ch := r.FormValue("slack.channel"); ch != "" {
				channel = ch
			} else if ch := c.Slack.Channel; ch != "" {
				channel = ch
			}

			payload := slackSvc.NewPayload(channel, message)
			if t := c.Slack.Username; t != "" {
				payload.Username = t
			}
			if t := c.Slack.IconEmoji; t != "" {
				payload.IconEmoji = t
			}

			if _, err := slackSvc.Post(payload); err != nil {
				log.Printf("error %v", err)
			}
		}

		{
			if _, err := lineSvc.Post(message); err != nil {
				log.Printf("error %v", err)
			}
		}

		{
			var id string
			if t := r.FormValue("facebook.id"); t != "" {
				id = t
			} else if t := c.Facebook.Messenger.ID; t != "" {
				id = t
			}
			payload := facebookSvc.NewPayload(id, message)
			payload.NotificationType = facebook.NotificationTypeRegular
			if _, err := facebookSvc.Post(payload); err != nil {
				log.Printf("error %v", err)
			}
		}
	}).Methods("POST")

	http.ListenAndServe(fmt.Sprintf(":%d", *port), r)
}
