package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type AddressCreateReq struct {
	g.Meta    `path:"/addresses" tags:"Address" method:"post" summary:"新增收货地址"`
	Consignee string `json:"consignee" v:"required|length:1,64" description:"收货人姓名"`
	Phone     string `json:"phone"     v:"required|length:1,20" description:"联系电话"`
	Country   string `json:"country"   description:"国家"`
	Province  string `json:"province"  v:"required|length:1,32" description:"省"`
	City      string `json:"city"      v:"required|length:1,32" description:"市"`
	District  string `json:"district"  v:"required|length:1,32" description:"区/县"`
	Detail    string `json:"detail"    v:"required|length:1,256" description:"详细地址"`
	ZipCode   string `json:"zip_code"  description:"邮编"`
	Tag       string `json:"tag"       description:"地址标签：home/office/company/other"`
	IsDefault int    `json:"is_default" description:"是否默认地址"`
}
type AddressCreateRes struct {
	Id int64 `json:"id"`
}

type AddressListReq struct {
	g.Meta `path:"/addresses" tags:"Address" method:"get" summary:"收货地址列表"`
}
type AddressListRes struct {
	List  []*entity.Addresses `json:"list"`
	Total int                 `json:"total"`
}

type AddressDetailReq struct {
	g.Meta `path:"/addresses/{id}" tags:"Address" method:"get" summary:"地址详情"`
	Id     int64 `json:"id"`
}
type AddressDetailRes struct {
	*entity.Addresses
}

type AddressGetDefaultReq struct {
	g.Meta `path:"/addresses/default" tags:"Address" method:"get" summary:"获取默认地址"`
}
type AddressGetDefaultRes struct {
	*entity.Addresses
}

type AddressUpdateReq struct {
	g.Meta    `path:"/addresses/{id}" tags:"Address" method:"put" summary:"更新收货地址"`
	Id        int64   `json:"id"          v:"required"`
	Consignee *string `json:"consignee" description:"收货人姓名"`
	Phone     *string `json:"phone"     description:"联系电话"`
	Country   *string `json:"country"   description:"国家"`
	Province  *string `json:"province"  description:"省"`
	City      *string `json:"city"      description:"市"`
	District  *string `json:"district"  description:"区/县"`
	Detail    *string `json:"detail"    description:"详细地址"`
	ZipCode   *string `json:"zip_code"  description:"邮编"`
	Tag       *string `json:"tag"       description:"地址标签"`
	IsDefault *int    `json:"is_default" description:"是否默认地址"`
}
type AddressUpdateRes struct{}

type AddressDeleteReq struct {
	g.Meta `path:"/addresses/{id}" tags:"Address" method:"delete" summary:"删除收货地址"`
	Id     int64 `json:"id"`
}
type AddressDeleteRes struct{}
