package transport

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type Producer interface {
	Start() error
}

type HTTPProducer struct {
	listenAddr     string
	produceChannel chan<- Message
}

func NewHTTPProducer(
	listenAddr string, produceChannel chan<- Message,
) *HTTPProducer {
	return &HTTPProducer{
		listenAddr:     listenAddr,
		produceChannel: produceChannel,
	}
}

// Implementing Handler interface.
// Refer to [net/http.Handler] interface.
func (p *HTTPProducer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
		if len(parts) != 2 {
			slog.Error("invalid route param count", "passed", len(parts))
		}

		protoName, topicName := parts[0], parts[1]
		switch protoName {
		case "publish":
			slog.Info("publishing", "topic", topicName)
			p.produceChannel <- Message{
				Topic: topicName,
				Data:  []byte("we don't know yet"),
			}
		default:
			w.WriteHeader(400)
			notSupportedErrMsg := fmt.Sprintf(
				"route %s not supported", protoName)
			w.Write([]byte(notSupportedErrMsg))
		}

		w.Write([]byte("message published"))

	case http.MethodDelete:

	default:
		w.WriteHeader(405)
	}
}

func (p *HTTPProducer) Start() error {
	slog.Info("HTTP transport started at producer", "port", p.listenAddr)

	// http.ListenAndServe is a blocking method unless there is an error
	return http.ListenAndServe(p.listenAddr, p)
}
