package payments

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"

	"gf-eshop/api/payments/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/entity"
	"gf-eshop/internal/service"
)

type sPayments struct{}

func init() {
	service.RegisterPayments(&sPayments{})
}

func generatePaymentNo() string {
	return fmt.Sprintf("PAY%s%04d", gtime.Now().Format("YmdHis"), grand.Intn(10000))
}

func generateRefundNo() string {
	return fmt.Sprintf("REF%s%04d", gtime.Now().Format("YmdHis"), grand.Intn(10000))
}

func (s *sPayments) CreatePayment(ctx context.Context, req *v1.PaymentsCreateReq) (res *v1.PaymentsCreateRes, err error) {
	paymentNo := generatePaymentNo()

	// 查询订单是否存在
	var order *entity.Orders
	err = dao.Orders.Ctx(ctx).Where(dao.Orders.Columns().OrderNo, req.OrderNo).Scan(&order)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "订单不存在")
	}

	var payment *entity.Payments
	err = dao.Payments.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 写入支付记录
		paymentId, err := tx.Model("tx_payments").InsertAndGetId(g.Map{
			"payment_no":      paymentNo,
			"order_no":        req.OrderNo,
			"order_id":        order.Id,
			"merchant_id":     0,
			"order_type":      "order",
			"amount":          req.Amount,
			"currency":        "CNY",
			"payment_method":  req.PaymentMethod,
			"channel":         req.Channel,
			"trade_type":      "native",
			"idempotency_key": paymentNo,
			"status":          "pending",
			"client_ip":       "",
			"channel_response": "",
			"created_at":      gtime.Now(),
		})
		if err != nil {
			return err
		}

		// 写入支付日志
		_, err = tx.Model("tx_payment_logs").Insert(g.Map{
			"payment_id":   paymentId,
			"payment_no":   paymentNo,
			"channel":      req.Channel,
			"action":       "create",
			"status":       "pending",
			"request_body":  "",
			"response_body": "",
			"created_at":   gtime.Now(),
		})
		if err != nil {
			return err
		}

		// 重新读取完整支付记录
		err = tx.Model("tx_payments").Where("id", paymentId).Scan(&payment)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &v1.PaymentsCreateRes{Payments: payment}, nil
}

func (s *sPayments) HandleCallback(ctx context.Context, req *v1.PaymentsCallbackReq) (res *v1.PaymentsCallbackRes, err error) {
	// 查询支付记录
	var payment *entity.Payments
	err = dao.Payments.Ctx(ctx).Where(dao.Payments.Columns().PaymentNo, req.PaymentNo).Scan(&payment)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "支付记录不存在")
	}

	err = dao.Payments.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		updateData := g.Map{
			"status":       req.Status,
			"transaction_id": req.TransactionID,
			"channel":      req.Channel,
			"updated_at":   gtime.Now(),
		}
		if req.FailureReason != "" {
			updateData["failure_reason"] = req.FailureReason
		}
		if req.Status == "success" {
			updateData["paid_at"] = gtime.Now()
		}

		_, err = tx.Model("tx_payments").Where("id", payment.Id).Data(updateData).Update()
		if err != nil {
			return err
		}

		// 写入支付日志
		_, err = tx.Model("tx_payment_logs").Insert(g.Map{
			"payment_id":     payment.Id,
			"payment_no":     payment.PaymentNo,
			"channel":        req.Channel,
			"transaction_id": req.TransactionID,
			"action":         "pay_callback",
			"request_body":   req.RawBody,
			"response_body":  "",
			"status":         req.Status,
			"created_at":     gtime.Now(),
		})
		if err != nil {
			return err
		}

		// 支付成功，更新关联订单的支付状态和状态
		if req.Status == "success" {
			_, err = tx.Model("tx_orders").Where("order_no", payment.OrderNo).Data(g.Map{
				"payment_status": "paid",
				"status":         "paid",
				"paid_at":        gtime.Now(),
				"updated_at":     gtime.Now(),
			}).Update()
			if err != nil {
				return err
			}
			// 同步更新子订单
			_, err = tx.Model("tx_sub_orders").Where("parent_order_no", payment.OrderNo).Data(g.Map{
				"status":     "paid",
				"paid_at":    gtime.Now(),
				"updated_at": gtime.Now(),
			}).Update()
			if err != nil {
				return err
			}
		}

		// 重新读取支付记录
		err = tx.Model("tx_payments").Where("id", payment.Id).Scan(&payment)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &v1.PaymentsCallbackRes{Payments: payment}, nil
}

func (s *sPayments) GetPayment(ctx context.Context, req *v1.PaymentsGetReq) (res *v1.PaymentsGetRes, err error) {
	m := dao.Payments.Ctx(ctx)
	if req.PaymentNo != "" {
		m = m.Where(dao.Payments.Columns().PaymentNo, req.PaymentNo)
	} else if req.OrderNo != "" {
		m = m.Where(dao.Payments.Columns().OrderNo, req.OrderNo)
	} else {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "请提供支付单号或订单号")
	}

	var payment *entity.Payments
	err = m.Scan(&payment)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "支付记录不存在")
	}
	return &v1.PaymentsGetRes{Payments: payment}, nil
}

func (s *sPayments) CreateRefund(ctx context.Context, req *v1.RefundsCreateReq) (res *v1.RefundsCreateRes, err error) {
	// 查询支付记录
	var payment *entity.Payments
	err = dao.Payments.Ctx(ctx).Where(dao.Payments.Columns().PaymentNo, req.PaymentNo).Scan(&payment)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "支付记录不存在")
	}

	refundNo := generateRefundNo()

	var refund *entity.Refunds
	err = dao.Payments.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 写入退款记录
		refundId, err := tx.Model("tx_refunds").InsertAndGetId(g.Map{
			"refund_no":         refundNo,
			"payment_id":        payment.Id,
			"payment_no":        req.PaymentNo,
			"order_no":          payment.OrderNo,
			"order_id":          payment.OrderId,
			"merchant_id":       0,
			"amount":            req.Amount,
			"reason":            req.Reason,
			"status":            "pending",
			"channel_refund_id": "",
			"channel_response":  "",
			"applied_at":        gtime.Now(),
			"created_at":        gtime.Now(),
		})
		if err != nil {
			return err
		}

		// 写入支付日志
		_, err = tx.Model("tx_payment_logs").Insert(g.Map{
			"payment_id":    payment.Id,
			"payment_no":    payment.PaymentNo,
			"channel":       "",
			"transaction_id": "",
			"action":        "refund_create",
			"status":        "pending",
			"request_body":  "",
			"response_body": "",
			"created_at":    gtime.Now(),
		})
		if err != nil {
			return err
		}

		// 重新读取退款记录
		err = tx.Model("tx_refunds").Where("id", refundId).Scan(&refund)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &v1.RefundsCreateRes{Refunds: refund}, nil
}

func (s *sPayments) ListRefund(ctx context.Context, req *v1.RefundsListReq) (res *v1.RefundsListRes, err error) {
	page := req.Page
	size := req.PageSize
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	m := dao.Refunds.Ctx(ctx)
	if req.PaymentNo != "" {
		m = m.Where(dao.Refunds.Columns().PaymentNo, req.PaymentNo)
	}
	if req.OrderNo != "" {
		m = m.Where(dao.Refunds.Columns().OrderNo, req.OrderNo)
	}
	if req.RefundNo != "" {
		m = m.Where(dao.Refunds.Columns().RefundNo, req.RefundNo)
	}
	if req.Status != "" {
		m = m.Where(dao.Refunds.Columns().Status, req.Status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &v1.RefundsListRes{
			List:  make([]*entity.Refunds, 0),
			Total: 0,
		}, nil
	}

	var list []*entity.Refunds
	err = m.Page(page, size).OrderDesc(dao.Refunds.Columns().Id).Scan(&list)
	if err != nil {
		return nil, err
	}
	return &v1.RefundsListRes{
		List:  list,
		Total: int(total),
	}, nil
}

func (s *sPayments) DetailRefund(ctx context.Context, req *v1.RefundsDetailReq) (res *v1.RefundsDetailRes, err error) {
	var refund *entity.Refunds
	err = dao.Refunds.Ctx(ctx).Where(dao.Refunds.Columns().RefundNo, req.RefundNo).Scan(&refund)
	if err != nil {
		return nil, err
	}
	if refund == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "退款记录不存在")
	}
	return &v1.RefundsDetailRes{Refunds: refund}, nil
}