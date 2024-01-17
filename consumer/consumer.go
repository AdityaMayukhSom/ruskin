package consumer

import (
	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Consumer struct {
	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Topics []string
}
