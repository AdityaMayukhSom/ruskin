package store

import (
	"fmt"
	"testing"
)

func TestMemoryStore(t *testing.T) {
	s := NewMemoryStore("mytopic")
	cnt := 10
	for idx := 0; idx < cnt; idx++ {
		key := fmt.Sprintf("foobarbaz_%d", idx)
		offset, err := s.Insert([]byte(key))
		if err != nil {
			t.Error(err)
		}

		data, err := s.Extract(offset)
		if err != nil {
			t.Error(err)
		}

		retKey := string(data)
		if key != retKey {
			t.Errorf("expected %s, got %s", key, retKey)
		}
	}
}
