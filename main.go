package main

import (
	"log/slog"
	"os"
	"time"

	server "github.com/AdityaMayukhSom/ruskin/server"
	store "github.com/AdityaMayukhSom/ruskin/store"

	tint "github.com/lmittmann/tint"
)

func main() {
	// create a new logger
	tintLogger := slog.New(
		tint.NewHandler(
			os.Stdout,
			&tint.Options{
				Level:      slog.LevelDebug,
				TimeFormat: time.DateTime,
			},
		),
	)

	// set global logger with custom options
	slog.SetDefault(tintLogger)

	serverconf := &server.ServerConfig{
		ProducerAddr: ":3000",
		ConsumerAddr: ":4000",
		StoreFactory: store.NewMemoryStoreFactory(nil),
	}

	server, err := server.NewServer(serverconf)
	if err != nil {
		slog.Error(err.Error())
	}

	err = server.Start()
	if err != nil {
		slog.Error(err.Error())
	}
}
