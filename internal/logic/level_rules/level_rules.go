package level_rules

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"gf-eshop/api/level_rules/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sLevelRules struct{}

func init() {
	service.RegisterLevelRules(&sLevelRules{})
}

func (s *sLevelRules) List(ctx context.Context, req *v1.LevelRulesListReq) (res *v1.LevelRulesListRes, err error) {
	page := req.Page
	size := req.PageSize
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	m := dao.LevelRules.Ctx(ctx)
	if req.RuleType != "" {
		m = m.Where(dao.LevelRules.Columns().RuleType, req.RuleType)
	}
	if req.Status != nil {
		m = m.Where(dao.LevelRules.Columns().Status, *req.Status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.LevelRulesListRes{
			List:  make([]*entity.LevelRules, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.LevelRules
	err = m.OrderAsc(dao.LevelRules.Columns().SortOrder).Page(page, size).Scan(&list)
	if err != nil {
		return nil, err
	}

	return &v1.LevelRulesListRes{List: list, Total: total}, nil
}

func (s *sLevelRules) Detail(ctx context.Context, req *v1.LevelRulesDetailReq) (res *v1.LevelRulesDetailRes, err error) {
	var rule entity.LevelRules
	err = dao.LevelRules.Ctx(ctx).Where(dao.LevelRules.Columns().Id, req.Id).Scan(&rule)
	if err != nil {
		return nil, err
	}
	return &v1.LevelRulesDetailRes{LevelRules: &rule}, nil
}

func (s *sLevelRules) Create(ctx context.Context, req *v1.LevelRulesCreateReq) (res *v1.LevelRulesCreateRes, err error) {
	result, err := dao.LevelRules.Ctx(ctx).Insert(do.LevelRules{
		Name:           req.Name,
		RuleType:       req.RuleType,
		FromLevelId:    req.FromLevelId,
		ToLevelId:      req.ToLevelId,
		ConditionType:  req.ConditionType,
		ConditionValue: req.ConditionValue,
		Description:    req.Description,
		SortOrder:      req.SortOrder,
		Status:         req.Status,
	})
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &v1.LevelRulesCreateRes{Id: int(id)}, nil
}

func (s *sLevelRules) Update(ctx context.Context, req *v1.LevelRulesUpdateReq) (res *v1.LevelRulesUpdateRes, err error) {
	count, err := dao.LevelRules.Ctx(ctx).Where(dao.LevelRules.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "规则不存在")
	}

	_, err = dao.LevelRules.Ctx(ctx).Where(dao.LevelRules.Columns().Id, req.Id).Update(do.LevelRules{
		Name:           req.Name,
		RuleType:       req.RuleType,
		FromLevelId:    req.FromLevelId,
		ToLevelId:      req.ToLevelId,
		ConditionType:  req.ConditionType,
		ConditionValue: req.ConditionValue,
		Description:    req.Description,
		SortOrder:      req.SortOrder,
		Status:         req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.LevelRulesUpdateRes{}, nil
}

func (s *sLevelRules) Delete(ctx context.Context, req *v1.LevelRulesDeleteReq) (res *v1.LevelRulesDeleteRes, err error) {
	count, err := dao.LevelRules.Ctx(ctx).Where(dao.LevelRules.Columns().Id, req.Id).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "规则不存在")
	}

	_, err = dao.LevelRules.Ctx(ctx).Where(dao.LevelRules.Columns().Id, req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &v1.LevelRulesDeleteRes{}, nil
}
