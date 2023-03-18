package config

import (
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
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

	return ""
}

func envValueInt(envKey string, fallback int) int {
	value := envValue(envKey, strconv.Itoa(fallback))
	if value == "" {
		return fallback
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal().Str("key", envKey).Msg("failed to parse int type environment variable")
		return fallback
	}

	return result
}

func envValueBool(envKey string, fallback bool) bool {
	value := envValue(envKey)
	if value == "" {
		return fallback
	}

	return value == "true" || value == "1"
}
