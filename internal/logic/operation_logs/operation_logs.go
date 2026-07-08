package operation_logs

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/service"
)

type sOperationLogs struct{}

func init() {
	service.RegisterOperationLogs(&sOperationLogs{})
}

func (s *sOperationLogs) Log(ctx context.Context, in *service.OperationLogInput) error {
	ip := ""
	if r := g.RequestFromCtx(ctx); r != nil {
		ip = r.GetClientIp()
	}
	_, err := dao.OperationLogs.Ctx(ctx).Insert(do.OperationLogs{
		StaffId:       in.StaffId,
		StaffName:     in.StaffName,
		Operation:     in.Operation,
		Resource:      in.Resource,
		ResourceId:    in.ResourceId,
		Detail:        in.Detail,
		Result:        in.Result,
		FailureReason: in.FailureReason,
		Ip:            ip,
		CreatedAt:     gtime.Now(),
	})
	return err
}
