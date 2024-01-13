/* The server must be running in another terminal on port 3000 */

package integratio_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"testing"

	transport "github.com/AdityaMayukhSom/ruskin/transport"
)

func TestPublishOverConnection(t *testing.T) {
	SERVER_PORT := 3000
	ROUTE_NAME := "/publish"

	BASE_URL := fmt.Sprintf(
		"http://localhost:%d%s",
		SERVER_PORT, ROUTE_NAME,
	)

	MIME_TYPE := "application/octet-stream"

	topics := []string{
		"food",
		"movies",
		"songs",
		"education",
	}

	topicsLen := len(topics)

	for i := 0; i < 1000; i++ {
		topic := topics[rand.Intn(topicsLen)]
		payload := []byte(fmt.Sprintf("message_%s_%d", topic, i))
		resp, err := http.Post(
			fmt.Sprintf("%s/%s", BASE_URL, topic),
			MIME_TYPE,
			bytes.NewReader(payload),
		)

		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		var responseBody transport.TransportResponse
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(responseBody)
	}
}
