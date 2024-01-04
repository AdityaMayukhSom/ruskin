package store

type StoreConfig struct {
}

type StoreFactory interface {
	Produce() Store
}

type MemoryStoreFactory struct {
	*StoreConfig
}

func NewMemoryStoreFactory(storeConfig *StoreConfig) *MemoryStoreFactory {
	return &MemoryStoreFactory{
		StoreConfig: storeConfig,
	}
}

func (m *MemoryStoreFactory) Produce() Store {
	return NewMemoryStore()
}
