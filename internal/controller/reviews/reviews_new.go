package reviews

import (
	"gf-eshop/api/reviews"
)

type ControllerV1 struct{}

func NewV1() reviews.IReviewsV1 {
	return &ControllerV1{}
}
