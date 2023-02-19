package server

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	r = chi.NewRouter()
)

func Init() chi.Router {
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware_ContentType)
	r.Use(middleware_CORS)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 30))

	r.Route("/api", func(r chi.Router) {
		r.Get("/video/metadata", handler_GetVideoMetadata)
		r.Get("/video", handler_GetVideoFile)
		r.Post("/video", handler_QueueVideoDownload)
	})

	return r
}
