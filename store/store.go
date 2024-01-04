package store

import (
	"fmt"
	"sync"
)

type Store interface {
	Insert([]byte) (int, error)
	Extract(int) ([]byte, error)
}

type MemoryStore struct {
	mut  sync.RWMutex
	data [][]byte
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make([][]byte, 0),
	}
}

func (s *MemoryStore) Insert(b []byte) (int, error) {
	s.mut.Lock()
	defer s.mut.Unlock()

	s.data = append(s.data, b)
	return len(s.data) - 1, nil
}

func (s *MemoryStore) Extract(offset int) ([]byte, error) {
	s.mut.RLock()
	defer s.mut.RUnlock()

	if offset < 0 {
		return nil, fmt.Errorf(
			"IndexOutOfBound :: offset must be non-negative, provided (%d)",
			offset,
		)
	}

	if offset >= len(s.data) {
		return nil, fmt.Errorf(
			"IndexOutOfBound :: provided (%d), length (%d)",
			offset, len(s.data),
		)
	}

	return s.data[offset], nil
}
