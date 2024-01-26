package consumer

import (
	websocket "github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Consumer struct {
	// The websocket connection.
	Conn *websocket.Conn
}
