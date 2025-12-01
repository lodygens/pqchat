package protocol

type HelloMessage struct {
	Type   string `json:"type"`
	Pseudo string `json:"pseudo"`
	UserID string `json:"user_id"`
	Pub    string `json:"pub"` // base64
	Sig    string `json:"sig"` // base64
}

type ChatMessage struct {
	Type      string   `json:"type"`
	From      string   `json:"from"`
	To        []string `json:"to"`
	Body      string   `json:"body"`
	Timestamp int64    `json:"timestamp"`
	Sig       string   `json:"sig"`
	Pub       string   `json:"pub"`
}
