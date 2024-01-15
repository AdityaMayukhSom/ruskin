/* Server must be running in another terminal */
/* Producer port must be 3000 */

package integration_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"testing"
)

func TestPublishOverConnection(t *testing.T) {
	const SERVER_PORT = 3000
	const ROUTE_NAME = "/publish"
	const MIME_TYPE = "application/octet-stream"
	BASE_URL := fmt.Sprintf("http://localhost:%d%s", SERVER_PORT, ROUTE_NAME)

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

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}
}
