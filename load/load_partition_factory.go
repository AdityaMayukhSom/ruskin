package load

import (
	"github.com/AdityaMayukhSom/ruskin/consumer"
	messagequeue "github.com/AdityaMayukhSom/ruskin/messagequeue"
)

type LoadPartitionFactoryConfig struct {
}

type LoadPartitionFactory struct {
	*LoadPartitionFactoryConfig
}

func NewLoadPartitionFactory(config *LoadPartitionFactoryConfig) *LoadPartitionFactory {
	return &LoadPartitionFactory{
		LoadPartitionFactoryConfig: config,
	}
}

func (lpf *LoadPartitionFactory) Produce(
	topicChannel <-chan string,
	consumerChannel <-chan consumer.Consumer,
) *LoadPartition {
	// mu.Lock()
	// defer mu.Unlock()

	// if f == nil {
	// 	panic("Factory is nil")
	// }
	// if _, exist := LoadPartitionFactories[pkgName]; exist {
	// 	panic("Factory already registered")
	// }

	// LoadPartitionFactories[pkgName] = f

	return &LoadPartition{
		topicChannel: make(<-chan messagequeue.TopicIdentifier),
		connnectionChannel: make(<-chan struct {
			consumer consumer.Consumer
			topic    messagequeue.TopicIdentifier
		}),
		relays: make(map[messagequeue.TopicIdentifier]consumer.ConsumerRelay),
	}
}
