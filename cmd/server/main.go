package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/axatol/jayd/pkg/config"
	"github.com/axatol/jayd/pkg/downloader"
	"github.com/axatol/jayd/pkg/downloader/miniodriver"
	"github.com/axatol/jayd/pkg/server"
	"github.com/rs/zerolog/log"
)

func main() {
	config.Print()

	if err := downloader.CreateCache(config.ServerBackupFile); err != nil {
		log.Error().
			Err(err).
			Str("file", config.ServerBackupFile).
			Msg("could not load cache from backup, will initialise from scratch")
	}

	defer func() {
		if err := downloader.SaveCache(config.ServerBackupFile); err != nil {
			log.Error().
				Err(err).
				Str("file", config.ServerBackupFile).
				Msg("could not save cache to backup")
		}
	}()

	if config.StorageEnabled {
		if _, err := miniodriver.AssertClient(context.Background()); err != nil {
			log.Fatal().
				Err(err).
				Msg("failed to initialise storage driver")
		}
	}

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
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT)
	signal := <-interrupt
	log.Info().Str("signal", signal.String()).Msg("caught interupt")
}
