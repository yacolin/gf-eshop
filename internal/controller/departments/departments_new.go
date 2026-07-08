package departments

import (
	"gf-eshop/api/departments"
)

type ControllerV1 struct{}

func NewV1() departments.IDepartmentsV1 {
	return &ControllerV1{}
}
