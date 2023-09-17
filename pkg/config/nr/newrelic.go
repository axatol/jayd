package nr

import (
	"os"

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/zerologWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	App     *newrelic.Application
	Enabled bool
)

func Configure() {
	Enabled = os.Getenv("NEW_RELIC_ENABLED") == "true"
	if !Enabled {
		return
	}

	app, err := newrelic.NewApplication(newrelic.ConfigFromEnvironment())
	if err != nil {
		log.Fatal().
			Err(err).
			Bool("enabled", Enabled).
			Str("app_name", os.Getenv("NEW_RELIC_APP_NAME")).
			Msg("could not configure newrelic")
	}

	writer := zerologWriter.New(os.Stdout, app)
	log.Logger = zerolog.New(writer)

	App = app
}
