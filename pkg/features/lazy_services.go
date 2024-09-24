package features

import (
	"github.com/matzefriedrich/parsley/internal"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
	"sync"
)

type lazy[T any] struct {
	instance      T
	activatorFunc func() T
	m             sync.RWMutex
}

// Value returns the instance of the lazy-initialized type, generating it if not already created. Ensures thread-safety.
func (l *lazy[T]) Value() T {
	l.m.RLock()
	if internal.IsNil(l.instance) == false {
		defer l.m.RUnlock()
	} else {
		l.m.RUnlock()
		l.m.Lock()
		defer l.m.Unlock()
		if internal.IsNil(l.instance) {
			instance := l.activatorFunc()
			l.instance = instance
		}
	}
	return l.instance
}

// Lazy represents a type whose value is initialized lazily upon first access, typically to improve performance or manage resources.
type Lazy[T any] interface {
	Value() T
}

var _ Lazy[any] = &lazy[any]{}

// RegisterLazy registers a lazily-activated service in the service registry using the provided activator function.
func RegisterLazy[T any](registry types.ServiceRegistry, activatorFunc func() T, _ types.LifetimeScope) error {

	lazyActivator := newLazyServiceFactory[T](activatorFunc)
	err := registration.RegisterInstance(registry, lazyActivator)
	if err != nil {
		return types.NewRegistryError("failed to register lazy service")
	}

	return nil
}

func newLazyServiceFactory[T any](activatorFunc func() T) Lazy[T] {
	return &lazy[T]{
		activatorFunc: activatorFunc,
	}
}
