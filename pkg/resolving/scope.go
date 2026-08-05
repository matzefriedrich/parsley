package resolving

import (
	"context"
	"github.com/matzefriedrich/parsley/internal/core"
	"github.com/matzefriedrich/parsley/pkg/types"
)

// NewScopedContext creates a new context with an associated service instance map, useful for managing service lifetimes within scope.
func NewScopedContext(ctx context.Context) context.Context {
	bag := core.NewInstancesBag(nil, types.LifetimeScoped)
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
