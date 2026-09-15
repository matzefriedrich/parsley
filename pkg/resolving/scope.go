package resolving

import (
	"context"
	"github.com/matzefriedrich/parsley/internal/core"
	"github.com/matzefriedrich/parsley/pkg/types"
)

type scopedContextOptions struct {
	teardownOrder []string
}

// ScopeOptionsFunc configures a scoped context created by NewScopedContextWithOptions.
type ScopeOptionsFunc func(*scopedContextOptions)

// WithTeardownOrder declares the lifecycle group teardown order for a scoped context. Empty group names are ignored.
func WithTeardownOrder(groups ...string) ScopeOptionsFunc {
	return func(o *scopedContextOptions) {
		order := make([]string, 0, len(groups))
		for _, group := range groups {
			if group == "" {
				continue
			}
			order = append(order, group)
		}
		o.teardownOrder = order
	}
}

// NewScopedContext creates a new context with an associated service instance map, useful for managing service lifetimes within scope.
func NewScopedContext(ctx context.Context) context.Context {
	return NewScopedContextWithOptions(ctx)
}

// NewScopedContextWithOptions creates a new context with an associated service instance map and the given scope options.
// The scoped services disposed via DisposeScope honour the declared lifecycle group teardown order.
func NewScopedContextWithOptions(ctx context.Context, options ...ScopeOptionsFunc) context.Context {
	opts := &scopedContextOptions{}
	for _, apply := range options {
		apply(opts)
	}
	bag := core.NewInstancesBag(nil, types.LifetimeScoped, opts.teardownOrder)
	return context.WithValue(ctx, core.ParsleyContext, bag)
}

// DisposeScope disposes all services in the given context that implement the Disposable interface.
func DisposeScope(ctx context.Context) error {
	bag, ok := ctx.Value(core.ParsleyContext).(*core.InstanceBag)
	if !ok {
		return nil
	}
	return bag.Dispose(ctx)
}
