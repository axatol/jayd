package nr

import (
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
)

var (
	App     *newrelic.Application
	Enabled bool
)

func Configure() {
	Enabled = os.Getenv("NEW_RELIC_ENABLED") == "true"

	app, err := newrelic.NewApplication(newrelic.ConfigFromEnvironment())
	if err != nil {
		log.Fatal().
			Err(err).
			Str("app_name", os.Getenv("NEW_RELIC_APP_NAME")).
			Msg("could not configure newrelic")
	}

	App = app
}
