package service

import (
	"context"

	"gf-eshop/api/product_descriptions/v1"
)

type IProductDescriptions interface {
	Detail(ctx context.Context, req *v1.ProductDescriptionsDetailReq) (res *v1.ProductDescriptionsDetailRes, err error)
	Save(ctx context.Context, req *v1.ProductDescriptionsSaveReq) (res *v1.ProductDescriptionsSaveRes, err error)
}

var localProductDescriptions IProductDescriptions

func ProductDescriptions() IProductDescriptions {
	if localProductDescriptions == nil {
		panic("implement not found for interface IProductDescriptions, forgot register?")
	}
	return localProductDescriptions
}

func RegisterProductDescriptions(i IProductDescriptions) {
	localProductDescriptions = i
}
