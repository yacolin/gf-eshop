package ws

import (
	"context"

	"gf-eshop/api/ws/v1"
	"gf-eshop/internal/service"
	"gf-eshop/utility"
)

func (c *ControllerV1) Stats(ctx context.Context, req *v1.WsStatsReq) (res *v1.WsStatsRes, err error) {
	userCount, connCount := service.WsHub().GetOnlineCount()
	return &v1.WsStatsRes{
		OnlineUsers: userCount,
		Connections: connCount,
	}, nil
}

func (c *ControllerV1) Session(ctx context.Context, req *v1.WsSessionReq) (res *v1.WsSessionRes, err error) {
	claims := utility.GetStaffClaims(ctx)
	if claims == nil {
		return &v1.WsSessionRes{Exists: false}, nil
	}
	session, err := service.WsHub().GetUserLastSeq(claims.StaffId)
	if err != nil {
		return nil, err
	}
	return &v1.WsSessionRes{
		Exists:  session > 0,
		LastSeq: session,
	}, nil
}
