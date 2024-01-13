package transport

import (
	"fmt"
	"io"
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
			errMsg := fmt.Sprintf(
				"invalid route param count : passed %d",
				len(parts),
			)

			slog.Error(errMsg)
			writeResponse(w, http.StatusBadRequest, errMsg)
			return
		}

		protoName, topicName := parts[0], parts[1]
		switch protoName {
		case "publish":
			data, err := io.ReadAll(r.Body)

			if err != nil {
				errMsg := "could not read request body for data to be published"

				slog.Error(errMsg)
				writeResponse(w, http.StatusBadRequest, errMsg)
				return
			}

			slog.Info("publishing under", "topic", topicName)

			p.produceChannel <- Message{
				Topic: topicName,
				Data:  data,
			}

			writeResponse(w, http.StatusBadRequest, "message published")

		default:
			errMsg := fmt.Sprintf("route %s not supported", protoName)
			writeResponse(w, http.StatusBadRequest, errMsg)
		}
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
