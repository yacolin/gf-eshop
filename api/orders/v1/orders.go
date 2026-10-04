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
	g.Meta        `path:"/orders" tags:"Orders" method:"get" summary:"订单列表"`
	Page          int    `json:"page"           description:"页码（与 cursor 二选一；offset 分页，仅不分表时精确）"`
	PageSize      int    `json:"page_size"      description:"每页条数（配合 page）"`
	Cursor        string `json:"cursor"         description:"游标：base64(上一页末位订单ID)，排序固定 id DESC。传 cursor 或 size 即走游标分页"`
	Size          int    `json:"size"           description:"游标分页每页条数（默认 20，上限 100）"`
	UserID        int64  `json:"user_id"        description:"用户ID"`
	Status        string `json:"status"          description:"订单状态"`
	PaymentStatus string `json:"payment_status"  description:"支付状态"`
	OrderNo       string `json:"order_no"       description:"订单号"`
}

type OrdersListRes struct {
	List  []*entity.Orders `json:"list"`
	Total int              `json:"total"` // 游标分页时为 -1，表示未统计（分表后 COUNT 需要跨片）
	// NextCursor 下一页游标；为空表示没有更多。仅在游标分页模式下返回。
	NextCursor string `json:"next_cursor,omitempty"`
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