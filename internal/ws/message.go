package ws

import "encoding/json"

const (
	messageCodeOK = 0
)

// WsEnvelope is the standard WS message wrapper, aligned with HTTP API response format.
type WsEnvelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func newEnvelope(data interface{}) []byte {
	payload, _ := json.Marshal(data)
	b, _ := json.Marshal(&WsEnvelope{Code: messageCodeOK, Message: "", Data: payload})
	return b
}

type PushPayload struct {
	Type       string      `json:"type"`
	SequenceID int64       `json:"sequence_id"`
	Timestamp  int64       `json:"timestamp"`
	Payload    interface{} `json:"payload"`
}

type StatsPayload struct {
	Type        string `json:"type"`
	OnlineUsers int    `json:"online_users"`
	Connections int    `json:"connections"`
}

type UserEventPayload struct {
	Type      string `json:"type"`
	Action    string `json:"action"`
	UserID    int64  `json:"user_id"`
	Timestamp int64  `json:"timestamp"`
}

type WelcomePayload struct {
	Type            string `json:"type"`
	SequenceID      int64  `json:"sequence_id"`
	RequireFullSync bool   `json:"require_full_sync"`
}

type SyncRequiredPayload struct {
	Type            string `json:"type"`
	RequireFullSync bool   `json:"require_full_sync"`
}

type PongPayload struct {
	Type    string `json:"type"`
	Seq     int64  `json:"seq,omitempty"`
	LastSeq int64  `json:"last_seq,omitempty"`
}

type ClientMessage struct {
	Type    string `json:"type"`
	Seq     int64  `json:"seq"`
	LastSeq int64  `json:"last_seq"`
}
