package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

// ---------- List ----------
type MerchantRolePermissionsListReq struct {
	g.Meta `path:"/merchant-role-permissions" tags:"MerchantRolePermissions" method:"get" summary:"角色权限列表"`

	RoleId int64 `json:"role_id" v:"required" description:"角色ID"`
}
type MerchantRolePermissionsListRes struct {
	List []*entity.MerchantRolePermissions `json:"list"`
}

// ---------- Update ----------
type MerchantRolePermissionsUpdateReq struct {
	g.Meta `path:"/merchant-role-permissions/{role_id}" tags:"MerchantRolePermissions" method:"put" summary:"批量更新角色权限"`

	RoleId          int64    `json:"role_id"       v:"required" description:"角色ID"`
	MerchantId      int64    `json:"merchant_id"   description:"商家ID"`
	PermissionNames []string `json:"permission_names" description:"权限标识列表"`
}
type MerchantRolePermissionsUpdateRes struct{}

// ---------- Create ----------
type MerchantRolePermissionsCreateReq struct {
	g.Meta `path:"/merchant-role-permissions" tags:"MerchantRolePermissions" method:"post" summary:"新增角色权限"`

	MerchantId     int64  `json:"merchant_id"     v:"required" description:"商家ID"`
	RoleId         int64  `json:"role_id"         v:"required" description:"角色ID"`
	PermissionName string `json:"permission_name" v:"required" description:"权限标识"`
}
type MerchantRolePermissionsCreateRes struct{}

// ---------- Delete ----------
type MerchantRolePermissionsDeleteReq struct {
	g.Meta `path:"/merchant-role-permissions/{id}" tags:"MerchantRolePermissions" method:"delete" summary:"删除角色权限"`
	Id     int64 `json:"id"`
}
type MerchantRolePermissionsDeleteRes struct{}
