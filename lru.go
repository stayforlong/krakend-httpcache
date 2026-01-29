package httpcache

import (
	"context"

	"github.com/krakend/lru"
)

type LruCache struct {
	*lru.LruCache
}

func NewLruCache(c *lru.LruCache) *LruCache {
	return &LruCache{LruCache: c}
}

func (c *LruCache) Get(_ context.Context, key string) (responseBytes []byte, ok bool) {
	return c.LruCache.Get(key)
}

func (c *LruCache) Set(_ context.Context, key string, responseBytes []byte) {
	c.LruCache.Set(key, responseBytes)
}

func (c *LruCache) Delete(_ context.Context, key string) {
	c.LruCache.Delete(key)
}
