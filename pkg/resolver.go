package pkg

import (
	"context"
	"github.com/matzefriedrich/parsley/internal"
	"github.com/matzefriedrich/parsley/pkg/types"
	"reflect"
)

type resolver struct {
	registry        types.ServiceRegistryAccessor
	globalInstances *internal.InstanceBag
}

func NewResolver(registry types.ServiceRegistryAccessor) types.Resolver {
	return &resolver{
		registry:        registry,
		globalInstances: internal.NewGlobalInstanceBag(),
	}
}

func (r *resolver) Resolve(ctx context.Context, serviceType reflect.Type) (interface{}, error) {

	registration, found := r.registry.TryGetServiceRegistration(serviceType)
	if !found {
		return nil, types.NewResolverError(types.ErrorServiceTypeNotRegistered, types.ForServiceType(serviceType.Name()))
	}

	makeDependencyInfo := func(sr types.ServiceRegistration, consumer types.DependencyInfo) types.DependencyInfo {
		instance, _ := r.globalInstances.TryResolveInstance(ctx, registration)
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

	instances := internal.NewInstancesBag(r.globalInstances, types.LifetimeTransient)

	for resolverStack.Any() {

		next := resolverStack.Pop()
		nextRegistration := next.Registration()
		instance, ok := instances.TryResolveInstance(ctx, nextRegistration)
		if ok {
			_ = next.SetInstance(instance)
			continue
		}

		instance, err := next.CreateInstance()
		if err != nil {
			return nil, types.NewResolverError(types.ErrorCannotResolveService, types.WithCause(err), types.ForServiceType(next.ServiceTypeName()))
		}

		instances.KeepInstance(ctx, nextRegistration, instance)
	}

	return root.Instance(), nil
}

var _ types.Resolver = &resolver{}
