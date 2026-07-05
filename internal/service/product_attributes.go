package service

import (
	"context"

	"gf-eshop/api/product_attributes/v1"
)

type IProductAttributes interface {
	List(ctx context.Context, req *v1.ProductAttributesListReq) (res *v1.ProductAttributesListRes, err error)
	Create(ctx context.Context, req *v1.ProductAttributesCreateReq) (res *v1.ProductAttributesCreateRes, err error)
	Delete(ctx context.Context, req *v1.ProductAttributesDeleteReq) (res *v1.ProductAttributesDeleteRes, err error)
}

var localProductAttributes IProductAttributes

func ProductAttributes() IProductAttributes {
	if localProductAttributes == nil {
		panic("implement not found for interface IProductAttributes, forgot register?")
	}
	return localProductAttributes
}

func RegisterProductAttributes(i IProductAttributes) {
	localProductAttributes = i
}
