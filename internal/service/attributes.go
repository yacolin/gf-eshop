package service

import (
	"context"

	"gf-eshop/api/attributes/v1"
)

type IAttributes interface {
	List(ctx context.Context, req *v1.AttributesListReq) (res *v1.AttributesListRes, err error)
	ListSearchable(ctx context.Context, req *v1.AttributesListSearchableReq) (res *v1.AttributesListSearchableRes, err error)
	ListSkuSpec(ctx context.Context, req *v1.AttributesListSkuSpecReq) (res *v1.AttributesListSkuSpecRes, err error)
	Detail(ctx context.Context, req *v1.AttributesDetailReq) (res *v1.AttributesDetailRes, err error)
	Create(ctx context.Context, req *v1.AttributesCreateReq) (res *v1.AttributesCreateRes, err error)
	Update(ctx context.Context, req *v1.AttributesUpdateReq) (res *v1.AttributesUpdateRes, err error)
	Delete(ctx context.Context, req *v1.AttributesDeleteReq) (res *v1.AttributesDeleteRes, err error)
}

var localAttributes IAttributes

func Attributes() IAttributes {
	if localAttributes == nil {
		panic("implement not found for interface IAttributes, forgot register?")
	}
	return localAttributes
}

func RegisterAttributes(i IAttributes) {
	localAttributes = i
}
