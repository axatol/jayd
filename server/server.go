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
		// TODO - r.Use(pagination)
		// r.Get("/events", getEvents)
		// r.Post("/events", createEvents)
	})

	return r
}
