package httpcache

import (
	"errors"
	"fmt"
	"time"

	"github.com/luraproject/lura/v2/config"
	"github.com/mitchellh/mapstructure"
)

const (
	BackendMemory = "memory"
	BackendRedis  = "redis"

	RedisModeRedis                 = "redis"
	RedisModeCluster               = "rediscluster"
	RedisDefaultDialTimeout        = 100 * time.Millisecond
	RedisDefaultReadTimeout        = 100 * time.Millisecond
	RedisDefaultWriteTimeout       = 200 * time.Millisecond
	RedisDefaultMaxRetries         = 0
	RedisDefaultIdleTimeout        = 5 * time.Minute
	RedisDefaultIdleCheckFrequency = 1 * time.Minute
	RedisDefaultPoolSize           = 10
	RedisDefaultPoolTimeout        = 10 * time.Millisecond
	RedisDefaultTtl                = 1 * time.Hour
)

var ErrMissingRequired = func(field string) error {
	return errors.New(fmt.Sprintf("Missing required krakend-httpcache config field [%s]", field))
}
var ErrInvalidValue = func(field string) error {
	return errors.New(fmt.Sprintf("Invalid value for krakend-httpcache config field [%s]", field))
}

type RedisConfig struct {
	Mode               string        `mapstructure:"mode"`
	Address            string        `mapstructure:"address"`
	DialTimeout        time.Duration `mapstructure:"dialTimeout"`
	ReadTimeout        time.Duration `mapstructure:"readTimeout"`
	WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
	MaxRetries         int           `mapstructure:"maxRetries"`
	IdleTimeout        time.Duration `mapstructure:"idleTimeout"`
	IdleCheckFrequency time.Duration `mapstructure:"idleCheckFrequency"`
	PoolSize           int           `mapstructure:"poolSize"`
	PoolTimeout        time.Duration `mapstructure:"poolTimeout"`
	Ttl                time.Duration `mapstructure:"ttl"`
}

type Config struct {
	Type        string
	Shared      bool
	RedisConfig RedisConfig `mapstructure:"redis"`
}

func ConfigGetter(cfg *config.Backend) (*Config, error) {
	raw, ok := cfg.ExtraConfig[Namespace]
	if !ok {
		return nil, nil
	}

	var cacheConfig = Config{
		RedisConfig: RedisConfig{
			Mode:               "",
			Address:            "",
			DialTimeout:        RedisDefaultDialTimeout,
			ReadTimeout:        RedisDefaultReadTimeout,
			WriteTimeout:       RedisDefaultWriteTimeout,
			MaxRetries:         RedisDefaultMaxRetries,
			IdleTimeout:        RedisDefaultIdleTimeout,
			IdleCheckFrequency: RedisDefaultIdleCheckFrequency,
			PoolSize:           RedisDefaultPoolSize,
			PoolTimeout:        RedisDefaultPoolTimeout,
			Ttl:                RedisDefaultTtl,
		},
	}
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		WeaklyTypedInput: false,
		Metadata:         nil,
		Result:           &cacheConfig,
	})
	err = d.Decode(raw)
	if err != nil {
		return nil, err
	}

	if cacheConfig.Type == "" {
		cacheConfig.Type = BackendMemory
	}

	switch cacheConfig.Type {
	case BackendMemory:
		return &cacheConfig, nil
	case BackendRedis:
		rc := cacheConfig.RedisConfig
		if rc.Mode != RedisModeRedis && rc.Mode != RedisModeCluster {
			return nil, ErrInvalidValue("redis.mode")
		}
		if rc.Address == "" {
			return nil, ErrMissingRequired("redis.address")
		}
	}

	return &cacheConfig, nil
}
