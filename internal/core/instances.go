package core

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
	ParsleyContext = "__parsley"
)

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

// TryResolveInstance attempts to locate an instance of a service identified by the given registration.
func (b *InstanceBag) TryResolveInstance(ctx context.Context, registration types.ServiceRegistration) (interface{}, bool) {
	id := registration.Id()
	instance, found := b.instances[id]
	if found {
		return instance, true
	}
	scopedInstances, hasParsleyContext := ctx.Value(ParsleyContext).(map[uint64]interface{})
	if hasParsleyContext {
		instance, found = scopedInstances[id]
		if found {
			return instance, true
		}
	}
	return nil, false
}

// KeepInstance stores an instance of a service based on the service's lifetime scope. Singleton instances are stored
// at the appropriate singleton level in the hierarchy. Scoped instances are stored in the context-specified scope.
// Transient instances are stored in the current instance bag.
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
		scopedInstances, hasParsleyContext := ctx.Value(ParsleyContext).(map[uint64]interface{})
		if hasParsleyContext {
			scopedInstances[id] = instance
		}
	case types.LifetimeTransient:
		fallthrough
	default:
		b.instances[id] = instance
	}
}
