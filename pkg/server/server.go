package server

import (
	"net/http"
	"time"

	"github.com/axatol/jayd/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	r = chi.NewRouter()
)

func Init() *http.Server {
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware_CORS)
	r.Use(middleware_Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 30))
	r.Use(middleware.Compress(5))
	r.Use(middleware_NewRelic)

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware_JWT)
		r.Use(middleware_ContentType)
		r.Get("/youtube/metadata", handler_GetVideoMetadata)
		r.Post("/youtube", handler_QueueVideoDownload)
		r.Get("/queue", handler_ListDownloadQueue)
		r.Delete("/queue", handler_DeleteQueueItem)

		r.Route("/content", func(r chi.Router) {
			r.Use(middleware_JWT)
			if config.StorageEnabled {
				r.Get("/{key}", handler_PresignedContent)
			} else {
				r.Get("/*", handler_StaticContent(config.DownloaderOutputDirectory))
			}
		})
	})

	r.Route("/ws", func(r chi.Router) {
		r.Use(middleware_JWT)
		r.Use(middleware_ContentType)
		r.Get("/queue", handler_QueueEvents)
	})

	if config.WebDirectory != "" {
		r.Get("/*", handler_StaticContent(config.WebDirectory))
	}

	return &http.Server{Addr: config.ServerAddress, Handler: r}
}
