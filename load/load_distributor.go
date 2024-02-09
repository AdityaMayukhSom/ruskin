package load

import (
	"fmt"
	"sync"

	consumer "github.com/AdityaMayukhSom/ruskin/consumer"
)

// Making a map of TopicId and LoadPartition
// The map values have to be mapped to LoadPartition uuid(will implement later)
// var TopicToPartitionMap map[string]*LoadPartition

type LoadDistributor struct {
	PartitionsCount int

	mutex sync.RWMutex

	TopicToPartitionMap map[string]*LoadPartition
}

func NewLoadDistributor() *LoadDistributor {
	return &LoadDistributor{}
}

/*
Function to get partition and check if the load partition is present
or not if its present it returns the load partition else it returns
after creating a new load partition.
*/
func (ld *LoadDistributor) GetPartition(topic string) *LoadPartition {
	ld.mutex.RLock()

	//Checking if loadPartition exist in the map
	partition, exist := ld.TopicToPartitionMap[topic]

	ld.mutex.RUnlock()

	if !exist {
		/*If the load partition is not present in the map then create a new load partition
		using the load partition factory and then use that factory using the config and then
		producing the channels*/

		//Adding the partition into the map and then returning the new partition

		ld.mutex.Lock()
		existingPartition, exist := ld.TopicToPartitionMap[topic]
		if !exist {
			//add the desired config to the partition
			config := LoadPartitionFactoryConfig{}

			//creating channels for load partition
			topicChannel := make(chan string)
			consumerChannel := make(chan consumer.Consumer)

			//creating a new load partition using load partiotion factory
			factory := NewLoadPartitionFactory(&config)
			partition = factory.Produce(topicChannel, consumerChannel)

			ld.TopicToPartitionMap[topic] = partition
		} else {
			partition = existingPartition
		}

		ld.mutex.Unlock()
	}

	return partition
}

func (ld *LoadDistributor) Start(chans []chan string) {
	agg := make(chan string)
	for _, ch := range chans {
		go func(c chan string) {
			for msg := range c {
				agg <- msg
			}
		}(ch)
	}

	select {
	case msg := <-agg:
		fmt.Println("received ", msg)
	}
}
