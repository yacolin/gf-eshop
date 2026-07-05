package ws

import (
	"context"

	"gf-eshop/api/ws/v1"
)

type IWsV1 interface {
	Stats(ctx context.Context, req *v1.WsStatsReq) (res *v1.WsStatsRes, err error)
	Session(ctx context.Context, req *v1.WsSessionReq) (res *v1.WsSessionRes, err error)
}
