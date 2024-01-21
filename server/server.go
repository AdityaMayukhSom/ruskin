package server

import (
	"fmt"
	"log/slog"

	store "github.com/AdityaMayukhSom/ruskin/store"
	transport "github.com/AdityaMayukhSom/ruskin/transport"
)

type ServerConfig struct {
	ProducerAddr string
	ConsumerAddr string
}

type Server struct {
	// We must use two mutex to synchronise the creation of new topic store.
	*ServerConfig

	producerHandlers     []transport.ProducerHandler
	subscriptionHandlers []transport.SubscriptionHandler
	streamProcessors     []transport.StreamProcessor

	// consumeChannel will be used to push the topicname to the client handler
	// for spawning appropriate client relay
	consumeChannel chan<- *store.Store

	// When we are willing to shutdown the server gracefully,
	// we need to signal this channel or close it.
	quitChannel chan struct{}
}

func NewServer(config *ServerConfig) (*Server, error) {
	// a channel shared between producers and the server
	// where producers push messages and server adds the message into
	// their respective topic stores
	var produceChannel = make(chan transport.Message)

	server := &Server{
		ServerConfig: config,
		topicStores:  make(map[string]store.Store),

		producerHandlers: []transport.ProducerHandler{
			transport.NewHTTPProducerHandler(
				config.ProducerAddr,
				produceChannel,
			),
		},
		produceChannel: produceChannel,

		// consumers: []transport.Consumer{
		// 	transport.NewWSConsumer(
		// 		config.ConsumerAddr,
		// 	),
		// },

		quitChannel: make(chan struct{}),
	}

	return server, nil
}

func (s *Server) notifySubscribers(topicName string) error {

	return nil
}

// Registers producers and consumers associated with the server and
// starts publishing messages to topics.
func (s *Server) Start() error {
	for _, producerHandler := range s.producerHandlers {
		go func(ph transport.ProducerHandler) {
			err := ph.Start()
			if err != nil {
				// if one producer is failing, doesn't mean whole
				// server has to be stopped, so print and move on
				fmt.Println(err)
			}
		}(producerHandler)
	}

	for _, subscriptionHandler := range s.subscriptionHandlers {
		go func(sh transport.SubscriptionHandler) {
			err := sh.Start()
			if err != nil {
				// if one consumer is failing, doesn't mean whole
				// server has to be stopped, so print and move on
				fmt.Println(err)
			}
		}(subscriptionHandler)
	}

	for _, streamProcessor := range s.streamProcessors {
		go func(sp transport.StreamProcessor) {
			err := sp.Start()
			if err != nil {
				// if one consumer is failing, doesn't mean whole
				// server has to be stopped, so print and move on
				fmt.Println(err)
			}
		}(streamProcessor)
	}

	for {
		select {
		case <-s.quitChannel:
			return nil
		case message := <-s.produceChannel:
			go func(s *Server, m transport.Message) {
				offset, err := s.publishMessage(m)
				if err != nil {
					slog.Error("could not publish message=%s to topic=%s",
						m.Topic, string(m.Data))
					return
				}
				slog.Info("produced", "message", m, "offset", offset)
			}(s, message)
		}
	}
}
