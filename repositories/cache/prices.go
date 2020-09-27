package cache

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/gommon/log"

	"github.com/ldegaetano/go-ddd-example/settings"
)

func (cr cacheRepository) GetPricesFor(itemsCode []string) (map[string]float64, error) {
	itemsPrice := map[string]float64{}

	cmd := cr.client.MGet(buildPricesKeys(itemsCode)...)
	if err := cmd.Err(); err != nil {
		log.Errorf("[process:get_redis][err:%s]", err.Error())
		return itemsPrice, errors.New("Redis get error")
	}

	errorList := []string{}
	for k, v := range cmd.Val() {
		if v == nil {
			errorList = append(errorList, fmt.Sprintf("Item %s do not exist", itemsCode[k]))
			continue
		}
		price, err := strconv.ParseFloat(v.(string), 64)
		if err != nil {
			errorList = append(errorList, fmt.Sprintf("Invalid value for %s", itemsCode[k]))
			continue
		}
		itemsPrice[itemsCode[k]] = price
	}

	if len(errorList) > 0 {
		return itemsPrice, errors.New(strings.Join(errorList, ","))
	}
	return itemsPrice, nil
}

func (cr cacheRepository) SetPricesFor(itemsPrice map[string]float64) error {
	for k, v := range itemsPrice {
		cmd := cr.client.Set(buildPriceKey(k), v, cr.defaultTimeout)
		if err := cmd.Err(); err != nil {
			log.Errorf("[process:set_redis][err:%s]", err.Error())
			return errors.New("Set cache error")
		}
	}
	return nil
}

func buildPricesKeys(itemsCode []string) (keys []string) {
	for _, i := range itemsCode {
		keys = append(keys, buildPriceKey(i))
	}
	return
}

func buildPriceKey(itemsCode string) string {
	return fmt.Sprintf(settings.Redis.PriceKey, itemsCode)
}
