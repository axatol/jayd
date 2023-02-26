package server

import (
	"net/http"
	"time"

	"github.com/axatol/jayd/config"
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

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware_JWT)
		r.Use(middleware_ContentType)
		r.Get("/youtube/metadata", handler_GetVideoMetadata)
		r.Post("/youtube", handler_QueueVideoDownload)
		r.Get("/queue", handler_ListDownloadQueue)
		r.Delete("/queue", handler_DeleteDownloadQueueItem)
	})

	r.Route("/static", func(r chi.Router) {
		r.Use(middleware_JWT)
		r.Get("/*", handler_StaticContent(config.DownloaderOutputDirectory))
	})

	r.Get("/*", handler_StaticContent(config.WebDirectory))

	return &http.Server{Addr: config.ServerAddress, Handler: r}
}
