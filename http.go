//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

package httpcache

import (
	"context"
	"net/http"

	"github.com/krakend/httpcache"
	"github.com/krakend/lru"
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/transport/http/client"
)

type Cache interface {
	// Get returns the []byte representation of a cached response and a bool
	// set to true if the value isn't empty
	Get(ctx context.Context, key string) (responseBytes []byte, ok bool)
	// Set stores the []byte representation of a response against a key
	Set(ctx context.Context, key string, responseBytes []byte)
	// Delete removes the value associated with the key
	Delete(ctx context.Context, key string)
}

// Namespace is the key to use to store and access the custom config data
const Namespace = "github.com/devopsfaith/krakend-httpcache"

// NewHTTPClient creates a HTTPClientFactory using an in-memory-cached http client
func NewHTTPClient(cfg *config.Backend, nextF client.HTTPClientFactory) client.HTTPClientFactory {
	cacheCfg, err := ConfigGetter(cfg)
	if err != nil || cacheCfg == nil {
		return nextF
	}

	switch cacheCfg.Type {
	case BackendMemory:
		return getCachedClient(nextF, selectMemoryCache(cacheCfg.MemoryConfig))
	case BackendRedis:
		var r Client
		switch cacheCfg.RedisConfig.Mode {
		case RedisModeRedis:
			r = NewRedis(cacheCfg.RedisConfig)
		case RedisModeCluster:
			r = NewRedisCluster(cacheCfg.RedisConfig)
		}
		return func(_ context.Context) *http.Client {
			return &http.Client{Transport: NewRedisCacheTransport(NewRedisCache(r, cacheCfg.RedisConfig.Ttl))}
		}
	}
	return defaultClient(nextF)
}

func defaultClient(nextF client.HTTPClientFactory) client.HTTPClientFactory {
	return getCachedClient(nextF, httpcache.NewMemoryCache())
}

func getCachedClient(nextF client.HTTPClientFactory, cache Cache) client.HTTPClientFactory {
	return func(ctx context.Context) *http.Client {
		httpClient := nextF(ctx)
		return &http.Client{
			Transport: &httpcache.Transport{
				Transport: httpClient.Transport,
				Cache:     cache,
			},
			CheckRedirect: httpClient.CheckRedirect,
			Jar:           httpClient.Jar,
			Timeout:       httpClient.Timeout,
		}
	}
}

func selectMemoryCache(opts MemoryConfig) Cache {
	if opts.MaxSize == 0 || opts.MaxItems == 0 {
		if opts.Shared {
			return globalCache
		}
		return httpcache.NewMemoryCache()
	}

	if !opts.Shared {
		cache, _ := lru.NewLruCache(opts.MaxSize, opts.MaxItems)
		return NewLruCache(cache)
	}

	if globalLruCache == nil {
		c, _ := lru.NewLruCache(opts.MaxSize, opts.MaxItems)
		globalLruCache = NewLruCache(c)
	}
	return globalLruCache
}

var (
	globalLruCache Cache
	globalCache    = httpcache.NewMemoryCache()
)
