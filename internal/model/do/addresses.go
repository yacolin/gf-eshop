// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Addresses is the golang structure of table usr_addresses for DAO operations like Where/Data.
type Addresses struct {
	g.Meta    `orm:"table:usr_addresses, do:true"`
	Id        interface{} // 主键
	UserId    interface{} // 用户ID
	Consignee interface{} // 收货人姓名
	Phone     interface{} // 联系电话
	Country   interface{} // 国家
	Province  interface{} // 省
	City      interface{} // 市
	District  interface{} // 区/县
	Detail    interface{} // 详细地址
	ZipCode   interface{} // 邮编
	Tag       interface{} // 地址标签：home/office/company/other
	IsDefault interface{} // 是否默认地址（NULL=非默认, 1=默认）
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
	DeletedAt *gtime.Time // 删除时间
}
