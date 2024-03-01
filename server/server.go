package server

import (
	"log/slog"

	"github.com/AdityaMayukhSom/ruskin/consumer"
	"github.com/AdityaMayukhSom/ruskin/load"
	"github.com/AdityaMayukhSom/ruskin/producer"
	"github.com/gorilla/websocket"
)

type Server struct {
	// We must use two mutex to synchronise the creation of new topic store.
	ProducerAddrs []string
	ConsumerAddrs []string

	consumerProxy   consumer.ConsumerProxy
	loadDistributor load.LoadDistributor
	producerBroker  *producer.ProducerBroker

	// When we are willing to shutdown the server gracefully,
	// we need to signal this channel or close it.
	quitChannel chan struct{}
}

type ServerOption func(*Server)

// will update the port at 3000 only if atleast one port it given,
// else will keep on running at 3000 if no addr is passed.
func WithProducerAddr(producerAddrs ...string) ServerOption {
	return func(server *Server) {
		if len(producerAddrs) > 0 {
			server.ProducerAddrs = producerAddrs
		}
	}
}

// will update the port at 4000 only if atleast one port it given,
// else will keep on running at 4000 if no addr is passed.
func WithConsumerAddr(consumerAddrs ...string) ServerOption {
	return func(server *Server) {
		server.ConsumerAddrs = consumerAddrs
	}
}

func NewServer(serverOpts ...ServerOption) (*Server, error) {
	// a channel shared between producers and the server where producers push
	// messages and server adds the message into their respective topic stores

	server := &Server{
		ProducerAddrs: []string{":3000"},
		ConsumerAddrs: []string{":4000"},
		quitChannel:   make(chan struct{}),
	}
	for _, opt := range serverOpts {
		opt(server)
	}

	consumerChannel := make(chan *websocket.Conn)

	server.loadDistributor = *load.NewLoadDistributor()
	server.producerBroker = producer.NewProducerBroker(&server.loadDistributor)

	server.consumerProxy = consumer.NewWSConsumerProxy(server.ConsumerAddrs, consumerChannel)

	return server, nil
}

// Registers producers and consumers associated with the server and
// starts publishing messages to topics.
// func (s *Server) Start() error {
// 	for _, producerHandler := range s.producerHandlers {
// 		go func(ph transport.ProducerHandler) {
// 			err := ph.Start()
// 			if err != nil {
// 				// if one producer is failing, doesn't mean whole
// 				// server has to be stopped, so print and move on
// 				fmt.Println(err)
// 			}
// 		}(producerHandler)
// 	}

// 	for _, subscriptionHandler := range s.subscriptionHandlers {
// 		go func(sh transport.SubscriptionHandler) {
// 			err := sh.Start()
// 			if err != nil {
// 				// if one consumer is failing, doesn't mean whole
// 				// server has to be stopped, so print and move on
// 				fmt.Println(err)
// 			}
// 		}(subscriptionHandler)
// 	}

// 	for _, streamProcessor := range s.streamProcessors {
// 		go func(sp transport.StreamProcessor) {
// 			err := sp.Start()
// 			if err != nil {
// 				// if one consumer is failing, doesn't mean whole
// 				// server has to be stopped, so print and move on
// 				fmt.Println(err)
// 			}
// 		}(streamProcessor)
// 	}

// 	for {
// 		select {
// 		case <-s.quitChannel:
// 			return nil
// 		case message := <-s.produceChannel:
// 			go func(s *Server, m transport.Message) {
// 				offset, err := s.publishMessage(m)
// 				if err != nil {
// 					slog.Error("could not publish message=%s to topic=%s",
// 						m.Topic, string(m.Data))
// 					return
// 				}
// 				slog.Info("produced", "message", m, "offset", offset)
// 			}(s, message)
// 		}
// 	}
// }

func (s *Server) Start() error {

	consumerPaths := make(chan string)
	// Start components
	if err := s.loadDistributor.Start(consumerPaths); err != nil {
		return err
	}
	if err := s.consumerProxy.Start(); err != nil {
		return err
	}

	producerPaths := make(chan string)
	if err := s.producerBroker.Start(producerPaths); err != nil {
		return err
	}

	slog.Info("🎉 Ruskin ready for writing... ✒️")
	return nil
}
