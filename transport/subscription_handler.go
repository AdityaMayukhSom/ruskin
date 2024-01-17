package transport

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"

	consumer "github.com/AdityaMayukhSom/ruskin/consumer"
	store "github.com/AdityaMayukhSom/ruskin/store"
)

var upgrader = websocket.Upgrader{}

type SubscriptionHandler interface {
	Start() error
}

type WSSubscriptionHandler struct {
	listenAddr     string
	consumerRelays map[*store.Store]consumer.ConsumerRelay
}

func NewWSSubscriptionHandler(listenAddr string) *WSSubscriptionHandler {
	return &WSSubscriptionHandler{
		listenAddr:     listenAddr,
		consumerRelays: make(map[*store.Store]consumer.ConsumerRelay),
	}
}

func (s *WSSubscriptionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// handle websocket
}

func (s *WSSubscriptionHandler) Start() error {
	slog.Info("websockets consumer started", "port", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, s)

}
