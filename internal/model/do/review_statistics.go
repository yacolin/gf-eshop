// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewStatistics is the golang structure of table rev_review_statistics for DAO operations like Where/Data.
type ReviewStatistics struct {
	g.Meta          `orm:"table:rev_review_statistics, do:true"`
	Id              interface{} //
	TargetType      interface{} // 统计目标类型 1-SPU 2-商家
	TargetId        interface{} // 目标ID（SPU_ID或Merchant_ID）
	Rating1Count    interface{} // 1星数量
	Rating2Count    interface{} // 2星数量
	Rating3Count    interface{} // 3星数量
	Rating4Count    interface{} // 4星数量
	Rating5Count    interface{} // 5星数量
	TotalCount      interface{} // 总评价数
	AvgRating       interface{} // 平均评分
	GoodRate        interface{} // 好评率（%）
	HasMediaCount   interface{} // 带图评价数
	HasContentCount interface{} // 有内容评价数
	LastUpdatedAt   *gtime.Time //
}
