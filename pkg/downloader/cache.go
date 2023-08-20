package downloader

import (
	"encoding/json"
	"os"
	"path"
	"sync"
	"time"

	"github.com/axatol/jayd/pkg/config"
	"github.com/axatol/jayd/pkg/data"
	"github.com/rs/zerolog/log"
)

type EventAction string

const (
	AddedEventAction     EventAction = "added"
	FailedEventAction    EventAction = "failed"
	CompletedEventAction EventAction = "completed"
	RemovedEventAction   EventAction = "removed"
)

type CacheItem struct {
	ID          string     `json:"id"`
	Info        InfoJSON   `json:"info"`
	AddedAt     *data.Time `json:"added_at"`
	FailedAt    *data.Time `json:"failed_at,omitempty"`
	CompletedAt *data.Time `json:"completed_at,omitempty"`
}

type CacheEvent struct {
	Action EventAction `json:"action"`
	Item   CacheItem   `json:"item"`
}

type CacheEventSubscriber func(CacheEvent)

type cache struct {
	mu          *sync.RWMutex
	Items       map[string]CacheItem
	subscribers map[string]CacheEventSubscriber
}

var Cache = cache{
	mu:          &sync.RWMutex{},
	subscribers: map[string]CacheEventSubscriber{},
	Items:       map[string]CacheItem{},
}

func (c cache) Entries() []CacheItem {
	items := make([]CacheItem, len(c.Items))
	i := 0
	for _, item := range c.Items {
		items[i] = item
		i += 1
	}

	return items
}

func (c cache) Subscribe(id string, subscriber CacheEventSubscriber) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.subscribers[id] = subscriber
}

func (c cache) Unsubscribe(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.subscribers, id)
}

func (c cache) emit(action EventAction, item CacheItem) {
	wg := sync.WaitGroup{}
	for _, subscriber := range c.subscribers {
		wg.Add(1)
		go func(subscriber CacheEventSubscriber) {
			subscriber(CacheEvent{Action: action, Item: item})
			wg.Done()
		}(subscriber)
	}

	wg.Wait()
}

func (c cache) Get(id string) *CacheItem {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if item, ok := c.Items[id]; ok {
		return &item
	}

	return nil
}

func (c cache) Set(id string, info InfoJSON) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := CacheItem{
		ID:      id,
		Info:    info,
		AddedAt: &data.Time{Time: time.Now()},
	}

	Cache.Items[id] = item
	c.emit(AddedEventAction, item)
}

func (c cache) Add(id string, info InfoJSON, overwrite ...bool) {
	allowOverwrite := len(overwrite) > 0 && overwrite[0]
	if c.Get(id) != nil && !allowOverwrite {
		return
	}

	c.Set(id, info)
}

func (c cache) Remove(id string) {
	if item, ok := c.Items[id]; ok {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.Items, id)
		c.emit(RemovedEventAction, item)
	}
}

func (c cache) SetFailed(id string) {
	if item, ok := c.Items[id]; ok {
		c.mu.Lock()
		defer c.mu.Unlock()
		item.FailedAt = &data.Time{Time: time.Now()}
		c.Items[id] = item
		c.emit(FailedEventAction, item)
	}
}

func (c cache) SetCompleted(id string) {
	if item, ok := c.Items[id]; ok {
		c.mu.Lock()
		defer c.mu.Unlock()
		item.CompletedAt = &data.Time{Time: time.Now()}
		c.Items[id] = item
		c.emit(CompletedEventAction, item)
	}
}

func createCache(initial ...map[string]CacheItem) {
	if len(initial) > 0 {
		Cache.Items = initial[0]
	}

	Cache.Subscribe("", func(event CacheEvent) {
		if event.Action != RemovedEventAction {
			return
		}

		filepath := path.Join(config.DownloaderOutputDirectory, event.Item.Info.Filename)
		if err := os.Remove(filepath); err != nil && os.IsNotExist(err) {
			log.Error().
				Err(err).
				Str("filepath", filepath).
				Str("video_id", event.Item.Info.VideoID).
				Str("format_id", event.Item.Info.FormatID).
				Msg("failed to delete file")
		}
	})
}

func CreateCache(filename string) error {
	items := map[string]CacheItem{}
	defer createCache(items)

	raw, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw, &items); err != nil {
		return err
	}

	for _, item := range items {
		items[item.ID] = item
	}

	return nil
}

func SaveCache(name string) error {
	raw, err := json.Marshal(Cache.Items)
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
