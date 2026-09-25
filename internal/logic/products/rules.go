package products

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"gf-eshop/api/products/v1"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/entity"
)

// 本文件收敛 UpdateFull 的业务规则：
//   Rule 1  商品类目创建后不可修改
//   Rule 2  不可删除旧规格组，新 SKU 的规格维度数必须与旧规格一致
//   Rule 3  不可移除旧的 SKU 组合，且请求中的 SKU ID 必须属于该商品

// validateUpdateFull 加载旧商品/旧 SKU/旧规格组并执行 Rule 1/2/3 校验，
// 返回旧 SKU 映射（ID → 实体）供后续事务复用。
func validateUpdateFull(ctx context.Context, req *v1.ProductsUpdateFullReq) (map[int64]*entity.Skus, error) {
	oldProduct, err := getByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if oldProduct == nil {
		return nil, errcode.ErrProductNotFound
	}

	// ── Rule 1: 禁止修改类目 ──────────────────────────────────────────
	if err := ruleCategoryImmutable(oldProduct, req); err != nil {
		return nil, err
	}

	// ── 加载旧 SKU ────────────────────────────────────────────────────
	oldSkus, err := listSkus(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	oldSKUMap := make(map[int64]*entity.Skus, len(oldSkus))
	for _, sku := range oldSkus {
		oldSKUMap[sku.Id] = sku
	}

	// ── 加载旧 sku_specs 规格组 ──────────────────────────────────────
	oldSkuIDs := make([]int64, len(oldSkus))
	for i, sku := range oldSkus {
		oldSkuIDs[i] = sku.Id
	}
	oldSpecAttrIDs, err := listSkuSpecAttrs(ctx, oldSkuIDs)
	if err != nil {
		return nil, err
	}

	// ── Rule 2: 禁止删除旧规格组 ─────────────────────────────────────
	if err := ruleSpecGroupPreserved(oldSkus, oldSKUMap, oldSpecAttrIDs, req); err != nil {
		return nil, err
	}

	// ── Rule 3: 禁止移除旧的 SKU 组合 ─────────────────────────────────
	if err := ruleSkuCombinationPreserved(oldSkus, oldSKUMap, req); err != nil {
		return nil, err
	}

	return oldSKUMap, nil
}

// ruleCategoryImmutable Rule 1：类目决定属性模板与规格体系，创建后不可修改。
func ruleCategoryImmutable(old *entity.Products, req *v1.ProductsUpdateFullReq) error {
	if old.CategoryId != req.CategoryId {
		return gerror.NewCode(errcode.Code(400), "商品类目不可修改")
	}
	return nil
}

// ruleSpecGroupPreserved Rule 2：无旧规格组时跳过；
// 已有 SKU 不允许改动规格文本快照，新增 SKU 的规格维度数必须与旧规格一致。
func ruleSpecGroupPreserved(oldSkus []*entity.Skus, oldSKUMap map[int64]*entity.Skus, oldSpecAttrIDs map[int64]string, req *v1.ProductsUpdateFullReq) error {
	if len(oldSpecAttrIDs) == 0 {
		return nil
	}
	for _, sku := range req.SKUs {
		if sku.Id > 0 {
			if old, ok := oldSKUMap[sku.Id]; ok && old.SpecSummary != sku.SpecSummary {
				return gerror.NewCode(errcode.Code(400),
					fmt.Sprintf("不允许修改已有 SKU(%d) 的规格文本", sku.Id))
			}
			continue
		}
		// 新 SKU 的规格维度数必须与旧规格匹配
		oldDim := 1
		if len(oldSkus) > 0 {
			oldDim = len(strings.Split(oldSkus[0].SpecSummary, " / "))
		}
		newDim := len(strings.Split(sku.SpecSummary, " / "))
		if newDim != oldDim {
			return gerror.NewCode(errcode.Code(400),
				fmt.Sprintf("新 SKU 规格维度(%d)与旧规格维度(%d)不匹配: %s", newDim, oldDim, sku.SpecSummary))
		}
	}
	return nil
}

// ruleSkuCombinationPreserved Rule 3：请求中的 SKU ID 必须属于该商品，且旧 SKU 组合一个都不能少。
func ruleSkuCombinationPreserved(oldSkus []*entity.Skus, oldSKUMap map[int64]*entity.Skus, req *v1.ProductsUpdateFullReq) error {
	dtoSkuIDSet := make(map[int64]struct{}, len(req.SKUs))
	for _, sku := range req.SKUs {
		if sku.Id > 0 {
			if _, exists := oldSKUMap[sku.Id]; !exists {
				return gerror.NewCode(gcode.CodeInvalidParameter,
					fmt.Sprintf("SKU ID %d 不属于该商品", sku.Id))
			}
			dtoSkuIDSet[sku.Id] = struct{}{}
		}
	}
	for _, oldSku := range oldSkus {
		if _, exists := dtoSkuIDSet[oldSku.Id]; !exists {
			return gerror.NewCode(errcode.Code(400),
				fmt.Sprintf("不允许移除旧的规格组合: %s", oldSku.SpecSummary))
		}
	}
	return nil
}
