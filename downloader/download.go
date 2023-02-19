package downloader

import (
	"github.com/rs/zerolog/log"
)

var (
	Jobs Queue
)

func Download(videoID string, format string) error {
	Jobs.Enqueue(videoID)
	defer Jobs.Dequeue(videoID)

	log.Debug().Str("video_id", videoID).
		Str("format", format).
		Int("queue_size", len(Jobs.items)).
		Msg("downloading")

	return execYoutubeDownloader(videoID, format)
}
