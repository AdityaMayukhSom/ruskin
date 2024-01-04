package main

import (
	"fmt"
	"log"
)

func main() {
	serverconf := &ServerConfig{
		port:         ":3000",
		storeFactory: NewMemoryStoreFactory(nil),
	}

	server, err := NewServer(serverconf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(server)
}
