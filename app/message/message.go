package message

type Message struct {
	Type    string      `json:"type"`
	Event   string      `json:"event,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Channel string      `json:"channel,omitempty"`
}
