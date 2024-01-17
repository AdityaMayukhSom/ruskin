package consumer

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

type ConsumerProxy interface {
	Start() error
}

type WSConsumerProxy struct {
	listenAddr        string
	connectionChannel chan<- *websocket.Conn
}

func (wscp *WSConsumerProxy) socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Info("Error during connection upgradation:", err)
		return
	}

	wscp.connectionChannel <- conn
}

func (wscp *WSConsumerProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (wscp *WSConsumerProxy) Start() error {
	slog.Info("web socket consumer proxy started", "port", wscp.listenAddr)

	// http.ListenAndServe is a blocking method unless there is an error
	return http.ListenAndServe(wscp.listenAddr, wscp)
}
