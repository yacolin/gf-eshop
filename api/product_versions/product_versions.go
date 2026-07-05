package productVersions

import (
	"context"

	"gf-eshop/api/product_versions/v1"
)

type IProductVersionsV1 interface {
	List(ctx context.Context, req *v1.ProductVersionsListReq) (res *v1.ProductVersionsListRes, err error)
	Detail(ctx context.Context, req *v1.ProductVersionsDetailReq) (res *v1.ProductVersionsDetailRes, err error)
}
