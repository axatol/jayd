package downloader

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/axatol/jayd/config"
	"github.com/rs/zerolog/log"
)

var (
	History Queue
)

func Download(info InfoJSON, formatID string) error {
	History.Add(info, formatID)
	defer History.SetCompleted(info.ID, formatID)

	log.Debug().
		Str("video_id", info.ID).
		Str("format_id", formatID).
		Msg("downloading")

	err := execYoutubeDownloader(info.ID, formatID)
	if err != nil {
		History.SetFailed(info.ID, formatID)
	}

	log.Debug().Err(err).
		Str("video_id", info.ID).
		Str("format_id", formatID).
		Msg("downloaded")
	return nil
}

const (
	FormatDefaultVideo = "defaultvideo"
	FormatDefaultAudio = "defaultaudio"
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
		// "--output",
		// "%(id)s.%(ext)s",
		"--no-overwrites",
		"--continue",
		"--no-write-info-json",
		"--cache-dir",
		config.DownloaderCacheDirectory,

		// verbosity and simulation options
		"--no-simulate",
		// "--dump-json",
		"--no-progress",
		// "--write-info-json",

		// workarounds
		"--no-check-certificates",
	}

	ytdlpExecAudioArguments = []string{
		// post-processing options
		"--extract-audio",
		"--audio-format",
		"opus",
		"--audio-quality",
		"0",
		"--no-keep-video",

		// video format options
		"--format-sort",
		"hasaud,quality,abr,asr,+size",
	}

	ytdlpExecVideoArguments = []string{
		// video format options
		"--format-sort",
		"hasvid,lang,quality,res:1080,+size",
	}
)

func execYoutubeDownloader(videoID string, formatID string) error {
	target := fmt.Sprintf("https://youtube.com/watch?v=%s", videoID)
	outputTemplate := fmt.Sprintf("%s_%s.%s", "%(id)s", formatID, "%(ext)s")
	args := []string{target, "--output", outputTemplate}
	args = append(args, ytdlpExecArguments...)
	switch formatID {
	case FormatDefaultAudio:
		args = append(args, ytdlpExecAudioArguments...)
	case FormatDefaultVideo:
		args = append(args, ytdlpExecVideoArguments...)
	default:
		args = append(args, "--format", formatID)
	}

	cmd := exec.Command(config.DownloaderExecutable, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().
			Err(err).
			Str("output", string(output)).
			Msg("youtube downloader execution failed")
		return err
	}

	return nil
}
