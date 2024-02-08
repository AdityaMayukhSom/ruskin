package producer

import (
	"fmt"

	"github.com/AdityaMayukhSom/ruskin/messagequeue"
)

type ProducerBroker interface {
	Start() error
	TopictoMessageQueueMap map[string] messagequeue.MessageQueue
}

// here for testing purposes only one channel is passed
// we have to actually pass multiple producers channels through it.
func Start(producerChannels []chan string) error {
	
	fmt.Println("Producer Broker is up")

	// Start a goroutine for each producer channel to produce messages to the message queue
	for _, ch := range producerChannels {
		go func(channel chan string) {
			for msg := range channel {
				err := messageQueue.Produce(msg)
				if err != nil {
					fmt.Println("Error producing message:", err)
				}
			}
		}(ch)
	}

	select {}

}