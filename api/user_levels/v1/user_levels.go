package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type LevelsListReq struct {
	g.Meta   `path:"/user-levels" tags:"UserLevels" method:"get" summary:"等级列表"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Name     string `json:"name"     description:"等级名称模糊搜索"`
	Status   *int   `json:"status"   description:"1-启用 0-禁用"`
}
type LevelsListRes struct {
	List  []*entity.Levels `json:"list"`
	Total int              `json:"total"`
}

type LevelsDetailReq struct {
	g.Meta `path:"/user-levels/{id}" tags:"UserLevels" method:"get" summary:"等级详情"`
	Id     int64 `json:"id"`
}
type LevelsDetailRes struct {
	*entity.Levels
}

type LevelsCreateReq struct {
	g.Meta           `path:"/user-levels" tags:"UserLevels" method:"post" summary:"新增等级"`
	Name             string  `json:"name"              v:"required" description:"等级名称"`
	Level            int     `json:"level"             v:"required" description:"等级数值"`
	MinPoints        int64   `json:"min_points"        description:"最低累计积分"`
	MaxPoints        int64   `json:"max_points"        description:"最高累计积分（0无上限）"`
	DiscountRate     int64   `json:"discount_rate"     description:"折扣率（千分比，1000=无折扣）"`
	FreeShipping     int     `json:"free_shipping"     description:"1-免运费"`
	PointsMultiplier float64 `json:"points_multiplier" description:"消费积分倍数"`
	Benefits         string  `json:"benefits"          description:"扩展权益JSON"`
	Status           int     `json:"status"            description:"1-启用 0-禁用"`
	SortOrder        int     `json:"sort_order"        description:"排序"`
}
type LevelsCreateRes struct {
	Id int64 `json:"id"`
}

type LevelsUpdateReq struct {
	g.Meta           `path:"/user-levels/{id}" tags:"UserLevels" method:"put" summary:"更新等级"`
	Id               int64   `json:"id"                v:"required"`
	Name             string  `json:"name"              description:"等级名称"`
	Level            int     `json:"level"             description:"等级数值"`
	MinPoints        int64   `json:"min_points"        description:"最低累计积分"`
	MaxPoints        int64   `json:"max_points"        description:"最高累计积分（0无上限）"`
	DiscountRate     int64   `json:"discount_rate"     description:"折扣率（千分比，1000=无折扣）"`
	FreeShipping     int     `json:"free_shipping"     description:"1-免运费"`
	PointsMultiplier float64 `json:"points_multiplier" description:"消费积分倍数"`
	Benefits         string  `json:"benefits"          description:"扩展权益JSON"`
	Status           int     `json:"status"            description:"1-启用 0-禁用"`
	SortOrder        int     `json:"sort_order"        description:"排序"`
}
type LevelsUpdateRes struct{}

type LevelsDeleteReq struct {
	g.Meta `path:"/user-levels/{id}" tags:"UserLevels" method:"delete" summary:"删除等级"`
	Id     int64 `json:"id"`
}
type LevelsDeleteRes struct{}
