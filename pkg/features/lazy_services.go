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

type Lazy[T any] interface {
	Value() T
}

var _ Lazy[any] = &lazy[any]{}

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
