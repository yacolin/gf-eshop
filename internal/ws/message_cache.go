package ws

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

const (
	msgCacheKeyPrefix = "ws:msg:"
	seqIDKeyPrefix    = "ws:seq:"
	defaultCacheSize  = 1000
	cacheExpire       = 24 * time.Hour
)

type MessageCache struct {
	cacheSize int
}

func NewMessageCache() *MessageCache {
	return &MessageCache{cacheSize: defaultCacheSize}
}

func (mc *MessageCache) NextSeqID(userID int64) (int64, error) {
	v, err := g.Redis().Do(context.Background(), "INCR", seqIDKeyPrefix+itoa(userID))
	if err != nil {
		return 0, err
	}
	return v.Int64(), nil
}

func (mc *MessageCache) GetCurrentSeqID(userID int64) (int64, error) {
	v, err := g.Redis().Do(context.Background(), "GET", seqIDKeyPrefix+itoa(userID))
	if err != nil {
		return 0, err
	}
	if v.IsNil() {
		return 0, nil
	}
	return v.Int64(), nil
}

func (mc *MessageCache) StoreMessage(userID int64, seqID int64, message []byte) error {
	key := msgCacheKeyPrefix + itoa(userID)
	_, err := g.Redis().Do(context.Background(), "ZADD", key, seqID, string(message))
	if err != nil {
		return err
	}
	g.Redis().Do(context.Background(), "EXPIRE", key, int(cacheExpire.Seconds()))
	g.Redis().Do(context.Background(), "ZREMRANGEBYRANK", key, 0, -int64(mc.cacheSize+1))
	return nil
}

func (mc *MessageCache) GetMessages(userID int64, lastSeq int64, currentSeq int64) ([][]byte, error) {
	key := msgCacheKeyPrefix + itoa(userID)
	v, err := g.Redis().Do(context.Background(), "ZRANGEBYSCORE", key, "("+itoa(lastSeq), itoa(currentSeq))
	if err != nil {
		return nil, err
	}
	if v.IsNil() {
		return nil, nil
	}
	msgs := v.Strings()
	result := make([][]byte, len(msgs))
	for i, msg := range msgs {
		result[i] = []byte(msg)
	}
	return result, nil
}

func (mc *MessageCache) GetCachedSeqRange(userID int64) (minSeq, maxSeq int64, err error) {
	key := msgCacheKeyPrefix + itoa(userID)
	minV, err := g.Redis().Do(context.Background(), "ZRANGE", key, 0, 0, "WITHSCORES")
	if err != nil {
		return 0, 0, err
	}
	if minV.IsNil() {
		return 0, 0, nil
	}
	vals := minV.Vars()
	if len(vals) < 2 {
		return 0, 0, nil
	}
	minSeq = vals[1].Int64()

	maxV, err := g.Redis().Do(context.Background(), "ZRANGE", key, -1, -1, "WITHSCORES")
	if err != nil {
		return 0, 0, err
	}
	if maxV.IsNil() {
		return 0, 0, nil
	}
	vals = maxV.Vars()
	if len(vals) < 2 {
		return 0, 0, nil
	}
	maxSeq = vals[1].Int64()
	return minSeq, maxSeq, nil
}
