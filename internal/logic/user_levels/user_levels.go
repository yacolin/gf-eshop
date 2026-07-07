package user_levels

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"gf-eshop/api/user_levels/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sUserLevels struct{}

func init() {
	service.RegisterUserLevels(&sUserLevels{})
}

func (s *sUserLevels) List(ctx context.Context, req *v1.LevelsListReq) (res *v1.LevelsListRes, err error) {
	page := req.Page
	size := req.PageSize
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	m := dao.Levels.Ctx(ctx)
	if req.Name != "" {
		m = m.WhereLike(dao.Levels.Columns().Name, "%"+req.Name+"%")
	}
	if req.Status != nil {
		m = m.Where(dao.Levels.Columns().Status, *req.Status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.LevelsListRes{
			List:  make([]*entity.Levels, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.Levels
	err = m.OrderAsc(dao.Levels.Columns().SortOrder).OrderAsc(dao.Levels.Columns().Level).Page(page, size).Scan(&list)
	if err != nil {
		return nil, err
	}

	return &v1.LevelsListRes{List: list, Total: total}, nil
}

func (s *sUserLevels) Detail(ctx context.Context, req *v1.LevelsDetailReq) (res *v1.LevelsDetailRes, err error) {
	var level *entity.Levels
	err = dao.Levels.Ctx(ctx).Where(dao.Levels.Columns().Id, req.Id).Scan(&level)
	if err != nil {
		return nil, err
	}
	if level == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "等级不存在")
	}
	return &v1.LevelsDetailRes{Levels: level}, nil
}

func (s *sUserLevels) Create(ctx context.Context, req *v1.LevelsCreateReq) (res *v1.LevelsCreateRes, err error) {
	result, err := dao.Levels.Ctx(ctx).Insert(do.Levels{
		Name:             req.Name,
		Level:            req.Level,
		MinPoints:        req.MinPoints,
		MaxPoints:        req.MaxPoints,
		DiscountRate:     req.DiscountRate,
		FreeShipping:     req.FreeShipping,
		PointsMultiplier: req.PointsMultiplier,
		Benefits:         req.Benefits,
		Status:           req.Status,
		SortOrder:        req.SortOrder,
		CreatedAt:        gtime.Now(),
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.LevelsCreateRes{Id: id}, nil
}

func (s *sUserLevels) Update(ctx context.Context, req *v1.LevelsUpdateReq) (res *v1.LevelsUpdateRes, err error) {
	count, err := dao.Levels.Ctx(ctx).Where(dao.Levels.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "等级不存在")
	}

	data := do.Levels{
		Name:             req.Name,
		Level:            req.Level,
		MinPoints:        req.MinPoints,
		MaxPoints:        req.MaxPoints,
		DiscountRate:     req.DiscountRate,
		FreeShipping:     req.FreeShipping,
		PointsMultiplier: req.PointsMultiplier,
		Benefits:         req.Benefits,
		Status:           req.Status,
		SortOrder:        req.SortOrder,
		UpdatedAt:        gtime.Now(),
	}
	_, err = dao.Levels.Ctx(ctx).Where(dao.Levels.Columns().Id, req.Id).Update(data)
	if err != nil {
		return nil, err
	}
	return &v1.LevelsUpdateRes{}, nil
}

func (s *sUserLevels) Delete(ctx context.Context, req *v1.LevelsDeleteReq) (res *v1.LevelsDeleteRes, err error) {
	count, err := dao.Levels.Ctx(ctx).Where(dao.Levels.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "等级不存在")
	}
	_, err = dao.Levels.Ctx(ctx).Where(dao.Levels.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &v1.LevelsDeleteRes{}, nil
}
