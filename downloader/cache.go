package downloader

import (
	"encoding/json"
	"os"
	"path"

	ds "github.com/axatol/go-utils/datastructures"
	"github.com/axatol/jayd/config"
	"github.com/rs/zerolog/log"
)

var (
	Cache       *ds.AsyncMap[InfoJSON]
	CacheEvents *ds.Fanout[ds.AsyncMapEvent[InfoJSON]]
)

func createCache(initial ...map[string]ds.AsyncMapItem[InfoJSON]) {
	Cache = ds.NewAsyncMap(initial...)
	CacheEvents = ds.NewFanout(Cache.Subscribe())

	self := make(chan ds.AsyncMapEvent[InfoJSON], 1)
	CacheEvents.Subscribe("self", self)

	go func() {
		for event := range self {
			log.Debug().
				Str("action", string(event.Action)).
				Str("format_id", event.Item.Data.FormatID).
				Str("video_id", event.Item.Data.VideoID).
				Msg("cache event")

			if event.Action != ds.RemovedEventAction && event.Action != ds.FailedEventAction {
				continue
			}

			filepath := path.Join(config.DownloaderOutputDirectory, event.Item.Data.Filename)
			if err := os.Remove(filepath); err != nil && os.IsNotExist(err) {
				log.Error().
					Err(err).
					Str("filepath", filepath).
					Str("video_id", event.Item.Data.VideoID).
					Str("format_id", event.Item.Data.FormatID).
					Msg("failed to delete file")
			}
		}

		log.Debug().Msg("goroutine loop end")
	}()
}

func CreateCache(filename string) error {
	raw, err := os.ReadFile(filename)
	if err != nil {
		log.Error().
			Err(err).
			Str("filename", filename).
			Msg("error reading cache file")
		createCache()
		return nil
	}

	var items []ds.AsyncMapItem[InfoJSON]
	if err := json.Unmarshal(raw, &items); err != nil {
		return err
	}

	mapping := make(map[string]ds.AsyncMapItem[InfoJSON], len(items))
	for _, item := range items {
		mapping[item.ID] = item
	}

	createCache(mapping)
	return nil
}

func SaveCache(name string) error {
	raw, err := json.Marshal(Cache.Entries())
	if err != nil {
		return err
	}

	if err := os.MkdirAll(path.Dir(name), 0777); err != nil {
		return err
	}

	if err := os.WriteFile(name, raw, 0777); err != nil {
		return err
	}

	return nil
}
