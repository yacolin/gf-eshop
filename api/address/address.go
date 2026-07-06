package address

import (
	"context"
	"gf-eshop/api/address/v1"
)

type IAddressV1 interface {
	Create(ctx context.Context, req *v1.AddressCreateReq) (res *v1.AddressCreateRes, err error)
	List(ctx context.Context, req *v1.AddressListReq) (res *v1.AddressListRes, err error)
	Detail(ctx context.Context, req *v1.AddressDetailReq) (res *v1.AddressDetailRes, err error)
	GetDefault(ctx context.Context, req *v1.AddressGetDefaultReq) (res *v1.AddressGetDefaultRes, err error)
	Update(ctx context.Context, req *v1.AddressUpdateReq) (res *v1.AddressUpdateRes, err error)
	Delete(ctx context.Context, req *v1.AddressDeleteReq) (res *v1.AddressDeleteRes, err error)
}
