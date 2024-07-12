package pkg

import (
	"context"
	"github.com/matzefriedrich/parsley/internal"
	"github.com/matzefriedrich/parsley/pkg/types"
	"reflect"
)

type resolver struct {
	registry        types.ServiceRegistry
	globalInstances *internal.InstanceBag
}

func ResolveRequiredService[T any](resolver types.Resolver, ctx context.Context) (T, error) {
	var nilInstance T
	t := reflect.TypeOf((*T)(nil)).Elem()
	switch t.Kind() {
	case reflect.Func:
	case reflect.Interface:
	default:
		return nilInstance, types.NewResolverError(types.ErrorActivatorFunctionInvalidReturnType)
	}
	resolve, err := resolver.Resolve(ctx, types.ServiceType[T]())
	if err != nil {
		return nilInstance, err
	}
	return resolve.(T), err
}

func NewResolver(registry types.ServiceRegistry) types.Resolver {
	return &resolver{
		registry:        registry,
		globalInstances: internal.NewGlobalInstanceBag(),
	}
}

func detectCircularDependency(sr types.ServiceRegistration, consumer types.DependencyInfo) error {
	stack := internal.MakeStack[types.DependencyInfo]()
	if consumer != nil {
		stack.Push(consumer)
		for stack.Any() {
			next := stack.Pop()
			if next.Registration().Id() == sr.Id() {
				return types.NewResolverError(types.ErrorCircularDependencyDetected, types.ForServiceType(next.ServiceTypeName()))
			}
			parent := consumer.Consumer()
			if parent != nil {
				stack.Push(parent)
			}
		}
	}
	return nil
}

func (r *resolver) createResolverRegistryAccessor(resolverOptions ...types.ResolverOptionsFunc) (types.ServiceRegistryAccessor, error) {
	if len(resolverOptions) > 0 {
		transientRegistry := r.registry.CreateLinkedRegistry()
		err := applyResolverOptions(transientRegistry, resolverOptions...)
		if err != nil {
			return nil, err
		}
		return NewMultiRegistryAccessor(r.registry, transientRegistry), nil
	}
	return r.registry, nil
}

func (r *resolver) Resolve(ctx context.Context, serviceType reflect.Type) (interface{}, error) {
	return r.ResolveWithOptions(ctx, serviceType)
}

func (r *resolver) ResolveWithOptions(ctx context.Context, serviceType reflect.Type, resolverOptions ...types.ResolverOptionsFunc) (interface{}, error) {

	registry, registryErr := r.createResolverRegistryAccessor(resolverOptions...)
	if registryErr != nil {
		return nil, types.NewResolverError("failed to create resolver service registry", types.WithCause(registryErr))
	}

	registration, found := registry.TryGetServiceRegistration(serviceType)
	if !found {
		return nil, types.NewResolverError(types.ErrorServiceTypeNotRegistered, types.ForServiceType(serviceType.Name()))
	}

	makeDependencyInfo := func(sr types.ServiceRegistration, consumer types.DependencyInfo) (types.DependencyInfo, error) {
		instance, _ := r.globalInstances.TryResolveInstance(ctx, registration)
		err := detectCircularDependency(sr, consumer)
		if err != nil {
			return nil, err
		}
		return types.NewDependencyInfo(sr, instance, consumer), nil
	}

	resolverStack := internal.MakeStack[types.DependencyInfo]()

	root, _ := makeDependencyInfo(registration, nil)
	stack := internal.MakeStack[types.DependencyInfo](root)
	for stack.Any() {
		next := stack.Pop()
		resolverStack.Push(next)
		requiredServices := next.RequiredServiceTypes()
		for _, requiredService := range requiredServices {
			requiredServiceRegistration, isRegistered := registry.TryGetServiceRegistration(requiredService)
			if isRegistered == false {
				return nil, types.NewResolverError(types.ErrorServiceTypeNotRegistered, types.ForServiceType(requiredService.Name()))
			}
			child, err := makeDependencyInfo(requiredServiceRegistration, next)
			if err != nil {
				return nil, types.NewResolverError(types.ErrorCannotBuildDependencyGraph, types.WithCause(err), types.ForServiceType(requiredService.Name()))
			}
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
