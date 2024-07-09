//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

// Package httpcache introduces an in-memory-cached http client into the KrakenD stack
package httpcache

import (
	"context"
	"net/http"

	"github.com/krakendio/httpcache"

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
	if err == nil {
		if cacheCfg == nil {
			return nextF
		}

		switch cacheCfg.Type {
		case BackendMemory:
			var cache Cache

			if cacheCfg.Shared {
				cache = globalCache
			}

			if cache == nil {
				cache = httpcache.NewMemoryCache()
			}

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
	}
	return client.NewHTTPClient
}

var globalCache = httpcache.NewMemoryCache()
