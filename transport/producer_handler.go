package transport

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type ProducerHandler interface {
	Start() error
}

type HTTPProducerHandler struct {
	listenAddr     string
	produceChannel chan<- Message
}

func NewHTTPProducerHandler(listenAddr string, produceChannel chan<- Message) *HTTPProducerHandler {
	return &HTTPProducerHandler{
		listenAddr:     listenAddr,
		produceChannel: produceChannel,
	}
}

// Implementing Handler interface.
// Refer to [net/http.Handler] interface.
func (p *HTTPProducerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello From Server!"))
	case http.MethodPost:
		if len(parts) != 2 {
			errMsg := fmt.Sprintf(
				"invalid route param count : passed %d",
				len(parts),
			)

			slog.Error(errMsg)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		protoName, topicName := parts[0], parts[1]
		switch protoName {
		case "publish":
			data, err := io.ReadAll(r.Body)

			if err != nil {
				errMsg := "could not read request body for data to be published"

				slog.Error(errMsg)
				http.Error(w, errMsg, http.StatusBadRequest)

				return
			}

			slog.Info("publishing under", "topic", topicName)

			p.produceChannel <- Message{
				Topic: topicName,
				Data:  data,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("message published"))

		default:
			errMsg := fmt.Sprintf("route %s not supported", protoName)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}
	case http.MethodDelete:

	default:
		w.WriteHeader(405)
	}
}

func (p *HTTPProducerHandler) Start() error {
	slog.Info("http producer started", "port", p.listenAddr)

	// http.ListenAndServe is a blocking method unless there is an error
	return http.ListenAndServe(p.listenAddr, p)
}
