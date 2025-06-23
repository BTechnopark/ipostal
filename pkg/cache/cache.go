package cache

import (
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
	"time"
)

func NewCache(d time.Duration) *cacheImpl {
	return &cacheImpl{
		items:    make(map[string]*CacheItem),
		duration: d,
	}
}

type Cache interface {
	Set(key string, value any) error
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

func (c *cacheImpl) Set(key string, value any) error {
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
		Expiration: now.Add(c.duration).Unix(),
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
