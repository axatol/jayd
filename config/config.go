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

var (
	Debug                     = envValue("DEBUG", "false") == "true"
	ServerBackupFile          = envValue("SERVER_BACKUP_FILE", "/data/history.json")
	ServerAddress             = envValue("SERVER_ADDRESS", ":8000")
	ServerCORSList            = envValue("SERVER_CORS_LIST", "")
	DownloaderExecutable      = envValue("DOWNLOADER_EXECUTABLE", "yt-dlp")
	DownloaderOutputDirectory = envValue("DOWNLOADER_OUTPUT_DIR", "/data/output")
	DownloaderCacheDirectory  = envValue("DOWNLOADER_CACHE_DIR", "/data/cache")
	DownloaderRetries         = envValueInt("DOWNLOADER_RETRIES", 1)
	DownloaderConcurrency     = envValueInt("DOWNLOADER_CONCURRENCY", 1)
	YoutubeAPIKey             = envValue("YOUTUBE_API_KEY")
	Auth0Domain               = envValue("AUTH0_DOMAIN", "")
	Auth0Audience             = envValue("AUTH0_AUDIENCE", "")
	WebDirectory              = envValue("WEB_DIRECTORY", "/web")
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
	flag.StringVar(&Auth0Domain, "auth0-domain", Auth0Domain, "auth0 domain")
	flag.StringVar(&Auth0Audience, "auth0-audience", Auth0Audience, "auth0 audience")
	flag.StringVar(&WebDirectory, "web-directory", WebDirectory, "web directory")
	flag.Parse()

	if Debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
}
