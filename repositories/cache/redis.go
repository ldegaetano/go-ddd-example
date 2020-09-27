package cache

import (
	"fmt"
	"time"

	"github.com/ldegaetano/go-ddd-example/settings"
	"github.com/go-redis/redis"
)

type cacheRepository struct {
	client         *redis.Client
	defaultTimeout time.Duration
}

func New(defaultTime time.Duration) cacheRepository {
	url := fmt.Sprintf("%s:%s", settings.Redis.Host, settings.Redis.Port)
	options := &redis.Options{
		Addr:     url,
		Password: "",
	}
	rc := redis.NewClient(options)

	return cacheRepository{rc, defaultTime}
}
