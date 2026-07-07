package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type PermissionListReq struct {
	g.Meta `path:"/permissions" tags:"Permissions" method:"get" summary:"权限列表"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	Resource   string `json:"resource"`
	Action     string `json:"action"`
	Category   string `json:"category"`
	Status     *int   `json:"status"`
}
type PermissionListRes struct {
	List  []*entity.Permissions `json:"list"`
	Total int                   `json:"total"`
}

type PermissionDetailReq struct {
	g.Meta `path:"/permissions/{id}" tags:"Permissions" method:"get" summary:"权限详情"`
	Id     int64 `json:"id"`
}
type PermissionDetailRes struct {
	*entity.Permissions
}

type PermissionCreateReq struct {
	g.Meta `path:"/permissions" tags:"Permissions" method:"post" summary:"新增权限"`
	Name        string `json:"name"         v:"required|length:1,100" description:"权限标识"`
	DisplayName string `json:"display_name" v:"required|length:1,100" description:"显示名称"`
	Description string `json:"description"                           description:"描述"`
	Resource    string `json:"resource"     v:"required|length:1,50"  description:"资源"`
	Action      string `json:"action"       v:"required|length:1,50"  description:"操作"`
	Category    string `json:"category"                               description:"分类"`
	SortOrder   int    `json:"sort_order"                             description:"排序"`
	Status      int    `json:"status"                                 description:"状态"`
}
type PermissionCreateRes struct {
	Id int64 `json:"id"`
}

type PermissionUpdateReq struct {
	g.Meta `path:"/permissions/{id}" tags:"Permissions" method:"put" summary:"更新权限"`
	Id          int64  `json:"id"          v:"required"`
	Name        string `json:"name"        v:"length:1,100"    description:"权限标识"`
	DisplayName string `json:"display_name" v:"length:1,100"   description:"显示名称"`
	Description string `json:"description"                    description:"描述"`
	Resource    string `json:"resource"    v:"length:1,50"     description:"资源"`
	Action      string `json:"action"      v:"length:1,50"     description:"操作"`
	Category    string `json:"category"                        description:"分类"`
	SortOrder   int    `json:"sort_order"                      description:"排序"`
	Status      int    `json:"status"                          description:"状态"`
}
type PermissionUpdateRes struct{}

type PermissionDeleteReq struct {
	g.Meta `path:"/permissions/{id}" tags:"Permissions" method:"delete" summary:"删除权限"`
	Id     int64 `json:"id"`
}
type PermissionDeleteRes struct{}

type PermissionCheckReq struct {
	g.Meta   `path:"/permissions/check" tags:"Permissions" method:"post" summary:"校验当前用户权限"`
	Permission string `json:"permission" v:"required" description:"权限标识，如 order:create"`
}
type PermissionCheckRes struct {
	HasPermission bool `json:"has_permission"`
}

type RolePermissionListReq struct {
	g.Meta `path:"/permissions/roles/{role_id}" tags:"Permissions" method:"get" summary:"角色权限列表"`
	RoleId int64 `json:"role_id"`
}
type RolePermissionListRes struct {
	List []*entity.Permissions `json:"list"`
}

type RolePermissionUpdateReq struct {
	g.Meta      `path:"/permissions/roles/{role_id}" tags:"Permissions" method:"put" summary:"替换角色权限"`
	RoleId      int64   `json:"role_id"`
	PermissionIds []int64 `json:"permission_ids" v:"required"`
}
type RolePermissionUpdateRes struct{}
