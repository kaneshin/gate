package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kaneshin/gate/cmd/internal"
	"github.com/kaneshin/gate/gate"
	"github.com/kaneshin/gate/gate/facebook"
	"github.com/kaneshin/gate/gate/slack"
)

var (
	slackSvc    *gate.SlackIncomingService
	lineSvc     *gate.LINENotifyService
	facebookSvc *gate.FacebookMessengerService
)

func main() {
	internal.ParseFlag()

	httpClient := http.DefaultClient

	{
		var str string
		if str = internal.Config.Slack.Incoming.URL; str == "" {
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
		if str = internal.Config.LINE.Notify.AccessToken; str == "" {
			str = os.Getenv("LINE_NOTIFY_ACCESS_TOKEN")
		}
		lineSvc = gate.NewLINENotifyService(
			gate.NewConfig().WithHTTPClient(httpClient).WithAccessToken(str),
		)
	}

	{
		var str string
		if str = internal.Config.Facebook.Messenger.AccessToken; str == "" {
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

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
				ch = internal.Config.Slack.Incoming.Channel
			}
			payload := slackSvc.NewPayload(ch, message)
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
				id = internal.Config.Facebook.Messenger.ID
			}
			payload := facebookSvc.NewPayload(id, message)
			payload.NotificationType = facebook.NotificationTypeRegular
			if _, err := facebookSvc.Post(payload); err != nil {
				log.Printf("error %v", err)
			}
		}
	})

	log.Panic(http.ListenAndServe(fmt.Sprintf(":%d", internal.Config.Gate.Port), http.DefaultServeMux))
	return 0
}
