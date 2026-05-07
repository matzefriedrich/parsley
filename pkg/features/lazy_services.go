package features

import (
	"context"
	"sync"

	"github.com/matzefriedrich/parsley/internal"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
)

type lazy[T any] struct {
	instance      T
	activatorFunc any
	resolver      types.Resolver
	m             sync.RWMutex
}

// Value returns the instance of the lazy-initialized type, generating it if not already created. Ensures thread-safety.
func (l *lazy[T]) Value(ctx context.Context) T {
	l.m.RLock()
	if !internal.IsNil(l.instance) {
		defer l.m.RUnlock()
	} else {
		l.m.RUnlock()
		l.m.Lock()
		defer l.m.Unlock()
		if internal.IsNil(l.instance) {
			instance, err := resolving.Activate[T](ctx, l.resolver, l.activatorFunc)
			if err != nil {
				return l.instance
			}
			l.instance = instance
		}
	}
	return l.instance
}

// Lazy represents a type whose value is initialized lazily upon first access, typically to improve performance or manage resources.
type Lazy[T any] interface {
	Value(ctx context.Context) T
}

var _ Lazy[any] = &lazy[any]{}

// RegisterLazy registers a lazily activated service in the service registry using the provided activator function.
func RegisterLazy[T any](registry types.ServiceRegistry, activatorFunc any, scope types.LifetimeScope) error {
	return registry.Register(func(resolver types.Resolver) Lazy[T] {
		return &lazy[T]{
			activatorFunc: activatorFunc,
			resolver:      resolver,
		}
	}, scope)
}
