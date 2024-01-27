package messagequeue

type MemoryStoreFactory struct {
	*StoreFactoryConfig
}

func NewMemoryStoreFactory(config *StoreFactoryConfig) *MemoryStoreFactory {
	return &MemoryStoreFactory{
		StoreFactoryConfig: config,
	}
}

func (m *MemoryStoreFactory) Produce(topicName string) Store {
	return NewMemoryStore(topicName)
}
