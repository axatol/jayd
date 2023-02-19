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

	log.Debug().
		Str("video_id", videoID).
		Str("format", format).
		Int("queue_size", len(Jobs.items)).
		Msg("downloading")

	err := execYoutubeDownloader(videoID, format)
	if err != nil {
		log.Error().Err(err).
			Str("video_id", videoID).
			Str("format", format).
			Int("queue_size", len(Jobs.items)).
			Msg("download failed")
		return err
	}

	log.Debug().Err(err).
		Str("video_id", videoID).
		Str("format", format).
		Int("queue_size", len(Jobs.items)-1).
		Msg("download complete")
	return nil
}
