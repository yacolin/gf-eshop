package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

// ---------- List ----------
type ListReq struct {
	g.Meta `path:"/categories" tags:"Categories" method:"get" summary:"类目列表（平铺）"`

	Page     int `json:"page"`     // 页码，默认1
	PageSize int `json:"pageSize"` // 每页条数，默认20
	ParentId int `json:"parentId"` // 按父级ID筛选
	Status   int `json:"status"`   // 按状态筛选
}
type ListRes struct {
	List  []*entity.Categories `json:"list"`
	Total int                  `json:"total"`
}

// ---------- All ----------
type AllReq struct {
	g.Meta `path:"/categories/all" tags:"Categories" method:"get" summary:"所有类目"`
}
type AllRes struct {
	List []*entity.Categories `json:"list"`
}

// ---------- Root ----------
type RootReq struct {
	g.Meta `path:"/categories/root" tags:"Categories" method:"get" summary:"根类目"`
}
type RootRes struct {
	List []*entity.Categories `json:"list"`
}


// ---------- Children ----------
type ChildrenReq struct {
	g.Meta `path:"/categories/{id}/children" tags:"Categories" method:"get" summary:"子类目"`
	Id     int64 `json:"id"`
}

type ChildrenRes struct {
	List []*entity.Categories `json:"list"`
}


// ---------- Level ----------
type LevelReq struct {
	g.Meta `path:"/categories/level/{level}" tags:"Categories" method:"get" summary:"层级类目"`
	Level     int64 `json:"level"`
}

type LevelRes struct {
	List []*entity.Categories `json:"list"`
}


// ---------- Tree ----------
type TreeReq struct {
	g.Meta `path:"/categories/tree" tags:"Categories" method:"get" summary:"类目树形结构"`
	Status int `json:"status"`
}
type TreeItem struct {
	*entity.Categories
	Children []*TreeItem `json:"children"`
}
type TreeRes struct {
	Tree []*TreeItem `json:"tree"`
}

// ---------- Detail ----------
type DetailReq struct {
	g.Meta `path:"/categories/{id}" tags:"Categories" method:"get" summary:"类目详情"`
	Id     int64 `json:"id"`
}
type DetailRes struct {
	*entity.Categories
}

// ---------- Create ----------
type CreateReq struct {
	g.Meta `path:"/categories" tags:"Categories" method:"post" summary:"新增类目"`

	Name      string `json:"name"      v:"required|length:1,100" description:"类目名称"`
	ParentId  int64  `json:"parentId"  description:"父级ID"`
	Level     int    `json:"level"     description:"层级"`
	Path      string `json:"path"      description:"路径"`
	IconUrl   string `json:"iconUrl"   description:"类目图标"`
	SortOrder int    `json:"sortOrder" description:"排序"`
	Status    int    `json:"status"    description:"状态"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

// ---------- Update ----------
type UpdateReq struct {
	g.Meta `path:"/categories/{id}" tags:"Categories" method:"put" summary:"更新类目"`

	Id        int64  `json:"id"        v:"required"`
	Name      string `json:"name"      v:"length:1,100" description:"类目名称"`
	ParentId  int64  `json:"parentId"  description:"父级ID"`
	Level     int    `json:"level"     description:"层级"`
	Path      string `json:"path"      description:"路径"`
	IconUrl   string `json:"iconUrl"   description:"类目图标"`
	SortOrder int    `json:"sortOrder" description:"排序"`
	Status    int    `json:"status"    description:"状态"`
}
type UpdateRes struct{}

// ---------- Delete ----------
type DeleteReq struct {
	g.Meta `path:"/categories/{id}" tags:"Categories" method:"delete" summary:"删除类目"`
	Id     int64 `json:"id"`
}
type DeleteRes struct{}
