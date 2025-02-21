package mock

import (
	"context"

	"github.com/InfluxCommunity/flux"
	"github.com/InfluxCommunity/flux/memory"
)

// Program is a mock program that can be returned by the mock compiler.
// It will construct a mock query that will then be passed to ExecuteFn.
type Program struct {
	StartFn   func(ctx context.Context, alloc memory.Allocator) (*Query, error)
	ExecuteFn func(ctx context.Context, q *Query, alloc memory.Allocator)
}

func (p *Program) Start(ctx context.Context, alloc memory.Allocator) (flux.Query, error) {
	startFn := p.StartFn
	if startFn == nil {
		var cancel func()
		ctx, cancel = context.WithCancel(ctx)
		startFn = func(ctx context.Context, alloc memory.Allocator) (*Query, error) {
			results := make(chan flux.Result)
			q := &Query{
				ResultsCh: results,
				CancelFn:  cancel,
				Canceled:  make(chan struct{}),
			}
			go func() {
				defer close(results)
				if p.ExecuteFn != nil {
					p.ExecuteFn(ctx, q, alloc)
				}
			}()
			return q, nil
		}
	}
	return startFn(ctx, alloc)
}
