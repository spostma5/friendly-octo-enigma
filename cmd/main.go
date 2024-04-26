package main

import (
	"log"
	"log/slog"

	"github.com/spostma5/friendly-octo-enigma/cmd/api"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	server := api.NewServer("localhost", 8080)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
