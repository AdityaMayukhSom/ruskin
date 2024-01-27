package connector

import (
	"github.com/AdityaMayukhSom/ruskin/messagequeue"
)

type PointerStoreConnector struct {
	storeAddr *messagequeue.Store
}

func NewPointerStoreConnector() {

}

func (pc *PointerStoreConnector) Fetch(offset int) ([]byte, error) {
	data, err := (*pc.storeAddr).Extract(offset)
	return data, err
}

// Fetches all the data starting from the beginning.
func (pc *PointerStoreConnector) FetchAll(offset int) ([][]byte, error) {
	return nil, nil
}

func (pc *PointerStoreConnector) FetchLatest() ([]byte, error) {
	data, err := (*pc.storeAddr).ExtractLatest()
	return data, err
}
