package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"gf-eshop/internal/model/entity"
)

type OrdersCreateReq struct {
	g.Meta `path:"/orders" tags:"Orders" method:"post" summary:"创建订单"`
	UserId         int64             `json:"user_id"      description:"用户ID（留空则从登录上下文获取）"`
	Items          []CreateOrderItem `json:"items"        v:"required" description:"商品列表"`
	CouponID       *int64            `json:"coupon_id"    description:"优惠券ID"`
	BuyerRemark    string            `json:"buyer_remark" description:"买家备注" v:"max-length:500"`
	Source         string            `json:"source"       description:"来源" v:"max-length:20"`
	Consignee      string            `json:"consignee"    v:"required" description:"收货人"`
	Phone          string            `json:"phone"        v:"required" description:"手机号"`
	Province       string            `json:"province"     description:"省"`
	City           string            `json:"city"         description:"市"`
	District       string            `json:"district"     description:"区"`
	DetailAddr     string            `json:"detail_addr"  description:"详细地址"`
	ZipCode        string            `json:"zip_code"     description:"邮编"`
}

type CreateOrderItem struct {
	SkuID    int64 `json:"sku_id"   v:"required" description:"SKU ID"`
	Quantity int   `json:"quantity" v:"required|min:1|max:99" description:"数量"`
}

type OrdersCreateRes struct {
	*entity.Orders
}

type OrdersListReq struct {
	g.Meta `path:"/orders" tags:"Orders" method:"get" summary:"订单列表"`
	// 列表只有游标分页一种形态（原先的 page/page_size offset 分页已移除：
	// 分表后 offset 无法跨片正确归并，且 COUNT 要跨片 fan-out）。
	// 传了已移除的 page/page_size 会直接报参数错误，不会静默忽略。
	Cursor string `json:"cursor"         description:"游标：base64(上一页末位订单ID)，排序固定 id DESC。首次不传，后续传上一页的 next_cursor"`
	Size   int    `json:"size"           description:"每页条数（默认 20，上限 100）"`

	UserID        int64  `json:"user_id"        description:"用户ID"`
	Status        string `json:"status"          description:"订单状态"`
	PaymentStatus string `json:"payment_status"  description:"支付状态"`
	OrderNo       string `json:"order_no"       description:"订单号（点查：不受月份窗口限制）"`
	// 时间窗口。分表后列表必须带时间范围，否则要跨全部活跃分片 fan-out
	//（实测一次翻页 = 活跃分片数次查询）。优先级：month > created_from/to > 默认近 N 个月。
	Month       string `json:"month"        description:"下单月份，形如 202610（路由到单个分片，翻页更省）"`
	CreatedFrom string `json:"created_from" description:"下单日期起（含），形如 2026-08-01"`
	CreatedTo   string `json:"created_to"   description:"下单日期止（含），形如 2026-10-31"`
}

// OrdersListRes 列表响应。游标分页字段与 products 列表完全一致：
// list + next_cursor + has_more，便于前端复用同一套翻页逻辑。
//
// 刻意**不返回 total**：游标分页不做 COUNT（分表后 COUNT 要跨 36 片 fan-out），
// 返回 -1 之类的占位值容易被误解成「总共就这么多单」。
type OrdersListRes struct {
	List []*entity.Orders `json:"list"`
	// NextCursor 下一页游标；为空表示没有更多。
	NextCursor string `json:"next_cursor,omitempty"`
	// HasMore 是否还有下一页（多取一条判定，末页不会误报）。
	HasMore bool `json:"has_more"`
	// 以下四个字段回显**本次实际生效**的时间窗口：客户端（尤其前端列表）
	// 应把它显示出来，否则「默认只看近 N 个月」会造成「怎么少了一单」的困惑。
	AppliedFrom     string   `json:"applied_from,omitempty"     description:"生效窗口起点（含）"`
	AppliedTo       string   `json:"applied_to,omitempty"       description:"生效窗口终点（含）"`
	AppliedMonths   []string `json:"applied_months,omitempty"   description:"生效窗口命中的月份（YYYYMM），用于说明本次查询了哪些分片"`
	WindowDefaulted bool     `json:"window_defaulted,omitempty" description:"true 表示未传时间参数、由后端默认窗口兜底"`
}

type OrdersDetailReq struct {
	g.Meta  `path:"/orders/{order_no}" tags:"Orders" method:"get" summary:"订单详情"`
	OrderNo string `json:"order_no"`
}

type OrdersDetailRes struct {
	Order     *entity.Orders       `json:"order"`
	SubOrders []*entity.SubOrders  `json:"sub_orders"`
	Items     []*entity.OrderItems `json:"items"`
}

type OrdersUpdateStatusReq struct {
	g.Meta  `path:"/orders/{order_no}/status" tags:"Orders" method:"put" summary:"更新订单状态"`
	OrderNo string `json:"order_no"`
	Status  string `json:"status" v:"required|in:paid,cancelled,shipped,delivered,completed" description:"目标状态"`
	Note    string `json:"note"   description:"备注" v:"max-length:500"`
}

type OrdersUpdateStatusRes struct {
	*entity.Orders
}