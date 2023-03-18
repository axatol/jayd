package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	envFile, _                = godotenv.Read()
	Debug                     = envValueBool("DEBUG", false)
	ServerBackupFile          = envValue("SERVER_BACKUP_FILE", "/data/cache.json")
	ServerAddress             = envValue("SERVER_ADDRESS", ":8000")
	ServerCORSList            = envValue("SERVER_CORS_LIST", "")
	DownloaderExecutable      = envValue("DOWNLOADER_EXECUTABLE", "yt-dlp")
	DownloaderOutputDirectory = envValue("DOWNLOADER_OUTPUT_DIR", "/data/output")
	DownloaderCacheDirectory  = envValue("DOWNLOADER_CACHE_DIR", "/data/cache")
	DownloaderRetries         = envValueInt("DOWNLOADER_RETRIES", 1)
	DownloaderConcurrency     = envValueInt("DOWNLOADER_CONCURRENCY", 1)
	YoutubeAPIKey             = envValue("YOUTUBE_API_KEY")
	Auth0Enabled              = envValueBool("AUTH0_ENABLED", true)
	Auth0Domain               = envValue("AUTH0_DOMAIN", "")
	Auth0Audience             = envValue("AUTH0_AUDIENCE", "")
	WebDirectory              = envValue("WEB_DIRECTORY", "/web")
	StorageEnabled            = envValueBool("STORAGE_ENABLED", true)
	StorageEndpoint           = envValue("STORAGE_ENDPOINT", "")
	StorageBucketName         = envValue("STORAGE_BUCKET_NAME", "jayd")
	StorageAccessKeyID        = envValue("STORAGE_ACCESS_KEY_ID", "")
	StorageSecretAccessKey    = envValue("STORAGE_SECRET_ACCESS_KEY", "")
)

func init() {
	flag.BoolVar(&Debug, "debug", Debug, "enable debug mode")
	flag.StringVar(&ServerBackupFile, "server-backup-file", ServerBackupFile, "server backup file")
	flag.StringVar(&ServerAddress, "server-address", ServerAddress, "enable debug mode")
	flag.StringVar(&ServerCORSList, "server-cors-list", ServerCORSList, "server cors list")
	flag.StringVar(&DownloaderExecutable, "downloader-executable", DownloaderExecutable, "downloader executable")
	flag.StringVar(&DownloaderOutputDirectory, "downloader-output-directory", DownloaderOutputDirectory, "downloader output directory")
	flag.StringVar(&DownloaderCacheDirectory, "downloader-cache-directory", DownloaderCacheDirectory, "downloader cache directory")
	flag.IntVar(&DownloaderRetries, "downloader-retries", DownloaderRetries, "downloader retries")
	flag.IntVar(&DownloaderConcurrency, "downloader-concurrency", DownloaderConcurrency, "downloader concurrency")
	flag.StringVar(&YoutubeAPIKey, "youtube-api-key", YoutubeAPIKey, "youtube api key")
	flag.BoolVar(&Auth0Enabled, "auth0-enabled", Auth0Enabled, "auth0 enabled")
	flag.StringVar(&Auth0Domain, "auth0-domain", Auth0Domain, "auth0 domain")
	flag.StringVar(&Auth0Audience, "auth0-audience", Auth0Audience, "auth0 audience")
	flag.StringVar(&WebDirectory, "web-directory", WebDirectory, "web directory")
	flag.BoolVar(&StorageEnabled, "storage-enabled", StorageEnabled, "storage enabled")
	flag.StringVar(&StorageEndpoint, "storage-endpoint", StorageEndpoint, "storage endpoint")
	flag.StringVar(&StorageBucketName, "storage-bucket-name", StorageBucketName, "storage bucket name")
	flag.StringVar(&StorageAccessKeyID, "storage-access-key-id", StorageAccessKeyID, "storage access key id")
	flag.StringVar(&StorageSecretAccessKey, "storage-secret-access-key", StorageSecretAccessKey, "storage secret access key")
	flag.Parse()

	if Debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	if Auth0Enabled && (Auth0Domain == "" || Auth0Audience == "") {
		log.Fatal().
			Str("auth0_audience", Auth0Audience).
			Str("auth0_domain", Auth0Domain).
			Msg("must provide auth0 domain and audience")
	}

	if StorageEnabled && (StorageEndpoint == "" || StorageAccessKeyID == "" || StorageSecretAccessKey == "") {
		log.Fatal().
			Str("storage_endpoint", StorageEndpoint).
			Msg("must provide storage endpoint and credentials")
	}

}
