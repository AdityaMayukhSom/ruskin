package messagequeue

type StoreConfig struct {
}

type StoreFactory interface {
	Produce(string) Store
}

type MemoryStoreFactory struct {
	*StoreConfig
}

func NewMemoryStoreFactory(storeConfig *StoreConfig) *MemoryStoreFactory {
	return &MemoryStoreFactory{
		StoreConfig: storeConfig,
	}
}

func (m *MemoryStoreFactory) Produce(topicName string) Store {
	return NewMemoryStore(topicName)
}
