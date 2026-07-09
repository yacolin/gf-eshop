package service

import (
	"context"

	"gf-eshop/api/attribute_values/v1"
)

type IAttributeValues interface {
	List(ctx context.Context, req *v1.AttributeValuesListReq) (res *v1.AttributeValuesListRes, err error)
	Detail(ctx context.Context, req *v1.AttributeValuesDetailReq) (res *v1.AttributeValuesDetailRes, err error)
	Create(ctx context.Context, req *v1.AttributeValuesCreateReq) (res *v1.AttributeValuesCreateRes, err error)
	Update(ctx context.Context, req *v1.AttributeValuesUpdateReq) (res *v1.AttributeValuesUpdateRes, err error)
	Delete(ctx context.Context, req *v1.AttributeValuesDeleteReq) (res *v1.AttributeValuesDeleteRes, err error)
	ListByAttr(ctx context.Context, req *v1.AttributeValuesListByAttrReq) (res *v1.AttributeValuesListByAttrRes, err error)
}

var localAttributeValues IAttributeValues

func AttributeValues() IAttributeValues {
	if localAttributeValues == nil {
		panic("implement not found for interface IAttributeValues, forgot register?")
	}
	return localAttributeValues
}

func RegisterAttributeValues(i IAttributeValues) {
	localAttributeValues = i
}
