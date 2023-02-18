package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	envFile, _ = godotenv.Read()
)

func envValue(envKey string, fallback ...string) string {
	value, ok := os.LookupEnv(envKey)
	if ok {
		return value
	}

	if value, ok := envFile[envKey]; ok {
		return value
	}

	if len(fallback) == 1 {
		return fallback[0]
	}

	log.Fatal().Str("key", envKey).Msg("required environment variable not available")
	return ""
}

var (
	Debug         = envValue("DEBUG", "false") == "true"
	ServerAddress = envValue("SERVER_ADDRESS", ":8000")
)

func init() {
	flag.BoolVar(&Debug, "debug", false, "enable debug mode")
	flag.StringVar(&ServerAddress, "server-address", ":8000", "enable debug mode")
	flag.Parse()

	if Debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
}
