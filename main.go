package main

import (
	"log/slog"
	"os"
	"time"

	server "github.com/AdityaMayukhSom/ruskin/server"
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

	producerAddrOpt := server.WithProducerAddr(":6000")
	consumerAddrOpt := server.WithConsumerAddr(":6900")

	server, err := server.NewServer(producerAddrOpt, consumerAddrOpt)
	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info("server created 🤖")

	// spawn components and start listeing for consumers and produvcers
	err = server.Start()
	if err != nil {
		slog.Error(err.Error())
	}

	quitCh := make(chan bool)
	quit := <-quitCh
	if quit {
		slog.Info("khatam, tata, bye bye...")
	}

}
