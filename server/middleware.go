package server

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/axatol/jayd/config"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func middleware_ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func middleware_CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func middleware_JWT(next http.Handler) http.Handler {
	if config.Auth0Audience == "" || config.Auth0Domain == "" {
		log.Warn().
			Str("auth0_audience", config.Auth0Audience).
			Str("auth0_domain", config.Auth0Domain).
			Msg("skipping JWT middleware configuration")
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
				event = event.
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
				Str("request_id", (requestID).(string)).
				Send()
		}()

		next.ServeHTTP(ww, r)
	})
}
