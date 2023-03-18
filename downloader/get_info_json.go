package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/axatol/jayd/config"
	"github.com/rs/zerolog/log"
)

type Format struct {
	Filesize int     `json:"filesize"`
	FormatID string  `json:"format_id"`
	Width    int     `json:"width"`
	Height   int     `json:"height"`
	FPS      float64 `json:"fps"`
	AudioExt string  `json:"audio_ext"`
	VideoExt string  `json:"video_ext"`
}

type InfoJSON struct {
	VideoID        string   `json:"id"`
	Title          string   `json:"title"`
	Formats        []Format `json:"formats"`
	Thumbnail      string   `json:"thumbnail"`
	Description    string   `json:"description"`
	Uploader       string   `json:"uploader"`
	UploaderID     string   `json:"uploader_id"`
	Duration       int      `json:"duration"`
	DurationString string   `json:"duration_string"`
	FormatID       string   `json:"format_id"`
}

var (
	infoCache = map[string]InfoJSON{}
)

func GetInfoJSON(ctx context.Context, videoID string) (*InfoJSON, error) {
	if value, ok := infoCache[videoID]; ok {
		log.Debug().Str("video_id", videoID).Msg("retrieved info json from cache")
		return &value, nil
	}

	target := fmt.Sprintf("https://youtube.com/watch?v=%s", videoID)
	cmd := exec.CommandContext(ctx, config.DownloaderExecutable, target, "--dump-json", "--skip-download")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Str("stderr", string(output)).Msg("youtube downloader execution failed")
		return nil, err
	}

	var info InfoJSON
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, err
	}

	var filtered []Format
	for _, format := range info.Formats {
		if format.Filesize > 0 || (format.AudioExt != "none" || format.VideoExt != "none") {
			filtered = append(filtered, format)
		}
	}
	info.Formats = filtered

	log.Debug().Str("video_id", videoID).Msg("fetched info json")
	infoCache[videoID] = info
	return &info, nil
}
