package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PublishedResponse struct {
	Message string
	Offset  int
}

func NewPublishedResponse(message string, offset int) *PublishedResponse {
	return &PublishedResponse{
		Message: message,
		Offset:  offset,
	}
}

func (pr PublishedResponse) String() string {
	return fmt.Sprintf("{\"Message\": %s, \"Offset\": %d}", pr.Message, pr.Offset)
}

type ErrorResponse struct {
	Message string
}

func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
	}
}

func (er ErrorResponse) String() string {
	return fmt.Sprintf("{\"Message\": %s}", er.Message)
}

func WriteResponse[T fmt.Stringer](w http.ResponseWriter, statusCode int, responseBody *T) {
	// you need to set headers before setting WriteHeader with status code
	// refer to https://stackoverflow.com/questions/39273490/adding-a-header-to-a-responsewriter
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if respBody, err := json.Marshal(responseBody); err == nil {
		w.Write(respBody)
	}
}
