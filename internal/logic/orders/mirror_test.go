package orders

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gtime"

	"gf-eshop/internal/model/entity"
)

// 本文件锁定 Phase 3 双写/影子读里**与数据库无关**的那部分契约：
// 分片纯计算、影子读抽样、指纹比对。

func TestShardPureComputation(t *testing.T) {
	// 单号 → 分片（只认内嵌时间，不做模式判断）
	sh, ok := shardOfOrderNo("ORD202609060029157196")
	if !ok || sh.suffix != "202609" {
		t.Errorf("shardOfOrderNo = (%q, %v), want (202609, true)", sh.suffix, ok)
	}
	if _, ok = shardOfOrderNo("ORD_NOT_EXIST"); ok {
		t.Error("无法解析的单号应返回 ok=false")
	}

	// 创建时间 → 分片（跨年边界：12 月不能算成次年 1 月）
	at := gtime.New(time.Date(2025, 12, 31, 23, 59, 59, 0, time.Local))
	if got := shardOfCreatedAt(at).suffix; got != "202512" {
		t.Errorf("跨年边界分片 = %q, want 202512", got)
	}

	// 编码主键 → 分片；老自增主键必须判为不可用（要查映射表）
	newID := encodeOrderID(time.Date(2027, 3, 15, 10, 30, 0, 0, time.Local), 42)
	if sh, ok = shardOfEncodedID(newID); !ok || sh.suffix != "202703" {
		t.Errorf("shardOfEncodedID(新主键) = (%q, %v), want (202703, true)", sh.suffix, ok)
	}
	if _, ok = shardOfEncodedID(2000); ok {
		t.Error("老自增主键不应被当作编码主键")
	}
}

func TestShouldShadowDeterministic(t *testing.T) {
	orderNo := "ORD20261004130509000001"
	if !shouldShadow(orderNo, 100) {
		t.Error("100% 抽样必须命中")
	}
	if shouldShadow(orderNo, 0) {
		t.Error("0% 抽样不应命中")
	}

	// 确定性 + 比例大致正确：影子读不能用随机数，否则同一条订单时而比对时而不比对
	const n = 2000
	hits := 0
	for i := 0; i < n; i++ {
		no := fmt.Sprintf("ORD20261004130509%06d", i)
		first := shouldShadow(no, 50)
		if first != shouldShadow(no, 50) {
			t.Fatalf("同一样本抽样结果不稳定: %s", no)
		}
		if first {
			hits++
		}
	}
	if hits < n/4 || hits > n*3/4 {
		t.Errorf("50%% 抽样命中 %d/%d，偏离过大（抽样函数可能退化成常量）", hits, n)
	}
}

func TestOrderFingerprintDetectsDifferences(t *testing.T) {
	base := &entity.Orders{
		OrderNo:       "ORD20261004130509000001",
		UserId:        7,
		TotalAmount:   1000,
		PayAmount:     900,
		Status:        "pending",
		PaymentStatus: "unpaid",
		CreatedAt:     gtime.New(time.Date(2026, 10, 4, 13, 5, 9, 0, time.Local)),
	}
	same := *base
	if orderFingerprint(base) != orderFingerprint(&same) {
		t.Error("同一订单应得到相同指纹")
	}

	cases := map[string]*entity.Orders{}
	c := *base
	c.Status = "paid"
	cases["状态"] = &c
	c2 := *base
	c2.PayAmount = 901
	cases["金额"] = &c2
	c3 := *base
	c3.PaymentStatus = "paid"
	cases["支付状态"] = &c3
	c4 := *base
	c4.CreatedAt = gtime.New(time.Date(2026, 10, 4, 13, 5, 10, 0, time.Local))
	cases["创建时间"] = &c4

	for name, o := range cases {
		if orderFingerprint(base) == orderFingerprint(o) {
			t.Errorf("%s 不同时必须产生不同指纹（否则影子读发现不了差异）", name)
		}
	}
	if orderFingerprint(nil) != "<nil>" {
		t.Error("nil 订单应有明确的指纹表示")
	}
}
