// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewStatistics is the golang structure for table review_statistics.
type ReviewStatistics struct {
	Id              int64       `json:"id"                description:""`
	TargetType      int         `json:"target_type"       description:"统计目标类型 1-SPU 2-商家"`
	TargetId        int64       `json:"target_id"         description:"目标ID（SPU_ID或Merchant_ID）"`
	Rating1Count    int         `json:"rating_1_count"    description:"1星数量"`
	Rating2Count    int         `json:"rating_2_count"    description:"2星数量"`
	Rating3Count    int         `json:"rating_3_count"    description:"3星数量"`
	Rating4Count    int         `json:"rating_4_count"    description:"4星数量"`
	Rating5Count    int         `json:"rating_5_count"    description:"5星数量"`
	TotalCount      int         `json:"total_count"       description:"总评价数"`
	AvgRating       float64     `json:"avg_rating"        description:"平均评分"`
	GoodRate        float64     `json:"good_rate"         description:"好评率（%）"`
	HasMediaCount   int         `json:"has_media_count"   description:"带图评价数"`
	HasContentCount int         `json:"has_content_count" description:"有内容评价数"`
	LastUpdatedAt   *gtime.Time `json:"last_updated_at"   description:""`
}
