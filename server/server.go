package server

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"

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
	// a channel shared between producers and the server
	// where producers push messages and server adds the message into
	// their respective topic stores
	var produceChannel = make(chan transport.Message)

	server := &Server{
		ServerConfig: config,
		producers: []transport.Producer{
			transport.NewHTTPProducer(
				config.ListenAddr,
				produceChannel,
			),
		},
		consumers:   []transport.Consumer{},
		topicStores: make(map[string]store.Store),

		produceChannel: produceChannel,
		quitChannel:    make(chan struct{}),
	}

	return server, nil
}

// Returns true if a store for the topic exists, otherwise returns false.
func (s *Server) checkStore(topicName string) (bool, error) {
	if len(topicName) == 0 {
		return false, errors.New("cannot check for store with empty name")
	}

	// Checks whether a topic with the name exists in the map or not.
	// ok will be true if it exists, if it doesn't, ok will be false
	s.getStoreStateMut.RLock()
	_, found := s.topicStores[topicName]
	s.getStoreStateMut.RUnlock()

	// errors.New("unable to check if store exists or not")
	return found, nil
}

// Creates a new topic if the topic does not already exists.
// Returns true if topic is successfully created, otherwise returns false.
func (s *Server) createStore(topicName string) (store.Store, error) {
	if len(topicName) == 0 {
		return nil, errors.New("cannot create topic store with empty name")
	}

	var topicStore store.Store

	// As we may require to modify the state of topicStores, we first lock
	// the changeStoreStateMut variable. Now no other function shall be able
	// to modify the state of topicStores (considering createStore is the
	// only way to modify topicStores)
	s.changeStoreStateMut.Lock()

	found, err := s.checkStore(topicName)
	if err == nil && !found {
		// err is nil here
		topicStore = s.StoreFactory.Produce()

		// we need to modify the state of the store
		// hence stop read access to check store
		s.getStoreStateMut.Lock()

		// add the newly generated store for the given topic
		s.topicStores[topicName] = topicStore

		s.getStoreStateMut.Unlock()
		slog.Info("created store", "topic", topicName)
	} else if err != nil {
		topicStore = nil
	} else if found {
		// err is nil here too
		topicStore = s.topicStores[topicName]
	}

	s.changeStoreStateMut.Unlock()

	return topicStore, err
}

// Returns the corresponding store
func (s *Server) getStore(topicName string) (store.Store, error) {
	if len(topicName) == 0 {
		return nil, errors.New("cannot get topic store with empty name")
	}

	var topicStore store.Store
	found, err := s.checkStore(topicName)
	if err == nil && !found {
		if topicStore, err = s.createStore(topicName); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else if found {
		topicStore = s.topicStores[topicName]
	}

	return topicStore, nil
}

// Publishes the message to the topic mentioned in `Message.Topic` field.
//
// If the mentioned topic does not exist, this function will create a
// new topic and then publish the data into that topic.
//
// Returns the offset of the message inside the topic store as first value.
//
// Returns an error if the topic name is empty string or unable
// to insert data into the topic store. If an error is thrown then
// offset value returned is -1.
func (s *Server) publishMessage(message transport.Message) (int, error) {
	store, err := s.getStore(message.Topic)
	if err != nil {
		return -1, fmt.Errorf(
			"cannot find store for topic %s", message.Topic)
	}

	offset, err := store.Insert(message.Data)
	if err != nil {
		return -1, fmt.Errorf("could not insert message=%s to topic=%s",
			string(message.Data), message.Topic)
	}

	return offset, nil
}

// Registers producers and consumers associated with the server and
// starts publishing messages to topics.
func (s *Server) Start() error {
	// for _, consumer := range s.consumers {
	// 	go func(c transport.Consumer) {
	// 		err := c.Start()
	// 		if err != nil {
	// 			// if one consumer is failing, doesn't mean whole
	// 			// server has to be stopped, so print and move on
	// 			fmt.Println(err)
	// 		}
	// 	}(consumer)
	// }

	for _, producer := range s.producers {
		go func(p transport.Producer) {
			err := p.Start()
			if err != nil {
				// if one producer is failing, doesn't mean whole
				// server has to be stopped, so print and move on
				fmt.Println(err)
			}
		}(producer)
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
