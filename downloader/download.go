package downloader

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/axatol/jayd/config"
	"github.com/rs/zerolog/log"
)

var (
	jobs Queue
)

func Jobs() []InfoJSON {
	return jobs.Entries()
}

func HasJob(id string) bool {
	return jobs.Has(id)
}

func Download(info InfoJSON, formatID string) error {
	jobs.Add(info)
	defer jobs.Remove(info.ID)

	log.Debug().
		Str("video_id", info.ID).
		Str("format", formatID).
		Msg("downloading")

	err := execYoutubeDownloader(info.ID, formatID)
	if err != nil {
		return err
	}

	log.Debug().Err(err).
		Str("video_id", info.ID).
		Str("format", formatID).
		Msg("download complete")
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
		"--output",
		"%(id)s.%(ext)s",
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

	args := []string{target}
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
		log.Error().Err(err).Str("stderr", string(output)).Msg("youtube downloader execution failed")
		return err
	}

	return nil
}
