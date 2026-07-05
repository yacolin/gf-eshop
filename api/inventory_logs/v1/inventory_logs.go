package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type InventoryLogsListReq struct {
	g.Meta `path:"/inventory_logs" tags:"InventoryLogs" method:"get" summary:"库存流水列表"`

	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	SkuId    int64  `json:"sku_id"`
	ChangeType string `json:"change_type"`
}
type InventoryLogsListRes struct {
	List  []*entity.InventoryLogs `json:"list"`
	Total int                     `json:"total"`
}

type InventoryLogsDetailReq struct {
	g.Meta `path:"/inventory_logs/{id}" tags:"InventoryLogs" method:"get" summary:"流水详情"`
	Id     int64 `json:"id"`
}
type InventoryLogsDetailRes struct {
	*entity.InventoryLogs
}
