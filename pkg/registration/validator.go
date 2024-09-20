package registration

import (
	"fmt"
	"github.com/matzefriedrich/parsley/internal"
	"github.com/matzefriedrich/parsley/internal/utils"
	"github.com/matzefriedrich/parsley/pkg/types"
)

const (
	ErrorFailedToRetrieveServiceRegistrations       = "failed to retrieve service registrations"
	ErrorRegistryMissesRequiredServiceRegistrations = "the registry misses required service registrations"
	ErrorCircularServiceRegistrationDetected        = "circular service registration detected"
)

var (

	// ErrFailedToRetrieveServiceRegistrations signifies an error encountered while attempting to retrieve service registrations.
	ErrFailedToRetrieveServiceRegistrations = types.NewRegistryError(ErrorFailedToRetrieveServiceRegistrations)

	// ErrRegistryMissesRequiredServiceRegistrations indicates that required service registrations are missing.
	ErrRegistryMissesRequiredServiceRegistrations = types.NewRegistryError(ErrorRegistryMissesRequiredServiceRegistrations)

	// ErrCircularServiceRegistrationDetected signifies that a circular service registration was encountered.
	ErrCircularServiceRegistrationDetected = types.NewResolverError(ErrorCircularServiceRegistrationDetected)
)

// Validator defines an interface to validate service registries..
type Validator interface {

	// Validate checks the provided ServiceRegistry for missing, invalid, or circular service dependencies. Returns an error if any issues are found.
	Validate(registry types.ServiceRegistry) error
}

type serviceRegistrationsValidator struct {
}

// Validate ensures that all required service registrations are present and service do not depend on them-selves (prevents circular dependencies).
func (s *serviceRegistrationsValidator) Validate(registry types.ServiceRegistry) error {

	registrations, err := registry.GetServiceRegistrations()
	if err != nil {
		return types.NewRegistryError(ErrorFailedToRetrieveServiceRegistrations, types.WithCause(err))
	}

	missingRegistrations := make([]types.ServiceType, 0)

	checkedServiceTypes := make(map[uint64]struct{})

	stack := internal.MakeStack[types.ServiceRegistration]()
	for _, registration := range registrations {
		stack.Push(registration)
	}

	for stack.Any() {
		next := stack.Pop()
		nextId := next.Id()
		_, seen := checkedServiceTypes[nextId]
		if seen {
			continue
		}
		dependencies := next.RequiredServiceTypes()
		for _, dependency := range dependencies {
			list, found := registry.TryGetServiceRegistrations(dependency)
			if found == false {
				missingRegistrations = append(missingRegistrations, dependency)
				continue
			}
			for _, item := range list.Registrations() {
				stack.Push(item)
			}
		}
		checkedServiceTypes[nextId] = struct{}{}
	}

	if len(missingRegistrations) > 0 {
		errors := utils.Map(missingRegistrations, func(serviceType types.ServiceType) error {
			return fmt.Errorf("missing service registration for service type %s", serviceType)
		})
		return types.NewRegistryError(ErrorRegistryMissesRequiredServiceRegistrations, types.WithAggregatedCause(errors...))
	}

	circularDependencyErrors := make([]error, 0)
	for _, registration := range registrations {
		if dependencyError := detectCircularDependency(registration, registry); dependencyError != nil {
			circularDependencyErrors = append(circularDependencyErrors, dependencyError)
		}
	}

	if len(circularDependencyErrors) > 0 {
		return types.NewRegistryError(ErrorCircularServiceRegistrationDetected, types.WithAggregatedCause(circularDependencyErrors...))
	}

	return nil
}

func detectCircularDependency(sr types.ServiceRegistration, registry types.ServiceRegistry) error {

	stack := internal.MakeStack[types.ServiceRegistration]()

	pushRequiredServices := func(r types.ServiceRegistration) {
		requiredServices := r.RequiredServiceTypes()
		for _, serviceType := range requiredServices {
			list, found := registry.TryGetServiceRegistrations(serviceType)
			if found == false {
				continue
			}
			for _, item := range list.Registrations() {
				stack.Push(item)
			}
		}
	}

	pushRequiredServices(sr)

	for stack.Any() {
		next := stack.Pop()
		if next.Id() == sr.Id() {
			serviceType := sr.ServiceType()
			return types.NewRegistryError(ErrorCircularServiceRegistrationDetected, types.ForServiceType(serviceType.Name()))
		}
		pushRequiredServices(next)
	}

	return nil
}

var _ Validator = (*serviceRegistrationsValidator)(nil)

// NewServiceRegistrationsValidator creates a new Validator instance.
func NewServiceRegistrationsValidator() Validator {
	return &serviceRegistrationsValidator{}
}
