package transport

import (
	"log/slog"
	"net/http"
)

type Consumer interface {
	Start() error
}

type WSConsumer struct {
	listenAddr string
}

func NewWSConsumer(listenAddr string) *WSConsumer {
	return &WSConsumer{
		listenAddr: listenAddr,
	}
}

func (c *WSConsumer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (c *WSConsumer) Start() error {
	slog.Info("websockets consumer started", "port", c.listenAddr)

	return http.ListenAndServe(c.listenAddr, c)
}
