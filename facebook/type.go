package facebook

// SendAPIURL is Send API for Facebook.
const SendAPIURL = "https://graph.facebook.com/v2.6/me/messages"

const (
	// NotificationTypeRegular is.
	NotificationTypeRegular = "REGULAR"
	// NotificationTypeSilentPush is.
	NotificationTypeSilentPush = "SILENT_PUSH"
	// NotificationTypeNoPush is.
	NotificationTypeNoPush = "NO_PUSH"
)

type (
	// Payload represents a facebook message request data.
	Payload struct {
		Recipient        Recipient `json:"recipient"`
		Message          Message   `json:"message,omitempty"`
		SenderAction     string    `json:"sender_action,omitempty"`
		NotificationType string    `json:"notification_type,omitempty"`
	}

	// Recipient is.
	Recipient struct {
		PhoneNumber string `json:"phone_number,omitempty"`
		ID          string `json:"id,omitempty"`
	}

	// Message is.
	Message struct {
		Text string `json:"text,omitempty"`
	}
)
