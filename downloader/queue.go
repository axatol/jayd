package downloader

import (
	"encoding/json"
	"os"
	"path"
	"sync"

	"github.com/axatol/jayd/config"
	"github.com/rs/zerolog/log"
)

type QueueItem struct {
	IsCompleted bool     `json:"is_completed"`
	IsFailed    bool     `json:"is_failed"`
	FormatID    string   `json:"selected_format_id"`
	Data        InfoJSON `json:"data"`
}

type Queue struct {
	mutex sync.RWMutex
	Items []QueueItem `json:"items"`
}

func (q *Queue) Add(info InfoJSON, formatID string) {
	if len(q.Get(info.ID, formatID)) > 0 {
		return
	}

	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.Items = append(q.Items, QueueItem{
		IsCompleted: false,
		IsFailed:    false,
		FormatID:    formatID,
		Data:        info,
	})
}

func (q *Queue) SetFailed(videoID string, formatID string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for i, info := range q.Items {
		if info.Data.ID == videoID && info.FormatID == formatID {
			q.Items[i].IsFailed = true
			return
		}
	}
}

func (q *Queue) SetCompleted(videoID string, formatID string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for i, info := range q.Items {
		if info.Data.ID == videoID && info.FormatID == formatID {
			q.Items[i].IsCompleted = true
			return
		}
	}
}

func (q *Queue) Remove(videoID string, formatID string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for i, item := range q.Items {
		if item.Data.ID == videoID && item.FormatID == formatID {
			q.Items = append(q.Items[:i], q.Items[i+1:]...)
			filename := itemFilename(item)
			filepath := path.Join(config.DownloaderOutputDirectory, filename)
			if err := os.Remove(filepath); err != nil {
				log.Error().
					Err(err).
					Str("filepath", filepath).
					Msg("failed to delete file")
			}
			return
		}
	}
}

func (q *Queue) Entries() []QueueItem {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	entries := make([]QueueItem, len(q.Items))
	copy(entries, q.Items)
	return entries
}

func (q *Queue) Get(videoID string, formatID string) []QueueItem {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	results := []QueueItem{}
	for _, item := range q.Items {
		if item.Data.ID == videoID && (formatID == "" || item.FormatID == formatID) {
			results = append(results, item)
		}
	}

	return results
}
func (q *Queue) Save(name string) error {
	raw, err := json.Marshal(q)
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

func LoadHistory(name string) error {
	raw, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw, &History); err != nil {
		return err
	}

	return nil
}
