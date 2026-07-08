package service

import (
	"context"

	"gf-eshop/api/departments/v1"
)

type IDepartments interface {
	List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error)
	All(ctx context.Context, req *v1.AllReq) (res *v1.AllRes, err error)
	Children(ctx context.Context, req *v1.ChildrenReq) (res *v1.ChildrenRes, err error)
	Tree(ctx context.Context, req *v1.TreeReq) (res *v1.TreeRes, err error)
	Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error)
	Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error)
	Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error)
	Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error)
}

var localDepartments IDepartments

func Departments() IDepartments {
	if localDepartments == nil {
		panic("implement not found for interface IDepartments, forgot register?")
	}
	return localDepartments
}

func RegisterDepartments(i IDepartments) {
	localDepartments = i
}
