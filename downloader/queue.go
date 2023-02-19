package downloader

import (
	"sync"
)

type Queue struct {
	mutex sync.RWMutex
	items []string
}

func (q *Queue) Enqueue(id string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.items = append(q.items, id)
}

func (q *Queue) Dequeue(id string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for i, cachedID := range q.items {
		if cachedID == id {
			q.items = append(q.items[:i], q.items[i+1:]...)
			return
		}
	}
}

func (q *Queue) Entries() []string {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	entries := make([]string, len(q.items))
	copy(entries, q.items)
	return entries
}
