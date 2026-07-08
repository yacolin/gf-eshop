// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Merchants is the golang structure of table mch_merchants for DAO operations like Where/Data.
type Merchants struct {
	g.Meta          `orm:"table:mch_merchants, do:true"`
	Id              interface{} // 商家ID
	MerchantName    interface{} // 商家名称（店铺名）
	MerchantCode    interface{} // 商家编码（系统生成，唯一）
	MerchantType    interface{} // 1-个人商家 2-企业商家 3-品牌直营
	MerchantLevel   interface{} // 商家等级 1-普通 2-银牌 3-金牌 4-钻石（影响佣金率/权限）
	BusinessScope   interface{} // 经营范围
	BusinessYears   interface{} // 经营年限（入驻年限）
	ContactPerson   interface{} // 主要联系人
	ContactPhone    interface{} // 联系电话
	ContactEmail    interface{} // 联系邮箱
	LogoUrl         interface{} // 店铺Logo
	BannerUrl       interface{} // 店铺Banner图
	ShopDescription interface{} // 店铺简介
	Status          interface{} // 0-待审核 1-正常 2-冻结 3-已注销
	AuditStatus     interface{} // 0-待审核 1-审核通过 2-审核拒绝
	AuditReason     interface{} // 审核拒绝原因
	AuditedAt       *gtime.Time // 审核时间
	FrozenReason    interface{} // 冻结原因
	CommissionRate  interface{} // 平台抽佣比例（千分比，如 50 表示5%）
	SettlementCycle interface{} // 结算周期 1-T+1 2-T+7 3-月结
	TotalOrders     interface{} // 历史总订单数
	TotalSales      interface{} // 历史总销售额（分）
	AvgRating       interface{} // 店铺平均评分
	ProductCount    interface{} // 在售商品数量
	SettledAt       *gtime.Time // 入驻时间
	ExpireAt        *gtime.Time // 合同到期时间
	CreatedBy       interface{} // 创建人
	UpdatedBy       interface{} // 更新人
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
	DeletedAt       *gtime.Time //
}
