//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE -self_package=$GOPACKAGE

package httpcache

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/krakendio/httpcache"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis.v8"
)

type Client interface {
	redis.Cmdable
}

func NewRedis(cfg RedisConfig) Client {
	c := redis.NewClient(&redis.Options{
		Addr:               cfg.Address,
		DialTimeout:        cfg.DialTimeout,
		ReadTimeout:        cfg.ReadTimeout,
		WriteTimeout:       cfg.WriteTimeout,
		MaxRetries:         cfg.MaxRetries,
		IdleTimeout:        cfg.IdleTimeout,
		IdleCheckFrequency: cfg.IdleCheckFrequency,
		PoolSize:           cfg.PoolSize,
		PoolTimeout:        cfg.PoolTimeout,
	})
	redistrace.WrapClient(c, redistrace.WithServiceName(serviceNameFromAddresses([]string{cfg.Address})))
	return c
}

func NewRedisCluster(cfg RedisConfig) Client {
	c := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:              []string{cfg.Address},
		DialTimeout:        cfg.DialTimeout,
		ReadTimeout:        cfg.ReadTimeout,
		WriteTimeout:       cfg.WriteTimeout,
		MaxRetries:         cfg.MaxRetries,
		IdleTimeout:        cfg.IdleTimeout,
		IdleCheckFrequency: cfg.IdleCheckFrequency,
		PoolSize:           cfg.PoolSize,
		PoolTimeout:        cfg.PoolTimeout,
	})
	redistrace.WrapClient(c, redistrace.WithServiceName(serviceNameFromAddresses([]string{cfg.Address})))
	return c
}

type RedisCache struct {
	client Client
	ttl    time.Duration
}

func NewRedisCache(client Client, ttl time.Duration) *RedisCache {
	return &RedisCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *RedisCache) Get(ctx context.Context, key string) (responseBytes []byte, ok bool) {
	r := c.client.Get(ctx, key)
	rb, err := r.Bytes()
	if err != nil {
		return []byte{}, false
	}
	return rb, true
}

func (c *RedisCache) Set(ctx context.Context, key string, responseBytes []byte) {
	c.client.Set(ctx, key, responseBytes, c.ttl)
}

func (c *RedisCache) Delete(ctx context.Context, key string) {
	c.client.Del(ctx, key)
}

func NewRedisCacheTransport(c Cache) *httpcache.Transport {
	t := httpcache.NewTransport(c)
	return t
}

func serviceNameFromAddresses(addr []string) string {
	prefix := "redis-"
	delimiter := "_"
	return prefix + strings.Join(addr, delimiter)
}
