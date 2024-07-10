package internal

import "sync"

type serviceId struct {
	n uint64
	m sync.Mutex
}

type ServiceIdSequence interface {
	Next() uint64
}

var _ ServiceIdSequence = &serviceId{}

func NewServiceId(n uint64) ServiceIdSequence {
	return &serviceId{n: n, m: sync.Mutex{}}
}

func (s *serviceId) Next() uint64 {
	s.m.Lock()
	defer s.m.Unlock()
	s.n++
	return s.n
}
