package transport

import (
	"log/slog"
	"net/http"
)

type Producer interface {
	Start() error
}

type HTTPProducer struct {
	listenAddr string
}

func NewHTTPProducer(listenAddr string) *HTTPProducer {
	return &HTTPProducer{
		listenAddr: listenAddr,
	}
}

// Implementing Handler interface.
//
// Refer to [net/http.Handler] interface.
func (p *HTTPProducer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handling Request for Producer", "PATH", r.URL.Path)
	w.Write([]byte("hello world"))
}

func (p *HTTPProducer) Start() error {
	slog.Info("HTTP transport started at producer", "port", p.listenAddr)

	// http.ListenAndServe is a blocking method unless there is an error
	return http.ListenAndServe(p.listenAddr, p)
}
