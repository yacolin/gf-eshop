package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/model/entity"
)

type PointsListReq struct {
	g.Meta    `path:"/user-points" tags:"UserPoints" method:"get" summary:"积分流水列表"`
	Page      int    `json:"page"`       // 页码
	PageSize  int    `json:"page_size"`  // 每页数量
	UserId    string `json:"user_id"`    // 用户ID筛选
	Source    string `json:"source"`     // 积分来源筛选：order-下单消费 review-评价 signin-签到 admin-管理员调整 refund-退款扣减 expire-过期清零
	StartTime string `json:"start_time"` // 开始时间
	EndTime   string `json:"end_time"`   // 结束时间
}

type PointsListRes struct {
	List  []*entity.Points `json:"list"`
	Total int              `json:"total"`
}

type PointsBalanceReq struct {
	g.Meta `path:"/user-points/balance" tags:"UserPoints" method:"get" summary:"用户积分余额"`
	UserId string `json:"user_id"` // 用户ID
}

type PointsBalanceRes struct {
	UserId  int64 `json:"user_id"`  // 用户ID
	Balance int64 `json:"balance"`  // 当前积分余额
	TotalIn int64 `json:"total_in"` // 累计获得积分
	TotalOut int64 `json:"total_out"` // 累计消耗积分
}

type PointsTrendReq struct {
	g.Meta    `path:"/user-points/trend" tags:"UserPoints" method:"get" summary:"积分趋势"`
	UserId    string `json:"user_id"`    // 用户ID
	StartTime string `json:"start_time"` // 开始时间
	EndTime   string `json:"end_time"`   // 结束时间
}

type PointsTrendItem struct {
	Date   string `json:"date"`   // 日期
	Points int64  `json:"points"` // 积分变动
	Balance int64 `json:"balance"` // 余额
}

type PointsTrendRes struct {
	List []*PointsTrendItem `json:"list"`
}

type PointsAdjustReq struct {
	g.Meta  `path:"/user-points/adjust" tags:"UserPoints" method:"post" summary:"手动调整积分"`
	UserId  int64  `json:"user_id"  v:"required#用户ID不能为空"` // 用户ID
	Points  int64  `json:"points"   v:"required#积分值不能为空"` // 积分变动（正=增加，负=扣减）
	Remark  string `json:"remark"`                            // 备注说明
}

type PointsAdjustRes struct {
	Id           int64 `json:"id"`            // 流水ID
	BalanceAfter int64 `json:"balance_after"` // 变动后积分余额
}
