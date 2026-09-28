package products

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/sync/errgroup"

	"gf-eshop/api/products/v1"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/search"
	"gf-eshop/internal/service"
)

// sProducts 商品服务实现，只做编排。
// 数据访问见 repo_products.go / repo_skus.go / repo_attrs.go，
// 缓存见 cache.go，响应组装见 assembler.go，业务规则见 rules.go。
type sProducts struct{}

func init() {
	service.RegisterProducts(&sProducts{})
}

func (s *sProducts) List(ctx context.Context, req *v1.ProductsListReq) (*v1.ProductsListRes, error) {
	size := req.Size
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	cursorId := decodeCursor(req.Cursor)

	// ZSET 路径：仅当没有文本搜索/价格筛选时可用
	if req.Name == "" && req.PriceMin == 0 && req.PriceMax == 0 {
		result, err := s.listFromZSET(ctx, req, cursorId, size)
		if err == nil {
			return result, nil
		}
		g.Log().Warningf(ctx, "ZSET list cache miss, fallback to DB: %v", err)
	}

	// 名称/价格筛选：优先走 ES（名称多字段子串检索 + 价格区间），
	// ES 未启用、处于熔断冷却期、索引未建立或查询失败时自动降级回 DB。
	// 无 name/price 的类目/品牌/状态组合仍由上面的 ZSET 路径承接，未被替换。
	if search.Available(ctx) {
		ids, hasMore, err := searchProductIDsES(ctx, req, cursorId, size)
		if err == nil {
			return buildListFromIDs(ctx, ids, hasMore)
		}
		g.Log().Warningf(ctx, "商品 ES 检索失败，降级 DB 查询: %v", err)
	}

	return s.listFromDB(ctx, req, cursorId, size)
}

// buildListFromIDs 按「已排好序且已裁剪」的 ID 列表组装列表响应。
// ES 与 DB 两条路径共用，保证组装逻辑完全一致。
func buildListFromIDs(ctx context.Context, ids []int64, hasMore bool) (*v1.ProductsListRes, error) {
	if len(ids) == 0 {
		return &v1.ProductsListRes{List: make([]*v1.ProductsListItem, 0)}, nil
	}
	products, err := listByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return buildListResponse(products, listProductStats(ctx, ids), hasMore), nil
}

// listFromZSET 通过 Redis ZSET 游标分页（仅筛选类目/品牌/状态）
func (s *sProducts) listFromZSET(ctx context.Context, req *v1.ProductsListReq, cursorId int64, size int) (*v1.ProductsListRes, error) {
	status := 0
	if req.Status != nil {
		status = *req.Status
	}
	key := cacheKeyProductListIDs(req.CategoryId, req.BrandId, status)

	// 缓存缺失时重建 ZSET
	if err := ensureProductListZSET(ctx, key, func() ([]int64, error) {
		return listIDs(ctx, req.CategoryId, req.BrandId, status)
	}); err != nil {
		return nil, err
	}

	// 从 ZSET 获取 ID 列表（多取一条判断 hasMore）
	ids, err := fetchProductListIDs(ctx, key, cursorId, size+1)
	if err != nil {
		return nil, err
	}
	hasMore := len(ids) > size
	if hasMore {
		ids = ids[:size]
	}
	return buildListFromIDs(ctx, ids, hasMore)
}

// listFromDB 通过数据库 keyset 游标分页（含名称/价格筛选）
func (s *sProducts) listFromDB(ctx context.Context, req *v1.ProductsListReq, cursorId int64, size int) (*v1.ProductsListRes, error) {
	ids, err := listIDsByFilter(ctx, req, cursorId, size)
	if err != nil {
		return nil, err
	}
	hasMore := len(ids) > size
	if hasMore {
		ids = ids[:size]
	}
	return buildListFromIDs(ctx, ids, hasMore)
}

func (s *sProducts) Detail(ctx context.Context, req *v1.ProductsDetailReq) (*v1.ProductsDetailRes, error) {
	var (
		product      *entity.Products
		skuEntities  []*entity.Skus
		description  *entity.ProductDescriptions
		productAttrs []v1.ProductAttrDetailResponse
	)

	eg, egCtx := errgroup.WithContext(ctx)

	// SPU（L1 → Bloom → L2 → DB 多级缓存）
	eg.Go(func() error {
		var err error
		product, err = getByIDCached(egCtx, req.Id)
		return err
	})

	// SKU 列表
	eg.Go(func() error {
		var err error
		skuEntities, err = listSkus(egCtx, req.Id)
		return err
	})

	// 图文描述（可选）
	eg.Go(func() error {
		d, err := getDescription(egCtx, req.Id)
		if err != nil {
			return err
		}
		description = d
		return nil
	})

	// 商品属性（含属性名称 join）
	eg.Go(func() error {
		var err error
		productAttrs, err = findAttrsWithName(egCtx, req.Id)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errcode.ErrProductNotFound
	}

	// SKU 详情项 + 库存
	skuItems := buildSkuItems(skuEntities)
	fillSkuInventory(skuItems, loadSkuInventory(ctx, skuItemIDs(skuItems)))

	// 规格维度：优先 sp_sku_specs EAV，无数据时从商品属性兜底
	specAttrs := aggregateSkuSpecs(ctx, skuEntities)
	if len(specAttrs) == 0 {
		specAttrs = fallbackSkuSpecAttrs(ctx, req.Id)
	}

	// 类目属性仅用于补充 attribute_id / sort_order
	idByName, orderByName := listCategoryAttrs(ctx, product.CategoryId)

	specs := buildSpecResponse(specAttrs, productAttrs, idByName, orderByName)
	return buildDetailResponse(product, description, skuItems, specs), nil
}

func (s *sProducts) DetailPure(ctx context.Context, req *v1.ProductsDetailPureReq) (*v1.ProductsDetailPureRes, error) {
	product, err := getByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errcode.ErrProductNotFound
	}
	return &v1.ProductsDetailPureRes{Products: product}, nil
}

func (s *sProducts) Create(ctx context.Context, req *v1.ProductsCreateReq) (*v1.ProductsCreateRes, error) {
	var productId int64

	err := withProductsTx(ctx, func(ctx context.Context, tx gdb.TX) error {
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

		id, err := insertProduct(ctx, tx, productData)
		if err != nil {
			return err
		}
		productId = id

		return replaceProductAttrs(ctx, tx, productId, req.Attributes)
	})
	if err != nil {
		return nil, err
	}

	if productId > 0 {
		productBloom.add(productId)
	}
	// 双写 ES（含最新价格区间）；失败只记日志，启动对账会修正
	syncProductDoc(context.Background(), productId)
	return &v1.ProductsCreateRes{Id: productId}, nil
}

func (s *sProducts) CreateFull(ctx context.Context, req *v1.ProductsCreateFullReq) (*v1.ProductsCreateFullRes, error) {
	var productId int64
	hasDesc := 0
	if req.Description != "" || req.MobileDesc != "" {
		hasDesc = 1
	}

	err := withProductsTx(ctx, func(ctx context.Context, tx gdb.TX) error {
		id, err := insertProduct(ctx, tx, do.Products{
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
		productId = id

		for _, sku := range req.SKUs {
			if _, err = insertSku(ctx, tx, do.Skus{
				ProductId:      productId,
				SkuCode:        sku.SkuCode,
				Barcode:        nullableString(sku.Barcode),
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
			}); err != nil {
				return err
			}
		}

		if hasDesc == 1 {
			if err = insertDescription(ctx, tx, productId, req.Description, req.MobileDesc); err != nil {
				return err
			}
		}

		if err = replaceProductAttrs(ctx, tx, productId, req.Attributes); err != nil {
			return err
		}

		return setHasDescription(ctx, tx, productId, hasDesc)
	})
	if err != nil {
		return nil, err
	}

	if productId > 0 {
		productBloom.add(productId)
	}
	delAllProductListCaches(context.Background())
	// 双写 ES（含最新价格区间）
	syncProductDoc(context.Background(), productId)
	return &v1.ProductsCreateFullRes{Id: productId}, nil
}

// GetAttributes 获取商品绑定的非销售属性列表
func (s *sProducts) GetAttributes(ctx context.Context, req *v1.ProductsGetAttributesReq) (*v1.ProductsGetAttributesRes, error) {
	list, err := listProductAttrsWithValue(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.ProductsGetAttributesRes{List: list}, nil
}

// UpdateAttributes 全量替换商品的非销售属性
func (s *sProducts) UpdateAttributes(ctx context.Context, req *v1.ProductsUpdateAttributesReq) (*v1.ProductsUpdateAttributesRes, error) {
	err := withProductsTx(ctx, func(ctx context.Context, tx gdb.TX) error {
		return replaceProductAttrs(ctx, tx, req.Id, req.Attributes)
	})
	if err != nil {
		return nil, err
	}
	return &v1.ProductsUpdateAttributesRes{}, nil
}

// EnrichedDetail 商品富化详情（同 Detail，独立路径供前端使用）
func (s *sProducts) EnrichedDetail(ctx context.Context, req *v1.ProductsEnrichedDetailReq) (*v1.ProductsEnrichedDetailRes, error) {
	detail, err := s.Detail(ctx, &v1.ProductsDetailReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &v1.ProductsEnrichedDetailRes{ProductsDetailRes: detail}, nil
}

// BatchCreateSKUs 批量生成 SKU（笛卡尔积）
func (s *sProducts) BatchCreateSKUs(ctx context.Context, req *v1.ProductsBatchCreateSKUsReq) (*v1.ProductsBatchCreateSKUsRes, error) {
	// 校验商品存在
	exists, err := existsByID(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errcode.ErrProductNotFound
	}

	// 校验所有属性均为 is_sku_spec = 1
	valueTextMap := make(map[int64]map[int64]string)
	for _, group := range req.SpecGroups {
		attr, err := getAttribute(ctx, group.AttributeId)
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

		if len(group.ValueIds) == 0 {
			return nil, gerror.NewCode(gcode.CodeInvalidParameter,
				fmt.Sprintf("属性 %s 至少选择一个值", attr.Name))
		}

		vals := make(map[int64]string, len(group.ValueIds))
		for _, vid := range group.ValueIds {
			v, err := getAttributeValue(ctx, vid)
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
	warehouses, err := listWarehouses(ctx)
	if err != nil {
		return nil, err
	}

	prefix := req.SkuCodePrefix
	if prefix == "" {
		prefix = fmt.Sprintf("SKU-%d-", req.ProductId)
	}

	// 已存在的 SKU 数量（含软删除）作为编码起始偏移，避免重新创建时唯一键冲突
	offset := countSkusIncludingDeleted(ctx, req.ProductId)

	type createdSku struct {
		id      int64
		code    string
		summary string
		price   int64
	}
	var createdSkus []createdSku

	err = withProductsTx(ctx, func(ctx context.Context, tx gdb.TX) error {
		for i, combo := range combinations {
			parts := make([]string, 0, len(combo))
			for _, opt := range combo {
				parts = append(parts, opt.ValueText)
			}
			summary := strings.Join(parts, " / ")

			skuCode := fmt.Sprintf("%s%03d", prefix, offset+i+1)

			skuId, err := insertSku(ctx, tx, do.Skus{
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

			for sortOrder, opt := range combo {
				if err = insertSkuSpec(ctx, tx, do.SkuSpecs{
					SkuId:            skuId,
					AttributeId:      opt.AttributeId,
					AttributeValueId: opt.ValueId,
					SortOrder:        sortOrder,
				}); err != nil {
					return err
				}
			}

			for _, w := range warehouses {
				if err = insertInventory(ctx, tx, do.Inventories{
					SkuId:       skuId,
					WarehouseId: w.Id,
					Quantity:    0,
					Reserved:    0,
					Threshold:   10,
				}); err != nil {
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
	// 新建 SKU 会改变商品的价格区间，需重新同步检索索引
	syncProductDoc(context.Background(), req.ProductId)
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

func (s *sProducts) Update(ctx context.Context, req *v1.ProductsUpdateReq) (*v1.ProductsUpdateRes, error) {
	exists, err := existsByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errcode.ErrProductNotFound
	}

	err = updateProduct(ctx, req.Id, do.Products{
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
	})
	if err != nil {
		return nil, err
	}

	delProductEntityCache(context.Background(), req.Id)
	delAllProductListCaches(context.Background())
	// 双写 ES：重新取库以保证文档与 DB 一致
	syncProductDoc(context.Background(), req.Id)
	return &v1.ProductsUpdateRes{}, nil
}

func (s *sProducts) UpdateFull(ctx context.Context, req *v1.ProductsUpdateFullReq) (*v1.ProductsUpdateFullRes, error) {
	// 校验 Rule 1/2/3，并取回旧 SKU 映射用于保留 sku_code
	oldSKUMap, err := validateUpdateFull(ctx, req)
	if err != nil {
		return nil, err
	}

	hasDesc := 0
	if req.Description != "" || req.MobileDesc != "" {
		hasDesc = 1
	}

	err = withProductsTx(ctx, func(ctx context.Context, tx gdb.TX) error {
		err := updateProductTx(ctx, tx, req.Id, do.Products{
			Name:      req.Name,
			Subtitle:  req.Subtitle,
			BrandId:   req.BrandId,
			Unit:      req.Unit,
			MainImage: req.MainImage,
			Images:    req.Images,
			VideoUrl:  req.VideoUrl,
			SortOrder: req.SortOrder,
			UpdatedBy: req.UpdatedBy,
		})
		if err != nil {
			return err
		}

		for _, sku := range req.SKUs {
			if sku.Id > 0 {
				// 已有 SKU：请求未带 sku_code 时沿用旧值
				skuCode := sku.SkuCode
				if skuCode == "" {
					skuCode = oldSKUMap[sku.Id].SkuCode
				}
				if err = updateSku(ctx, tx, sku.Id, do.Skus{
					SkuCode:        skuCode,
					Barcode:        nullableString(sku.Barcode),
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
				}); err != nil {
					return err
				}
				continue
			}

			// 新增 SKU
			if _, err = insertSku(ctx, tx, do.Skus{
				ProductId:      req.Id,
				SkuCode:        sku.SkuCode,
				Barcode:        nullableString(sku.Barcode),
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
			}); err != nil {
				return err
			}
		}

		if hasDesc == 1 {
			if err = upsertDescription(ctx, tx, req.Id, req.Description, req.MobileDesc); err != nil {
				return err
			}
		} else if err = deleteDescription(ctx, tx, req.Id); err != nil {
			return err
		}

		// 属性为空时视为不修改，避免误删
		if len(req.Attributes) > 0 {
			if err = replaceProductAttrs(ctx, tx, req.Id, req.Attributes); err != nil {
				return err
			}
		}

		return setHasDescription(ctx, tx, req.Id, hasDesc)
	})
	if err != nil {
		return nil, err
	}

	delProductEntityCache(context.Background(), req.Id)
	delAllProductListCaches(context.Background())
	// 双写 ES：全量更新可能改了名称/类目/品牌/状态，重新取库同步
	syncProductDoc(context.Background(), req.Id)
	return &v1.ProductsUpdateFullRes{}, nil
}

func (s *sProducts) Delete(ctx context.Context, req *v1.ProductsDeleteReq) (*v1.ProductsDeleteRes, error) {
	if err := deleteProduct(ctx, req.Id); err != nil {
		return nil, err
	}
	delProductEntityCache(context.Background(), req.Id)
	delAllProductListCaches(context.Background())
	// 双写 ES：从检索索引中移除
	removeProductDoc(context.Background(), req.Id)
	return &v1.ProductsDeleteRes{}, nil
}
