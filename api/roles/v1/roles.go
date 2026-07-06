package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type RoleListReq struct {
	g.Meta `path:"/roles" tags:"Roles" method:"get" summary:"角色列表"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Name     string `json:"name"`
	RoleType string `json:"role_type"`
	Status   int    `json:"status"`
}
type RoleListRes struct {
	List  []*entity.Roles `json:"list"`
	Total int             `json:"total"`
}

type RoleDetailReq struct {
	g.Meta `path:"/roles/{id}" tags:"Roles" method:"get" summary:"角色详情"`
	Id     int64 `json:"id"`
}
type RoleDetailRes struct {
	*entity.Roles
}

type RoleCreateReq struct {
	g.Meta `path:"/roles" tags:"Roles" method:"post" summary:"新增角色"`
	Name        string `json:"name"         v:"required|length:1,50" description:"角色名称"`
	DisplayName string `json:"display_name" v:"required|length:1,100" description:"显示名称"`
	Description string `json:"description"                          description:"描述"`
	RoleType    string `json:"role_type"                             description:"类型"`
	SortOrder   int    `json:"sort_order"                            description:"排序"`
	Status      int    `json:"status"                                description:"状态"`
}
type RoleCreateRes struct {
	Id int64 `json:"id"`
}

type RoleUpdateReq struct {
	g.Meta `path:"/roles/{id}" tags:"Roles" method:"put" summary:"更新角色"`
	Id          int64  `json:"id"          v:"required"`
	Name        string `json:"name"        v:"length:1,50"  description:"角色名称"`
	DisplayName string `json:"display_name" v:"length:1,100" description:"显示名称"`
	Description string `json:"description"                  description:"描述"`
	RoleType    string `json:"role_type"                     description:"类型"`
	SortOrder   int    `json:"sort_order"                    description:"排序"`
	Status      int    `json:"status"                        description:"状态"`
}
type RoleUpdateRes struct{}

type RoleDeleteReq struct {
	g.Meta `path:"/roles/{id}" tags:"Roles" method:"delete" summary:"删除角色"`
	Id     int64 `json:"id"`
}
type RoleDeleteRes struct{}
