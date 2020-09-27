package settings

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type redisSettings struct {
	Host              string `envconfig:"REDIS_HOST" required:"true"`
	Port              string `envconfig:"REDIS_PORT" required:"true"`
	PriceKey          string
	DefaultExpiration time.Duration
}

var Redis redisSettings

func init() {
	if err := envconfig.Process("", &Redis); err != nil {
		panic(err.Error())
	}

	Redis.PriceKey = "price:%s"
	Redis.DefaultExpiration = 1 * time.Minute
}
