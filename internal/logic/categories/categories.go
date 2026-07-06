package categories

import (
	"context"

	"github.com/bytedance/sonic"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/api/categories/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sCategories struct{}

func init() {
	service.RegisterCategories(&sCategories{})
}

// List 类目平铺列表（支持分页）
func (s *sCategories) List(ctx context.Context, req *v1.CategoryListReq) (res *v1.CategoryListRes, err error) {
	var (
		page = req.Page
		size = req.PageSize
	)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	// 无筛选条件时走 Lua 脚本
	if req.ParentId <= 0 && req.Status <= 0 && req.Name == "" && req.Level <= 0 {
		list, total, err := getCategoryPage(ctx, page, size)
		if err == nil && total > 0 {
			return &v1.CategoryListRes{List: list, Total: total}, nil
		}
		if ctx.Err() != nil {
			return nil, gerror.NewCode(errcode.Code(57), "请求已取消")
		}
		ensureCategoryCache(ctx)
		list, total, err = getCategoryPage(ctx, page, size)
		if err == nil && total > 0 {
			return &v1.CategoryListRes{List: list, Total: total}, nil
		}
		var dbAll []*entity.Categories
		if err := dao.Categories.Ctx(ctx).OrderAsc(dao.Categories.Columns().SortOrder).Scan(&dbAll); err != nil {
			return nil, err
		}
		if len(dbAll) == 0 {
			return &v1.CategoryListRes{List: make([]*entity.Categories, 0), Total: 0}, nil
		}
		return paginateCategories(page, size, dbAll), nil
	}

	// 有筛选条件时直接查询数据库
	var (
		m    = dao.Categories.Ctx(ctx)
		list []*entity.Categories
	)
	if req.ParentId > 0 {
		m = m.Where(dao.Categories.Columns().ParentId, req.ParentId)
	}
	if req.Status > 0 {
		m = m.Where(dao.Categories.Columns().Status, req.Status)
	}
	if req.Name != "" {
		m = m.Where(dao.Categories.Columns().Name+" LIKE ?", "%"+req.Name+"%")
	}
	if req.Level > 0 {
		m = m.Where(dao.Categories.Columns().Level, req.Level)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.CategoryListRes{
			List:  make([]*entity.Categories, 0),
			Total: 0,
		}, nil
	}

	err = m.Page(page, size).OrderAsc(dao.Categories.Columns().SortOrder).OrderDesc(dao.Categories.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.CategoryListRes{
		List:  list,
		Total: total,
	}, nil
}

func paginateCategories(page, size int, all []*entity.Categories) *v1.CategoryListRes {
	total := len(all)
	start := (page - 1) * size
	if start >= total {
		return &v1.CategoryListRes{List: make([]*entity.Categories, 0), Total: total}
	}
	end := start + size
	if end > total {
		end = total
	}
	return &v1.CategoryListRes{List: all[start:end], Total: total}
}

// All 所有类目列表（Lua 全量读取）
func (s *sCategories) All(ctx context.Context, req *v1.CategoryAllReq) (res *v1.CategoryAllRes, err error) {
	// ZRANGE 0 -1 一次性拿全部 ID
	v, err := g.Redis().Do(ctx, "ZRANGE", categoryIdsKey, 0, -1)
	if err == nil && !v.IsNil() && len(v.Vars()) > 0 {
		args := make([]interface{}, len(v.Vars()))
		for i, idv := range v.Vars() {
			args[i] = cacheKeyCategory(idv.Int64())
		}
		mvals, err := g.Redis().Do(ctx, "MGET", args...)
		if err == nil && !mvals.IsNil() {
			var list []*entity.Categories
			for _, mv := range mvals.Vars() {
				if mv.IsNil() {
					list = nil
					break
				}
				var c entity.Categories
				if err := sonic.Unmarshal(mv.Bytes(), &c); err != nil {
					list = nil
					break
				}
				list = append(list, &c)
			}
			if list != nil {
				return &v1.CategoryAllRes{List: list}, nil
			}
		}
	}
	if ctx.Err() != nil {
		return nil, gerror.NewCode(errcode.Code(57), "请求已取消")
	}
	ensureCategoryCache(ctx)
	// 兜底：查 DB
	var list []*entity.Categories
	if err := dao.Categories.Ctx(ctx).OrderAsc(dao.Categories.Columns().SortOrder).OrderDesc(dao.Categories.Columns().Id).Scan(&list); err != nil {
		return nil, err
	}
	return &v1.CategoryAllRes{List: list}, nil
}

// NonRoot 非根类目列表（parent_id != 0，用于品牌绑定联动）
func (s *sCategories) NonRoot(ctx context.Context, req *v1.CategoryNonRootReq) (res *v1.CategoryNonRootRes, err error) {
	var list []*entity.Categories
	err = dao.Categories.Ctx(ctx).
		Where(dao.Categories.Columns().ParentId+" > ?", 0).

		OrderAsc(dao.Categories.Columns().SortOrder).
		OrderDesc(dao.Categories.Columns().Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = make([]*entity.Categories, 0)
	}
	return &v1.CategoryNonRootRes{List: list}, nil
}

// Root 根类目列表
func (s *sCategories) Root(ctx context.Context, req *v1.CategoryRootReq) (res *v1.CategoryRootRes, err error) {
	var (
		m    = dao.Categories.Ctx(ctx)
		list []*entity.Categories
	)
	err = m.Where(dao.Categories.Columns().ParentId, 0).OrderAsc(dao.Categories.Columns().SortOrder).OrderDesc(dao.Categories.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.CategoryRootRes{List: list}, nil
}

// Children 子类目列表
func (s *sCategories) Children(ctx context.Context, req *v1.CategoryChildrenReq) (res *v1.CategoryChildrenRes, err error) {
	var (
		m    = dao.Categories.Ctx(ctx)
		list []*entity.Categories
	)
	err = m.Where(dao.Categories.Columns().ParentId, req.Id).OrderAsc(dao.Categories.Columns().SortOrder).OrderDesc(dao.Categories.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.CategoryChildrenRes{List: list}, nil
}

// Level 层级类目列表
func (s *sCategories) Level(ctx context.Context, req *v1.CategoryLevelReq) (res *v1.CategoryLevelRes, err error) {
	var (
		m    = dao.Categories.Ctx(ctx)
		list []*entity.Categories
	)
	err = m.Where(dao.Categories.Columns().Level, req.Level).OrderAsc(dao.Categories.Columns().SortOrder).OrderDesc(dao.Categories.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.CategoryLevelRes{List: list}, nil
}

// Tree 类目树形结构
func (s *sCategories) Tree(ctx context.Context, req *v1.CategoryTreeReq) (res *v1.CategoryTreeRes, err error) {
	var (
		m   = dao.Categories.Ctx(ctx)
		all []*entity.Categories
	)
	if req.Status > 0 {
		m = m.Where(dao.Categories.Columns().Status, req.Status)
	}
	err = m.OrderAsc(dao.Categories.Columns().SortOrder).OrderAsc(dao.Categories.Columns().Id).Scan(&all)
	if err != nil {
		return nil, err
	}
	return &v1.CategoryTreeRes{Tree: buildTree(all, 0)}, nil
}

// buildTree 递归构建树形结构
func buildTree(nodes []*entity.Categories, parentId int64) []*v1.CategoryTreeItem {
	var tree []*v1.CategoryTreeItem
	for _, n := range nodes {
		if n.ParentId == parentId {
			item := &v1.CategoryTreeItem{
				Categories: n,
				Children:   buildTree(nodes, n.Id),
			}
			tree = append(tree, item)
		}
	}
	return tree
}

// Detail 类目详情（缓存旁路）
func (s *sCategories) Detail(ctx context.Context, req *v1.CategoryDetailReq) (res *v1.CategoryDetailRes, err error) {
	cached, err := getCategoryEntityCache(ctx, req.Id)
	if err == nil && cached != nil {
		return &v1.CategoryDetailRes{Categories: cached}, nil
	}
	if ctx.Err() != nil {
		return nil, gerror.NewCode(errcode.Code(57), "请求已取消")
	}

	var entity *entity.Categories
	err = dao.Categories.Ctx(ctx).Where(dao.Categories.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errcode.ErrCategoryNotFound
	}
	if err := setCategoryEntityCache(context.Background(), entity); err != nil {
		g.Log().Warning(ctx, "setCategoryEntityCache failed: %v", err)
	}
	return &v1.CategoryDetailRes{Categories: entity}, nil
}

// Create 新增类目
func (s *sCategories) Create(ctx context.Context, req *v1.CategoryCreateReq) (res *v1.CategoryCreateRes, err error) {
	result, err := dao.Categories.Ctx(ctx).Insert(do.Categories{
		Name:      req.Name,
		ParentId:  req.ParentId,
		Level:     req.Level,
		Path:      req.Path,
		IconUrl:   req.IconUrl,
		SortOrder: req.SortOrder,
		Status:    req.Status,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	addCategoryToIndex(context.Background(), id, req.SortOrder)
	return &v1.CategoryCreateRes{Id: id}, nil
}

// Update 更新类目
func (s *sCategories) Update(ctx context.Context, req *v1.CategoryUpdateReq) (res *v1.CategoryUpdateRes, err error) {
	count, err := dao.Categories.Ctx(ctx).Where(dao.Categories.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrCategoryNotFound
	}
	_, err = dao.Categories.Ctx(ctx).Data(do.Categories{
		Name:      req.Name,
		ParentId:  req.ParentId,
		Level:     req.Level,
		Path:      req.Path,
		IconUrl:   req.IconUrl,
		SortOrder: req.SortOrder,
		Status:    req.Status,
	}).Where(dao.Categories.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	g.Redis().Do(context.Background(), "ZADD", categoryIdsKey, encodeCategoryScore(req.SortOrder, req.Id), req.Id)
	delCategoryEntityCache(context.Background(), req.Id)
	return &v1.CategoryUpdateRes{}, nil
}

// Delete 删除类目（软删除）
func (s *sCategories) Delete(ctx context.Context, req *v1.CategoryDeleteReq) (res *v1.CategoryDeleteRes, err error) {
	_, err = dao.Categories.Ctx(ctx).Where(dao.Categories.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	removeCategoryFromIndex(context.Background(), req.Id)
	delCategoryEntityCache(context.Background(), req.Id)
	return &v1.CategoryDeleteRes{}, nil
}
