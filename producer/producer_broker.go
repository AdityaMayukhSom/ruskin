package producer

import (
	"log/slog"

	load "github.com/AdityaMayukhSom/ruskin/load"
	mq "github.com/AdityaMayukhSom/ruskin/messagequeue"
	transport "github.com/AdityaMayukhSom/ruskin/transport"
)

type ProducerBrokerIdentifier *ProducerBroker

type ProducerBroker struct {
	// later we can use multiple channels and distribute the load here too
	messageChannel  chan transport.Message
	loadDistributor *load.LoadDistributor
	topicQueueMap   map[mq.TopicIdentifier]*mq.MessageQueue
}

func NewProducerBroker(loadDistributor *load.LoadDistributor) *ProducerBroker {
	return &ProducerBroker{
		messageChannel:  make(chan transport.Message),
		loadDistributor: loadDistributor,
		topicQueueMap:   make(map[mq.TopicIdentifier]*mq.MessageQueue),
	}
}

func (pb *ProducerBroker) addMessageToQueue(message transport.Message) {
	ti := mq.TopicIdentifier(message.Topic)

	queue, found := pb.topicQueueMap[ti]

	if !found {

	}

	store, err := queue.GetStore(message.Topic)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	store.Insert(message.Data)
}

// here for testing purposes only one channel is passed
// we have to actually pass multiple producers channels through it.
func (pb *ProducerBroker) Start(listenAddr string, producerChannels ...chan<- string) error {

	slog.Info("Producer Broker is up")

	// for initial version a producer handler is spawned
	// initially when a Producer Broker is created by the server
	pb.SpawnProducerHandler(listenAddr, pb.messageChannel)

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

func (pb *ProducerBroker) SpawnProducerHandler(listenAddr string, producerChannels ...chan<- transport.Message) *HTTPProducerHandler {
	producerHandler := NewHTTPProducerHandler(listenAddr, pb.messageChannel, pb)
	go func() {
		if err := producerHandler.Start(); err != nil {
			slog.Error("Failed to start HTTP producer handler", "error", err)
		}
	}()
	return producerHandler
}
