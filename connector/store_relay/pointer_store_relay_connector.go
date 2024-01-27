package store_relay_connector

import (
	"github.com/AdityaMayukhSom/ruskin/messagequeue"
)

type PointerStoreConnector struct {
	storeAddr *messagequeue.Store
}

func NewPointerStoreConnector(topicStoreAddr *messagequeue.Store) *PointerStoreConnector {
	return &PointerStoreConnector{
		storeAddr: topicStoreAddr,
	}
}

func (pc *PointerStoreConnector) Fetch(offset int) ([]byte, error) {
	data, err := (*pc.storeAddr).Extract(offset)
	return data, err
}

// Fetches all the data starting from the beginning.
func (pc *PointerStoreConnector) FetchAll() ([][]byte, error) {
	return nil, nil
}

func (pc *PointerStoreConnector) FetchLatest() ([]byte, error) {
	data, err := (*pc.storeAddr).ExtractLatest()
	return data, err
}
