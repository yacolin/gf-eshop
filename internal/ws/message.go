package ws

import "encoding/json"

type PushMessage struct {
	Type       string      `json:"type"`
	SequenceID int64       `json:"sequence_id"`
	Timestamp  int64       `json:"timestamp"`
	Data       interface{} `json:"data"`
}

func (m *PushMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *PushMessage) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}

type RealtimeMessage struct {
	Seq     int64       `json:"seq"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (m *RealtimeMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

type SystemMessage struct {
	Type            string `json:"type"`
	SequenceID      int64  `json:"sequence_id"`
	Message         string `json:"message"`
	RequireFullSync bool   `json:"require_full_sync"`
}

func (m *SystemMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func NewSystemMessage(msgType string, message string, requireFullSync bool) *SystemMessage {
	return &SystemMessage{
		Type:            msgType,
		SequenceID:      0,
		Message:         message,
		RequireFullSync: requireFullSync,
	}
}

type ClientMessage struct {
	Type    string `json:"type"`
	Seq     int64  `json:"seq"`
	LastSeq int64  `json:"last_seq"`
}
