package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/ldegaetano/go-ddd-example/settings"

	"github.com/stretchr/testify/assert"
)

func TestPriceFor_RedisNil(t *testing.T) {
	cache := New(time.Second)
	_, err := cache.GetPricesFor([]string{"c1"})

	assert.Contains(t, err.Error(), "Item c1 do not exist")
}

func TestPriceFor_RedisGetError(t *testing.T) {
	aux := settings.Redis.Host
	settings.Redis.Host = "invalid_host"
	cache := New(time.Second)
	_, err := cache.GetPricesFor([]string{"c1"})

	assert.Contains(t, err.Error(), "Redis get error")
	settings.Redis.Host = aux
}

func TestPriceFor_RedisSetError(t *testing.T) {
	aux := settings.Redis.Host
	settings.Redis.Host = "invalid_host"
	cache := New(time.Second)
	itemsPrices := map[string]float64{
		"c3": 1,
		"c5": 3,
	}
	err := cache.SetPricesFor(itemsPrices)

	assert.NotNil(t, err)
	settings.Redis.Host = aux
}

func TestPriceFor_InvalidFormat(t *testing.T) {
	cache := New(time.Second)
	cache.client.Set(fmt.Sprintf(settings.Redis.PriceKey, "c3"), "invalid_format", time.Second)
	_, err := cache.GetPricesFor([]string{"c3"})

	assert.Contains(t, err.Error(), "Invalid value for c3")
}

func TestPriceFor_ValueExpired(t *testing.T) {
	cache := New(time.Millisecond * 100)
	delay := time.Millisecond * 200
	itemsPrices := map[string]float64{
		"c3": 10.5,
		"c5": 3,
	}
	cache.SetPricesFor(itemsPrices)
	price, err := cache.GetPricesFor([]string{"c3"})

	assert.Nil(t, err)
	assert.Equal(t, 10.5, price["c3"])

	time.Sleep(delay)
	itemsPrices = map[string]float64{
		"c4": 9,
		"c7": 3,
	}
	cache.SetPricesFor(itemsPrices)

	prices, err := cache.GetPricesFor([]string{"c4", "c3"})
	assert.Contains(t, "Item c3 do not exist", err.Error())
	assert.Equal(t, float64(9), prices["c4"])
}
