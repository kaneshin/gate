package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/kaneshin/gate"
	"github.com/kaneshin/gate/cmd/internal"
	"github.com/kaneshin/gate/facebook"
	"github.com/kaneshin/gate/slack"
)

var (
	slackIncomingServices    map[string]*gate.SlackIncomingService
	lineNotifyService        *gate.LINENotifyService
	facebookMessengerService *gate.FacebookMessengerService
)

const slackIncomingDefaultKey = "github.com/kaneshin/gate/incoming_default"

func handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	if r.Method != http.MethodPost {
		return
	}

	var err error
	message := r.FormValue("message")
	color := r.FormValue("color")
	image := r.FormValue("image")

	if slackIncomingServices != nil {
		ch := r.FormValue("slack.channel")
		name := r.FormValue("slack.username")
		emoji := r.FormValue("slack.emoji")

		key := ""
		if ch == "" {
			if len(internal.Config.Slack.App.Incoming) > 0 {
				ch = internal.Config.Slack.App.Incoming[0].Channel
				key = ch
			} else {
				ch = internal.Config.Slack.Incoming.Channel
				key = slackIncomingDefaultKey
			}
		} else {
			for _, v := range internal.Config.Slack.App.Incoming {
				if ch == v.Channel {
					key = ch
					break
				}
			}
			if key == "" {
				key = slackIncomingDefaultKey
			}
		}
		if name != "" || emoji != "" {
			key = slackIncomingDefaultKey
		}
		svc, ok := slackIncomingServices[key]
		if !ok {
			goto finish_slack
		}

		payload := svc.NewPayload(ch, message)
		payload.Username = internal.Config.Slack.Incoming.Username
		payload.IconEmoji = internal.Config.Slack.Incoming.IconEmoji

		if name != "" {
			payload.Username = name
		}
		if emoji != "" {
			payload.IconEmoji = emoji
		}

		if color != "" {
			att := slack.Attachment{
				Color: color,
				Text:  message,
			}
			payload.Text = ""
			payload.Attachments = append(payload.Attachments, att)
		}

		if image != "" {
			att := slack.Attachment{
				ImageURL: image,
			}
			payload.Attachments = append(payload.Attachments, att)
		}

		_, err = svc.Post(payload)
		if err != nil {
			log.Printf("[ERROR] %v", err)
		}
	}
finish_slack:

	if lineNotifyService != nil {
		_, err = lineNotifyService.Post(message)
		if err != nil {
			log.Printf("[ERROR] %v", err)
		}
	}

	if facebookMessengerService != nil {
		id := internal.Config.Facebook.Messenger.ID
		payload := facebookMessengerService.NewPayload(id, message)
		payload.NotificationType = facebook.NotificationTypeRegular
		_, err = facebookMessengerService.Post(payload)
		if err != nil {
			log.Printf("[ERROR] %v", err)
		}
	}
}

func newSlackIncomingServices(hcl *http.Client) map[string]*gate.SlackIncomingService {
	svc := map[string]*gate.SlackIncomingService{}
	url := internal.Config.Slack.Incoming.URL
	if url == "" {
		url = os.Getenv("SLACK_INCOMING_URL")
	}
	if url != "" {
		svc[slackIncomingDefaultKey] = gate.NewSlackIncomingService(
			gate.NewConfig().WithHTTPClient(hcl),
		).WithBaseURL(url)
	}

	for _, v := range internal.Config.Slack.App.Incoming {
		if v.URL == "" || v.Channel == "" {
			continue
		}
		svc[v.Channel] = gate.NewSlackIncomingService(
			gate.NewConfig().WithHTTPClient(hcl),
		).WithBaseURL(v.URL)
	}
	return svc
}

func newLINENotifyService(hcl *http.Client) *gate.LINENotifyService {
	token := internal.Config.LINE.Notify.AccessToken
	if token == "" {
		token = os.Getenv("LINE_NOTIFY_ACCESS_TOKEN")
	}
	if token == "" {
		return nil
	}
	conf := gate.NewConfig().WithHTTPClient(hcl).WithAccessToken(token)
	return gate.NewLINENotifyService(conf)
}

func newFacebookMessengerService(hcl *http.Client) *gate.FacebookMessengerService {
	token := internal.Config.Facebook.Messenger.AccessToken
	if token == "" {
		token = os.Getenv("FACEBOOK_MESSENGER_ACCESS_TOKEN")
	}
	if token == "" {
		return nil
	}
	conf := gate.NewConfig().WithHTTPClient(hcl).WithAccessToken(token)
	return gate.NewFacebookMessengerService(conf)
}

func main() {
	internal.ParseFlag()

	sigc := make(chan os.Signal)
	internal.Trap(sigc, map[syscall.Signal]func(os.Signal){
		syscall.SIGINT: func(sig os.Signal) {
			log.Printf("interrupt: %v", sig)
			os.Exit(1)
		},
		syscall.SIGTERM: func(sig os.Signal) {
			log.Printf("interrupt: %v", sig)
			os.Exit(1)
		},
	})

	hcl := &http.Client{}
	slackIncomingServices = newSlackIncomingServices(hcl)
	lineNotifyService = newLINENotifyService(hcl)
	facebookMessengerService = newFacebookMessengerService(hcl)

	http.HandleFunc("/", handler)

	port := fmt.Sprintf(":%d", internal.Config.Gate.Port)
	log.Fatal(http.ListenAndServe(port, http.DefaultServeMux))
}
