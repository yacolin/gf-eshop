package notification

import (
	"gf-eshop/api/notification"
)

type ControllerV1 struct{}

func NewV1() notification.INotificationV1 {
	return &ControllerV1{}
}
