package server

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/axatol/jayd/pkg/config"
	"github.com/axatol/jayd/pkg/config/nr"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
)

func middleware_ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func middleware_CORS(next http.Handler) http.Handler {
	options := cors.Options{
		AllowedOrigins:   strings.Split(config.ServerCORSList, ","),
		AllowedMethods:   []string{http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodDelete},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
	}

	return cors.Handler(options)(next)
}

func middleware_JWT(next http.Handler) http.Handler {
	if !config.Auth0Enabled {
		return next
	}

	issuer := fmt.Sprintf("https://%s/", config.Auth0Domain)
	issuerURL, err := url.Parse(issuer)
	if err != nil {
		log.Fatal().Err(err).Str("issuer", issuer).Msg("invalid issuer url")
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	validator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{config.Auth0Audience},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to configure jwt validator")
	}

	errhandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Error().Err(err).Msg("error validating jwt")
		responseErr(w, err_Unauthorised, http.StatusUnauthorized)
	}

	middleware := jwtmiddleware.New(
		validator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errhandler),
	)

	return middleware.CheckJWT(next)
}

func middleware_NewRelic(next http.Handler) http.Handler {
	if !nr.Enabled {
		return next
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		route, ok := getRoutePattern(r)
		if !ok {
			route = r.RequestURI
		}

		_, handler := newrelic.WrapHandleFunc(nr.App, route, next.ServeHTTP)
		handler(w, r)
	}

	return http.HandlerFunc(fn)
}

func middleware_Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := r.Context().Value(middleware.RequestIDKey)
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			event := log.Info()

			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}

			url := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

			if rec := recover(); rec != nil {
				event = log.Error().
					Interface("recovery", rec).
					Bytes("stack", debug.Stack())
			}

			event.
				Int("status", ww.Status()).
				Str("method", r.Method).
				Str("url", url).
				Str("proto", r.Proto).
				Dur("duration", time.Since(start)).
				Int("bytes_written", ww.BytesWritten()).
				Str("remote_addr", r.RemoteAddr).
				Str("origin", r.Header.Get("Origin")).
				Str("request_id", (requestID).(string)).
				Send()
		}()

		next.ServeHTTP(ww, r)
	})
}
