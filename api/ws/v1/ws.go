package v1

import "github.com/gogf/gf/v2/frame/g"

type WsStatsReq struct {
	g.Meta `path:"/stats" tags:"WS" method:"get" summary:"WebSocket 在线统计"`
}
type WsStatsRes struct {
	OnlineUsers int `json:"online_users"`
	Connections int `json:"connections"`
}

type WsSessionReq struct {
	g.Meta `path:"/session" tags:"WS" method:"get" summary:"当前用户会话信息"`
}
type WsSessionRes struct {
	Exists         bool   `json:"exists"`
	LastSeq        int64  `json:"last_seq"`
	ConnectedAt    string `json:"connected_at"`
	LastActiveAt   string `json:"last_active_at"`
	ReconnectCount int    `json:"reconnect_count"`
}
