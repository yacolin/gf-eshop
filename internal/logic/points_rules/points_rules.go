package points_rules

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"gf-eshop/api/points_rules/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sPointsRules struct{}

func init() {
	service.RegisterPointsRules(&sPointsRules{})
}

func (s *sPointsRules) List(ctx context.Context, req *v1.PointsRulesListReq) (res *v1.PointsRulesListRes, err error) {
	page := req.Page
	size := req.PageSize
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	m := dao.PointsRules.Ctx(ctx)
	if req.RuleKey != "" {
		m = m.Where(dao.PointsRules.Columns().RuleKey, req.RuleKey)
	}
	if req.Status != nil {
		m = m.Where(dao.PointsRules.Columns().Status, *req.Status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.PointsRulesListRes{
			List:  make([]*entity.PointsRules, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.PointsRules
	err = m.OrderAsc(dao.PointsRules.Columns().SortOrder).Page(page, size).Scan(&list)
	if err != nil {
		return nil, err
	}

	return &v1.PointsRulesListRes{List: list, Total: total}, nil
}

func (s *sPointsRules) Detail(ctx context.Context, req *v1.PointsRulesDetailReq) (res *v1.PointsRulesDetailRes, err error) {
	var rule entity.PointsRules
	err = dao.PointsRules.Ctx(ctx).Where(dao.PointsRules.Columns().Id, req.Id).Scan(&rule)
	if err != nil {
		return nil, err
	}
	return &v1.PointsRulesDetailRes{PointsRules: &rule}, nil
}

func (s *sPointsRules) Create(ctx context.Context, req *v1.PointsRulesCreateReq) (res *v1.PointsRulesCreateRes, err error) {
	result, err := dao.PointsRules.Ctx(ctx).Insert(do.PointsRules{
		Name:        req.Name,
		RuleKey:     req.RuleKey,
		RuleValue:   req.RuleValue,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
	})
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &v1.PointsRulesCreateRes{Id: int(id)}, nil
}

func (s *sPointsRules) Update(ctx context.Context, req *v1.PointsRulesUpdateReq) (res *v1.PointsRulesUpdateRes, err error) {
	count, err := dao.PointsRules.Ctx(ctx).Where(dao.PointsRules.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "积分规则不存在")
	}

	_, err = dao.PointsRules.Ctx(ctx).Where(dao.PointsRules.Columns().Id, req.Id).Update(do.PointsRules{
		Name:        req.Name,
		RuleKey:     req.RuleKey,
		RuleValue:   req.RuleValue,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.PointsRulesUpdateRes{}, nil
}

func (s *sPointsRules) Delete(ctx context.Context, req *v1.PointsRulesDeleteReq) (res *v1.PointsRulesDeleteRes, err error) {
	count, err := dao.PointsRules.Ctx(ctx).Where(dao.PointsRules.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "积分规则不存在")
	}

	_, err = dao.PointsRules.Ctx(ctx).Where(dao.PointsRules.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &v1.PointsRulesDeleteRes{}, nil
}
