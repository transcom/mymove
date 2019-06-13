package migrate

import (
	"io"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

var (
	ErrWait = errors.New("wait for input")
)

// Buffer is a wrapper around strings.Builder for concurrent loading and reading.
type Buffer struct {
	*sync.RWMutex
	buffer strings.Builder
	closed bool
}

func (b *Buffer) Close() {
	b.Lock()
	b.closed = true
	b.Unlock()
}

func (b *Buffer) Closed() bool {
	closed := false
	b.RLock()
	closed = b.closed
	b.RUnlock()
	return closed
}

func (b *Buffer) Len() int {
	return b.buffer.Len()
}

func (b *Buffer) WriteString(x string) (int, error) {
	b.Lock()
	n, err := b.buffer.WriteString(x)
	b.Unlock()
	return n, err
}

func (b *Buffer) WriteByte(x byte) error {
	b.Lock()
	err := b.buffer.WriteByte(x)
	b.Unlock()
	return err
}

func (b *Buffer) WriteRune(x rune) (int, error) {
	b.Lock()
	n, err := b.buffer.WriteRune(x)
	b.Unlock()
	return n, err
}

func (b *Buffer) Index(i int) (byte, error) {
	b.RLock()
	if i >= b.buffer.Len() {
		if b.closed {
			b.RUnlock()
			return byte(0), io.EOF
		}
		b.RUnlock()
		return byte(0), ErrWait
	}
	x := b.buffer.String()[i]
	b.RUnlock()
	return x, nil
}

func (b *Buffer) Range(start int, end int) (string, error) {
	if start >= end {
		return "", errors.New("start should be less than end")
	}
	b.RLock()
	if end >= b.buffer.Len() {
		if b.closed {
			b.RUnlock()
			return "", io.EOF
		}
		b.RUnlock()
		return "", ErrWait
	}
	str := b.buffer.String()[start:end]
	b.RUnlock()
	return str, nil
}

// NewBuffer returns a new Buffer.
func NewBuffer() *Buffer {
	return &Buffer{
		RWMutex: &sync.RWMutex{},
		buffer:  strings.Builder{},
		closed:  false,
	}
}
