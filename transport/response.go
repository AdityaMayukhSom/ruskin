package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TransportResponse struct {
	Message string
}

func (tr TransportResponse) String() string {
	return fmt.Sprintf("{\"Message\": %s}", tr.Message)
}

func writeResponse(w http.ResponseWriter, statusCode int, responseMessage string) {
	w.WriteHeader(http.StatusBadRequest)
	if respBody, err := json.Marshal(TransportResponse{
		Message: responseMessage,
	}); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(respBody)
	}
}
