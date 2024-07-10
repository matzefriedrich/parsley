package internal

import (
	"context"
	"github.com/matzefriedrich/parsley/pkg/types"
)

type InstanceBag struct {
	parent    *InstanceBag
	instances map[uint64]interface{}
	scope     types.LifetimeScope
}

const (
	parsleyContext = "__parsley"
)

func NewScopedContext(ctx context.Context) context.Context {
	instances := make(map[uint64]interface{})
	return context.WithValue(ctx, parsleyContext, instances)
}

// NewGlobalInstanceBag Creates a new InstanceBag object with global scope.
func NewGlobalInstanceBag() *InstanceBag {
	return &InstanceBag{
		instances: make(map[uint64]interface{}),
		scope:     types.LifetimeSingleton,
	}
}

// NewInstancesBag Creates a new InstanceBag object.
func NewInstancesBag(parent *InstanceBag, scope types.LifetimeScope) *InstanceBag {
	bag := &InstanceBag{
		parent:    parent,
		scope:     scope,
		instances: make(map[uint64]interface{}),
	}
	for k, v := range parent.instances {
		bag.instances[k] = v
	}
	return bag
}

func (b *InstanceBag) TryResolveInstance(ctx context.Context, registration types.ServiceRegistration) (interface{}, bool) {
	id := registration.Id()
	instance, found := b.instances[id]
	if found {
		return instance, true
	}
	scopedInstances, hasParsleyContext := ctx.Value(parsleyContext).(map[uint64]interface{})
	if hasParsleyContext {
		instance, found = scopedInstances[id]
		if found {
			return instance, true
		}
	}
	return nil, false
}

func (b *InstanceBag) KeepInstance(ctx context.Context, registration types.ServiceRegistration, instance interface{}) {
	id := registration.Id()
	switch registration.LifetimeScope() {
	case types.LifetimeSingleton:
		if b.scope == types.LifetimeSingleton {
			b.instances[id] = instance
		} else {
			if b.parent != nil {
				b.parent.KeepInstance(ctx, registration, instance)
			}
		}
	case types.LifetimeScoped:
		scopedInstances, hasParsleyContext := ctx.Value(parsleyContext).(map[uint64]interface{})
		if hasParsleyContext {
			scopedInstances[id] = instance
		}
	case types.LifetimeTransient:
		fallthrough
	default:
		b.instances[id] = instance
	}
}
