/* Server must be running in another terminal */
/* Consumer port must be 000 */

package integration_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/gorilla/websocket"
)

func TestConsumeOverConnection(t *testing.T) {
	SERVER_PORT := 4000
	ROUTE_NAME := "/publish"
	// MIME_TYPE := "application/octet-stream"
	BASE_URL := fmt.Sprintf("ws://localhost:%d%s", SERVER_PORT, ROUTE_NAME)

	topics := []string{
		"food",
		"movies",
		"songs",
		"education",
	}

	for _, topic := range topics {
		_, _, err := websocket.DefaultDialer.Dial(
			fmt.Sprintf("%s/%s", BASE_URL, topic),
			nil,
		)

		if err != nil {
			log.Fatal(err)
		}
	}

}
