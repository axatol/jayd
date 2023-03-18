package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/axatol/jayd/config"
	"github.com/axatol/jayd/downloader"
	"github.com/axatol/jayd/server"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().
		Str("file", config.ServerBackupFile).
		Msg("loading cache from file")

	if err := downloader.CreateCache(config.ServerBackupFile); err != nil {
		log.Error().
			Err(err).
			Str("file", config.ServerBackupFile).
			Msg("could not load cache from backup")
	}

	defer func() {
		log.Info().
			Str("file", config.ServerBackupFile).
			Msg("saving cache to file")

		if err := downloader.SaveCache(config.ServerBackupFile); err != nil {
			log.Error().
				Err(err).
				Str("file", config.ServerBackupFile).
				Msg("could not save cache to backup")
		}
	}()

	server := server.Init()
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("server exited")
		}
	}()

	log.Info().
		Bool("debug", config.Debug).
		Str("server_address", config.ServerAddress).
		Msg("started server")

	waitForInterrupt()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error().
			Err(err).
			Msg("error shutting down server")
	}
}

func waitForInterrupt() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}
