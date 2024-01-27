package load

import (
	"fmt"

	consumer "github.com/AdityaMayukhSom/ruskin/consumer"
)

// Making a map of TopicId and LoadPartition
// The map values have to be mapped to LoadPartition uid(will implement later)
var TopicToPartitionMap = make(map[string]*LoadPartition)

type LoadDistributor struct {
	PartitionsCount int
	Partitions      []LoadPartition
}

func (ld *LoadDistributor) getOrCreatePartition(topic string) *LoadPartition {

	//Checking if loadPartition exist in the map
	if partition, exist := TopicToPartitionMap[topic]; exist {
		return partition
	}

	/*If the load partition is not present in the map then create a new load partition
	using the load partition factory and then use that factory using the config and then
	producing the channels*/

	config := LoadPartitionFactoryConfig{}

	topicChannel := make(chan string)
	consumerChannel := make(chan consumer.Consumer)

	factory := NewLoadPartitionFactory(&config)
	newPartition := factory.Produce(topicChannel, consumerChannel)

	//Adding the partition into the map and then returning
	TopicToPartitionMap[topic] = newPartition
	return newPartition
}

func (ld *LoadDistributor) Listen() {
	agg := make(chan string)
	for _, ch := range chans {
		go func(c chan string) {
			for msg := range c {
				agg <- msg
			}
		}(ch)
	}

	select {
	case msg <- agg:
		fmt.Println("received ", msg)
	}
}
