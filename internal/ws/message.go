package ws

import "encoding/json"

// RealtimeMessage 实时推送消息体，与 monolith RealtimeMessage 对齐。
// 前端 useWebSocket hook 期望 { type, payload } 格式。
type RealtimeMessage struct {
	Seq     int64       `json:"seq,omitempty"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// newEnvelope 创建扁平 { type, payload } 消息
func newEnvelope(msgType string, payload interface{}) []byte {
	b, _ := json.Marshal(&RealtimeMessage{
		Type:    msgType,
		Payload: payload,
	})
	return b
}

// newEnvelopeWithSeq 创建带 seq 的扁平消息（用于 PushToUser 等需要序列号的场景）
func newEnvelopeWithSeq(msgType string, payload interface{}, seq int64) []byte {
	b, _ := json.Marshal(&RealtimeMessage{
		Seq:     seq,
		Type:    msgType,
		Payload: payload,
	})
	return b
}

// ── 广播/实时推送 Payload（不含 type 字段，type 在 RealtimeMessage.Type） ──

type StatsPayload struct {
	OnlineUsers int `json:"online_users"`
	Connections int `json:"connections"`
}

type UserEventPayload struct {
	Action    string `json:"action"`
	UserID    int64  `json:"user_id"`
	Timestamp int64  `json:"timestamp"`
}

type WelcomePayload struct {
	SequenceID      int64 `json:"sequence_id"`
	RequireFullSync bool  `json:"require_full_sync"`
}

type SyncRequiredPayload struct {
	RequireFullSync bool `json:"require_full_sync"`
}

type NotificationPayload struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Level   string `json:"level"`
}

type PongPayload struct {
	Seq     int64 `json:"seq,omitempty"`
	LastSeq int64 `json:"last_seq,omitempty"`
}

// ── 客户端上行消息 ──

type ClientMessage struct {
	Type    string `json:"type"`
	Seq     int64  `json:"seq"`
	LastSeq int64  `json:"last_seq"`
}
