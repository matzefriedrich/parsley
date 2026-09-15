package core

import (
	"context"
	"maps"
	"sync"

	"github.com/matzefriedrich/parsley/pkg/types"
)

type InstanceBag struct {
	parent        *InstanceBag
	instances     map[uint64]any
	disposables   []*disposableRef
	scope         types.LifetimeScope
	teardownOrder []string
	mu            sync.Mutex
}

type disposableRef struct {
	id       uint64
	group    string
	disposed bool
}

type ContextKey string

const (
	ParsleyContext ContextKey = "__parsley"
)

// NewGlobalInstanceBag Creates a new InstanceBag object with global scope.
func NewGlobalInstanceBag(teardownOrder []string) *InstanceBag {
	return &InstanceBag{
		instances:     make(map[uint64]any),
		disposables:   make([]*disposableRef, 0),
		scope:         types.LifetimeSingleton,
		teardownOrder: copyTeardownOrder(teardownOrder),
	}
}

// NewInstancesBag Creates a new InstanceBag object.
func NewInstancesBag(parent *InstanceBag, scope types.LifetimeScope, teardownOrder []string) *InstanceBag {
	bag := &InstanceBag{
		parent:        parent,
		scope:         scope,
		instances:     make(map[uint64]any),
		disposables:   make([]*disposableRef, 0),
		teardownOrder: copyTeardownOrder(teardownOrder),
	}
	if parent != nil {
		parent.mu.Lock()
		maps.Copy(bag.instances, parent.instances)
		parent.mu.Unlock()
	}
	return bag
}

func copyTeardownOrder(order []string) []string {
	copied := make([]string, len(order))
	copy(copied, order)
	return copied
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
				b.disposables = append(b.disposables, &disposableRef{id: id, group: registration.LifecycleGroup()})
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
				scoped.disposables = append(scoped.disposables, &disposableRef{id: id, group: registration.LifecycleGroup()})
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

// Dispose disposes the tracked disposable instances. When a teardown order is declared, instances are disposed group
// by group in that order, with instances within a group disposed in reverse resolution order; instances whose group is
// not listed are treated as part of the DefaultLifecycleGroup and are disposed at the "default" position (its declared
// position, or last when "default" is not listed). Without a declared order, all disposables are disposed in reverse
// resolution order.
func (b *InstanceBag) Dispose(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	var errs []error
	for _, ref := range disposalOrder(b.teardownOrder, b.disposables) {
		if ref.disposed {
			continue
		}
		if disposable, ok := b.instances[ref.id].(types.Disposable); ok {
			err := disposable.Dispose(ctx)
			if err != nil {
				errs = append(errs, err)
			}
			ref.disposed = true
		}
	}
	if len(errs) > 0 {
		return types.NewResolverError("failed to dispose one or more services", types.WithAggregatedCause(errs...))
	}
	return nil
}

// disposalOrder computes the disposal order over the given disposables. Without a declared teardown order the
// disposables are returned in reverse resolution order (the default behaviour). With a declared order, each listed
// group is visited once in declaration order (deduplicated, empty names skipped) and its members are appended in
// reverse resolution order; a "default" pass collects explicit DefaultLifecycleGroup members together with members of
// any group not listed in the order, and is placed at the "default" position (its declared position, or last when
// "default" is not listed).
func disposalOrder(teardownOrder []string, disposables []*disposableRef) []*disposableRef {
	if len(teardownOrder) == 0 {
		order := make([]*disposableRef, 0, len(disposables))
		for i := len(disposables) - 1; i >= 0; i-- {
			order = append(order, disposables[i])
		}
		return order
	}
	listed := make(map[string]struct{}, len(teardownOrder))
	sequence := make([]string, 0, len(teardownOrder)+1)
	for _, group := range teardownOrder {
		if group == "" {
			continue
		}
		if _, exists := listed[group]; exists {
			continue
		}
		listed[group] = struct{}{}
		sequence = append(sequence, group)
	}
	if needsDefaultPass(disposables, sequence, listed) {
		sequence = append(sequence, types.DefaultLifecycleGroup)
	}
	order := make([]*disposableRef, 0, len(disposables))
	for _, group := range sequence {
		for i := len(disposables) - 1; i >= 0; i-- {
			ref := disposables[i]
			if group == types.DefaultLifecycleGroup {
				if isDefaultMatch(ref.group, listed) {
					order = append(order, ref)
				}
			} else if ref.group == group {
				order = append(order, ref)
			}
		}
	}
	return order
}

// isDefaultMatch reports whether the given group is handled by the "default" pass: the explicit
// DefaultLifecycleGroup or any group not listed in the declared teardown order.
func isDefaultMatch(group string, listed map[string]struct{}) bool {
	if group == "" {
		group = types.DefaultLifecycleGroup
	}
	if group == types.DefaultLifecycleGroup {
		return true
	}
	_, ok := listed[group]
	return !ok
}

// needsDefaultPass reports whether a "default" pass has to be appended to the disposal sequence: a ref whose group is
// handled by the "default" pass exists while "default" is not yet part of the sequence.
func needsDefaultPass(disposables []*disposableRef, sequence []string, listed map[string]struct{}) bool {
	for _, group := range sequence {
		if group == types.DefaultLifecycleGroup {
			return false
		}
	}
	for _, ref := range disposables {
		if isDefaultMatch(ref.group, listed) {
			return true
		}
	}
	return false
}
