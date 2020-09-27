package prices

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/ldegaetano/go-ddd-example/errors"
	"github.com/ldegaetano/go-ddd-example/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type serviceMock struct {
	mock.Mock
}

func (_m *serviceMock) GetPricesFor(itemCode ...string) (map[string]float64, *errors.CustomError) {
	_va := make([]interface{}, len(itemCode))
	for _i := range itemCode {
		_va[_i] = itemCode[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 map[string]float64
	if rf, ok := ret.Get(0).(func(...string) map[string]float64); ok {
		r0 = rf(itemCode...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]float64)
		}
	}

	var r1 *errors.CustomError
	if rf, ok := ret.Get(1).(func(...string) *errors.CustomError); ok {
		r1 = rf(itemCode...)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*errors.CustomError)
		}
	}

	return r0, r1
}

func (_m *serviceMock) SetPriceFor(itemCode string, price float64) *errors.CustomError {
	ret := _m.Called(itemCode, price)

	var r0 *errors.CustomError
	if rf, ok := ret.Get(0).(func(string, float64) *errors.CustomError); ok {
		r0 = rf(itemCode, price)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*errors.CustomError)
		}
	}

	return r0
}

func TestGetPricesFor_InvalidItems(t *testing.T) {
	handler := StartHandler()
	path := handler.BasePath + handler.PricesPath
	w := utils.ServeTestRequest("GET", path, nil, handler.GetPricesFor, "items_codes=ppppppp")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid items: ppppppp.")
}

func TestGetPricesFor_AtLeastOneItem(t *testing.T) {
	handler := StartHandler()
	path := handler.BasePath + handler.PricesPath
	w := utils.ServeTestRequest("GET", path, nil, handler.GetPricesFor, "item")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "You must provide at least one item code.")
}

func TestGetPricesFor_MaxItemsExceded(t *testing.T) {
	handler := StartHandler()
	path := handler.BasePath + handler.PricesPath
	query := "items_codes=q,w,e,r,t,y,u,i,o,p,a"
	w := utils.ServeTestRequest("GET", path, nil, handler.GetPricesFor, query)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Max items quantity exceded.")
}

func TestGetPricesFor_InternalErr(t *testing.T) {
	service := serviceMock{}
	handler := StartHandler()
	handler.ItemsService = &service
	service.On("GetPricesFor", "p1").Return(map[string]float64{}, errors.InternalError)

	path := handler.BasePath + handler.PricesPath
	w := utils.ServeTestRequest("GET", path, nil, handler.GetPricesFor, "items_codes=p1")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Internal server error.")
}

func TestGetPricesFor_NotFound(t *testing.T) {
	service := serviceMock{}
	handler := StartHandler()
	handler.ItemsService = &service
	service.On("GetPricesFor", "p2").Return(map[string]float64{}, errors.NotFoundItems)

	path := handler.BasePath + handler.PricesPath
	w := utils.ServeTestRequest("GET", path, nil, handler.GetPricesFor, "items_codes=p2")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Items not found")
}

func TestGetPricesFor_ReturnPrices(t *testing.T) {
	service := serviceMock{}
	handler := StartHandler()
	handler.ItemsService = &service
	service.On("GetPricesFor", "p2").Return(map[string]float64{"p2": 10}, nil)

	path := handler.BasePath + handler.PricesPath
	w := utils.ServeTestRequest("GET", path, nil, handler.GetPricesFor, "items_codes=p2")

	response := pricesResponse{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, float64(10), response.Items[0].ItemPrice)
	assert.Equal(t, "p2", response.Items[0].ItemCode)
}

func TestPostPricesFor_InvalidFormat(t *testing.T) {
	service := serviceMock{}
	handler := StartHandler()
	handler.ItemsService = &service

	path := handler.BasePath + handler.PricesPath
	w := utils.ServeTestRequest("POST", path, strings.NewReader("{"), handler.SetPricesFor, "")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Request invalid format")
}

func TestPostPricesFor_InternalErr(t *testing.T) {
	service := serviceMock{}
	handler := StartHandler()
	handler.ItemsService = &service
	service.On("SetPriceFor", "p14", float64(15)).Return(errors.InternalError)

	body := `{"item_code": "p14","item_price": 15}`

	path := handler.BasePath + handler.PricesPath
	w := utils.ServeTestRequest("POST", path, strings.NewReader(body), handler.SetPricesFor, "")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Internal server error")
}

func TestPostPricesFor_StatusOK(t *testing.T) {
	service := serviceMock{}
	handler := StartHandler()
	handler.ItemsService = &service
	service.On("SetPriceFor", "p14", float64(15)).Return(nil)

	body := `{"item_code": "p14","item_price": 15}`

	path := handler.BasePath + handler.PricesPath
	w := utils.ServeTestRequest("POST", path, strings.NewReader(body), handler.SetPricesFor, "")

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "", w.Body.String())
}
