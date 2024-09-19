package registration

import (
	"github.com/matzefriedrich/parsley/internal/core"
	"github.com/matzefriedrich/parsley/pkg/types"
	"sync"
)

type serviceRegistrationList struct {
	id               uint64
	identifierSource core.ServiceIdSequence
	registrations    []types.ServiceRegistration
	m                sync.RWMutex
}

// IsEmpty Returns true if the list contains no registrations, otherwise false.
func (s *serviceRegistrationList) IsEmpty() bool {
	return len(s.registrations) == 0
}

// Registrations returns a slice of ServiceRegistration representing the current registrations in the list.
func (s *serviceRegistrationList) Registrations() []types.ServiceRegistration {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.registrations
}

// AddRegistration adds a new service registration to the list.
// It returns an ErrorTypeAlreadyRegistered error if the registration already exists.
func (s *serviceRegistrationList) AddRegistration(registration types.ServiceRegistrationSetup) error {

	s.m.Lock()
	defer s.m.Unlock()

	for _, reg := range s.registrations {
		if reg.IsSame(registration) {
			return types.NewRegistryError(types.ErrorTypeAlreadyRegistered)
		}
	}

	registrationId := s.identifierSource.Next()
	err := registration.SetId(registrationId)
	if err != nil {
		return types.NewRegistryError(types.ErrorServiceAlreadyLinkedWithAnotherList, types.WithCause(err))
	}

	s.registrations = append(s.registrations, registration)
	return nil
}

// Id returns the unique identifier of the service registration list.
func (s *serviceRegistrationList) Id() uint64 {
	return s.id
}

var _ types.ServiceRegistrationList = &serviceRegistrationList{}

// NewServiceRegistrationList creates a new service registration list instance.
func NewServiceRegistrationList(sequence core.ServiceIdSequence) types.ServiceRegistrationList {
	id := sequence.Next()
	return &serviceRegistrationList{
		id:               id,
		identifierSource: sequence,
		registrations:    make([]types.ServiceRegistration, 0),
		m:                sync.RWMutex{},
	}
}
