package products

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/sync/errgroup"
)

// WarmupStage 预热阶段接口。后续新增模块实现此接口即可加入预热管线。
type WarmupStage interface {
	Name() string
	Warmup(ctx context.Context) (int, error)
}

// WarmupPipeline 串联多个预热阶段，各阶段并行执行。
type WarmupPipeline struct {
	stages []WarmupStage
}

func NewWarmupPipeline(stages ...WarmupStage) *WarmupPipeline {
	return &WarmupPipeline{stages: stages}
}

func (p *WarmupPipeline) Add(stage WarmupStage) {
	p.stages = append(p.stages, stage)
}

// Run 并行执行所有已注册的预热阶段
func (p *WarmupPipeline) Run(ctx context.Context) map[string]int {
	results := make(map[string]int, len(p.stages))
	eg, egCtx := errgroup.WithContext(ctx)
	for _, stage := range p.stages {
		stage := stage
		eg.Go(func() error {
			n, err := stage.Warmup(egCtx)
			if err != nil {
				g.Log().Warningf(ctx, "warmup stage %q failed: %v", stage.Name(), err)
				return nil // 不阻断其他阶段
			}
			results[stage.Name()] = n
			g.Log().Infof(ctx, "warmup stage %q done: %d items", stage.Name(), n)
			return nil
		})
	}
	_ = eg.Wait()
	g.Log().Infof(ctx, "warmup pipeline finished: %v", results)
	return results
}

// FuncStage 将函数适配为 WarmupStage
type FuncStage struct {
	name string
	fn   func(context.Context) (int, error)
}

func NewFuncStage(name string, fn func(context.Context) (int, error)) *FuncStage {
	return &FuncStage{name: name, fn: fn}
}

func (s *FuncStage) Name() string              { return s.name }
func (s *FuncStage) Warmup(ctx context.Context) (int, error) { return s.fn(ctx) }
