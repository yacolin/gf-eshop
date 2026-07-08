package departments

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"gf-eshop/internal/dao"
	"gf-eshop/internal/model/entity"
)

const (
	departmentIdsKey     = "department:ids"
	departmentIdsTTL     = 10 * time.Minute
	departmentEntityTTL  = 10 * time.Minute
	departmentScoreScale = 1_000_000_000_000
)

const departmentListLua = `
local zkey  = KEYS[1]
local pfx   = ARGV[1]   -- "department:"
local start = tonumber(ARGV[2])
local stop  = tonumber(ARGV[3])

local total = redis.call('ZCARD', zkey)
if total == 0 then
    return {0}
end

local ids = redis.call('ZRANGE', zkey, start, stop)
if #ids == 0 then
    return {total}
end

local keys = {}
for i, id in ipairs(ids) do
    keys[i] = pfx .. id
end
if redis.call('EXISTS', unpack(keys)) < #ids then
    return {-1}
end

local vals = redis.call('MGET', unpack(keys))
local result = {total}
for i = 1, #vals do
    result[#result + 1] = vals[i]
end
return result
`

func cacheKeyDepartment(id int64) string { return fmt.Sprintf("department:%d", id) }

func encodeDepartmentScore(sortOrder int, id int64) float64 {
	return float64(int64(sortOrder)*departmentScoreScale + (departmentScoreScale - id))
}

// --- Lua 分页读取 ---

func getDepartmentPage(ctx context.Context, page, size int) (list []*entity.Departments, total int, err error) {
	start := (page - 1) * size
	stop := start + size - 1

	v, err := g.Redis().Do(ctx, "EVAL", departmentListLua, 1, departmentIdsKey, "department:", start, stop)
	if err != nil || v.IsNil() || len(v.Vars()) == 0 {
		return nil, 0, err
	}
	elems := v.Vars()
	total = elems[0].Int()
	if total <= 0 {
		return nil, 0, nil
	}
	list = make([]*entity.Departments, 0, len(elems)-1)
	for _, e := range elems[1:] {
		var d entity.Departments
		if err := sonic.Unmarshal(e.Bytes(), &d); err != nil {
			return nil, 0, err
		}
		list = append(list, &d)
	}
	return list, total, nil
}

// --- Singleflight 重建 ---

var (
	rebuildMu  sync.Mutex
	rebuilding = map[string]chan struct{}{
		"departments": nil,
	}
)

func ensureDepartmentCache(ctx context.Context) {
	card, _ := g.Redis().Do(ctx, "ZCARD", departmentIdsKey)
	if !card.IsNil() && card.Int() > 0 {
		return
	}

	rebuildMu.Lock()
	ch, ok := rebuilding["departments"]
	if !ok || ch == nil {
		ch = make(chan struct{})
		rebuilding["departments"] = ch
		rebuildMu.Unlock()

		_ = rebuildDepartmentCache(context.Background())

		rebuildMu.Lock()
		delete(rebuilding, "departments")
		rebuildMu.Unlock()
		close(ch)
		return
	}
	rebuildMu.Unlock()
	<-ch
}

// --- 索引 ---

func addDepartmentToIndex(ctx context.Context, id int64, sortOrder int) {
	g.Redis().Do(ctx, "ZADD", departmentIdsKey, encodeDepartmentScore(sortOrder, id), id)
}

func removeDepartmentFromIndex(ctx context.Context, id int64) {
	g.Redis().Do(ctx, "ZREM", departmentIdsKey, id)
}

// --- 缓存重建 ---

func rebuildDepartmentCache(ctx context.Context) error {
	var list []*entity.Departments
	if err := dao.Departments.Ctx(ctx).OrderAsc(dao.Departments.Columns().SortOrder).Scan(&list); err != nil {
		return err
	}
	if len(list) == 0 {
		return nil
	}
	g.Redis().Do(ctx, "DEL", departmentIdsKey)
	for _, d := range list {
		g.Redis().Do(ctx, "ZADD", departmentIdsKey, encodeDepartmentScore(d.SortOrder, d.Id), d.Id)
		data, _ := sonic.Marshal(d)
		g.Redis().Do(ctx, "SETEX", cacheKeyDepartment(d.Id), int(departmentEntityTTL.Seconds()), string(data))
	}
	g.Redis().Do(ctx, "EXPIRE", departmentIdsKey, int(departmentIdsTTL.Seconds()))
	return nil
}

// Warmup 启动预热
func Warmup(ctx context.Context) {
	if err := rebuildDepartmentCache(ctx); err != nil {
		g.Log().Warningf(ctx, "department cache warmup failed: %v", err)
		return
	}
	g.Log().Infof(ctx, "department cache warmed up")
}

// --- 单条缓存 ---

func getDepartmentEntityCache(ctx context.Context, id int64) (*entity.Departments, error) {
	v, err := g.Redis().Do(ctx, "GET", cacheKeyDepartment(id))
	if err != nil || v.IsNil() {
		return nil, err
	}
	var d entity.Departments
	if err := sonic.Unmarshal(v.Bytes(), &d); err != nil {
		return nil, err
	}
	return &d, nil
}

func setDepartmentEntityCache(ctx context.Context, d *entity.Departments) error {
	data, err := sonic.Marshal(d)
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "SETEX", cacheKeyDepartment(d.Id), int(departmentEntityTTL.Seconds()), string(data))
	return err
}

func delDepartmentEntityCache(ctx context.Context, id int64) {
	g.Redis().Do(ctx, "DEL", cacheKeyDepartment(id))
}
