package service

import (
	"context"

	"gf-eshop/api/category_attributes/v1"
)

type ICategoryAttributes interface {
	List(ctx context.Context, req *v1.CategoryAttributesListReq) (res *v1.CategoryAttributesListRes, err error)
	Create(ctx context.Context, req *v1.CategoryAttributesCreateReq) (res *v1.CategoryAttributesCreateRes, err error)
	BatchCreate(ctx context.Context, req *v1.CategoryAttributesBatchCreateReq) (res *v1.CategoryAttributesBatchCreateRes, err error)
	Delete(ctx context.Context, req *v1.CategoryAttributesDeleteReq) (res *v1.CategoryAttributesDeleteRes, err error)
	ListByCat(ctx context.Context, req *v1.CategoryAttributesListByCatReq) (res *v1.CategoryAttributesListByCatRes, err error)
}

var localCategoryAttributes ICategoryAttributes

func CategoryAttributes() ICategoryAttributes {
	if localCategoryAttributes == nil {
		panic("implement not found for interface ICategoryAttributes, forgot register?")
	}
	return localCategoryAttributes
}

func RegisterCategoryAttributes(i ICategoryAttributes) {
	localCategoryAttributes = i
}
