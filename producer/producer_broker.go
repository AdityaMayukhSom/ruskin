package producer

import (
	"log/slog"

	load "github.com/AdityaMayukhSom/ruskin/load"
	mq "github.com/AdityaMayukhSom/ruskin/messagequeue"
	transport "github.com/AdityaMayukhSom/ruskin/transport"
)

type ProducerBroker struct {
	// later we can use multiple channels and distribute the load here too
	messageChannel  chan transport.Message
	loadDistributor *load.LoadDistributor
	topicQueueMap   map[mq.TopicIdentifier]mq.MessageQueue
}

func NewProducerBroker(loadDistributor *load.LoadDistributor) *ProducerBroker {
	return &ProducerBroker{
		messageChannel:  make(chan transport.Message),
		loadDistributor: loadDistributor,
		topicQueueMap:   make(map[mq.TopicIdentifier]mq.MessageQueue),
	}
}

func (pb *ProducerBroker) addMessageToQueue(message transport.Message) {

}

// here for testing purposes only one channel is passed
// we have to actually pass multiple producers channels through it.
func (pb *ProducerBroker) Start(producerChannels ...chan<- string) error {

	slog.Info("Producer Broker is up")

	go func(pb *ProducerBroker) {
		for msg := range pb.messageChannel {
			go pb.addMessageToQueue(msg)
		}
	}(pb)

	return nil

}

func (pb *ProducerBroker) handleIncomingProducer(producer <-chan transport.Message) {
	for msg := range producer {
		pb.messageChannel <- msg
	}

}

func (pb *ProducerBroker) AddProducer(producerChannels ...chan transport.Message) error {
	for _, producerChannel := range producerChannels {
		go pb.handleIncomingProducer(producerChannel)
	}
	return nil
}
