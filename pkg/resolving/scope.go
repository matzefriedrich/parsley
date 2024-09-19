package resolving

import (
	"context"
	"github.com/matzefriedrich/parsley/internal/core"
)

// NewScopedContext creates a new context with an associated service instance map, useful for managing service lifetimes within scope.
func NewScopedContext(ctx context.Context) context.Context {
	instances := make(map[uint64]interface{})
	return context.WithValue(ctx, core.ParsleyContext, instances)
}
