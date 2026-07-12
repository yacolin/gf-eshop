package products

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/singleflight"

	"gf-eshop/api/products/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sProducts struct {
	sf singleflight.Group
}

func init() {
	service.RegisterProducts(&sProducts{})
}

func (s *sProducts) List(ctx context.Context, req *v1.ProductsListReq) (res *v1.ProductsListRes, err error) {
	size := req.Size
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	var cursorId int64
	if req.Cursor != "" {
		b, err := base64.StdEncoding.DecodeString(req.Cursor)
		if err == nil {
			cursorId, _ = strconv.ParseInt(string(b), 10, 64)
		}
	}

	// ZSET 路径：仅当没有文本搜索/价格筛选时可用
	useZSET := req.Name == "" && req.PriceMin == 0 && req.PriceMax == 0
	if useZSET {
		result, err := s.listFromZSET(ctx, req, cursorId, size)
		if err == nil {
			return result, nil
		}
		g.Log().Warning(ctx, "ZSET list cache miss, fallback to DB: %v", err)
	}

	return s.listFromDB(ctx, req, cursorId, size)
}

// listFromZSET 通过 Redis ZSET 游标分页（仅筛选类目/品牌/状态）
func (s *sProducts) listFromZSET(ctx context.Context, req *v1.ProductsListReq, cursorId int64, size int) (*v1.ProductsListRes, error) {
	status := 0
	if req.Status != nil {
		status = *req.Status
	}
	key := cacheKeyProductListIDs(req.CategoryId, req.BrandId, status)

	// Singleflight：同一组合仅一个 goroutine 构建 ZSET
	_, err, _ := s.sf.Do(key, func() (interface{}, error) {
		exists, err := g.Redis().Do(ctx, "EXISTS", key)
		if err != nil || exists.Int() == 0 {
			ids, err := buildProductListIDs(ctx, req.CategoryId, req.BrandId, status)
			if err != nil {
				return nil, err
			}
			_ = setProductListZSET(ctx, key, ids)
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	// 从 ZSET 获取 ID 列表
	ids, err := fetchProductListIDs(ctx, key, cursorId, size+1)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return &v1.ProductsListRes{List: make([]*v1.ProductsListItem, 0)}, nil
	}

	hasMore := len(ids) > size
	if hasMore {
		ids = ids[:size]
	}

	// 按 ID 取完整实体
	var products []*entity.Products
	err = dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id+" IN (?)", ids).OrderAsc(dao.Products.Columns().Id).Scan(&products)
	if err != nil {
		return nil, err
	}
	if products == nil {
		products = make([]*entity.Products, 0)
	}

	// 聚合价格&库存
	stats := enrichProductStats(ctx, ids)

	list := make([]*v1.ProductsListItem, len(products))
	for i, p := range products {
		s := stats[p.Id]
		list[i] = &v1.ProductsListItem{
			Products:   p,
			PriceMin:   s.PriceMin,
			PriceMax:   s.PriceMax,
			TotalStock: s.TotalStock,
		}
	}

	cursor := ""
	if len(products) > 0 {
		cursor = base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(products[len(products)-1].Id, 10)))
	}
	return &v1.ProductsListRes{List: list, Cursor: cursor, HasMore: hasMore}, nil
}

// listFromDB 通过数据库 keyset 游标分页（含名称/价格筛选）
func (s *sProducts) listFromDB(ctx context.Context, req *v1.ProductsListReq, cursorId int64, size int) (*v1.ProductsListRes, error) {
	m := dao.Products.Ctx(ctx)
	if req.Name != "" {
		m = m.WhereLike(dao.Products.Columns().Name, "%"+req.Name+"%")
	}
	if req.CategoryId > 0 {
		m = m.Where(dao.Products.Columns().CategoryId+" IN (SELECT id FROM sp_categories WHERE id = ? OR path LIKE CONCAT((SELECT IFNULL(path,'') FROM sp_categories WHERE id = ?), ?, '/%'))",
			req.CategoryId, req.CategoryId, strconv.FormatInt(req.CategoryId, 10))
	}
	if req.BrandId > 0 {
		m = m.Where(dao.Products.Columns().BrandId, req.BrandId)
	}
	if req.Status != nil {
		m = m.Where(dao.Products.Columns().Status, *req.Status)
	}
	if req.PriceMin > 0 {
		m = m.Where(dao.Products.Columns().Id+" IN (SELECT DISTINCT product_id FROM sp_skus WHERE price >= ? AND deleted_at IS NULL)", req.PriceMin)
	}
	if req.PriceMax > 0 {
		m = m.Where(dao.Products.Columns().Id+" IN (SELECT DISTINCT product_id FROM sp_skus WHERE price <= ? AND deleted_at IS NULL)", req.PriceMax)
	}
	if cursorId > 0 {
		m = m.WhereGT(dao.Products.Columns().Id, cursorId)
	}

	values, err := m.Fields(dao.Products.Columns().Id).OrderAsc(dao.Products.Columns().Id).Limit(size + 1).Array()
	if err != nil {
		return nil, err
	}
	ids := make([]int64, len(values))
	for i, v := range values {
		ids[i] = v.Int64()
	}
	if len(ids) == 0 {
		return &v1.ProductsListRes{List: make([]*v1.ProductsListItem, 0)}, nil
	}

	hasMore := len(ids) > size
	if hasMore {
		ids = ids[:size]
	}

	var products []*entity.Products
	err = dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id+" IN (?)", ids).OrderAsc(dao.Products.Columns().Id).Scan(&products)
	if err != nil {
		return nil, err
	}
	if products == nil {
		products = make([]*entity.Products, 0)
	}

	// 聚合价格&库存
	stats := enrichProductStats(ctx, ids)

	list := make([]*v1.ProductsListItem, len(products))
	for i, p := range products {
		s := stats[p.Id]
		list[i] = &v1.ProductsListItem{
			Products:   p,
			PriceMin:   s.PriceMin,
			PriceMax:   s.PriceMax,
			TotalStock: s.TotalStock,
		}
	}

	cursor := ""
	if len(products) > 0 {
		cursor = base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(products[len(products)-1].Id, 10)))
	}
	return &v1.ProductsListRes{List: list, Cursor: cursor, HasMore: hasMore}, nil
}

func (s *sProducts) Detail(ctx context.Context, req *v1.ProductsDetailReq) (res *v1.ProductsDetailRes, err error) {
	var (
		product      *entity.Products
		skuEntities  []*entity.Skus
		description  *entity.ProductDescriptions
		productAttrs []v1.ProductAttrDetailResponse
	)

	g, egCtx := errgroup.WithContext(ctx)

	// SPU（L1 → Bloom → L2 → DB 多级缓存）
	g.Go(func() error {
		// 1. L1 本地缓存（最快）
		if p, ok := productLocalCache.get(req.Id); ok {
			product = p
			return nil
		}

		// 2. Bloom Filter 快速拦截（未预热时 count==0 放行）
		if !productBloom.mayExist(req.Id) {
			return errcode.ErrProductNotFound
		}

		// 3. L2 Redis
		if p, err := getProductEntityCache(egCtx, req.Id); err == nil && p != nil {
			productLocalCache.set(p.Id, p)
			product = p
			return nil
		}

		// 4. DB 兜底
		var p *entity.Products
		err := dao.Products.Ctx(egCtx).Where(dao.Products.Columns().Id, req.Id).Scan(&p)
		if err != nil {
			return err
		}
		if p == nil {
			return errcode.ErrProductNotFound
		}
		product = p
		// 回填所有缓存层级
		_ = setProductEntityCache(context.Background(), product)
		return nil
	})

	// SKU 列表
	g.Go(func() error {
		var list []*entity.Skus
		err := dao.Skus.Ctx(egCtx).Where(dao.Skus.Columns().ProductId, req.Id).OrderAsc(dao.Skus.Columns().Id).Scan(&list)
		if err != nil {
			return err
		}
		if list == nil {
			list = make([]*entity.Skus, 0)
		}
		skuEntities = list
		return nil
	})

	// 图文描述（可选）
	g.Go(func() error {
		var d *entity.ProductDescriptions
		err := dao.ProductDescriptions.Ctx(egCtx).Where(dao.ProductDescriptions.Columns().ProductId, req.Id).Scan(&d)
		if err != nil {
			return err
		}
		if d != nil {
			description = d
		}
		return nil
	})

	// 商品属性（含属性名称 join）
	g.Go(func() error {
		var err error
		productAttrs, err = findProductAttrsWithName(egCtx, req.Id)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errcode.ErrProductNotFound
	}

	// 从 SKU 实体映射为 SkuDetailItem，加载库存
	skuItems := make([]*v1.SkuDetailItem, len(skuEntities))
	for i, sku := range skuEntities {
		skuItems[i] = &v1.SkuDetailItem{Skus: sku}
	}
	enrichSKUInventory(ctx, skuItems)

	// 从 sku_specs EAV 表聚合规格维度
	specAttrs := aggregateSpecAttrs(ctx, skuEntities)

	// 聚合商品级价格&库存
	var priceMin, priceMax, totalStock int64
	for _, item := range skuItems {
		if item.Price > 0 {
			if priceMin == 0 || item.Price < priceMin {
				priceMin = item.Price
			}
			if item.Price > priceMax {
				priceMax = item.Price
			}
		}
		totalStock += item.AvailableQuantity
	}

	// 查询类目下所有属性，用于补充 attribute_id / sort_order
	attrNameMap := make(map[string]int64)
	attrOrder := make(map[string]int)
	if product.CategoryId > 0 {
		var catAttrs []*entity.Attributes
		err := dao.Attributes.Ctx(ctx).
			Where(dao.Attributes.Columns().CategoryId, product.CategoryId).
			Scan(&catAttrs)
		if err == nil {
			for i, a := range catAttrs {
				attrNameMap[a.Name] = a.Id
				attrOrder[a.Name] = i
			}
		}
	}

	// 富化可选的规格（attribute_id / sort_order）
	for i := range specAttrs {
		specAttrs[i].AttributeId = attrNameMap[specAttrs[i].AttributeName]
		specAttrs[i].SortOrder = attrOrder[specAttrs[i].AttributeName]
	}
	if specAttrs == nil {
		specAttrs = make([]v1.ProductAttrDetailResponse, 0)
	}

	// 不可选的规格：去重（排除已在可选规格中出现的）
	seen := make(map[string]bool, len(specAttrs))
	for _, s := range specAttrs {
		seen[s.AttributeName] = true
	}
	nonSelectable := make([]v1.ProductAttrDetailResponse, 0, len(productAttrs))
	for _, a := range productAttrs {
		if seen[a.AttributeName] {
			continue
		}
		seen[a.AttributeName] = true
		a.AttributeId = attrNameMap[a.AttributeName]
		a.SortOrder = attrOrder[a.AttributeName]
		nonSelectable = append(nonSelectable, a)
	}
	if nonSelectable == nil {
		nonSelectable = make([]v1.ProductAttrDetailResponse, 0)
	}

	specs := &v1.ProductSpecResponse{
		Selectable:    specAttrs,
		NonSelectable: nonSelectable,
	}

	return &v1.ProductsDetailRes{
		Products:    product,
		PriceMin:    priceMin,
		PriceMax:    priceMax,
		TotalStock:  totalStock,
		Description: description,
		SKUs:        skuItems,
		Specs:       specs,
	}, nil
}

func (s *sProducts) DetailPure(ctx context.Context, req *v1.ProductsDetailPureReq) (res *v1.ProductsDetailPureRes, err error) {
	var entity *entity.Products
	err = dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errcode.ErrProductNotFound
	}
	return &v1.ProductsDetailPureRes{Products: entity}, nil
}

func (s *sProducts) Create(ctx context.Context, req *v1.ProductsCreateReq) (res *v1.ProductsCreateRes, err error) {
	var productId int64

	err = dao.Products.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		productData := do.Products{
			Name:       req.Name,
			Subtitle:   req.Subtitle,
			CategoryId: req.CategoryId,
			BrandId:    req.BrandId,
			Unit:       req.Unit,
			MainImage:  req.MainImage,
			VideoUrl:   req.VideoUrl,
			SortOrder:  req.SortOrder,
			Status:     req.Status,
			CreatedBy:  req.CreatedBy,
		}
		if req.Images != "" {
			productData.Images = req.Images
		}
		result, err := tx.Model("sp_products").Insert(productData)
		if err != nil {
			return err
		}
		id, _ := result.LastInsertId()
		productId = id

		// 写入非销售属性（is_sku_spec = 0）
		for _, attr := range req.Attributes {
			data := do.ProductAttributes{
				ProductId:   productId,
				AttributeId: attr.AttributeId,
				Value:       attr.Value,
			}
			if attr.AttributeValueId > 0 {
				data.AttributeValueId = attr.AttributeValueId
			}
			_, err = tx.Model("sp_product_attributes").Insert(data)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if productId > 0 {
		productBloom.add(productId)
	}
	return &v1.ProductsCreateRes{Id: productId}, nil
}

func (s *sProducts) CreateFull(ctx context.Context, req *v1.ProductsCreateFullReq) (res *v1.ProductsCreateFullRes, err error) {
	var productId int64

	err = dao.Products.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, err := dao.Products.Ctx(ctx).TX(tx).Insert(do.Products{
			Name:       req.Name,
			Subtitle:   req.Subtitle,
			CategoryId: req.CategoryId,
			BrandId:    req.BrandId,
			Unit:       req.Unit,
			MainImage:  req.MainImage,
			Images:     req.Images,
			VideoUrl:   req.VideoUrl,
			SortOrder:  req.SortOrder,
			Status:     0,
			CreatedBy:  req.CreatedBy,
		})
		if err != nil {
			return err
		}
		id, _ := result.LastInsertId()
		productId = id

		for _, sku := range req.SKUs {
			barcode := interface{}(sku.Barcode)
			if sku.Barcode == "" {
				barcode = nil
			}
			_, err = tx.Model("sp_skus").Insert(do.Skus{
				ProductId:      productId,
				SkuCode:        sku.SkuCode,
				Barcode:        barcode,
				SpecSummary:    sku.SpecSummary,
				Price:          sku.Price,
				MarketPrice:    sku.MarketPrice,
				CostPrice:      sku.CostPrice,
				Weight:         sku.Weight,
				Volume:         sku.Volume,
				Length:         sku.Length,
				Width:          sku.Width,
				Height:         sku.Height,
				MinPurchaseQty: sku.MinPurchaseQty,
				MaxPurchaseQty: sku.MaxPurchaseQty,
				Image:          sku.Image,
			})
			if err != nil {
				return err
			}
		}

		if req.Description != "" || req.MobileDesc != "" {
			_, err = tx.Model("sp_product_descriptions").Insert(do.ProductDescriptions{
				ProductId:         productId,
				Description:       req.Description,
				MobileDescription: req.MobileDesc,
			})
			if err != nil {
				return err
			}
		}

		for _, attr := range req.Attributes {
			data := do.ProductAttributes{
				ProductId:   productId,
				AttributeId: attr.AttributeId,
				Value:       attr.Value,
			}
			if attr.AttributeValueId > 0 {
				data.AttributeValueId = attr.AttributeValueId
			}
			_, err = tx.Model("sp_product_attributes").Insert(data)
			if err != nil {
				return err
			}
		}

		hasDesc := 0
		if req.Description != "" || req.MobileDesc != "" {
			hasDesc = 1
		}
		_, err = tx.Model("sp_products").Where("id", productId).Update(g.Map{
			"has_description": hasDesc,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if productId > 0 {
		productBloom.add(productId)
	}
	delAllProductListCaches(context.Background())
	return &v1.ProductsCreateFullRes{Id: productId}, nil
}
// GetAttributes 获取商品绑定的非销售属性列表
func (s *sProducts) GetAttributes(ctx context.Context, req *v1.ProductsGetAttributesReq) (res *v1.ProductsGetAttributesRes, err error) {
	type attrRow struct {
		Id               int64  `orm:"id"`
		AttributeId      int64  `orm:"attribute_id"`
		AttributeName    string `orm:"attribute_name"`
		AttributeValueId int64  `orm:"attribute_value_id"`
		Value            string `orm:"value"`
	}
	var rows []attrRow
	err = g.DB().Model("sp_product_attributes pa").
		Fields("pa.id", "pa.attribute_id", "a.name AS attribute_name", "pa.attribute_value_id", "pa.value").
		LeftJoin("sp_attributes a", "a.id = pa.attribute_id").
		Where("pa.product_id", req.Id).
		Where("pa.deleted_at IS NULL").
		Order("pa.sort_order ASC, pa.id ASC").
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = make([]attrRow, 0)
	}
	list := make([]*v1.ProductsAttributeItem, len(rows))
	for i, r := range rows {
		list[i] = &v1.ProductsAttributeItem{
			Id:               r.Id,
			AttributeId:      r.AttributeId,
			AttributeName:    r.AttributeName,
			AttributeValueId: r.AttributeValueId,
			Value:            r.Value,
		}
	}
	return &v1.ProductsGetAttributesRes{List: list}, nil
}

// UpdateAttributes 全量替换商品的非销售属性
func (s *sProducts) UpdateAttributes(ctx context.Context, req *v1.ProductsUpdateAttributesReq) (res *v1.ProductsUpdateAttributesRes, err error) {
	err = dao.Products.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model("sp_product_attributes").
			Where("product_id", req.Id).
			Unscoped().Delete()
		if err != nil {
			return err
		}
		for _, attr := range req.Attributes {
			data := do.ProductAttributes{
				ProductId:   req.Id,
				AttributeId: attr.AttributeId,
				Value:       attr.Value,
			}
			if attr.AttributeValueId > 0 {
				data.AttributeValueId = attr.AttributeValueId
			}
			_, err = tx.Model("sp_product_attributes").Insert(data)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.ProductsUpdateAttributesRes{}, nil
}

// EnrichedDetail 商品富化详情（同 Detail，独立路径供前端使用）
func (s *sProducts) EnrichedDetail(ctx context.Context, req *v1.ProductsEnrichedDetailReq) (res *v1.ProductsEnrichedDetailRes, err error) {
	detail, err := s.Detail(ctx, &v1.ProductsDetailReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &v1.ProductsEnrichedDetailRes{ProductsDetailRes: detail}, nil
}

// BatchCreateSKUs 批量生成 SKU（笛卡尔积）
func (s *sProducts) BatchCreateSKUs(ctx context.Context, req *v1.ProductsBatchCreateSKUsReq) (res *v1.ProductsBatchCreateSKUsRes, err error) {
	// 校验商品存在
	count, err := dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, req.ProductId).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrProductNotFound
	}

	// 校验所有属性均为 is_sku_spec = 1
	attrNameMap := make(map[int64]string)
	valueTextMap := make(map[int64]map[int64]string)
	for _, group := range req.SpecGroups {
		var attr *entity.Attributes
		err := dao.Attributes.Ctx(ctx).Where(dao.Attributes.Columns().Id, group.AttributeId).Scan(&attr)
		if err != nil {
			return nil, err
		}
		if attr == nil {
			return nil, gerror.NewCode(gcode.CodeInvalidParameter,
				fmt.Sprintf("属性不存在: %d", group.AttributeId))
		}
		if attr.IsSkuSpec != 1 {
			return nil, gerror.NewCode(gcode.CodeInvalidParameter,
				fmt.Sprintf("非销售属性不可用于生成SKU: %s", attr.Name))
		}
		attrNameMap[group.AttributeId] = attr.Name

		if len(group.ValueIds) == 0 {
			return nil, gerror.NewCode(gcode.CodeInvalidParameter,
				fmt.Sprintf("属性 %s 至少选择一个值", attr.Name))
		}

		vals := make(map[int64]string, len(group.ValueIds))
		for _, vid := range group.ValueIds {
			var v *entity.AttributeValues
			err := dao.AttributeValues.Ctx(ctx).Where(dao.AttributeValues.Columns().Id, vid).Scan(&v)
			if err != nil {
				return nil, err
			}
			if v == nil {
				return nil, gerror.NewCode(gcode.CodeInvalidParameter,
					fmt.Sprintf("属性值不存在: %d", vid))
			}
			vals[vid] = v.Value
		}
		valueTextMap[group.AttributeId] = vals
	}

	// 笛卡尔积计算
	combinations := computeCartesianProduct(req.SpecGroups, valueTextMap)
	if len(combinations) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "没有有效的属性值组合")
	}

	// 查询所有仓库
	var warehouses []*entity.Warehouses
	err = dao.Warehouses.Ctx(ctx).Scan(&warehouses)
	if err != nil {
		return nil, err
	}

	prefix := req.SkuCodePrefix
	if prefix == "" {
		prefix = fmt.Sprintf("SKU-%d-", req.ProductId)
	}

	// 查询当前已存在的 SKU 数量（含软删除），避免重新创建时唯一键冲突
	var existingCount int
	cnt, cntErr := dao.Skus.Ctx(ctx).Unscoped().Where(dao.Skus.Columns().ProductId, req.ProductId).Count()
	if cntErr == nil {
		existingCount = cnt
	}
	offset := existingCount

	type createdSku struct {
		id      int64
		code    string
		summary string
		price   int64
	}
	var createdSkus []createdSku

	err = dao.Products.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for i, combo := range combinations {
			parts := make([]string, 0, len(combo))
			for _, opt := range combo {
				parts = append(parts, opt.ValueText)
			}
			summary := strings.Join(parts, " / ")

			skuCode := fmt.Sprintf("%s%03d", prefix, offset+i+1)

			result, err := tx.Model("sp_skus").Insert(do.Skus{
				ProductId:   req.ProductId,
				SkuCode:     skuCode,
				SpecSummary: summary,
				Price:       req.BasePrice,
				MarketPrice: req.BasePrice,
				CostPrice:   int64(float64(req.BasePrice) * 0.6),
			})
			if err != nil {
				return err
			}
			skuId, _ := result.LastInsertId()

			for sortOrder, opt := range combo {
				_, err = tx.Model("sp_sku_specs").Insert(do.SkuSpecs{
					SkuId:            skuId,
					AttributeId:      opt.AttributeId,
					AttributeValueId:  opt.ValueId,
					SortOrder:        sortOrder,
				})
				if err != nil {
					return err
				}
			}

			for _, w := range warehouses {
				_, err = tx.Model("sp_inventories").Insert(do.Inventories{
					SkuId:       skuId,
					WarehouseId: w.Id,
					Quantity:    0,
					Reserved:    0,
					Threshold:   10,
				})
				if err != nil {
					return err
				}
			}

			createdSkus = append(createdSkus, createdSku{
				id:      skuId,
				code:    skuCode,
				summary: summary,
				price:   req.BasePrice,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	items := make([]*v1.BatchCreateSkuItem, len(createdSkus))
	for i, sku := range createdSkus {
		items[i] = &v1.BatchCreateSkuItem{
			Id:          sku.id,
			SkuCode:     sku.code,
			SpecSummary: sku.summary,
			Price:       sku.price,
		}
	}

	delAllProductListCaches(context.Background())
	return &v1.ProductsBatchCreateSKUsRes{
		Total: len(items),
		SKUs:  items,
	}, nil
}

// specOption 笛卡尔积中的一个维度选项
type specOption struct {
	AttributeId int64
	ValueId     int64
	ValueText   string
}

// computeCartesianProduct 计算属性值的笛卡尔积，返回所有规格组合
func computeCartesianProduct(groups []v1.SpecGroupItem, valueTextMap map[int64]map[int64]string) [][]specOption {
	if len(groups) == 0 {
		return nil
	}
	result := [][]specOption{{}}
	for _, group := range groups {
		var newResult [][]specOption
		for _, combo := range result {
			for _, vid := range group.ValueIds {
				text := valueTextMap[group.AttributeId][vid]
				newCombo := make([]specOption, len(combo), len(combo)+1)
				copy(newCombo, combo)
				newCombo = append(newCombo, specOption{
					AttributeId: group.AttributeId,
					ValueId:     vid,
					ValueText:   text,
				})
				newResult = append(newResult, newCombo)
			}
		}
		result = newResult
	}
	return result
}

func (s *sProducts) Update(ctx context.Context, req *v1.ProductsUpdateReq) (res *v1.ProductsUpdateRes, err error) {
	count, err := dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrProductNotFound
	}

	_, err = dao.Products.Ctx(ctx).Data(do.Products{
		Name:       req.Name,
		Subtitle:   req.Subtitle,
		CategoryId: req.CategoryId,
		BrandId:    req.BrandId,
		Unit:       req.Unit,
		MainImage:  req.MainImage,
		Images:     req.Images,
		VideoUrl:   req.VideoUrl,
		SortOrder:  req.SortOrder,
		Status:     req.Status,
		UpdatedBy:  req.UpdatedBy,
	}).Where(dao.Products.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	delProductEntityCache(context.Background(), req.Id)
	delAllProductListCaches(context.Background())
	return &v1.ProductsUpdateRes{}, nil
}

func (s *sProducts) Delete(ctx context.Context, req *v1.ProductsDeleteReq) (res *v1.ProductsDeleteRes, err error) {
	_, err = dao.Products.Ctx(ctx).Where(dao.Products.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	delProductEntityCache(context.Background(), req.Id)
	delAllProductListCaches(context.Background())
	return &v1.ProductsDeleteRes{}, nil
}

// ── Helper: Build Product List IDs ────────────────────────────────────

// buildProductListIDs 查询数据库获取符合条件的商品 ID 列表（用于重建 ZSET）
func buildProductListIDs(ctx context.Context, categoryId, brandId int64, status int) ([]int64, error) {
	m := dao.Products.Ctx(ctx).Fields(dao.Products.Columns().Id)
	if categoryId > 0 {
		m = m.Where(dao.Products.Columns().CategoryId+" IN (SELECT id FROM sp_categories WHERE id = ? OR path LIKE CONCAT((SELECT IFNULL(path,'') FROM sp_categories WHERE id = ?), ?, '/%'))",
			categoryId, categoryId, strconv.FormatInt(categoryId, 10))
	}
	if brandId > 0 {
		m = m.Where(dao.Products.Columns().BrandId, brandId)
	}
	if status > 0 {
		m = m.Where(dao.Products.Columns().Status, status)
	}

	values, err := m.OrderAsc(dao.Products.Columns().Id).Array()
	if err != nil {
		return nil, err
	}
	ids := make([]int64, len(values))
	for i, v := range values {
		ids[i] = v.Int64()
	}
	return ids, nil
}

// ── Helper: Product Attributes with Name ────────────────────────────────

// findProductAttrsWithName 查询商品属性并 left join 属性表获取名称
func findProductAttrsWithName(ctx context.Context, productId int64) ([]v1.ProductAttrDetailResponse, error) {
	type attrRow struct {
		AttributeId   int64  `orm:"attribute_id"`
		AttributeName string `orm:"attribute_name"`
		Value         string `orm:"value"`
		SortOrder     int    `orm:"sort_order"`
	}
	var rows []attrRow
	err := g.DB().Model("sp_product_attributes pa").
		Fields("pa.attribute_id", "a.name AS attribute_name", "pa.value", "pa.sort_order").
		LeftJoin("sp_attributes a", "a.id = pa.attribute_id").
		Where("pa.product_id", productId).
		Where("pa.deleted_at IS NULL").
		Order("pa.sort_order ASC, pa.id ASC").
		Scan(&rows)
	if err != nil {
		return nil, err
	}

	// 按 attribute_id 分组聚合 values，保留原始顺序
	seen := make(map[int64]int) // attribute_id → result index
	result := make([]v1.ProductAttrDetailResponse, 0)
	for _, r := range rows {
		if idx, ok := seen[r.AttributeId]; ok {
			result[idx].Values = append(result[idx].Values, r.Value)
		} else {
			seen[r.AttributeId] = len(result)
			result = append(result, v1.ProductAttrDetailResponse{
				AttributeId:   r.AttributeId,
				AttributeName: r.AttributeName,
				Values:        []string{r.Value},
				SortOrder:     r.SortOrder,
			})
		}
	}
	if result == nil {
		result = make([]v1.ProductAttrDetailResponse, 0)
	}
	return result, nil
}

// ── Helper: SKU Spec Aggregation ────────────────────────────────────────

// aggregateSpecAttrs 从 sku_specs EAV 表聚合规格维度（如 颜色→[红,蓝]）
func aggregateSpecAttrs(ctx context.Context, skus []*entity.Skus) []v1.ProductAttrDetailResponse {
	if len(skus) == 0 {
		return make([]v1.ProductAttrDetailResponse, 0)
	}

	skuIDs := make([]int64, 0, len(skus))
	for _, sku := range skus {
		skuIDs = append(skuIDs, sku.Id)
	}

	type specRow struct {
		AttributeName  string `orm:"attribute_name"`
		AttributeValue string `orm:"attribute_value"`
	}
	var rows []specRow
	err := g.DB().Model("sp_sku_specs ss").
		Fields("a.name AS attribute_name", "av.value AS attribute_value").
		LeftJoin("sp_attributes a", "a.id = ss.attribute_id").
		LeftJoin("sp_attribute_values av", "av.id = ss.attribute_value_id").
		Where("ss.sku_id IN (?)", skuIDs).
		Order("ss.sort_order ASC, av.sort_order ASC").
		Scan(&rows)
	if err != nil || len(rows) == 0 {
		return make([]v1.ProductAttrDetailResponse, 0)
	}

	type attrValues struct {
		set   map[string]struct{}
		order []string
	}
	attrs := make(map[string]*attrValues)
	keyOrder := make([]string, 0)

	for _, r := range rows {
		av, ok := attrs[r.AttributeName]
		if !ok {
			av = &attrValues{set: make(map[string]struct{})}
			attrs[r.AttributeName] = av
			keyOrder = append(keyOrder, r.AttributeName)
		}
		if _, seen := av.set[r.AttributeValue]; !seen {
			av.set[r.AttributeValue] = struct{}{}
			av.order = append(av.order, r.AttributeValue)
		}
	}

	result := make([]v1.ProductAttrDetailResponse, len(keyOrder))
	for i, name := range keyOrder {
		result[i] = v1.ProductAttrDetailResponse{
			AttributeName: name,
			Values:        attrs[name].order,
		}
	}
	return result
}




type productStats struct {
	PriceMin   int64
	PriceMax   int64
	TotalStock int64
}

// enrichProductStats 批量查询商品级价格区间和总库存
func enrichProductStats(ctx context.Context, productIDs []int64) map[int64]productStats {
	result := make(map[int64]productStats, len(productIDs))
	if len(productIDs) == 0 {
		return result
	}
	for _, id := range productIDs {
		result[id] = productStats{}
	}

	// 价格区间：从 sp_skus 聚合
	type priceRow struct {
		ProductId int64
		PriceMin  int64
		PriceMax  int64
	}
	var prices []priceRow
	err := dao.Skus.Ctx(ctx).
		Fields("product_id", "MIN(price) AS price_min", "MAX(price) AS price_max").
		Where("product_id IN (?)", productIDs).
		Where("deleted_at IS NULL").
		Group("product_id").
		Scan(&prices)
	if err == nil {
		for _, r := range prices {
			s := result[r.ProductId]
			s.PriceMin = r.PriceMin
			s.PriceMax = r.PriceMax
			result[r.ProductId] = s
		}
	}

	// 总库存：SUM(quantity - reserved) per product
	type stockRow struct {
		ProductId  int64
		TotalStock int64
	}
	var stocks []stockRow
	err = g.DB().Model("sp_skus s").
		Fields("s.product_id", "SUM(COALESCE(i.quantity,0) - COALESCE(i.reserved,0)) AS total_stock").
		LeftJoin("sp_inventories i", "i.sku_id = s.id AND i.deleted_at IS NULL").
		Where("s.product_id IN (?)", productIDs).
		Where("s.deleted_at IS NULL").
		Group("s.product_id").
		Scan(&stocks)
	if err == nil {
		for _, r := range stocks {
			s := result[r.ProductId]
			if r.TotalStock > 0 {
				s.TotalStock = r.TotalStock
			}
			result[r.ProductId] = s
		}
	}

	return result
}

// enrichSKUInventory 批量加载 SKU 库存并写入 SkuDetailItem
func enrichSKUInventory(ctx context.Context, items []*v1.SkuDetailItem) {
	if len(items) == 0 {
		return
	}
	ids := make([]int64, len(items))
	for i, item := range items {
		ids[i] = item.Id
	}

	type invRow struct {
		SkuId    int64 `orm:"sku_id"`
		Quantity int64
		Reserved int64
	}
	var rows []invRow
	err := dao.Inventories.Ctx(ctx).
		Fields("sku_id", "SUM(quantity) AS quantity", "SUM(reserved) AS reserved").
		Where("sku_id IN (?)", ids).
		Where("deleted_at IS NULL").
		Group("sku_id").
		Scan(&rows)
	if err != nil {
		return
	}

	invMap := make(map[int64]struct{ qty, rsv int64 }, len(rows))
	for _, r := range rows {
		invMap[r.SkuId] = struct{ qty, rsv int64 }{qty: r.Quantity, rsv: r.Reserved}
	}

	for _, item := range items {
		if inv, ok := invMap[item.Id]; ok {
			avail := inv.qty - inv.rsv
			if avail < 0 {
				avail = 0
			}
			item.AvailableQuantity = avail
			if avail > 0 {
				item.InventoryStatus = "充足"
			} else {
				item.InventoryStatus = "缺货"
			}
		} else {
			item.AvailableQuantity = 0
			item.InventoryStatus = "缺货"
		}
	}
}
