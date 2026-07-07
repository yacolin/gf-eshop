package user_points

import (
	"context"
	"database/sql"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/api/user_points/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
	"gf-eshop/utility"
)

type sUserPoints struct{}

func init() {
	service.RegisterUserPoints(&sUserPoints{})
}

// List 积分流水列表（支持按用户、时间、来源筛选）
func (s *sUserPoints) List(ctx context.Context, req *v1.PointsListReq) (res *v1.PointsListRes, err error) {
	page := req.Page
	size := req.PageSize
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	m := dao.Points.Ctx(ctx)
	if req.UserId != "" {
		m = m.Where(dao.Points.Columns().UserId, req.UserId)
	}
	if req.Source != "" {
		m = m.Where(dao.Points.Columns().Source, req.Source)
	}
	if req.StartTime != "" {
		m = m.WhereGTE(dao.Points.Columns().CreatedAt, req.StartTime)
	}
	if req.EndTime != "" {
		m = m.WhereLTE(dao.Points.Columns().CreatedAt, req.EndTime)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.PointsListRes{
			List:  make([]*entity.Points, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.Points
	err = m.OrderDesc(dao.Points.Columns().CreatedAt).Page(page, size).Scan(&list)
	if err != nil {
		return nil, err
	}

	return &v1.PointsListRes{List: list, Total: total}, nil
}

// Balance 查询用户当前积分余额
func (s *sUserPoints) Balance(ctx context.Context, req *v1.PointsBalanceReq) (res *v1.PointsBalanceRes, err error) {
	// 查询最近一条已确认的记录获取余额
	var latest entity.Points
	err = dao.Points.Ctx(ctx).
		Where(dao.Points.Columns().UserId, req.UserId).
		Where(dao.Points.Columns().Status, 1).
		OrderDesc(dao.Points.Columns().Id).
		Limit(1).
		Scan(&latest)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	balance := latest.BalanceAfter

	// 统计累计获得积分（points > 0 且 status=1）
	totalIn, err := dao.Points.Ctx(ctx).
		Where(dao.Points.Columns().UserId, req.UserId).
		Where(dao.Points.Columns().Status, 1).
		WhereGT(dao.Points.Columns().Points, 0).
		Sum(dao.Points.Columns().Points)
	if err != nil {
		return nil, err
	}

	// 统计累计消耗积分（points < 0 且 status=1，取绝对值）
	totalOut, err := dao.Points.Ctx(ctx).
		Where(dao.Points.Columns().UserId, req.UserId).
		Where(dao.Points.Columns().Status, 1).
		WhereLT(dao.Points.Columns().Points, 0).
		Sum(dao.Points.Columns().Points)
	if err != nil {
		return nil, err
	}

	return &v1.PointsBalanceRes{
		UserId:   int64(latest.UserId),
		Balance:  balance,
		TotalIn:  int64(totalIn),
		TotalOut: int64(-totalOut), // 转为正值
	}, nil
}

// Trend 查询用户积分趋势（按日统计）
func (s *sUserPoints) Trend(ctx context.Context, req *v1.PointsTrendReq) (res *v1.PointsTrendRes, err error) {
	m := dao.Points.Ctx(ctx).
		Where(dao.Points.Columns().UserId, req.UserId).
		Where(dao.Points.Columns().Status, 1)
	if req.StartTime != "" {
		m = m.WhereGTE(dao.Points.Columns().CreatedAt, req.StartTime)
	}
	if req.EndTime != "" {
		m = m.WhereLTE(dao.Points.Columns().CreatedAt, req.EndTime)
	}

	// 按日期分组统计积分变动
	type TrendRow struct {
		Date    string
		Points  int64
		Balance int64
	}
	var rows []*TrendRow
	err = m.Fields(
		"DATE(created_at) AS date",
		"SUM(points) AS points",
		"MAX(balance_after) AS balance",
	).Group("DATE(created_at)").OrderAsc("DATE(created_at)").Scan(&rows)
	if err != nil {
		return nil, err
	}

	list := make([]*v1.PointsTrendItem, 0)
	for _, row := range rows {
		list = append(list, &v1.PointsTrendItem{
			Date:    row.Date,
			Points:  row.Points,
			Balance: row.Balance,
		})
	}

	return &v1.PointsTrendRes{List: list}, nil
}

// Adjust 手动调整用户积分（需要管理员权限）
func (s *sUserPoints) Adjust(ctx context.Context, req *v1.PointsAdjustReq) (res *v1.PointsAdjustRes, err error) {
	// 校验管理员权限
	claims := utility.GetStaffClaims(ctx)
	if claims == nil {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	isAdmin, err := service.Roles().IsAdmin(ctx, claims.StaffId)
	if err != nil || !isAdmin {
		return nil, gerror.NewCode(gcode.CodeOperationFailed, "无权限，需要管理员角色")
	}

	if req.Points == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "积分变动值不能为0")
	}

	// 查询当前积分余额
	var latest entity.Points
	err = dao.Points.Ctx(ctx).
		Where(dao.Points.Columns().UserId, req.UserId).
		Where(dao.Points.Columns().Status, 1).
		OrderDesc(dao.Points.Columns().Id).
		Limit(1).
		Scan(&latest)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	balanceAfter := latest.BalanceAfter + req.Points
	if balanceAfter < 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "积分扣减后余额不能为负数")
	}

	// 插入积分变动记录
	result, err := dao.Points.Ctx(ctx).Insert(do.Points{
		UserId:       req.UserId,
		Points:       req.Points,
		BalanceAfter: balanceAfter,
		Source:       "admin",
		Status:       1, // 管理员调整直接确认
		Remark:       req.Remark,
	})
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeOperationFailed, "获取流水ID失败")
	}

	// 记录操作日志
	g.Log().Info(ctx, "管理员调整积分: user_id=%d, points=%d, balance_after=%d, remark=%s", req.UserId, req.Points, balanceAfter, req.Remark)

	return &v1.PointsAdjustRes{
		Id:           id,
		BalanceAfter: balanceAfter,
	}, nil
}
