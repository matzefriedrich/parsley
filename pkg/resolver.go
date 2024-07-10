package pkg

import (
	"context"
	"github.com/matzefriedrich/parsley/internal"
	"github.com/matzefriedrich/parsley/pkg/types"
	"reflect"
)

type resolver struct {
	registry  types.ServiceRegistryAccessor
	instances map[uint64]interface{}
}

func NewResolver(registry types.ServiceRegistryAccessor) types.Resolver {
	return &resolver{
		registry:  registry,
		instances: make(map[uint64]interface{}),
	}
}

func (r *resolver) Resolve(ctx context.Context, serviceType reflect.Type) (interface{}, error) {

	registration, found := r.registry.TryGetServiceRegistration(serviceType)
	if !found {
		return nil, types.NewResolverError(types.ErrorServiceTypeNotRegistered, types.ForServiceType(serviceType.Name()))
	}

	makeDependencyInfo := func(sr types.ServiceRegistration, consumer types.DependencyInfo) types.DependencyInfo {
		instance, _ := r.TryResolveInstance(ctx, registration)
		return types.NewDependencyInfo(sr, instance, consumer)
	}

	resolverStack := internal.MakeStack[types.DependencyInfo]()

	root := makeDependencyInfo(registration, nil)
	stack := internal.MakeStack[types.DependencyInfo](root)
	for stack.Any() {
		next := stack.Pop()
		resolverStack.Push(next)
		requiredServices := next.RequiredServiceTypes()
		for _, requiredService := range requiredServices {
			requiredServiceRegistration, isRegistered := r.registry.TryGetServiceRegistration(requiredService)
			if isRegistered == false {
				return nil, types.NewResolverError(types.ErrorServiceTypeNotRegistered, types.ForServiceType(requiredService.Name()))
			}
			child := makeDependencyInfo(requiredServiceRegistration, next)
			if child.HasInstance() { // skip traversal, if instance is already present
				continue
			}
			next.AddRequiredServiceInfo(child)
			stack.Push(child)
		}
	}

	for resolverStack.Any() {
		next := resolverStack.Pop()
		_, err := next.CreateInstance()
		if err != nil {
			return nil, types.NewResolverError(types.ErrorCannotResolveService, types.WithCause(err), types.ForServiceType(next.ServiceTypeName()))
		}
		// TODO: check instance lifetime; if configred per scope, or singleton, we need to keep the instance (currently everything is transient)
	}

	// TODO:
	// * detect circular dependency; abort if any of the parameters matches the requested type
	// * if a registered service is a function, it is a factory; call it to resolve the service
	// * evaluate the lifetime of service registration; if singleton register it depending on its scope
	//		* means: per-context, transient (no registration), singleton (global)

	return root.Instance(), nil
}

func (r *resolver) TryResolveInstance(ctx context.Context, registration types.ServiceRegistration) (interface{}, bool) {
	id := registration.Id()
	instance, found := r.instances[id]
	if found {
		f := func() interface{} {
			return instance
		}
		return f, true
	}
	// TODO: try to find the instance on the context
	return nil, false
}

var _ types.Resolver = &resolver{}
