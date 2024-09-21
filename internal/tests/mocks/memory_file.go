package mocks

import (
	"fmt"
	"io"
	"strings"
)

type memoryFileTarget struct {
	buffer strings.Builder
}

type MemoryFile interface {
	io.WriteCloser
	fmt.Stringer
}

func (m *memoryFileTarget) String() string {
	return m.buffer.String()
}

func (m *memoryFileTarget) Write(p []byte) (n int, err error) {
	return m.buffer.Write(p)
}

func (m *memoryFileTarget) Close() error {
	return nil
}

var _ MemoryFile = (*memoryFileTarget)(nil)

// NewMemoryFile creates a new MemoryFile instance.
// The returned MemoryFile acts as a WriteCloser and a Stringer, storing data in a string buffer.
func NewMemoryFile() MemoryFile {
	return &memoryFileTarget{
		buffer: strings.Builder{},
	}
}
