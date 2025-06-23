package cache

import (
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
	"time"
)

func NewCache(d time.Duration) *cacheImpl {
	cache := cacheImpl{
		items:    make(map[string]*CacheItem),
		duration: d,
	}

	go cache.cleanUp()

	return &cache
}

type Cache interface {
	Set(key string, value any, exp time.Duration) error
	Get(key string, resp any) error
	Delete(key string)
}

var ErrCacheNotFound = errors.New("cache not found")

type CacheItem struct {
	Key        string
	Value      []byte
	Expiration int64
}

type cacheImpl struct {
	items    map[string]*CacheItem
	mu       sync.RWMutex
	duration time.Duration
}

func (c *cacheImpl) Set(key string, value any, exp time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	cache, found := c.items[key]
	if found && now.Unix() < cache.Expiration {
		slog.Info("Cache Expired", slog.String("key", key))
		return nil
	}

	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	cacheItem := CacheItem{
		Key:        key,
		Value:      raw,
		Expiration: now.Add(exp).Unix(),
	}

	c.items[key] = &cacheItem

	slog.Info("Set Cache", slog.String("key", key))
	return nil
}

func (c *cacheImpl) Get(key string, resp any) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found || time.Now().Unix() > item.Expiration {
		return ErrCacheNotFound
	}

	err := json.Unmarshal(item.Value, resp)
	if err != nil {
		return err
	}

	slog.Info("Cache Found", slog.String("key", key))
	return nil
}

func (c *cacheImpl) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *cacheImpl) cleanUp() {
	t := time.NewTicker(time.Minute)

	for range t.C {
		slog.Info("Clean Up Cache", slog.Int("count", len(c.items)))

		for key, item := range c.items {
			if time.Now().Unix() > item.Expiration {
				slog.Info("Delete Cache", slog.String("key", key))
				c.Delete(key)
			}
		}
	}
}
