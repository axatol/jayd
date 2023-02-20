package downloader

import (
	"sync"
)

type Queue struct {
	mutex sync.RWMutex
	items []InfoJSON
}

func (q *Queue) Add(info InfoJSON) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.items = append(q.items, info)
}

func (q *Queue) Remove(id string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	// TODO defer notify

	for i, info := range q.items {
		if info.ID == id {
			q.items = append(q.items[:i], q.items[i+1:]...)
			return
		}
	}
}

func (q *Queue) Entries() []InfoJSON {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	entries := make([]InfoJSON, len(q.items))
	copy(entries, q.items)
	return entries
}

func (q *Queue) Has(id string) bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	for _, item := range q.items {
		if item.ID == id {
			return true
		}
	}

	return false
}
