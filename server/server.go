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
	// We must use two mutex to synchronise the creation of new topic store.
	*ServerConfig

	// It is used when a goroutine is trying to check if a store
	// corresponding to a topic already exists or not.
	getStoreStateMut sync.RWMutex

	// It is used in when a goroutine is willing to update topicStores.
	changeStoreStateMut sync.Mutex

	producers []transport.Producer
	consumers []transport.Consumer

	// To update topic store in a goroutine safe manner, the developer must
	// lock the changeStoreStateMut before testing for the condition and
	// must write lock getStoreStateMut before actually updating the state
	// and then they must uplock them in the reverse manner.
	topicStores map[string]store.Store

	produceChannel <-chan transport.Message

	// When we are willing to shutdown the server gracefully,
	// we need to signal this channel or close it.
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
