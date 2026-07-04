package categories

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/api/categories/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sCategories struct{}

func init() {
	service.RegisterCategories(&sCategories{})
}

// List 类目平铺列表（支持分页）
func (s *sCategories) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	var (
		m     = dao.Categories.Ctx(ctx)
		list  []*entity.Categories
		total int
		page  = req.Page
		size  = req.PageSize
	)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	if req.ParentId > 0 {
		m = m.Where(dao.Categories.Columns().ParentId, req.ParentId)
	}
	if req.Status > 0 {
		m = m.Where(dao.Categories.Columns().Status, req.Status)
	}

	total, err = m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.ListRes{
			List:  make([]*entity.Categories, 0),
			Total: 0,
		}, nil
	}

	err = m.Page(page, size).OrderAsc(dao.Categories.Columns().SortOrder).OrderDesc(dao.Categories.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.ListRes{
		List:  list,
		Total: total,
	}, nil
}

// All 所有类目列表（缓存旁路）
func (s *sCategories) All(ctx context.Context, req *v1.AllReq) (res *v1.AllRes, err error) {
	// 尝试从缓存读取
	cached, err := getCategoryAllCache(ctx)
	if err == nil && cached != nil {
		return &v1.AllRes{List: cached}, nil
	}

	var (
		m    = dao.Categories.Ctx(ctx)
		list []*entity.Categories
	)
	err = m.OrderAsc(dao.Categories.Columns().SortOrder).OrderDesc(dao.Categories.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}

	// 回写缓存
	if err := setCategoryAllCache(ctx, list); err != nil {
		g.Log().Warning(ctx, "setCategoryAllCache failed: %v", err)
	}

	return &v1.AllRes{List: list}, nil
}

// Root 根类目列表
func (s *sCategories) Root(ctx context.Context, req *v1.RootReq) (res *v1.RootRes, err error) {
	var (
		m    = dao.Categories.Ctx(ctx)
		list []*entity.Categories
	)
	err = m.Where(dao.Categories.Columns().ParentId, 0).OrderAsc(dao.Categories.Columns().SortOrder).OrderDesc(dao.Categories.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.RootRes{List: list}, nil
}

// Children 子类目列表
func (s *sCategories) Children(ctx context.Context, req *v1.ChildrenReq) (res *v1.ChildrenRes, err error) {
	var (
		m    = dao.Categories.Ctx(ctx)
		list []*entity.Categories
	)
	err = m.Where(dao.Categories.Columns().ParentId, req.Id).OrderAsc(dao.Categories.Columns().SortOrder).OrderDesc(dao.Categories.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.ChildrenRes{List: list}, nil
}

// Level 层级类目列表
func (s *sCategories) Level(ctx context.Context, req *v1.LevelReq) (res *v1.LevelRes, err error) {
	var (
		m    = dao.Categories.Ctx(ctx)
		list []*entity.Categories
	)
	err = m.Where(dao.Categories.Columns().Level, req.Level).OrderAsc(dao.Categories.Columns().SortOrder).OrderDesc(dao.Categories.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.LevelRes{List: list}, nil
}

// Tree 类目树形结构
func (s *sCategories) Tree(ctx context.Context, req *v1.TreeReq) (res *v1.TreeRes, err error) {
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
	return &v1.TreeRes{Tree: buildTree(all, 0)}, nil
}

// buildTree 递归构建树形结构
func buildTree(nodes []*entity.Categories, parentId int64) []*v1.TreeItem {
	var tree []*v1.TreeItem
	for _, n := range nodes {
		if n.ParentId == parentId {
			item := &v1.TreeItem{
				Categories: n,
				Children:   buildTree(nodes, n.Id),
			}
			tree = append(tree, item)
		}
	}
	return tree
}

// Detail 类目详情（缓存旁路）
func (s *sCategories) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	// 尝试从缓存读取
	cached, err := getCategoryEntityCache(ctx, req.Id)
	if err == nil && cached != nil {
		return &v1.DetailRes{Categories: cached}, nil
	}

	var entity *entity.Categories
	err = dao.Categories.Ctx(ctx).Where(dao.Categories.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "类目不存在")
	}

	// 回写缓存
	if err := setCategoryEntityCache(ctx, entity); err != nil {
		g.Log().Warning(ctx, "setCategoryEntityCache failed: %v", err)
	}

	return &v1.DetailRes{Categories: entity}, nil
}

// Create 新增类目
func (s *sCategories) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
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
	// 缓存失效
	delCategoryAllCache(ctx)
	return &v1.CreateRes{Id: id}, nil
}

// Update 更新类目
func (s *sCategories) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	count, err := dao.Categories.Ctx(ctx).Where(dao.Categories.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "类目不存在")
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
	// 缓存失效
	delCategoryEntityCache(ctx, req.Id)
	delCategoryAllCache(ctx)
	return &v1.UpdateRes{}, nil
}

// Delete 删除类目（软删除）
func (s *sCategories) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	_, err = dao.Categories.Ctx(ctx).Where(dao.Categories.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	// 缓存失效
	delCategoryEntityCache(ctx, req.Id)
	delCategoryAllCache(ctx)
	return &v1.DeleteRes{}, nil
}
