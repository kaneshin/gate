package line

// NotifyAPIURL is Notify API URL for LINE.
const NotifyAPIURL = "https://notify-api.line.me/api/notify"

type (
	// Payload represents a line request data.
	Payload struct {
		Message string `json:"message"`
	}
)
