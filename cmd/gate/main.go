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
	"github.com/kaneshin/gate/gate/slack"
)

var (
	configPath = flag.String("config", "$HOME/.gate.tml", "")
	port       = flag.Int("port", 8080, "")
)

var config = struct {
	Slack struct {
		Incoming struct {
			URL       string `toml:"url"`
			Channel   string `toml:"channel"`
			Username  string `toml:"username"`
			IconEmoji string `toml:"icon_emoji"`
		} `toml:"incoming"`
	} `toml:"slack"`
	LINE struct {
		Notify struct {
			AccessToken string `toml:"access_token"`
		} `toml:"notify"`
	} `toml:"line"`
	Facebook struct {
		Messenger struct {
			ID          string `toml:"id"`
			AccessToken string `toml:"access_token"`
		} `toml:"messenger"`
	} `toml:"facebook"`
}{}

var (
	slackSvc    *gate.SlackIncomingService
	lineSvc     *gate.LINENotifyService
	facebookSvc *gate.FacebookMessengerService
)

func main() {
	flag.Parse()

	fp := os.ExpandEnv(*configPath)
	if _, err := toml.DecodeFile(fp, &config); err != nil {
		log.Fatal(err)
		return
	}

	httpClient := http.DefaultClient

	{
		var str string
		if str = config.Slack.Incoming.URL; str == "" {
			str = os.Getenv("SLACK_INCOMING_URL")
		}
		if str != "" {
			slackSvc = gate.NewSlackIncomingService(
				gate.NewConfig().WithHTTPClient(httpClient),
			).WithBaseURL(str)
		}
	}

	{
		var str string
		if str = config.LINE.Notify.AccessToken; str == "" {
			str = os.Getenv("LINE_NOTIFY_ACCESS_TOKEN")
		}
		lineSvc = gate.NewLINENotifyService(
			gate.NewConfig().WithHTTPClient(httpClient).WithAccessToken(str),
		)
	}

	{
		var str string
		if str = config.Facebook.Messenger.AccessToken; str == "" {
			str = os.Getenv("FACEBOOK_MESSENGER_ACCESS_TOKEN")
		}
		facebookSvc = gate.NewFacebookMessengerService(
			gate.NewConfig().WithHTTPClient(httpClient).WithAccessToken(str),
		)
	}

	os.Exit(run())
}

func run() int {

	const (
		slackIncoming = iota
		lineNotify
		facebookMessenger
	)

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()

		var (
			message = r.FormValue("message")
			color   = r.FormValue("color")
			image   = r.FormValue("image")
		)

		{
			var (
				ch    = r.FormValue("slack.channel")
				name  = r.FormValue("slack.username")
				emoji = r.FormValue("slack.emoji")
			)

			if ch == "" {
				ch = config.Slack.Incoming.Channel
			}
			payload := slackSvc.NewPayload(ch, message)
			payload.Username = config.Slack.Incoming.Username
			payload.IconEmoji = config.Slack.Incoming.IconEmoji

			if name != "" {
				payload.Username = name
			}
			if emoji != "" {
				payload.IconEmoji = emoji
			}

			if color != "" {
				att := slack.Attachment{
					Color:   color,
					Pretext: message,
					Text:    message,
				}
				payload.Attachments = append(payload.Attachments, att)
			}

			if image != "" {
				att := slack.Attachment{
					ImageURL: image,
				}
				payload.Attachments = append(payload.Attachments, att)
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
			var (
				id = r.FormValue("facebook.id")
			)
			if id == "" {
				id = config.Facebook.Messenger.ID
			}
			payload := facebookSvc.NewPayload(id, message)
			payload.NotificationType = facebook.NotificationTypeRegular
			if _, err := facebookSvc.Post(payload); err != nil {
				log.Printf("error %v", err)
			}
		}
	}).Methods("POST")

	log.Panic(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
	return 0
}
