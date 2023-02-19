package config

import (
	"os"
	"strconv"
	"strings"

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

func envValueInt(envKey string, fallback ...int) int {
	fallbackValue := 0
	if len(fallback) == 1 {
		fallbackValue = fallback[0]
	}

	value := envValue(envKey, strconv.Itoa(fallbackValue))
	if value == "" {
		return fallbackValue
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal().Str("key", envKey).Msg("failed to parse int type environment variable")
		return fallbackValue
	}

	return result
}

func envValueList(envKey string) []string {
	value := envValue(envKey, "")
	if value == "" {
		return []string{}
	}

	return strings.Split(value, ",")
}
