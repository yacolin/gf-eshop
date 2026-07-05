package ws

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

const (
	sessionKeyPrefix = "ws:session:"
	sessionExpire    = 7 * 24 * time.Hour
)

type Session struct {
	UserID         int64     `json:"user_id"`
	LastSeq        int64     `json:"last_seq"`
	ConnectedAt    time.Time `json:"connected_at"`
	LastActiveAt   time.Time `json:"last_active_at"`
	ReconnectCount int       `json:"reconnect_count"`
}

type SessionManager struct{}

func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

func (sm *SessionManager) SaveSession(session *Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(context.Background(), "SETEX", sessionKeyPrefix+itoa(session.UserID), int(sessionExpire.Seconds()), string(data))
	return err
}

func (sm *SessionManager) GetSession(userID int64) (*Session, error) {
	v, err := g.Redis().Do(context.Background(), "GET", sessionKeyPrefix+itoa(userID))
	if err != nil {
		return nil, err
	}
	if v.IsNil() {
		return nil, nil
	}
	var session Session
	if err := json.Unmarshal(v.Bytes(), &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (sm *SessionManager) UpdateLastSeq(userID int64, lastSeq int64) error {
	session, err := sm.GetSession(userID)
	if err != nil {
		return err
	}
	if session == nil {
		session = &Session{
			UserID:       userID,
			LastSeq:      lastSeq,
			ConnectedAt:  time.Now(),
			LastActiveAt: time.Now(),
		}
		return sm.SaveSession(session)
	}
	session.LastSeq = lastSeq
	session.LastActiveAt = time.Now()
	return sm.SaveSession(session)
}

func (sm *SessionManager) IncrementReconnectCount(userID int64) error {
	session, err := sm.GetSession(userID)
	if err != nil {
		return err
	}
	if session == nil {
		return nil
	}
	session.ReconnectCount++
	session.LastActiveAt = time.Now()
	return sm.SaveSession(session)
}

func (sm *SessionManager) DeleteSession(userID int64) error {
	_, err := g.Redis().Do(context.Background(), "DEL", sessionKeyPrefix+itoa(userID))
	return err
}

func itoa(n int64) string {
	return strconv.FormatInt(n, 10)
}
