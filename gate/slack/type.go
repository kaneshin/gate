package slack

import (
	"encoding/json"
)

type (
	// Payload represents a slack incoming webhook request data.
	Payload struct {
		Channel     string       `json:"channel"`
		Username    string       `json:"username"`
		Text        string       `json:"text"`
		IconEmoji   string       `json:"icon_emoji"`
		Attachments []Attachment `json:"attachments,omitempty"`
	}

	// AttachmentField contains information for an attachment field
	// An Attachment can contain multiple of these
	AttachmentField struct {
		Title string `json:"title"`
		Value string `json:"value"`
		Short bool   `json:"short"`
	}

	// AttachmentAction is a button to be included in the attachment. Required when
	// using message buttons and otherwise not useful. A maximum of 5 actions may be
	// provided per attachment.
	AttachmentAction struct {
		Name    string              `json:"name"`
		Text    string              `json:"text"`
		Style   string              `json:"style,omitempty"`
		Type    string              `json:"type"`
		Value   string              `json:"value,omitempty"`
		Confirm []ConfirmationField `json:"confirm,omitempty"`
	}

	// ConfirmationField are used to ask users to confirm actions
	ConfirmationField struct {
		Title       string `json:"title,omitempty"`
		Text        string `json:"text"`
		OkText      string `json:"ok_text,omitempty"`
		DismissText string `json:"dismiss_text,omitempty"`
	}

	// Attachment contains all the information for an attachment
	Attachment struct {
		Color    string `json:"color,omitempty"`
		Fallback string `json:"fallback"`

		CallbackID string `json:"callback_id,omitempty"`

		AuthorName    string `json:"author_name,omitempty"`
		AuthorSubname string `json:"author_subname,omitempty"`
		AuthorLink    string `json:"author_link,omitempty"`
		AuthorIcon    string `json:"author_icon,omitempty"`

		Title     string `json:"title,omitempty"`
		TitleLink string `json:"title_link,omitempty"`
		Pretext   string `json:"pretext,omitempty"`
		Text      string `json:"text"`

		ImageURL string `json:"image_url,omitempty"`
		ThumbURL string `json:"thumb_url,omitempty"`

		Fields     []AttachmentField  `json:"fields,omitempty"`
		Actions    []AttachmentAction `json:"actions,omitempty"`
		MarkdownIn []string           `json:"mrkdwn_in,omitempty"`

		Footer     string `json:"footer,omitempty"`
		FooterIcon string `json:"footer_icon,omitempty"`

		Ts json.Number `json:"ts,omitempty"`
	}
)
