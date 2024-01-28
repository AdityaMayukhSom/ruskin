package load

import (
	"fmt"
	"sync"

	consumer "github.com/AdityaMayukhSom/ruskin/consumer"
)

// Making a map of TopicId and LoadPartition
// The map values have to be mapped to LoadPartition uid(will implement later)
var TopicToPartitionMap map[string]*LoadPartition

type LoadDistributor struct {
	PartitionsCount int
	Partitions      []LoadPartition
	mutex           sync.RWMutex
}

/*Function to get partition and check if the load partition is present or not if
its present it returns the load partition else it returns after creating a new load partition*/

func (ld *LoadDistributor) GetPartition(topic string) *LoadPartition {

	ld.mutex.Lock()
	defer ld.mutex.Unlock()

	//Checking if loadPartition exist in the map
	if partition, exist := TopicToPartitionMap[topic]; exist {
		return partition
	}

	/*If the load partition is not present in the map then create a new load partition
	using the load partition factory and then use that factory using the config and then
	producing the channels*/

	//add the desired config to the partition
	config := LoadPartitionFactoryConfig{}

	//creating channels for load partition
	topicChannel := make(chan string)
	consumerChannel := make(chan consumer.Consumer)

	//creating a new load partition using load partiotion factory
	factory := NewLoadPartitionFactory(&config)
	newPartition := factory.Produce(topicChannel, consumerChannel)

	//Adding the partition into the map and then returning the new partition
	TopicToPartitionMap[topic] = newPartition
	return newPartition
}

func (ld *LoadDistributor) Start() {
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
