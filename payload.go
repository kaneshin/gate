package gate

type (
	// TextPayload represents a text payload.
	TextPayload struct {
		Text string `json:"text"`
	}

	// MessagePayload represents a message payload.
	MessagePayload struct {
		Message string `json:"message"`
	}

	// GraphPayload represents a graph payload.
	GraphPayload struct {
		ID       string `json:"id"`
		Date     string `json:"date"`
		Quantity string `json:"quantity"`
	}
)
