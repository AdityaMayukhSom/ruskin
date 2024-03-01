package consumer

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/AdityaMayukhSom/ruskin/messagequeue"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

type ConsumerProxy interface {
	Start() error
}

type WSConsumerProxy struct {
	listenAddr                       string
	consumerProxyLoadBalancerChannel chan<- *Consumer
}

func NewWSConsumerProxy(listenAddrs []string, consumerProxyLoadBalancerChannel chan<- *Consumer) *WSConsumerProxy {
	return &WSConsumerProxy{
		// TODO: handle consumers at multiple listen addrs too
		listenAddr:                       listenAddrs[0],
		consumerProxyLoadBalancerChannel: consumerProxyLoadBalancerChannel,
	}
}

func (wscp *WSConsumerProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Info("Error during connection upgradation:", err)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	// TODO: check if parts[0] is consume and parts[1] is topic name

	// puts the connection inside consumer channel
	wscp.consumerProxyLoadBalancerChannel <- &Consumer{
		Conn:  conn,
		Topic: messagequeue.TopicIdentifier(parts[1]),
	}
}

func (wscp *WSConsumerProxy) Start() error {
	slog.Info("web socket consumer proxy started", "port", wscp.listenAddr)

	// http.ListenAndServe is a blocking method unless there is an error
	go http.ListenAndServe(wscp.listenAddr, wscp)
	return nil
}
