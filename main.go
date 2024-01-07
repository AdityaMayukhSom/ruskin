package main

import (
	"fmt"
	"log"

	server "github.com/AdityaMayukhSom/ruskin/server"
	store "github.com/AdityaMayukhSom/ruskin/store"
)

func main() {
	serverconf := &server.ServerConfig{
		Port:         ":3000",
		StoreFactory: store.NewMemoryStoreFactory(nil),
	}

	server, err := server.NewServer(serverconf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(server)
}
