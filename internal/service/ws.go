package service

import (
	"gf-eshop/internal/ws"
)

var localWsHub *ws.Hub

func WsHub() *ws.Hub {
	if localWsHub == nil {
		panic("implement not found for WsHub, forgot register?")
	}
	return localWsHub
}

func RegisterWsHub(hub *ws.Hub) {
	localWsHub = hub
}
