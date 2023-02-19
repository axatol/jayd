package downloader

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/axatol/jayd/config"
	"github.com/rs/zerolog/log"
)

var (
	ytdlpExecArguments = []string{
		// general
		"--abort-on-error",
		"--no-mark-watched",

		// video selection
		"--no-playlist",

		// download options
		"--retries",
		strconv.Itoa(config.DownloaderRetries),

		// filesystem options
		"--paths",
		config.DownloaderOutputDirectory,
		"--output",
		"%(id)s.%(ext)s",
		"--no-overwrites",
		"--continue",
		"--cache-dir",
		config.DownloaderCacheDirectory,

		// verbosity and simulation options
		"--no-simulate",
		// "--dump-json",
		"--no-progress",
		// "--write-info-json",

		// workarounds
		"--no-check-certificates",

		// post-processing options
		"--extract-audio",
		"--audio-format",
		"opus",
		"--audio-quality",
		"0",
		"--no-keep-video",
	}

	ytdlpExecAudioArguments = []string{
		"--extract-audio",
		"--audio-format",
		"opus",
		"--audio-quality",
		"0",
		"--no-keep-video",
	}

	ytdlpExecVideoArguments = []string{}
)

func execYoutubeDownloader(videoID string, format string) error {
	target := fmt.Sprintf("https://youtube.com/watch?v=%s", videoID)

	args := []string{target}
	switch format {
	case FormatAudio:
		args = append(args, ytdlpExecArguments...)
		args = append(args, ytdlpExecAudioArguments...)
	case FormatVideo:
		args = append(args, ytdlpExecArguments...)
		args = append(args, ytdlpExecVideoArguments...)
	default:
		return fmt.Errorf("invalid format: %s", format)
	}

	log.Debug().Str("target", target).Msg("executing youtube downloader")
	cmd := exec.Command(config.DownloaderExecutable, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Str("stderr", string(output)).Msg("youtube downloader execution failed")
		return err
	}

	return nil
}
