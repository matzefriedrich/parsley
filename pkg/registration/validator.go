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
)

var (
	ErrFailedToRetrieveServiceRegistrations       = types.NewRegistryError(ErrorFailedToRetrieveServiceRegistrations)
	ErrRegistryMissesRequiredServiceRegistrations = types.NewRegistryError(ErrorRegistryMissesRequiredServiceRegistrations)
)

type Validator interface {
	Validate(registry types.ServiceRegistry) error
}

type serviceRegistrationsValidator struct {
}

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

	for stack.IsEmpty() == false {
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

	return nil
}

var _ Validator = (*serviceRegistrationsValidator)(nil)

func NewServiceRegistrationsValidator() Validator {
	return &serviceRegistrationsValidator{}
}
