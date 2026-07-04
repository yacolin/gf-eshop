package service

import (
	"context"

	"gf-eshop/api/categories/v1"
)

type ICategories interface {
	List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error)
	All(ctx context.Context, req *v1.AllReq) (res *v1.AllRes, err error)
	Root(ctx context.Context, req *v1.RootReq) (res *v1.RootRes, err error)
	Children(ctx context.Context, req *v1.ChildrenReq) (res *v1.ChildrenRes, err error)
	Level(ctx context.Context, req *v1.LevelReq) (res *v1.LevelRes, err error)
	Tree(ctx context.Context, req *v1.TreeReq) (res *v1.TreeRes, err error)
	Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error)
	Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error)
	Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error)
	Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error)
}

var localCategories ICategories

func Categories() ICategories {
	if localCategories == nil {
		panic("implement not found for interface ICategories, forgot register?")
	}
	return localCategories
}

func RegisterCategories(i ICategories) {
	localCategories = i
}
