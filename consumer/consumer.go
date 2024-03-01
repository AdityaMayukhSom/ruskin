package consumer

import (
	"github.com/AdityaMayukhSom/ruskin/messagequeue"
	websocket "github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Consumer struct {
	Topic messagequeue.TopicIdentifier
	// The websocket connection.
	Conn *websocket.Conn
}
