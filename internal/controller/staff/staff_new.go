package staff

import (
	"gf-eshop/api/staff"
)

type ControllerV1 struct{}

func NewV1() staff.IStaffV1 {
	return &ControllerV1{}
}
