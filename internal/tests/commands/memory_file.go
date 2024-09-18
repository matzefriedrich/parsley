package commands

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

func newMemoryFile() MemoryFile {
	return &memoryFileTarget{
		buffer: strings.Builder{},
	}
}
