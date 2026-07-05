package productVersions

import (
	"context"

	"gf-eshop/api/product_versions/v1"
	"gf-eshop/internal/service"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.ProductVersionsListReq) (res *v1.ProductVersionsListRes, err error) {
	return service.ProductVersions().List(ctx, req)
}

func (c *ControllerV1) Detail(ctx context.Context, req *v1.ProductVersionsDetailReq) (res *v1.ProductVersionsDetailRes, err error) {
	return service.ProductVersions().Detail(ctx, req)
}
