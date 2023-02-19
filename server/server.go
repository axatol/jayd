package server

import (
	"time"

	"github.com/axatol/jayd/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	r = chi.NewRouter()
)

func Init() chi.Router {
	r.Use(middleware_ContentType)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware_CORS)
	r.Use(middleware_JWT)
	r.Use(middleware_Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 30))

	r.Route("/api", func(r chi.Router) {
		r.Get("/youtube/metadata", handler_GetVideoMetadata)
		r.Post("/youtube", handler_QueueVideoDownload)
		r.Get("/queue", handler_ListDownloadQueue)
	})

	r.Route("/static", func(r chi.Router) {
		r.Get("/*", handler_StaticContent(config.DownloaderOutputDirectory))
	})

	return r
}
