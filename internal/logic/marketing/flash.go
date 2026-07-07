package marketing

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gf-eshop/api/marketing/v1"
	"gf-eshop/internal/dao"
	"gf-eshop/internal/errcode"
	"gf-eshop/internal/model/do"
	"gf-eshop/internal/model/entity"
	"gf-eshop/utility"
)

const (
	flashStockRedisPrefix = "flash:stock:"
	flashUserRedisPrefix  = "flash:users:"
	flashMaxRetries       = 3
	flashJobChanSize      = 1024
)

type FlashBuyJob struct {
	UserID      int64
	PromotionID int64
	QueueToken  string
	AcquireTime time.Time
}

const flashBuyLua = `
local stockKey = KEYS[1]
local userKey  = KEYS[2]
local userID   = ARGV[1]

local bought = redis.call("SISMEMBER", userKey, userID)
if bought == 1 then
	return -1
end

local total = redis.call("HGET", stockKey, "total")
local sold  = redis.call("HGET", stockKey, "sold")
if total and tonumber(total) > 0 then
	if not sold then sold = 0 end
	if tonumber(sold) >= tonumber(total) then
		return 0
	end
end

redis.call("HINCRBY", stockKey, "sold", 1)
redis.call("SADD", userKey, userID)
return 1
`

func flashStockKey(promotionID int64) string {
	return flashStockRedisPrefix + itoa(promotionID)
}

func flashUserKey(promotionID int64) string {
	return flashUserRedisPrefix + itoa(promotionID)
}

func LoadFlashStock(ctx context.Context, p *entity.Promotions) error {
	_, err := g.Redis().Do(ctx, "HSET", flashStockKey(p.Id), "total", p.TotalQuantity, "sold", p.UsedQuantity)
	return err
}

func (s *sMarketing) FlashBuy(ctx context.Context, req *v1.FlashBuyReq) (res *v1.FlashBuyRes, err error) {
	userID := utility.GetUserClaims(ctx).UserId

	p, err := s.promotionByID(ctx, req.PromotionId)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, gerror.NewCode(errcode.Code(errcode.CodeNotFound), "promotion not found")
	}
	if p.PromoType != 3 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "not a flash sale")
	}
	if p.Status != 2 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "flash sale not active")
	}

	now := gtime.Now()
	if now.Before(p.StartTime) || now.After(p.EndTime) {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "not in flash sale time window")
	}

	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(tokenBytes)

	stockKey := flashStockKey(req.PromotionId)
	userKey := flashUserKey(req.PromotionId)

	v, err := g.Redis().Do(ctx, "EVAL", flashBuyLua, 2, stockKey, userKey, userID)
	if err != nil {
		return nil, err
	}
	result := v.Int()
	switch result {
	case -1:
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "already participated")
	case 0:
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "sold out")
	}

	job := FlashBuyJob{
		UserID:      userID,
		PromotionID: req.PromotionId,
		QueueToken:  token,
		AcquireTime: time.Now(),
	}

	select {
	case s.jobChan <- job:
	default:
		if err := s.syncFlashBuyToDB(context.Background(), job); err != nil {
			g.Log().Warning(ctx, "sync flash buy to DB failed: %v", err)
		}
	}

	return &v1.FlashBuyRes{Token: token}, nil
}

func (s *sMarketing) FlashConfirm(ctx context.Context, req *v1.FlashConfirmReq) (res *v1.FlashConfirmRes, err error) {
	userID := utility.GetUserClaims(ctx).UserId

	_, err = dao.UserPromotions.Ctx(ctx).
		Data(do.UserPromotions{
			Status:   2,
			UsedTime: gtime.Now(),
		}).
		Where(dao.UserPromotions.Columns().QueueToken, req.Token).
		Where(dao.UserPromotions.Columns().UserId, userID).
		Where(dao.UserPromotions.Columns().Status, 1).
		Update()
	if err != nil {
		return nil, err
	}

	return &v1.FlashConfirmRes{}, nil
}

func (s *sMarketing) asyncWorker() {
	for {
		select {
		case job := <-s.jobChan:
			s.processFlashJob(job)
		case <-s.stopCh:
			return
		}
	}
}

func (s *sMarketing) processFlashJob(job FlashBuyJob) {
	var err error
	for i := 0; i < flashMaxRetries; i++ {
		if i > 0 {
			time.Sleep(time.Duration(50*(1<<(i-1))) * time.Millisecond)
		}
		err = s.syncFlashBuyToDB(context.Background(), job)
		if err == nil {
			return
		}
	}
	g.Log().Warningf(context.Background(), "flash sync failed after %d retries: %v", flashMaxRetries, err)
}

func (s *sMarketing) syncFlashBuyToDB(ctx context.Context, job FlashBuyJob) error {
	upt := dao.UserPromotions.Table()
	pt := dao.Promotions.Table()

	return dao.Promotions.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model(upt).Data(do.UserPromotions{
			UserPromotionNo: generateUserPromotionNo(job.UserID, job.PromotionID),
			UserId:          job.UserID,
			PromotionId:     job.PromotionID,
			AcquireTime:     gtime.New(job.AcquireTime),
			Status:          1,
			QueueToken:      job.QueueToken,
		}).Insert()
		if err != nil {
			return err
		}
		_, err = tx.Model(pt).
			Data(g.Map{"used_quantity": gdb.Raw("used_quantity + 1")}).
			Where("id", job.PromotionID).
			Update()
		return err
	})
}
