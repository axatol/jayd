package config

import (
	"flag"
	"os"
	"strings"

	"github.com/axatol/jayd/pkg/config/nr"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	_ = godotenv.Load()

	Debug     = envValueBool("DEBUG", false)
	LogFormat = envValue("LOG_FORMAT", "json")

	Auth0Enabled  = envValueBool("AUTH0_ENABLED", false)
	Auth0Domain   = envValue("AUTH0_DOMAIN", "")
	Auth0Audience = envValue("AUTH0_AUDIENCE", "")

	DownloaderExecutable      = envValue("DOWNLOADER_EXECUTABLE", "yt-dlp")
	DownloaderOutputDirectory = envValue("DOWNLOADER_OUTPUT_DIR", "/data/output")
	DownloaderCacheDirectory  = envValue("DOWNLOADER_CACHE_DIR", "/data/cache")
	DownloaderRetries         = envValueInt("DOWNLOADER_RETRIES", 1)
	DownloaderConcurrency     = envValueInt("DOWNLOADER_CONCURRENCY", 1)

	ServerBackupFile = envValue("SERVER_BACKUP_FILE", "/data/cache.json")
	ServerAddress    = envValue("SERVER_ADDRESS", ":8000")
	ServerCORSList   = envValue("SERVER_CORS_LIST", "")

	StorageEnabled         = envValueBool("STORAGE_ENABLED", false)
	StorageEndpoint        = envValue("STORAGE_ENDPOINT", "")
	StorageBucketName      = envValue("STORAGE_BUCKET_NAME", "jayd")
	StorageAccessKeyID     = envValue("STORAGE_ACCESS_KEY_ID", "")
	StorageSecretAccessKey = envValue("STORAGE_SECRET_ACCESS_KEY", "")

	WebDirectory = envValue("WEB_DIRECTORY", "/web")

	YoutubeAPIKey = envValue("YOUTUBE_API_KEY")
)

func init() {
	flag.BoolVar(&Debug, "debug", Debug, "enable debug mode")
	flag.StringVar(&LogFormat, "log-format", LogFormat, "set log format")

	flag.BoolVar(&Auth0Enabled, "auth0-enabled", Auth0Enabled, "auth0 enabled")
	flag.StringVar(&Auth0Domain, "auth0-domain", Auth0Domain, "auth0 domain")
	flag.StringVar(&Auth0Audience, "auth0-audience", Auth0Audience, "auth0 audience")

	flag.StringVar(&DownloaderExecutable, "downloader-executable", DownloaderExecutable, "downloader executable")
	flag.StringVar(&DownloaderOutputDirectory, "downloader-output-directory", DownloaderOutputDirectory, "downloader output directory")
	flag.StringVar(&DownloaderCacheDirectory, "downloader-cache-directory", DownloaderCacheDirectory, "downloader cache directory")
	flag.IntVar(&DownloaderRetries, "downloader-retries", DownloaderRetries, "downloader retries")
	flag.IntVar(&DownloaderConcurrency, "downloader-concurrency", DownloaderConcurrency, "downloader concurrency")

	flag.StringVar(&ServerBackupFile, "server-backup-file", ServerBackupFile, "server backup file")
	flag.StringVar(&ServerAddress, "server-address", ServerAddress, "enable debug mode")
	flag.StringVar(&ServerCORSList, "server-cors-list", ServerCORSList, "server cors list")

	flag.BoolVar(&StorageEnabled, "storage-enabled", StorageEnabled, "storage enabled")
	flag.StringVar(&StorageEndpoint, "storage-endpoint", StorageEndpoint, "storage endpoint")
	flag.StringVar(&StorageBucketName, "storage-bucket-name", StorageBucketName, "storage bucket name")
	flag.StringVar(&StorageAccessKeyID, "storage-access-key-id", StorageAccessKeyID, "storage access key id")
	flag.StringVar(&StorageSecretAccessKey, "storage-secret-access-key", StorageSecretAccessKey, "storage secret access key")

	flag.StringVar(&WebDirectory, "web-directory", WebDirectory, "web directory")

	flag.StringVar(&YoutubeAPIKey, "youtube-api-key", YoutubeAPIKey, "youtube api key")

	flag.Parse()

	if Debug {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	if LogFormat == "text" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	if Auth0Enabled && (Auth0Domain == "" || Auth0Audience == "") {
		log.Fatal().
			Str("auth0_audience", Auth0Audience).
			Str("auth0_domain", Auth0Domain).
			Msg("must provide auth0 domain and audience")
	}

	nr.Configure()

	if StorageEnabled && (StorageEndpoint == "" || StorageAccessKeyID == "" || StorageSecretAccessKey == "") {
		log.Fatal().
			Str("storage_endpoint", StorageEndpoint).
			Msg("must provide storage endpoint and credentials")
	}
}

func obscure(input string, visiblePrefix int) string {
	if len(input) < 1 {
		return ""
	}

	visible := input[0:visiblePrefix]
	obscured := strings.Repeat("*", len(input)-visiblePrefix)
	return visible + obscured
}

func Print() {
	log.Debug().
		Bool("Debug", Debug).
		Str("LogFormat", LogFormat).
		Bool("Auth0Enabled", Auth0Enabled).
		Str("Auth0Domain", Auth0Domain).
		Str("Auth0Audience", Auth0Audience).
		Str("DownloaderExecutable", DownloaderExecutable).
		Str("DownloaderOutputDirectory", DownloaderOutputDirectory).
		Str("DownloaderCacheDirectory", DownloaderCacheDirectory).
		Int("DownloaderRetries", DownloaderRetries).
		Int("DownloaderConcurrency", DownloaderConcurrency).
		Str("ServerBackupFile", ServerBackupFile).
		Str("ServerAddress", ServerAddress).
		Str("ServerCORSList", ServerCORSList).
		Bool("StorageEnabled", StorageEnabled).
		Str("StorageEndpoint", StorageEndpoint).
		Str("StorageBucketName", StorageBucketName).
		Str("StorageAccessKeyID", obscure(StorageAccessKeyID, 3)).
		Str("StorageSecretAccessKey", obscure(StorageSecretAccessKey, 3)).
		Str("WebDirectory", WebDirectory).
		Str("YoutubeAPIKey", obscure(YoutubeAPIKey, 3)).
		Send()
}
