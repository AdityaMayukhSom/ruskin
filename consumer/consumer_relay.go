package consumer

import (
	messagequeue "github.com/AdityaMayukhSom/ruskin/messagequeue"

	"github.com/gorilla/websocket"
)

type ConsumerRelay interface {
	Relay() error
}

type WSConsumerRelay struct {
	topicStore *messagequeue.Store

	// Registered clients.
	clients map[*websocket.Conn]bool

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
		clients:    map[*websocket.Conn]bool{},
		broadcast:  make(chan []byte),
		register:   make(chan *Consumer),
		unregister: make(chan *Consumer),
	}
}

func (wscr *WSConsumerRelay) Relay() error {
	// this should relay messages to all the connected consumers
	// before that consume message from

	msg, err := (*wscr.topicStore).Extract(10)
	if err != nil {
		return err
	}

	// sends the message to the corresponding consumer
	for consumer, isValid := range wscr.clients {
		if isValid {
			consumer.WriteMessage(websocket.BinaryMessage, msg)
		}
	}
	return nil
}
