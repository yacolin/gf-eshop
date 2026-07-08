package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

// ---------- List ----------
type ListReq struct {
	g.Meta `path:"/departments" tags:"Departments" method:"get" summary:"部门列表（平铺）"`

	Page     int    `json:"page"`      // 页码，默认1
	PageSize int    `json:"page_size"` // 每页条数，默认20
	ParentId int64  `json:"parent_id"` // 按父级ID筛选
	Name     string `json:"name"`      // 按名称模糊搜索
	Status   *int   `json:"status"`    // 按状态筛选
}
type ListRes struct {
	List  []*entity.Departments `json:"list"`
	Total int                   `json:"total"`
}

// ---------- All ----------
type AllReq struct {
	g.Meta `path:"/departments/all" tags:"Departments" method:"get" summary:"所有部门"`
}
type AllRes struct {
	List []*entity.Departments `json:"list"`
}

// ---------- Children ----------
type ChildrenReq struct {
	g.Meta `path:"/departments/{id}/children" tags:"Departments" method:"get" summary:"子部门"`
	Id     int64 `json:"id"`
}
type ChildrenRes struct {
	List []*entity.Departments `json:"list"`
}

// ---------- Tree ----------
type TreeReq struct {
	g.Meta `path:"/departments/tree" tags:"Departments" method:"get" summary:"部门树形结构"`
	Status *int `json:"status"`
}
type TreeItem struct {
	*entity.Departments
	Children []*TreeItem `json:"children"`
}
type TreeRes struct {
	Tree []*TreeItem `json:"tree"`
}

// ---------- Detail ----------
type DetailReq struct {
	g.Meta `path:"/departments/{id}" tags:"Departments" method:"get" summary:"部门详情"`
	Id     int64 `json:"id"`
}
type DetailRes struct {
	*entity.Departments
}

// ---------- Create ----------
type CreateReq struct {
	g.Meta `path:"/departments" tags:"Departments" method:"post" summary:"新增部门"`

	Name      string `json:"name"      v:"required|length:1,100" description:"部门名称"`
	ParentId  int64  `json:"parent_id"  description:"上级部门ID（0=根部门）"`
	SortOrder int    `json:"sort_order" description:"排序值"`
	Status    int    `json:"status"    description:"1-启用 0-禁用"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

// ---------- Update ----------
type UpdateReq struct {
	g.Meta `path:"/departments/{id}" tags:"Departments" method:"put" summary:"更新部门"`

	Id        int64  `json:"id"        v:"required"`
	Name      string `json:"name"      v:"length:1,100" description:"部门名称"`
	ParentId  int64  `json:"parent_id"  description:"上级部门ID（0=根部门）"`
	SortOrder int    `json:"sort_order" description:"排序值"`
	Status    int    `json:"status"    description:"1-启用 0-禁用"`
}
type UpdateRes struct{}

// ---------- Delete ----------
type DeleteReq struct {
	g.Meta `path:"/departments/{id}" tags:"Departments" method:"delete" summary:"删除部门"`
	Id     int64 `json:"id"`
}
type DeleteRes struct{}
