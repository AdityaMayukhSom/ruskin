package consumer

import (
	messagequeue "github.com/AdityaMayukhSom/ruskin/messagequeue"

	websocket "github.com/gorilla/websocket"

	mapset "github.com/deckarep/golang-set/v2"
)

type WSConsumerRelay struct {
	topicStore *messagequeue.Store

	// Registered clients, we will use a threadsafe Set implementation.
	clients mapset.Set[*websocket.Conn]

	// Inbound messages from the clients.
	// broadcast chan []byte
	broadcast chan []byte

	// Register requests from the clients.
	// register chan *Consumer
	register chan *Consumer

	// Unregister requests from clients.
	// unregister chan *Consumer
	unregister chan *Consumer
}

func NewWSConsumerRelay(topicStore *messagequeue.Store) *WSConsumerRelay {
	return &WSConsumerRelay{
		topicStore: topicStore,
		clients:    mapset.NewSet[*websocket.Conn](),
		broadcast:  make(chan []byte),
		register:   make(chan *Consumer),
		unregister: make(chan *Consumer),
	}
}

// Adds a new connection for the messages to be relayed.
func (wscr *WSConsumerRelay) AddConsumer(consumer *Consumer) bool {
	isadded := wscr.clients.Add(consumer.Conn)
	return isadded
}

// Sends the message to all
func (wscr *WSConsumerRelay) Relay() error {
	// this should relay messages to all the connected consumers
	// before that consume message from

	msg, err := (*wscr.topicStore).Extract(10)
	if err != nil {
		return err
	}

	// sends the message to the corresponding consumer

	for consumer := range wscr.clients.Iter() {
		consumer.WriteMessage(websocket.BinaryMessage, msg)
	}

	return nil
}
