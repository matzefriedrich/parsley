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

func (s *serviceRegistrationList) IsEmpty() bool {
	return len(s.registrations) == 0
}

func (s *serviceRegistrationList) Registrations() []types.ServiceRegistration {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.registrations
}

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

func (s *serviceRegistrationList) Id() uint64 {
	return s.id
}

var _ types.ServiceRegistrationList = &serviceRegistrationList{}

func NewServiceRegistrationList(sequence core.ServiceIdSequence) types.ServiceRegistrationList {
	id := sequence.Next()
	return &serviceRegistrationList{
		id:               id,
		identifierSource: sequence,
		registrations:    make([]types.ServiceRegistration, 0),
		m:                sync.RWMutex{},
	}
}
