package departments

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/api/departments/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sDepartments struct{}

func init() {
	service.RegisterDepartments(&sDepartments{})
}

// List 部门平铺列表（支持分页）
func (s *sDepartments) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
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
	if req.ParentId <= 0 && (req.Status == nil || *req.Status <= 0) && req.Name == "" {
		list, total, err := getDepartmentPage(ctx, page, size)
		if err == nil && total > 0 {
			return &v1.ListRes{List: list, Total: total}, nil
		}
		if ctx.Err() != nil {
			return nil, err
		}
		ensureDepartmentCache(ctx)
		list, total, err = getDepartmentPage(ctx, page, size)
		if err == nil && total > 0 {
			return &v1.ListRes{List: list, Total: total}, nil
		}
		// 兜底：直接查库
		var dbAll []*entity.Departments
		if err := dao.Departments.Ctx(ctx).OrderAsc(dao.Departments.Columns().SortOrder).Scan(&dbAll); err != nil {
			return nil, err
		}
		if len(dbAll) == 0 {
			return &v1.ListRes{List: make([]*entity.Departments, 0), Total: 0}, nil
		}
		return paginateDepartments(page, size, dbAll), nil
	}

	// 有筛选条件时直接查询数据库
	var (
		m    = dao.Departments.Ctx(ctx)
		list []*entity.Departments
	)
	if req.ParentId > 0 {
		m = m.Where(dao.Departments.Columns().ParentId, req.ParentId)
	}
	if req.Status != nil {
		m = m.Where(dao.Departments.Columns().Status, *req.Status)
	}
	if req.Name != "" {
		m = m.Where(dao.Departments.Columns().Name+" LIKE ?", "%"+req.Name+"%")
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.ListRes{
			List:  make([]*entity.Departments, 0),
			Total: 0,
		}, nil
	}

	err = m.Page(page, size).OrderAsc(dao.Departments.Columns().SortOrder).OrderDesc(dao.Departments.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.ListRes{
		List:  list,
		Total: total,
	}, nil
}

func paginateDepartments(page, size int, all []*entity.Departments) *v1.ListRes {
	total := len(all)
	start := (page - 1) * size
	if start >= total {
		return &v1.ListRes{List: make([]*entity.Departments, 0), Total: total}
	}
	end := start + size
	if end > total {
		end = total
	}
	return &v1.ListRes{List: all[start:end], Total: total}
}

// All 所有部门列表
func (s *sDepartments) All(ctx context.Context, req *v1.AllReq) (res *v1.AllRes, err error) {
	v, err := g.Redis().Do(ctx, "ZRANGE", departmentIdsKey, 0, -1)
	if err == nil && !v.IsNil() && len(v.Vars()) > 0 {
		args := make([]interface{}, len(v.Vars()))
		for i, idv := range v.Vars() {
			args[i] = cacheKeyDepartment(idv.Int64())
		}
		mvals, err := g.Redis().Do(ctx, "MGET", args...)
		if err == nil && !mvals.IsNil() {
			var list []*entity.Departments
			for _, mv := range mvals.Vars() {
				if mv.IsNil() {
					list = nil
					break
				}
				var d entity.Departments
				if err := sonic.Unmarshal(mv.Bytes(), &d); err != nil {
					list = nil
					break
				}
				list = append(list, &d)
			}
			if list != nil {
				return &v1.AllRes{List: list}, nil
			}
		}
	}
	if ctx.Err() != nil {
		return nil, err
	}
	ensureDepartmentCache(ctx)
	var list []*entity.Departments
	if err := dao.Departments.Ctx(ctx).OrderAsc(dao.Departments.Columns().SortOrder).Scan(&list); err != nil {
		return nil, err
	}
	return &v1.AllRes{List: list}, nil
}

// Children 子部门列表
func (s *sDepartments) Children(ctx context.Context, req *v1.ChildrenReq) (res *v1.ChildrenRes, err error) {
	var list []*entity.Departments
	err = dao.Departments.Ctx(ctx).
		Where(dao.Departments.Columns().ParentId, req.Id).
		OrderAsc(dao.Departments.Columns().SortOrder).
		OrderDesc(dao.Departments.Columns().Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.ChildrenRes{List: list}, nil
}

// Tree 部门树形结构
func (s *sDepartments) Tree(ctx context.Context, req *v1.TreeReq) (res *v1.TreeRes, err error) {
	var (
		m   = dao.Departments.Ctx(ctx)
		all []*entity.Departments
	)
	if req.Status != nil {
		m = m.Where(dao.Departments.Columns().Status, *req.Status)
	}
	err = m.OrderAsc(dao.Departments.Columns().SortOrder).OrderAsc(dao.Departments.Columns().Id).Scan(&all)
	if err != nil {
		return nil, err
	}
	return &v1.TreeRes{Tree: buildTree(all, 0)}, nil
}

// buildTree 递归构建树形结构
func buildTree(nodes []*entity.Departments, parentId int64) []*v1.TreeItem {
	var tree []*v1.TreeItem
	for _, n := range nodes {
		if n.ParentId == parentId {
			item := &v1.TreeItem{
				Departments: n,
				Children:    buildTree(nodes, n.Id),
			}
			tree = append(tree, item)
		}
	}
	return tree
}

// Detail 部门详情（缓存旁路）
func (s *sDepartments) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	cached, err := getDepartmentEntityCache(ctx, req.Id)
	if err == nil && cached != nil {
		return &v1.DetailRes{Departments: cached}, nil
	}
	if ctx.Err() != nil {
		return nil, err
	}

	var entity *entity.Departments
	err = dao.Departments.Ctx(ctx).Where(dao.Departments.Columns().Id, req.Id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errcode.ErrDepartmentNotFound
	}
	if err := setDepartmentEntityCache(context.Background(), entity); err != nil {
		g.Log().Warning(ctx, "setDepartmentEntityCache failed: %v", err)
	}
	return &v1.DetailRes{Departments: entity}, nil
}

// Create 新增部门
func (s *sDepartments) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	result, err := dao.Departments.Ctx(ctx).Insert(do.Departments{
		Name:      req.Name,
		ParentId:  req.ParentId,
		SortOrder: req.SortOrder,
		Status:    req.Status,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	addDepartmentToIndex(context.Background(), id, req.SortOrder)
	return &v1.CreateRes{Id: id}, nil
}

// Update 更新部门
func (s *sDepartments) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	count, err := dao.Departments.Ctx(ctx).Where(dao.Departments.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errcode.ErrDepartmentNotFound
	}

	_, err = dao.Departments.Ctx(ctx).Data(do.Departments{
		Name:      req.Name,
		ParentId:  req.ParentId,
		SortOrder: req.SortOrder,
		Status:    req.Status,
	}).Where(dao.Departments.Columns().Id, req.Id).Update()
	if err != nil {
		return nil, err
	}
	g.Redis().Do(context.Background(), "ZADD", departmentIdsKey, encodeDepartmentScore(req.SortOrder, req.Id), req.Id)
	delDepartmentEntityCache(context.Background(), req.Id)
	return &v1.UpdateRes{}, nil
}

// Delete 删除部门（软删除）
func (s *sDepartments) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	_, err = dao.Departments.Ctx(ctx).Where(dao.Departments.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	removeDepartmentFromIndex(context.Background(), req.Id)
	delDepartmentEntityCache(context.Background(), req.Id)
	return &v1.DeleteRes{}, nil
}
