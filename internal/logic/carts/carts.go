package carts

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gf-eshop/api/carts/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
	"gf-eshop/utility"
)

type sCarts struct {
	mu          sync.Mutex
	pendingSync map[int64]*time.Timer
}

func init() {
	service.RegisterCarts(&sCarts{
		pendingSync: make(map[int64]*time.Timer),
	})
}

const (
	cartCacheTTL = 86400 // 24h in seconds
)

func cartItemsKey(userID int64) string { return fmt.Sprintf("cart:%d:items", userID) }
func cartMetaKey(userID int64) string  { return fmt.Sprintf("cart:%d:meta", userID) }

type cartMeta struct {
	ItemCount   int   `json:"item_count"`
	TotalAmount int64 `json:"total_amount"`
}

type cachedCartItem struct {
	ProductID   int64  `json:"product_id"`
	ProductName string `json:"product_name"`
	SkuSpec     string `json:"sku_spec"`
	Image       string `json:"image"`
	Price       int64  `json:"price"`
	Quantity    int    `json:"quantity"`
}

// ── Redis CRUD ─────────────────────────────────

func readCartFromRedis(ctx context.Context, userID int64) (*v1.CartResponse, error) {
	metaVar, err := g.Redis().Do(ctx, "GET", cartMetaKey(userID))
	if err != nil {
		return nil, err
	}
	if metaVar.IsNil() {
		return nil, gerror.NewCode(gcode.CodeNotFound, "cart not in redis")
	}

	var meta cartMeta
	if err := sonic.Unmarshal(metaVar.Bytes(), &meta); err != nil {
		return nil, err
	}

	itemsVar, err := g.Redis().Do(ctx, "HGETALL", cartItemsKey(userID))
	if err != nil {
		return nil, err
	}

	itemsMap := itemsVar.MapStrStr()

	items := make([]v1.CartItemResponse, 0, len(itemsMap))
	for skuIDStr, itemJSON := range itemsMap {
		var ci cachedCartItem
		if err := sonic.Unmarshal([]byte(itemJSON), &ci); err != nil {
			continue
		}
		skuID, _ := strconv.ParseInt(skuIDStr, 10, 64)
		items = append(items, v1.CartItemResponse{
			SkuID:       skuID,
			ProductID:   ci.ProductID,
			ProductName: ci.ProductName,
			SkuSpec:     ci.SkuSpec,
			Image:       ci.Image,
			Price:       ci.Price,
			Quantity:    ci.Quantity,
			Subtotal:    ci.Price * int64(ci.Quantity),
		})
	}

	resp := &v1.CartResponse{
		ID:          userID,
		ItemCount:   meta.ItemCount,
		TotalAmount: meta.TotalAmount,
		Items:       items,
	}
	return resp, nil
}

func writeCartToRedis(ctx context.Context, userID int64, resp *v1.CartResponse) error {
	metaData, _ := sonic.Marshal(cartMeta{ItemCount: resp.ItemCount, TotalAmount: resp.TotalAmount})
	if _, err := g.Redis().Do(ctx, "SETEX", cartMetaKey(userID), cartCacheTTL, metaData); err != nil {
		return err
	}
	if _, err := g.Redis().Do(ctx, "DEL", cartItemsKey(userID)); err != nil {
		return err
	}
	for _, item := range resp.Items {
		data, _ := sonic.Marshal(cachedCartItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			SkuSpec:     item.SkuSpec,
			Image:       item.Image,
			Price:       item.Price,
			Quantity:    item.Quantity,
		})
		if _, err := g.Redis().Do(ctx, "HSET", cartItemsKey(userID), strconv.FormatInt(item.SkuID, 10), data); err != nil {
			return err
		}
	}
	if _, err := g.Redis().Do(ctx, "EXPIRE", cartItemsKey(userID), cartCacheTTL); err != nil {
		return err
	}
	return nil
}

func delCartFromRedis(ctx context.Context, userID int64) {
	_, _ = g.Redis().Do(ctx, "DEL", cartItemsKey(userID), cartMetaKey(userID))
}

// ── Helper ──────────────────────────────────────

func resolveCartUserID(ctx context.Context, reqUserID int64) int64 {
	if reqUserID > 0 {
		return reqUserID
	}
	if claims := utility.GetStaffClaims(ctx); claims != nil {
		return claims.StaffId
	}
	return 0
}

func buildCartResp(cart *entity.Carts, items []entity.CartItems) *v1.CartResponse {
	resp := &v1.CartResponse{
		ID:    cart.Id,
		Items: make([]v1.CartItemResponse, 0, len(items)),
	}
	for _, item := range items {
		resp.Items = append(resp.Items, v1.CartItemResponse{
			SkuID:       item.SkuId,
			ProductID:   item.ProductId,
			ProductName: item.ProductName,
			SkuSpec:     item.SkuSpec,
			Image:       item.Image,
			Price:       item.Price,
			Quantity:    item.Quantity,
			Subtotal:    item.Price * int64(item.Quantity),
		})
	}
	recalcCartResp(resp)
	return resp
}

func recalcCartResp(resp *v1.CartResponse) {
	resp.ItemCount = len(resp.Items)
	resp.TotalAmount = 0
	for _, item := range resp.Items {
		resp.TotalAmount += item.Subtotal
	}
}

// ── Service Methods ─────────────────────────────

func (s *sCarts) GetCart(ctx context.Context, req *v1.CartsGetReq) (res *v1.CartsGetRes, err error) {
	userID := resolveCartUserID(ctx, 0)
	if userID == 0 {
		return &v1.CartsGetRes{CartResponse: &v1.CartResponse{Items: make([]v1.CartItemResponse, 0)}}, nil
	}

	resp, err := readCartFromRedis(ctx, userID)
	if err == nil {
		return &v1.CartsGetRes{CartResponse: resp}, nil
	}

	// DB fallback
	cart, err := dao.Carts.Ctx(ctx).Where(dao.Carts.Columns().UserId, userID).One()
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return &v1.CartsGetRes{CartResponse: &v1.CartResponse{Items: make([]v1.CartItemResponse, 0)}}, nil
	}

	var items []entity.CartItems
	err = dao.CartItems.Ctx(ctx).Where(dao.CartItems.Columns().CartId, cart["id"].Int64()).Scan(&items)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = make([]entity.CartItems, 0)
	}

	var cartEnt entity.Carts
	if err := cart.Struct(&cartEnt); err != nil {
		return nil, err
	}

	resp = buildCartResp(&cartEnt, items)
	_ = writeCartToRedis(ctx, userID, resp)
	return &v1.CartsGetRes{CartResponse: resp}, nil
}

func (s *sCarts) AddItem(ctx context.Context, req *v1.CartsAddItemReq) (res *v1.CartsAddItemRes, err error) {
	userID := resolveCartUserID(ctx, 0)
	if userID == 0 {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "未登录")
	}

	// 校验 SKU
	var sku entity.Skus
	err = dao.Skus.Ctx(ctx).Where(dao.Skus.Columns().Id, req.SkuID).Scan(&sku)
	if err != nil {
		return nil, err
	}
	if sku.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "SKU不存在")
	}

	resp, err := s.loadOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 获取商品名称
	var product entity.Products
	_ = dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, sku.ProductId).Scan(&product)
	productName := sku.Spec
	if product.Id > 0 && product.Name != "" {
		productName = product.Name
	}

	found := false
	for i, item := range resp.Items {
		if item.SkuID == req.SkuID {
			resp.Items[i].Quantity += req.Quantity
			resp.Items[i].Price = sku.Price
			resp.Items[i].Subtotal = resp.Items[i].Price * int64(resp.Items[i].Quantity)
			found = true
			break
		}
	}
	if !found {
		resp.Items = append(resp.Items, v1.CartItemResponse{
			SkuID:       sku.Id,
			ProductID:   sku.ProductId,
			ProductName: productName,
			SkuSpec:     sku.Spec,
			Image:       sku.Image,
			Price:       sku.Price,
			Quantity:    req.Quantity,
			Subtotal:    sku.Price * int64(req.Quantity),
		})
	}

	recalcCartResp(resp)
	_ = writeCartToRedis(ctx, userID, resp)
	s.debounceSync(userID)
	return &v1.CartsAddItemRes{CartResponse: resp}, nil
}

func (s *sCarts) UpdateItem(ctx context.Context, req *v1.CartsUpdateItemReq) (res *v1.CartsUpdateItemRes, err error) {
	userID := resolveCartUserID(ctx, 0)
	if userID == 0 {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "未登录")
	}

	resp, err := s.loadOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.Quantity <= 0 {
		for i, item := range resp.Items {
			if item.SkuID == req.SkuID {
				resp.Items = append(resp.Items[:i], resp.Items[i+1:]...)
				break
			}
		}
	} else {
		for i, item := range resp.Items {
			if item.SkuID == req.SkuID {
				resp.Items[i].Quantity = req.Quantity
				resp.Items[i].Subtotal = item.Price * int64(req.Quantity)
				break
			}
		}
	}

	recalcCartResp(resp)
	_ = writeCartToRedis(ctx, userID, resp)
	s.debounceSync(userID)
	return &v1.CartsUpdateItemRes{CartResponse: resp}, nil
}

func (s *sCarts) RemoveItem(ctx context.Context, req *v1.CartsRemoveItemReq) (res *v1.CartsRemoveItemRes, err error) {
	userID := resolveCartUserID(ctx, 0)
	if userID == 0 {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "未登录")
	}

	resp, err := s.loadOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i, item := range resp.Items {
		if item.SkuID == req.SkuID {
			resp.Items = append(resp.Items[:i], resp.Items[i+1:]...)
			break
		}
	}

	recalcCartResp(resp)
	_ = writeCartToRedis(ctx, userID, resp)
	s.debounceSync(userID)
	return &v1.CartsRemoveItemRes{CartResponse: resp}, nil
}

func (s *sCarts) ClearCart(ctx context.Context, req *v1.CartsClearReq) (res *v1.CartsClearRes, err error) {
	userID := resolveCartUserID(ctx, 0)
	if userID == 0 {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "未登录")
	}

	delCartFromRedis(ctx, userID)

	cart, err := dao.Carts.Ctx(ctx).Where(dao.Carts.Columns().UserId, userID).One()
	if err != nil {
		return nil, err
	}
	if cart != nil {
		cartId := cart["id"].Int64()
		_, _ = dao.CartItems.Ctx(ctx).Where(dao.CartItems.Columns().CartId, cartId).Delete()
		_, _ = dao.Carts.Ctx(ctx).Where(dao.Carts.Columns().Id, cartId).Data(g.Map{
			"item_count":   0,
			"total_amount": 0,
			"updated_at":   gtime.Now(),
		}).Update()
	}
	return &v1.CartsClearRes{}, nil
}

// ── 辅助方法 ────────────────────────────────────

func (s *sCarts) loadOrCreateCart(ctx context.Context, userID int64) (*v1.CartResponse, error) {
	resp, err := readCartFromRedis(ctx, userID)
	if err == nil {
		return resp, nil
	}

	// Redis miss → DB 兜底
	cart, err := dao.Carts.Ctx(ctx).Where(dao.Carts.Columns().UserId, userID).One()
	if err != nil {
		return nil, err
	}

	if cart == nil {
		id, err := dao.Carts.Ctx(ctx).Data(g.Map{
			"user_id":      userID,
			"item_count":   0,
			"total_amount": 0,
			"created_at":   gtime.Now(),
		}).InsertAndGetId()
		if err != nil {
			return nil, err
		}
		resp := &v1.CartResponse{
			ID:          id,
			ItemCount:   0,
			TotalAmount: 0,
			Items:       make([]v1.CartItemResponse, 0),
		}
		_ = writeCartToRedis(ctx, userID, resp)
		return resp, nil
	}

	var items []entity.CartItems
	err = dao.CartItems.Ctx(ctx).Where(dao.CartItems.Columns().CartId, cart["id"].Int64()).Scan(&items)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = make([]entity.CartItems, 0)
	}

	var cartEnt entity.Carts
	if err := cart.Struct(&cartEnt); err != nil {
		return nil, err
	}

	resp = buildCartResp(&cartEnt, items)
	_ = writeCartToRedis(ctx, userID, resp)
	return resp, nil
}

// debounceSync 防抖延迟同步：5s 内无新写入再落库
func (s *sCarts) debounceSync(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if t, ok := s.pendingSync[userID]; ok {
		t.Stop()
	}
	s.pendingSync[userID] = time.AfterFunc(5*time.Second, func() {
		s.syncToDB(context.Background(), userID)
	})
}

func (s *sCarts) syncToDB(ctx context.Context, userID int64) {
	defer func() {
		s.mu.Lock()
		delete(s.pendingSync, userID)
		s.mu.Unlock()
	}()

	resp, err := readCartFromRedis(ctx, userID)
	if err != nil {
		return
	}

	cart, err := dao.Carts.Ctx(ctx).Where(dao.Carts.Columns().UserId, userID).One()
	if err != nil {
		return
	}

	var cartId int64
	if cart == nil {
		id, err := dao.Carts.Ctx(ctx).Data(g.Map{
			"user_id":      userID,
			"item_count":   0,
			"total_amount": 0,
			"created_at":   gtime.Now(),
		}).InsertAndGetId()
		if err != nil {
			return
		}
		cartId = id
	} else {
		cartId = cart["id"].Int64()
	}

	_ = dao.Carts.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, _ = tx.Model("tx_cart_items").Where("cart_id", cartId).Delete()
		for _, item := range resp.Items {
			_, err := tx.Model("tx_cart_items").Insert(g.Map{
				"cart_id":      cartId,
				"sku_id":       item.SkuID,
				"product_id":   item.ProductID,
				"product_name": item.ProductName,
				"sku_spec":     item.SkuSpec,
				"image":        item.Image,
				"price":        item.Price,
				"quantity":     item.Quantity,
			})
			if err != nil {
				return err
			}
		}
		_, err := tx.Model("tx_carts").Where("id", cartId).Data(g.Map{
			"item_count":   resp.ItemCount,
			"total_amount": resp.TotalAmount,
			"updated_at":   gtime.Now(),
		}).Update()
		return err
	})
}
