package resolving

import (
	"context"
	"github.com/matzefriedrich/parsley/internal/core"
)

func NewScopedContext(ctx context.Context) context.Context {
	instances := make(map[uint64]interface{})
	return context.WithValue(ctx, core.ParsleyContext, instances)
}
