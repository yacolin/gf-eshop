package ws

import (
	"encoding/json"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

const (
	writeWait       = 10 * time.Second
	pingPeriod      = 30 * time.Second
	pongWait        = 90 * time.Second
	sendBufferSize  = 64
	maxPingFailures = 2
)

type Client struct {
	Hub           *Hub
	Conn          *ghttp.WebSocket
	Send          chan []byte
	UserID        int64
	LastSeq       int64
	closed        bool
	pingFailCount int
	lastPingTime  time.Time
	pingTicker    *time.Ticker
}

func NewClient(hub *Hub, conn *ghttp.WebSocket, userID int64, lastSeq int64) *Client {
	return &Client{
		Hub:         hub,
		Conn:        conn,
		Send:        make(chan []byte, sendBufferSize),
		UserID:      userID,
		LastSeq:     lastSeq,
		pingTicker:  time.NewTicker(pingPeriod),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.pingTicker.Stop()
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(appData string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		c.pingFailCount = 0
		c.updateLastSeqFromPong(appData)
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		c.handleClientMessage(message)
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.pingTicker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(ghttp.WsMsgClose, []byte{})
				return
			}
			w, err := c.Conn.NextWriter(ghttp.WsMsgText)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.Send)
			}
			if err := w.Close(); err != nil {
				return
			}

		case <-c.pingTicker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			payload := struct {
				Timestamp int64 `json:"timestamp"`
				LastSeq   int64 `json:"last_seq"`
			}{
				Timestamp: time.Now().UnixMilli(),
				LastSeq:   c.LastSeq,
			}
			data, _ := json.Marshal(payload)
			if err := c.Conn.WriteMessage(ghttp.WsMsgPing, data); err != nil {
				c.pingFailCount++
				if c.pingFailCount >= maxPingFailures {
					return
				}
			} else {
				c.pingFailCount = 0
				c.lastPingTime = time.Now()
			}
		}
	}
}

func (c *Client) handleClientMessage(message []byte) {
	var msg ClientMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		return
	}
	switch msg.Type {
	case "ping":
		c.sendPong(msg.Seq)
	case "sync":
		c.Hub.SyncMessages(c.UserID, msg.LastSeq)
	case "reconnect":
		c.LastSeq = msg.LastSeq
		c.Hub.SyncMessages(c.UserID, msg.LastSeq)
	}
}

func (c *Client) sendPong(seq int64) {
	data := newEnvelope(&PongPayload{
		Type:    "pong",
		Seq:     seq,
		LastSeq: c.LastSeq,
	})
	select {
	case c.Send <- data:
	default:
	}
}

func (c *Client) updateLastSeqFromPong(appData string) {
	if appData == "" {
		return
	}
	var payload struct {
		LastSeq int64 `json:"last_seq"`
	}
	if err := json.Unmarshal([]byte(appData), &payload); err == nil && payload.LastSeq > c.LastSeq {
		c.LastSeq = payload.LastSeq
	}
}
