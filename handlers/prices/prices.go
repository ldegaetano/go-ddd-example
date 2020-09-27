package prices

import (
	"net/http"
	"strings"

	"github.com/ldegaetano/go-ddd-example/errors"
	"github.com/ldegaetano/go-ddd-example/repositories/cache"
	"github.com/ldegaetano/go-ddd-example/repositories/storage"
	"github.com/ldegaetano/go-ddd-example/services/prices"
	"github.com/ldegaetano/go-ddd-example/settings"

	"github.com/gin-gonic/gin"
)

const (
	maxItemsLength  = 5
	maxItems        = 10
	itemsCodesParam = "items_codes"
)

type PricesHandler struct {
	BasePath      string
	PricesPath    string
	PricesService prices.Service
}

func StartHandler() PricesHandler {
	return PricesHandler{
		BasePath:   "/api/items",
		PricesPath: "/prices",
		PricesService: prices.NewService(
			storage.New(),
			cache.New(settings.Redis.DefaultExpiration),
		),
	}
}

func (i PricesHandler) GetPricesFor(c *gin.Context) {
	itemsStr := c.Query(itemsCodesParam)

	itemsCodes, validateErr := validateItems(itemsStr)
	if validateErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, validateErr)
		return
	}

	itemsPrices, err := i.PricesService.GetPricesFor(itemsCodes...)
	if err != nil {
		if err.Code == errors.NotFoundCode {
			c.AbortWithStatusJSON(http.StatusNotFound, err)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
		return
	}

	c.JSON(http.StatusOK, buildPricesResponse(itemsPrices))
}

//SetPricesFor set price to item_code, if exists update the price
func (i PricesHandler) SetPricesFor(c *gin.Context) {
	p := priceCreate{}

	if err := c.BindJSON(&p); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.InvalidFormat)
		return
	}

	if err := i.PricesService.SetPriceFor(p.ItemCode, p.ItemPrice); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalError)
		return
	}

	c.Status(http.StatusNoContent)
}

func buildPricesResponse(itemsPrices map[string]float64) (response pricesResponse) {
	for itemCode, price := range itemsPrices {
		response.Items = append(response.Items, item{
			ItemCode:  itemCode,
			ItemPrice: price,
		})
	}
	return
}

func validateItems(itemsString string) ([]string, error) {
	itemsCodes := strings.Split(itemsString, ",")
	totalItems := len(itemsCodes)
	if len(itemsString) == 0 || totalItems == 0 {
		return itemsCodes, errors.AtLeastOneItem
	}
	if totalItems > maxItems {
		return itemsCodes, errors.MaxItemsExceded
	}
	if invalids := getInvalidItems(itemsCodes); len(invalids) > 0 {
		return itemsCodes, errors.InvalidItems.WithParams(invalids)
	}
	return itemsCodes, nil
}

func getInvalidItems(items []string) (res []string) {
	for _, i := range items {
		if len(i) > maxItemsLength {
			res = append(res, i)
		}
	}
	return
}
