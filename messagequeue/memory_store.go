package messagequeue

import (
	"fmt"
	"sync"
)

type MemoryStore struct {
	mut       sync.RWMutex
	topicName string
	data      [][]byte
}

func NewMemoryStore(topicName string) *MemoryStore {
	return &MemoryStore{
		data:      make([][]byte, 0),
		topicName: topicName,
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

func (s *MemoryStore) ExtractLatest() ([]byte, error) {
	s.mut.RLock()
	defer s.mut.RUnlock()
	return s.data[len(s.data)-1], nil
}
