package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

// ---------- List ----------
type MerchantRolesListReq struct {
	g.Meta `path:"/merchant-roles" tags:"MerchantRoles" method:"get" summary:"商家角色列表"`

	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	MerchantId int64  `json:"merchant_id" description:"商家ID，可选筛选"`
	Name       string `json:"name"        description:"按角色名称模糊搜索"`
	Status     *int   `json:"status"      description:"按状态筛选"`
}
type MerchantRolesListRes struct {
	List  []*entity.MerchantRoles `json:"list"`
	Total int                `json:"total"`
}

// ---------- Detail ----------
type MerchantRolesDetailReq struct {
	g.Meta `path:"/merchant-roles/{id}" tags:"MerchantRoles" method:"get" summary:"角色详情"`
	Id     int64 `json:"id"`
}
type MerchantRolesDetailRes struct {
	*entity.MerchantRoles
}

// ---------- Create ----------
type MerchantRolesCreateReq struct {
	g.Meta `path:"/merchant-roles" tags:"MerchantRoles" method:"post" summary:"新增角色"`

	MerchantId  int64  `json:"merchant_id"  v:"required" description:"商家ID"`
	Name        string `json:"name"         v:"required|length:1,64" description:"角色名称标识"`
	DisplayName string `json:"display_name" v:"required|length:1,64" description:"角色显示名称"`
	Description string `json:"description"  description:"角色描述"`
	RoleType    string `json:"role_type"    description:"builtin-系统预置 custom-商家自定义"`
	SortOrder   int    `json:"sort_order"   description:"排序值"`
	Status      int    `json:"status"       description:"1-启用 0-禁用"`
}
type MerchantRolesCreateRes struct {
	Id int64 `json:"id"`
}

// ---------- Update ----------
type MerchantRolesUpdateReq struct {
	g.Meta `path:"/merchant-roles/{id}" tags:"MerchantRoles" method:"put" summary:"更新角色"`

	Id          int64  `json:"id"          v:"required"`
	Name        string `json:"name"        v:"length:1,64" description:"角色名称标识"`
	DisplayName string `json:"display_name" v:"length:1,64" description:"角色显示名称"`
	Description string `json:"description" description:"角色描述"`
	RoleType    string `json:"role_type"   description:"builtin-系统预置 custom-商家自定义"`
	SortOrder   int    `json:"sort_order"  description:"排序值"`
	Status      int    `json:"status"      description:"1-启用 0-禁用"`
}
type MerchantRolesUpdateRes struct{}

// ---------- Delete ----------
type MerchantRolesDeleteReq struct {
	g.Meta `path:"/merchant-roles/{id}" tags:"MerchantRoles" method:"delete" summary:"删除角色"`
	Id     int64 `json:"id"`
}
type MerchantRolesDeleteRes struct{}
