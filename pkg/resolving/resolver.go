package resolving

import (
	"context"
	"github.com/matzefriedrich/parsley/internal"
	"github.com/matzefriedrich/parsley/internal/core"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
	"reflect"
)

type resolver struct {
	registry        types.ServiceRegistry
	globalInstances *core.InstanceBag
}

// ResolveRequiredServices resolves all registered services of a specified type T using the given resolver and context.
func ResolveRequiredServices[T any](resolver types.Resolver, ctx context.Context) ([]T, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	switch t.Kind() {
	case reflect.Func:
	case reflect.Interface:
	case reflect.Pointer:
	case reflect.Slice:
	case reflect.Struct:
	default:
		return []T{}, types.NewResolverError(types.ErrorActivatorFunctionInvalidReturnType)
	}
	resolvedInstances, err := resolver.Resolve(ctx, types.MakeServiceType[T]())
	if err != nil {
		return []T{}, err
	}

	result := make([]T, 0, len(resolvedInstances))
	for _, instance := range resolvedInstances {
		result = append(result, instance.(T))
	}
	return result, err
}

// ResolveRequiredService resolves a single service instance of the specified type using the given resolver and context.
// The method can return the following errors: ErrorCannotResolveService, ErrorAmbiguousServiceInstancesResolved.
func ResolveRequiredService[T any](resolver types.Resolver, ctx context.Context) (T, error) {
	var nilInstance T
	services, err := ResolveRequiredServices[T](resolver, ctx)
	if err != nil {
		return nilInstance, err
	}
	if len(services) == 1 {
		return services[0], nil
	} else if len(services) > 1 {
		return nilInstance, types.NewResolverError(types.ErrorAmbiguousServiceInstancesResolved)
	}
	return nilInstance, types.NewResolverError(types.ErrorCannotResolveService)
}

// NewResolver creates and returns a new Resolver instance based on the provided ServiceRegistry.
func NewResolver(registry types.ServiceRegistry) types.Resolver {
	r := &resolver{
		registry:        registry,
		globalInstances: core.NewGlobalInstanceBag(),
	}
	_ = registration.RegisterInstance[types.Resolver](registry, r)
	return r
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
			parent := next.Consumer()
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
		return registration.NewMultiRegistryAccessor(transientRegistry, r.registry), nil
	}
	return r.registry, nil
}

// Resolve returns a list of instances associated with the specified service type.
func (r *resolver) Resolve(ctx context.Context, serviceType types.ServiceType) ([]interface{}, error) {
	return r.ResolveWithOptions(ctx, serviceType)
}

// ResolveWithOptions resolves instances for the given service type with the provided resolver options.
func (r *resolver) ResolveWithOptions(scope context.Context, serviceType types.ServiceType, resolverOptions ...types.ResolverOptionsFunc) ([]interface{}, error) {

	ctx, span := newResolveWithOptionsSpan(scope, serviceType)
	defer span.End()

	registry, registryErr := r.createResolverRegistryAccessor(resolverOptions...)
	if registryErr != nil {
		return nil, types.NewResolverError("failed to create resolver service registry", types.WithCause(registryErr))
	}

	serviceRegistrationList, found := registry.TryGetServiceRegistrations(serviceType)
	if !found {
		return nil, types.NewResolverError(types.ErrorServiceTypeNotRegistered, types.ForServiceType(serviceType.Name()))
	}

	resolvedInstances := make([]interface{}, 0)

	for _, serviceRegistration := range serviceRegistrationList.Registrations() {

		makeDependencyInfo := func(sr types.ServiceRegistration, consumer types.DependencyInfo) (types.DependencyInfo, error) {
			instance, _ := r.globalInstances.TryResolveInstance(ctx, serviceRegistration)
			err := detectCircularDependency(sr, consumer)
			if err != nil {
				return nil, err
			}
			return registration.NewDependencyInfo(sr, instance, consumer), nil
		}

		resolverStack := internal.MakeStack[types.DependencyInfo]()

		root, _ := makeDependencyInfo(serviceRegistration, nil)
		stack := internal.MakeStack[types.DependencyInfo](root)
		for stack.Any() {
			next := stack.Pop()
			resolverStack.Push(next)
			requiredServices := next.RequiredServiceTypes()
			for _, requiredService := range requiredServices {
				requiredServiceRegistration, isRegistered := registry.TryGetSingleServiceRegistration(requiredService)
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

		instances := core.NewInstancesBag(r.globalInstances, types.LifetimeTransient)

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

		resolvedInstances = append(resolvedInstances, root.Instance())
	}

	return resolvedInstances, nil
}

var _ types.Resolver = &resolver{}
