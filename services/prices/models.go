package prices

import "github.com/ldegaetano/go-ddd-example/errors"

type (
	// Service implements a transparent cache for returning prices
	Service interface {
		GetPricesFor(itemCode ...string) (map[string]float64, *errors.CustomError)
		SetPriceFor(itemCode string, price float64) *errors.CustomError
	}

	cacheRepository interface {
		GetPricesFor(itemsCode []string) (map[string]float64, error)
		SetPricesFor(prices map[string]float64) error
	}

	storageRepository interface {
		GetPricesFor(itemsCode []string) (map[string]float64, error)
		SetPriceFor(itemCode string, price float64) error
	}

	// Service is a service that allow interact with items
	// Implements a transparent cache for returning prices
	service struct {
		storage storageRepository
		cache   cacheRepository
	}
)
