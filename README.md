Krakend HTTP Cache
====

A cached http client for the [KrakenD](github.com/devopsfaith/krakend) framework

## Using it

This package exposes two simple factories capable to create a instances of the `proxy.HTTPClientFactory` and the `proxy.BackendFactory` interfaces, respectively, embedding an in-memory-cached http client using the package [github.com/gregjones/httpcache](https://github.com/gregjones/httpcache). The client will cache the responses honoring the defined Cache HTTP header.

	import 	(
		"context"
		"net/http"
		"github.com/devopsfaith/krakend/config"
		"github.com/devopsfaith/krakend/proxy"
		"github.com/devopsfaith/krakend-httpcache"
	)

	requestExecutorFactory := func(cfg *config.Backend) proxy.HTTPRequestExecutor {
		clientFactory := httpcache.NewHTTPClient(cfg)
		return func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return clientFactory(ctx).Do(req.WithContext(ctx))
		}
	}

You can create your own proxy.HTTPRequestExecutor and inject it into your BackendFactory


## Stayforlong changes

Added support to redis & redis-cluster as cache backends.

Backend extra configuration block example with default values:

```
"backend": [
{
    "extra_config": {
        "github.com/devopsfaith/krakend-httpcache": {
            "type": "redis",
            "redis": {
                "mode": "redis",
                "address": "localhost:6379",
                "dialTimeout": "100ms",
                "readTimeout": "100ms",
                "writeTimeout": "200ms",
                "maxRetries": 0,
                "idleTimeout": "5m",
                "idleCheckFrequency": "1m",
                "poolSize": 10,
                "poolTimeout": "10ms",
                "ttl": "1h",
            }
        },
    }
}
```

Config options:
- **type** (optional, if missing `memory` will be used for an inmemory cache): `memory` | `redis`

**redis** block is required if type is set to `redis`:
- **mode** (required): `redis` or `rediscluster`
- **address** (required): Address to the redis node or one of the redis nodes when using a cluster. Format: `ip:port`
- **ttl**: TTL to be used on all redis keys being set. Redis will expire the keys after this time. This value is independent 
from the Cache-Control header that the library will use to decide if the cached value is still valid or not.

Check the description of the following options in the go-redis package documentation: https://github.com/go-redis/redis
- **dialTimeout**: 
- **readTimeout**: 
- **writeTimeout**: 
- **maxRetries**: 
- **idleTimeout**: 
- **idleCheckFrequency**: 
- **poolSize**: 
- **poolTimeout**: 
