package products

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/api/products/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
)

// 本文件是属性 EAV（sp_product_attributes、sp_attributes、sp_attribute_values、sp_category_attributes）
// 与图文描述（sp_product_descriptions）的数据访问层。

// ── 属性定义 ────────────────────────────────────────────────────────────

// getAttribute 按 ID 查询属性定义；不存在时返回 (nil, nil)。
func getAttribute(ctx context.Context, id int64) (*entity.Attributes, error) {
	var attr *entity.Attributes
	if err := dao.Attributes.Ctx(ctx).Where(dao.Attributes.Columns().Id, id).Scan(&attr); err != nil {
		return nil, err
	}
	return attr, nil
}

// getAttributeValue 按 ID 查询属性值字典项；不存在时返回 (nil, nil)。
func getAttributeValue(ctx context.Context, id int64) (*entity.AttributeValues, error) {
	var v *entity.AttributeValues
	if err := dao.AttributeValues.Ctx(ctx).Where(dao.AttributeValues.Columns().Id, id).Scan(&v); err != nil {
		return nil, err
	}
	return v, nil
}

// listCategoryAttrs 查询类目关联属性，返回 属性名→属性ID 与 属性名→排序 两个映射。
// 属性可能被多类目共享（category_id 仅标记归属类目），故经 sp_category_attributes 关联查询。
// 该结果仅用于详情富化，查询失败返回空映射而不阻断详情。
func listCategoryAttrs(ctx context.Context, categoryId int64) (map[string]int64, map[string]int) {
	idByName := make(map[string]int64)
	orderByName := make(map[string]int)
	if categoryId <= 0 {
		return idByName, orderByName
	}

	var catAttrs []*entity.Attributes
	err := dao.Attributes.Ctx(ctx).
		Fields("sp_attributes.*").
		InnerJoin("sp_category_attributes ca", "ca.attribute_id = sp_attributes.id").
		Where("ca.category_id", categoryId).
		OrderAsc("ca.sort_order").
		Scan(&catAttrs)
	if err != nil {
		return idByName, orderByName
	}
	for i, a := range catAttrs {
		idByName[a.Name] = a.Id
		orderByName[a.Name] = i
	}
	return idByName, orderByName
}

// ── 商品属性读取 ────────────────────────────────────────────────────────

// findAttrsWithName 查询商品属性并 left join 属性表/属性值表获取名称，按属性聚合多值。
// 空结果返回空切片而非 nil。
func findAttrsWithName(ctx context.Context, productId int64) ([]v1.ProductAttrDetailResponse, error) {
	type attrRow struct {
		AttributeId      int64  `orm:"attribute_id"`
		AttributeName    string `orm:"attribute_name"`
		AttributeValueId int64  `orm:"attribute_value_id"`
		AttributeValue   string `orm:"attribute_value"`
		SortOrder        int    `orm:"sort_order"`
	}
	var rows []attrRow
	err := g.DB().Model("sp_product_attributes pa").
		Fields("pa.attribute_id", "a.name AS attribute_name", "pa.attribute_value_id", "COALESCE(av.value, pa.value) AS attribute_value", "pa.sort_order").
		LeftJoin("sp_attributes a", "a.id = pa.attribute_id").
		LeftJoin("sp_attribute_values av", "av.id = pa.attribute_value_id").
		Where("pa.product_id", productId).
		Where("pa.deleted_at IS NULL").
		Order("pa.sort_order ASC, pa.id ASC").
		Scan(&rows)
	if err != nil {
		return nil, err
	}

	seen := make(map[int64]int)
	result := make([]v1.ProductAttrDetailResponse, 0)
	for _, r := range rows {
		var vals []*v1.ProductAttrValueResponse
		if r.AttributeValueId > 0 || r.AttributeValue != "" {
			vals = []*v1.ProductAttrValueResponse{{Id: r.AttributeValueId, Value: r.AttributeValue}}
		}
		if idx, ok := seen[r.AttributeId]; ok {
			result[idx].Values = append(result[idx].Values, vals...)
		} else {
			seen[r.AttributeId] = len(result)
			result = append(result, v1.ProductAttrDetailResponse{
				AttributeId:   r.AttributeId,
				AttributeName: r.AttributeName,
				Values:        vals,
				SortOrder:     r.SortOrder,
			})
		}
	}
	if result == nil {
		result = make([]v1.ProductAttrDetailResponse, 0)
	}
	return result, nil
}

// listProductAttrsWithValue 查询商品绑定的属性原始行（含关联 ID 与值），供属性管理接口使用。
func listProductAttrsWithValue(ctx context.Context, productId int64) ([]*v1.ProductsAttributeItem, error) {
	type attrRow struct {
		Id               int64  `orm:"id"`
		AttributeId      int64  `orm:"attribute_id"`
		AttributeName    string `orm:"attribute_name"`
		AttributeValueId int64  `orm:"attribute_value_id"`
		Value            string `orm:"value"`
	}
	var rows []attrRow
	err := g.DB().Model("sp_product_attributes pa").
		Fields("pa.id", "pa.attribute_id", "a.name AS attribute_name", "pa.attribute_value_id", "pa.value").
		LeftJoin("sp_attributes a", "a.id = pa.attribute_id").
		Where("pa.product_id", productId).
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
	return list, nil
}

// ── 商品属性写入 ────────────────────────────────────────────────────────

// buildAttrRows 将扁平属性列表展开为单行记录，每行一个 (product_id, attribute_id, attribute_value_id)。
// 适配 uk_product_attribute(product_id, attribute_id, attribute_value_id) 唯一索引。
func buildAttrRows(ctx context.Context, tx gdb.TX, productId int64, list []v1.CreateProductAttrItem) []do.ProductAttributes {
	var rows []do.ProductAttributes

	allValueIDs := make([]int64, 0)
	for _, attr := range list {
		for _, vid := range attr.AttributeValueId {
			if vid > 0 {
				allValueIDs = append(allValueIDs, vid)
			}
		}
	}

	valMap := make(map[int64]string)
	if len(allValueIDs) > 0 {
		type valRow struct {
			Id    int64  `orm:"id"`
			Value string `orm:"value"`
		}
		var vals []valRow
		err := tx.Model("sp_attribute_values").Where("id IN (?)", allValueIDs).Scan(&vals)
		if err == nil {
			for _, v := range vals {
				valMap[v.Id] = v.Value
			}
		}
	}

	for _, attr := range list {
		if len(attr.AttributeValueId) > 0 {
			for _, vid := range attr.AttributeValueId {
				if vid <= 0 {
					continue
				}
				v := valMap[vid]
				if v == "" {
					v = attr.Value
				}
				rows = append(rows, do.ProductAttributes{
					ProductId:        productId,
					AttributeId:      attr.AttributeId,
					AttributeValueId: vid,
					Value:            v,
				})
			}
		} else if attr.Value != "" {
			rows = append(rows, do.ProductAttributes{
				ProductId:   productId,
				AttributeId: attr.AttributeId,
				Value:       attr.Value,
			})
		}
	}
	return rows
}

// replaceProductAttrs 在事务中全量替换商品属性：硬删旧行后重插，避免软删除行占用唯一索引。
func replaceProductAttrs(ctx context.Context, tx gdb.TX, productId int64, list []v1.CreateProductAttrItem) error {
	_, err := tx.Model("sp_product_attributes").
		Where("product_id", productId).
		Unscoped().Delete()
	if err != nil {
		return err
	}
	for _, data := range buildAttrRows(ctx, tx, productId, list) {
		if _, err = tx.Model("sp_product_attributes").Insert(data); err != nil {
			return err
		}
	}
	return nil
}

// ── 图文描述 ────────────────────────────────────────────────────────────

// getDescription 查询商品图文描述；不存在时返回 (nil, nil)。
func getDescription(ctx context.Context, productId int64) (*entity.ProductDescriptions, error) {
	var d *entity.ProductDescriptions
	err := dao.ProductDescriptions.Ctx(ctx).
		Where(dao.ProductDescriptions.Columns().ProductId, productId).
		Scan(&d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// insertDescription 在事务中新增商品图文描述。
func insertDescription(ctx context.Context, tx gdb.TX, productId int64, description, mobileDescription string) error {
	_, err := tx.Model("sp_product_descriptions").Insert(do.ProductDescriptions{
		ProductId:         productId,
		Description:       description,
		MobileDescription: mobileDescription,
	})
	return err
}

// upsertDescription 在事务中写入商品图文描述：已有记录则更新，否则新增。
func upsertDescription(ctx context.Context, tx gdb.TX, productId int64, description, mobileDescription string) error {
	count, err := tx.Model("sp_product_descriptions").Where("product_id", productId).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		_, err = tx.Model("sp_product_descriptions").Where("product_id", productId).Update(do.ProductDescriptions{
			Description:       description,
			MobileDescription: mobileDescription,
		})
		return err
	}
	return insertDescription(ctx, tx, productId, description, mobileDescription)
}

// deleteDescription 在事务中删除商品图文描述（表含 deleted_at，自动软删除）。
func deleteDescription(ctx context.Context, tx gdb.TX, productId int64) error {
	_, err := tx.Model("sp_product_descriptions").Where("product_id", productId).Delete()
	return err
}
