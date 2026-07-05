package ws

import (
	"gf-eshop/api/ws"
)

type ControllerV1 struct{}

func NewV1() ws.IWsV1 {
	return &ControllerV1{}
}
