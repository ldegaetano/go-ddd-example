package prices

import "github.com/ldegaetano/go-ddd-example/errors"

// NewService return a items service for consult prices
func NewService(storage storageRepository, cache cacheRepository) Service {
	return &service{
		storage: storage,
		cache:   cache,
	}
}

// GetPriceFor gets the price for the item, either from the cache or the actual service if it was not cached or too old
func (s *service) GetPricesFor(itemsCode ...string) (map[string]float64, *errors.CustomError) {
	storagePrices := map[string]float64{}

	cachePrices, err := s.cache.GetPricesFor(itemsCode)
	if err == nil {
		return cachePrices, nil
	}

	if missingItems := getMissingItems(itemsCode, cachePrices); len(missingItems) > 0 {
		storagePrices, err = s.storage.GetPricesFor(missingItems)
		if err != nil {
			return storagePrices, errors.InternalError
		}
		s.cache.SetPricesFor(storagePrices)
	}

	storagePrices = getItemsUnion(cachePrices, storagePrices)
	if missingItems := getMissingItems(itemsCode, storagePrices); len(missingItems) > 0 {
		return storagePrices, errors.NotFoundItems.WithParams(missingItems)
	}

	return storagePrices, nil
}

func getMissingItems(itemsCode []string, prices map[string]float64) []string {
	missingItems := []string{}
	for _, item := range itemsCode {
		if _, ok := prices[item]; !ok {
			missingItems = append(missingItems, item)
		}
	}
	return missingItems
}

func getItemsUnion(cache, storage map[string]float64) map[string]float64 {
	for k, v := range cache {
		storage[k] = v
	}
	return storage
}

func (s *service) SetPriceFor(itemCode string, price float64) *errors.CustomError {

	if err := s.storage.SetPriceFor(itemCode, price); err != nil {
		return errors.InternalError
	}

	s.cache.SetPricesFor(map[string]float64{
		itemCode: price,
	})

	return nil
}
