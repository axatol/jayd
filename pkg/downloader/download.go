package downloader

import (
	"context"
	"fmt"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/axatol/jayd/pkg/config"
	"github.com/axatol/jayd/pkg/config/nr"
	"github.com/axatol/jayd/pkg/downloader/miniodriver"
	"github.com/rs/zerolog/log"
)

type FormatType string

const (
	AudioVideoFormatType FormatType = "audio_video"
	AudioOnlyFormatType  FormatType = "audio_only"
)

func CacheItemID(videoID string, formatID string) string {
	return fmt.Sprintf("%s#%s", videoID, formatID)
}

func Download(ctx context.Context, info InfoJSON, formatID string, overwrite bool) error {
	defer nr.Segment(ctx, "downloader.Download", nr.Attrs{
		"video_id":  info.VideoID,
		"format_id": formatID,
	}).End()

	id := CacheItemID(info.VideoID, formatID)
	info.FormatID = formatID
	info.Formats = selectItemFormats(formatID, info.Formats)
	info.Ext = selectItemExt(info.Formats)
	info.Filename = renderItemFilename(info)
	formatType := selectFormatType(info.Formats)

	Cache.Add(id, info, overwrite)
	defer Cache.SetCompleted(id)

	if info.Ext == "" {
		Cache.SetFailed(id)
		return fmt.Errorf("could not determine format extension")
	}

	log.Debug().
		Str("ext", info.Ext).
		Str("filename", info.Filename).
		Str("format_type", string(formatType)).
		Bool("overwrite", overwrite).
		Msg("downloading")

	if err := executeYTDL(info.VideoID, formatID, formatType); err != nil {
		Cache.SetFailed(id)
		return err
	}

	if config.StorageEnabled {
		client, err := miniodriver.AssertClient(ctx)
		if err != nil {
			Cache.SetFailed(id)
			return err
		}

		filepath := path.Join(config.DownloaderOutputDirectory, info.Filename)
		if err := client.FPutObject(ctx, info.Filename, filepath, miniodriver.Tags{}); err != nil {
			Cache.SetFailed(id)
			return err
		}
	}

	return nil
}

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
		// "--format-sort",
		// "hasaud,quality,abr,asr,+size",
	}

	ytdlpExecVideoArguments = []string{
		// post-processing options
		// TODO

		// video format options
		// "--format-sort",
		// "hasvid,lang,quality,res:1080,+size",
	}
)

func executeYTDL(videoID string, formatID string, formatType FormatType) error {
	target := fmt.Sprintf("https://youtube.com/watch?v=%s", videoID)
	outputTemplate := fmt.Sprintf("%s_%s.%s", "%(id)s", formatID, "%(ext)s")

	args := append(
		ytdlpExecArguments,
		target,
		"--output", outputTemplate,
		"--format", formatID,
	)

	switch formatType {
	case AudioVideoFormatType:
		args = append(args, ytdlpExecAudioArguments...)
	case AudioOnlyFormatType:
		args = append(args, ytdlpExecVideoArguments...)
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

func selectFormatType(formats []Format) FormatType {
	for _, format := range formats {
		if format.VideoExt != "none" {
			return AudioVideoFormatType
		}
	}

	return AudioOnlyFormatType
}

func selectItemFormats(id string, formats []Format) []Format {
	ids := strings.Split(id, "+")
	results := []Format{}

	for _, format := range formats {
		for _, id := range ids {
			if id == format.FormatID {
				results = append(results, format)
			}
		}
	}

	return results
}

func selectItemExt(formats []Format) string {
	for _, format := range formats {
		if format.VideoExt != "none" {
			return format.VideoExt
		}

		if format.AudioExt != "none" {
			return format.AudioExt
		}
	}

	return ""
}

func renderItemFilename(info InfoJSON) string {
	return fmt.Sprintf("%s_%s.%s", info.VideoID, info.FormatID, info.Ext)
}
