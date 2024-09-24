package core

import "sync"

type serviceId struct {
	n uint64
	m sync.Mutex
}

// ServiceIdSequence is an interface for generating sequential service identifiers.
type ServiceIdSequence interface {
	Next() uint64
}

var _ ServiceIdSequence = &serviceId{}

// NewServiceId creates a new instance of ServiceIdSequence starting from the specified initial value. This is useful for generating unique service identifiers in a thread-safe manner.
func NewServiceId(n uint64) ServiceIdSequence {
	return &serviceId{n: n, m: sync.Mutex{}}
}

func (s *serviceId) Next() uint64 {
	s.m.Lock()
	defer s.m.Unlock()
	s.n++
	return s.n
}
