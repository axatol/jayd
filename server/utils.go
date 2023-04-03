package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// https://github.com/rl404/fairy/blob/v0.22.0/monitoring/newrelic/middleware/http.go
func getRoutePattern(r *http.Request) (string, bool) {
	routePath := r.URL.Path
	if r.URL.RawPath != "" {
		routePath = r.URL.RawPath
	}

	rctx := chi.RouteContext(r.Context())
	tctx := chi.NewRouteContext()
	if rctx.Routes.Match(tctx, r.Method, routePath) {
		return tctx.RoutePattern(), true
	}

	return "", false
}
