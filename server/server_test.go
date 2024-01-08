package server

import (
	"log"
	"testing"

	store "github.com/AdityaMayukhSom/ruskin/store"
)

// Function to return a mock server for testing.
// TODO: Need to mock the ServerConfig otherwise we are not able
// to isolate the Server completely for unit testing.
func getNewServer() *Server {
	serverconf := &ServerConfig{
		StoreFactory: store.NewMemoryStoreFactory(nil),
	}

	server, err := NewServer(serverconf)
	if err != nil {
		log.Fatal(err)
	}

	return server
}

func TestIfStoreExistsInServer(t *testing.T) {
	t.Run("test non-existing key", func(t *testing.T) {
		server := getNewServer()

		server.createStore("test_store")
		found, err := server.checkStore("test_store")

		if err != nil {
			logger := log.Default()
			logger.Println(err)
			t.Fail()
		}

		if !found {
			logger := log.Default()
			logger.Println("topic inserted but not found")
			t.Fail()
		}
	})

	t.Run("test existing key", func(t *testing.T) {
		server := getNewServer()

		found, err := server.checkStore("test_store")

		if err != nil {
			logger := log.Default()
			logger.Println(err)
			t.Fail()
		}

		if found {
			logger := log.Default()
			logger.Println("topic not inserted but found")
			t.Fail()
		}
	})

}
