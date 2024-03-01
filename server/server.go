package server

import (
	"log/slog"

	"github.com/AdityaMayukhSom/ruskin/consumer"
	"github.com/AdityaMayukhSom/ruskin/load"
	"github.com/AdityaMayukhSom/ruskin/producer"
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

	consumerProxyLoadBalancerChannel := make(chan *consumer.Consumer)

	server.loadDistributor = load.NewLoadDistributor(consumerProxyLoadBalancerChannel)
	server.producerBroker = producer.NewProducerBroker(&server.loadDistributor)
	server.consumerProxy = consumer.NewWSConsumerProxy(server.ConsumerAddrs, consumerProxyLoadBalancerChannel)

	return server, nil
}

func (s *Server) Start() error {
	producerPaths := make(chan string)

	slog.Info("load distributor starting")
	if err := s.loadDistributor.Start(); err != nil {
		return err
	}
	slog.Info("load distributor started")

	slog.Info("consumer proxy starting")
	if err := s.consumerProxy.Start(); err != nil {
		return err
	}
	slog.Info("consumer proxy started")

	slog.Info("producer broker starting")
	if err := s.producerBroker.Start(producerPaths); err != nil {
		return err
	}
	slog.Info("producer broker started")

	slog.Info("🎉 Ruskin ready for writing... ✒️")
	return nil
}
