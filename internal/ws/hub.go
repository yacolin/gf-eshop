package ws

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

type UserBrief struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type UserEventPayload struct {
	Action    string    `json:"action"`
	User      UserBrief `json:"user"`
	Timestamp int64     `json:"timestamp"`
}

type StatsPayload struct {
	OnlineUsers int `json:"online_users"`
	Connections int `json:"connections"`
}

type Hub struct {
	mu         sync.RWMutex
	clients    map[int64]map[*Client]bool
	register   chan *Client
	unregister chan *Client

	msgCache   *MessageCache
	sessionMgr *SessionManager
	globalSeq  atomic.Int64

	shutdownCh chan struct{}
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]map[*Client]bool),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
		msgCache:   NewMessageCache(),
		sessionMgr: NewSessionManager(),
		shutdownCh: make(chan struct{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.handleRegister(client)
		case client := <-h.unregister:
			h.handleUnregister(client)
		case <-h.shutdownCh:
			h.handleShutdown()
			return
		}
	}
}

func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	wasOffline := h.clients[client.UserID] == nil
	if h.clients[client.UserID] == nil {
		h.clients[client.UserID] = make(map[*Client]bool)
	}
	h.clients[client.UserID][client] = true
	count := len(h.clients[client.UserID])
	h.mu.Unlock()

	go h.sendWelcomeMessage(client)
	go h.broadcastStats()
	if wasOffline {
		go h.broadcastUserEvent(client, "online")
	}
	_ = count
}

func (h *Hub) handleUnregister(client *Client) {
	h.mu.Lock()
	isOffline := false
	if _, ok := h.clients[client.UserID]; ok {
		if _, exists := h.clients[client.UserID][client]; exists {
			delete(h.clients[client.UserID], client)
			close(client.Send)
			if len(h.clients[client.UserID]) == 0 {
				delete(h.clients, client.UserID)
				isOffline = true
			}
		}
	}
	h.mu.Unlock()

	if isOffline {
		go h.broadcastUserEvent(client, "offline")
	}
	go h.broadcastStats()
	go h.sessionMgr.UpdateLastSeq(client.UserID, client.LastSeq)
}

func (h *Hub) handleShutdown() {
	h.mu.Lock()
	for userID, conns := range h.clients {
		for client := range conns {
			close(client.Send)
			client.Conn.Close()
		}
		delete(h.clients, userID)
	}
	h.mu.Unlock()
}

func (h *Hub) snapshotClients() []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	total := 0
	for _, clients := range h.clients {
		total += len(clients)
	}
	snapshot := make([]*Client, 0, total)
	for _, clients := range h.clients {
		for c := range clients {
			snapshot = append(snapshot, c)
		}
	}
	return snapshot
}

func (h *Hub) snapshotUserClients(userID int64) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clients := h.clients[userID]
	snapshot := make([]*Client, 0, len(clients))
	for c := range clients {
		snapshot = append(snapshot, c)
	}
	return snapshot
}

func (h *Hub) sendToClient(client *Client, data []byte) {
	select {
	case client.Send <- data:
		if id := extractSeqID(data); id > 0 {
			client.LastSeq = id
		}
	default:
		select {
		case h.unregister <- client:
		default:
		}
	}
}

func extractSeqID(data []byte) int64 {
	var msg PushMessage
	if err := msg.Unmarshal(data); err == nil {
		return msg.SequenceID
	}
	return 0
}

func (h *Hub) SendToUser(userID int64, data []byte) {
	for _, client := range h.snapshotUserClients(userID) {
		h.sendToClient(client, data)
	}
}

func (h *Hub) Broadcast(data []byte) {
	clients := h.snapshotClients()
	g, _ := errgroup.WithContext(context.Background())
	for _, client := range clients {
		client := client
		g.Go(func() error {
			select {
			case client.Send <- data:
			default:
				select {
				case h.unregister <- client:
				default:
				}
			}
			return nil
		})
	}
	g.Wait()
}

func (h *Hub) broadcastSafe(data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, clients := range h.clients {
		for client := range clients {
			select {
			case client.Send <- data:
			default:
			}
		}
	}
}

func (h *Hub) GetOnlineCount() (userCount int, connCount int) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	connCount = 0
	for _, clients := range h.clients {
		connCount += len(clients)
	}
	return len(h.clients), connCount
}

func (h *Hub) broadcastUserEvent(client *Client, action string) {
	msg := &RealtimeMessage{
		Seq:  h.globalSeq.Add(1),
		Type: "user",
		Payload: UserEventPayload{
			Action: action,
			User: UserBrief{
				ID:       client.UserID,
				Username: "",
				Nickname: "",
			},
			Timestamp: time.Now().UnixMilli(),
		},
	}
	data, err := msg.Marshal()
	if err != nil {
		return
	}
	h.broadcastSafe(data)
}

func (h *Hub) broadcastStats() {
	userCount, connCount := h.GetOnlineCount()
	msg := &RealtimeMessage{
		Seq:  h.globalSeq.Add(1),
		Type: "stats",
		Payload: StatsPayload{
			OnlineUsers: userCount,
			Connections: connCount,
		},
	}
	data, err := msg.Marshal()
	if err != nil {
		return
	}
	h.broadcastSafe(data)
}

func (h *Hub) SyncMessages(userID int64, lastSeq int64) {
	currentSeq, err := h.msgCache.GetCurrentSeqID(userID)
	if err != nil {
		return
	}
	if lastSeq >= currentSeq {
		return
	}
	cachedMinSeq, cachedMaxSeq, err := h.msgCache.GetCachedSeqRange(userID)
	if err != nil {
		return
	}
	_ = cachedMaxSeq
	if lastSeq < cachedMinSeq {
		h.sendFullSyncRequired(userID)
		return
	}
	messages, err := h.msgCache.GetMessages(userID, lastSeq, currentSeq)
	if err != nil {
		return
	}
	for _, msg := range messages {
		h.SendToUser(userID, msg)
	}
	go h.sessionMgr.IncrementReconnectCount(userID)
	go h.sessionMgr.UpdateLastSeq(userID, currentSeq)
}

func (h *Hub) GetUserLastSeq(userID int64) (int64, error) {
	session, err := h.sessionMgr.GetSession(userID)
	if err != nil {
		return 0, err
	}
	if session == nil {
		return 0, nil
	}
	return session.LastSeq, nil
}

func (h *Hub) PushToUser(userID int64, eventType string, data interface{}) error {
	seqID, err := h.msgCache.NextSeqID(userID)
	if err != nil {
		return err
	}
	msg := &PushMessage{
		Type:       eventType,
		SequenceID: seqID,
		Timestamp:  time.Now().UnixMilli(),
		Data:       data,
	}
	msgJSON, err := msg.Marshal()
	if err != nil {
		return err
	}
	if err := h.msgCache.StoreMessage(userID, seqID, msgJSON); err != nil {
		return err
	}
	h.SendToUser(userID, msgJSON)
	go h.sessionMgr.UpdateLastSeq(userID, seqID)
	return nil
}

func (h *Hub) sendWelcomeMessage(client *Client) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	currentSeq, _ := h.msgCache.GetCurrentSeqID(client.UserID)
	msg := NewSystemMessage("welcome", "连接成功", false)
	msg.SequenceID = currentSeq
	data, _ := msg.Marshal()
	select {
	case client.Send <- data:
	default:
	}
}

func (h *Hub) sendFullSyncRequired(userID int64) {
	msg := NewSystemMessage("sync_required", "缺失消息超过缓存窗口，请先获取全量数据", true)
	data, _ := msg.Marshal()
	h.SendToUser(userID, data)
}

func (h *Hub) Register() chan<- *Client {
	return h.register
}

func (h *Hub) NextSeq() int64 {
	return h.globalSeq.Add(1)
}

func (h *Hub) Shutdown(ctx context.Context) error {
	select {
	case h.shutdownCh <- struct{}{}:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
