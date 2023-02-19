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
	ServerAddress             = envValue("SERVER_ADDRESS", ":8000")
	DownloaderExecutable      = envValue("DOWNLOADER_EXECUTABLE", "yt-dlp")
	DownloaderOutputDirectory = envValue("DOWNLOADER_OUTPUT_DIR", "/data/output")
	DownloaderCacheDirectory  = envValue("DOWNLOADER_CACHE_DIR", "/data/cache")
	DownloaderRetries         = envValueInt("DOWNLOADER_RETRIES", 3)
	DownloaderConcurrency     = envValueInt("DOWNLOADER_CONCURRENCY", 1)
	YoutubeAPIKey             = envValue("YOUTUBE_API_KEY")
)

func init() {
	flag.BoolVar(&Debug, "debug", Debug, "enable debug mode")
	flag.StringVar(&ServerAddress, "server-address", ServerAddress, "enable debug mode")
	flag.StringVar(&DownloaderExecutable, "downloader-executable", DownloaderExecutable, "downloader executable")
	flag.StringVar(&DownloaderOutputDirectory, "downloader-output-directory", DownloaderOutputDirectory, "downloader output directory")
	flag.StringVar(&DownloaderCacheDirectory, "downloader-cache-directory", DownloaderCacheDirectory, "downloader cache directory")
	flag.IntVar(&DownloaderRetries, "downloader-retries", DownloaderRetries, "downloader retries")
	flag.IntVar(&DownloaderConcurrency, "downloader-concurrency", DownloaderConcurrency, "downloader concurrency")
	flag.StringVar(&YoutubeAPIKey, "youtube-api-key", YoutubeAPIKey, "youtube api key")
	flag.Parse()

	if Debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
}
