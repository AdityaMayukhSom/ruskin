package server

import (
	"fmt"

	store "github.com/AdityaMayukhSom/ruskin/store"
	transport "github.com/AdityaMayukhSom/ruskin/transport"
)

type ServerConfig struct {
	ListenAddr   string
	StoreFactory store.StoreFactory
}

type Server struct {
	*ServerConfig
	topics      map[string]store.Store
	producers   []transport.Producer
	consumers   []transport.Consumer
	quitChannel chan struct{}
}

func NewServer(config *ServerConfig) (*Server, error) {
	return &Server{
		ServerConfig: config,
		producers: []transport.Producer{
			transport.NewHTTPProducer(
				config.ListenAddr,
			),
		},
		topics:      make(map[string]store.Store),
		quitChannel: make(chan struct{}),
	}, nil
}

func (s *Server) CreateTopic(name string) bool {
	_, ok := s.topics[name]
	if !ok {
		s.topics[name] = s.StoreFactory.Produce()
		return true
	}
	return false
}

func (s *Server) Start() error {

	for _, consumer := range s.consumers {
		err := consumer.Start()
		if err != nil {
			// if one consumer is failing, doesn't mean whole
			// server has to be stopped, so print and move on
			fmt.Println(err)
		}
	}

	for _, producer := range s.producers {
		err := producer.Start()
		if err != nil {
			// if one producer is failing, doesn't mean whole
			// server has to be stopped, so print and move on
			fmt.Println(err)
		}
	}

	<-s.quitChannel
	return nil
}
