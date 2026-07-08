// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Addresses is the golang structure for table addresses.
type Addresses struct {
	Id        int64       `json:"id"         description:"主键"`
	UserId    int64       `json:"user_id"    description:"用户ID"`
	Consignee string      `json:"consignee"  description:"收货人姓名"`
	Phone     string      `json:"phone"      description:"联系电话"`
	Country   string      `json:"country"    description:"国家"`
	Province  string      `json:"province"   description:"省"`
	City      string      `json:"city"       description:"市"`
	District  string      `json:"district"   description:"区/县"`
	Detail    string      `json:"detail"     description:"详细地址"`
	ZipCode   string      `json:"zip_code"   description:"邮编"`
	Tag       string      `json:"tag"        description:"地址标签：home/office/company/other"`
	IsDefault int         `json:"is_default" description:"是否默认地址（NULL=非默认, 1=默认）"`
	CreatedAt *gtime.Time `json:"created_at" description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updated_at" description:"更新时间"`
	DeletedAt *gtime.Time `json:"deleted_at" description:"删除时间"`
}
