package downloader

import (
	"sync"
)

type QueueItem struct {
	IsCompleted bool     `json:"is_completed"`
	IsFailed    bool     `json:"is_failed"`
	FormatID    string   `json:"selected_format_id"`
	Data        InfoJSON `json:"data"`
}

type Queue struct {
	mutex sync.RWMutex
	items []QueueItem
}

func (q *Queue) Add(info InfoJSON, formatID string) {
	if len(q.Get(info.ID, formatID)) > 0 {
		return
	}

	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.items = append(q.items, QueueItem{
		IsCompleted: false,
		IsFailed:    false,
		FormatID:    formatID,
		Data:        info,
	})
}

func (q *Queue) SetFailed(videoID string, formatID string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for i, info := range q.items {
		if info.Data.ID == videoID && info.FormatID == formatID {
			q.items[i].IsFailed = true
			return
		}
	}
}

func (q *Queue) SetCompleted(videoID string, formatID string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for i, info := range q.items {
		if info.Data.ID == videoID && info.FormatID == formatID {
			q.items[i].IsCompleted = true
			return
		}
	}
}

func (q *Queue) Remove(videoID string, formatID string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for i, info := range q.items {
		if info.Data.ID == videoID && info.FormatID == formatID {
			q.items = append(q.items[:i], q.items[i+1:]...)
			return
		}
	}
}

func (q *Queue) Entries() []QueueItem {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	entries := make([]QueueItem, len(q.items))
	copy(entries, q.items)
	return entries
}

func (q *Queue) Get(videoID string, formatID string) []QueueItem {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	results := []QueueItem{}
	for _, item := range q.items {
		if item.Data.ID == videoID && (formatID == "" || item.FormatID == formatID) {
			results = append(results, item)
		}
	}

	return results
}
