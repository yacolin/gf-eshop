// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AfterSaleEvidences is the golang structure of table tx_after_sale_evidences for DAO operations like Where/Data.
type AfterSaleEvidences struct {
	g.Meta      `orm:"table:tx_after_sale_evidences, do:true"`
	Id          interface{} // 主键
	AfterSaleId interface{} // 售后单ID
	MediaType   interface{} // 1-图片 2-视频
	MediaUrl    interface{} // 凭证URL
	SortOrder   interface{} // 排序
	CreatedAt   *gtime.Time //
}
