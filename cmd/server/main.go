package main

import (
	"net/http"

	"github.com/axatol/jayd/config"
	"github.com/axatol/jayd/server"
	"github.com/rs/zerolog/log"
)

func main() {
	router := server.Init()

	log.Info().
		Bool("debug", config.Debug).
		Str("server_address", config.ServerAddress).
		Msg("starting server")

	if err := http.ListenAndServe(config.ServerAddress, router); err != nil {
		log.Fatal().Err(err).Send()
	}
}
