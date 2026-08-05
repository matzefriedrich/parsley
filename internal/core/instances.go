package core

import (
	"context"
	"maps"
	"sync"

	"github.com/matzefriedrich/parsley/pkg/types"
)

type InstanceBag struct {
	parent      *InstanceBag
	instances   map[uint64]any
	disposables []*disposableRef
	scope       types.LifetimeScope
	mu          sync.Mutex
}

type disposableRef struct {
	id       uint64
	disposed bool
}

type ContextKey string

const (
	ParsleyContext ContextKey = "__parsley"
)

// NewGlobalInstanceBag Creates a new InstanceBag object with global scope.
func NewGlobalInstanceBag() *InstanceBag {
	return &InstanceBag{
		instances:   make(map[uint64]any),
		disposables: make([]*disposableRef, 0),
		scope:       types.LifetimeSingleton,
	}
}

// NewInstancesBag Creates a new InstanceBag object.
func NewInstancesBag(parent *InstanceBag, scope types.LifetimeScope) *InstanceBag {
	bag := &InstanceBag{
		parent:      parent,
		scope:       scope,
		instances:   make(map[uint64]any),
		disposables: make([]*disposableRef, 0),
	}
	if parent != nil {
		parent.mu.Lock()
		maps.Copy(bag.instances, parent.instances)
		parent.mu.Unlock()
	}
	return bag
}

// TryResolveInstance attempts to locate an instance of a service identified by the given registration.
func (b *InstanceBag) TryResolveInstance(ctx context.Context, registration types.ServiceRegistration) (any, bool) {
	id := registration.Id()
	b.mu.Lock()
	instance, found := b.instances[id]
	b.mu.Unlock()
	if found {
		return instance, true
	}
	scoped, hasParsleyContext := ctx.Value(ParsleyContext).(*InstanceBag)
	if hasParsleyContext {
		scoped.mu.Lock()
		instance, found = scoped.instances[id]
		scoped.mu.Unlock()
		if found {
			return instance, true
		}
	}
	return nil, false
}

// KeepInstance stores an instance of a service based on the service's lifetime scope. Singleton instances are stored
// at the appropriate singleton level in the hierarchy. Scoped instances are stored in the context-specified scope.
// Transient instances are stored in the current instance bag.
func (b *InstanceBag) KeepInstance(ctx context.Context, registration types.ServiceRegistration, instance any) {
	id := registration.Id()
	switch registration.LifetimeScope() {
	case types.LifetimeSingleton:
		if b.scope == types.LifetimeSingleton {
			b.mu.Lock()
			b.instances[id] = instance
			if _, ok := instance.(types.Disposable); ok {
				b.disposables = append(b.disposables, &disposableRef{id: id})
			}
			b.mu.Unlock()
		} else {
			if b.parent != nil {
				b.parent.KeepInstance(ctx, registration, instance)
			}
		}
	case types.LifetimeScoped:
		scoped, hasParsleyContext := ctx.Value(ParsleyContext).(*InstanceBag)
		if hasParsleyContext {
			scoped.mu.Lock()
			scoped.instances[id] = instance
			if _, ok := instance.(types.Disposable); ok {
				scoped.disposables = append(scoped.disposables, &disposableRef{id: id})
			}
			scoped.mu.Unlock()
		}
	case types.LifetimeTransient:
		fallthrough
	default:
		b.mu.Lock()
		b.instances[id] = instance
		b.mu.Unlock()
	}
}

func (b *InstanceBag) Dispose(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	var errs []error
	for i := len(b.disposables) - 1; i >= 0; i-- {
		ref := b.disposables[i]
		if disposable, ok := b.instances[ref.id].(types.Disposable); ok {
			err := disposable.Dispose(ctx)
			if err != nil {
				errs = append(errs, err)
			}
			ref.disposed = true
		}
	}
	if len(errs) > 0 {
		return types.NewAggregateError("failed to dispose one or more services", errs)
	}
	return nil
}
