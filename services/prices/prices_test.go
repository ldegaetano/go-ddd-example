package prices

import (
	"errors"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockResult has the float64 and err to return
type mockResult struct {
	price      float64
	err        error
	expiration time.Time
}

type mockStorage struct {
	numCalls    int
	mockResults map[string]mockResult // what price and err to return for a particular itemCode
	callDelay   time.Duration         // how long to sleep on each call so that we can simulate calls to be expensive
}

func (m *mockStorage) GetPricesFor(itemsCode []string) (map[string]float64, error) {

	m.numCalls++            // increase the number of calls
	time.Sleep(m.callDelay) // sleep to simulate expensive call

	result := map[string]float64{}
	var resultErr error
	for _, i := range itemsCode {
		p, ok := m.mockResults[i]
		if !ok {
			return result, nil
		}
		result[i] = p.price
		if p.err != nil {
			resultErr = p.err
		}
	}
	return result, resultErr
}

func (m *mockStorage) getNumCalls() int {
	return m.numCalls
}

func (m *mockStorage) SetPriceFor(itemCode string, price float64) error {

	m.numCalls++ // increase the number of calls
	if m.mockResults[itemCode].err != nil {
		return m.mockResults[itemCode].err
	}
	return nil
}

type mockCache struct {
	numCalls int
	maxAge   time.Duration
	prices   map[string]mockResult // what price and err to return for a particular itemCode
}

func (m *mockCache) GetPricesFor(itemsCode []string) (map[string]float64, error) {

	m.numCalls++ // increase the number of calls

	result := map[string]float64{}
	var resultErr error
	for _, i := range itemsCode {
		p, ok := m.prices[i]
		if ok && p.expiration.After(time.Now()) {
			result[i] = p.price
			resultErr = p.err
		} else {
			resultErr = errors.New("not found in cache")
		}
	}
	return result, resultErr
}

func (m *mockCache) SetPricesFor(prices map[string]float64) error {

	m.numCalls++ // increase the number of calls
	if m.prices == nil {
		m.prices = make(map[string]mockResult)
	}
	for k, p := range prices {
		m.prices[k] = mockResult{p, nil, time.Now().Add(m.maxAge)}
	}
	return nil
}

func (m *mockCache) getNumCalls() int {
	return m.numCalls
}

func getPriceWithNoErr(t *testing.T, service Service, itemCode string) float64 {
	prices, err := service.GetPricesFor(itemCode)
	if err != nil {
		t.Error("error getting prices for", itemCode)
	}
	return prices[itemCode]
}

func getPricesWithNoErr(t *testing.T, service Service, itemCodes ...string) []float64 {
	prices, err := service.GetPricesFor(itemCodes...)
	if err != nil {
		t.Error("error getting prices for", itemCodes)
	}
	result := []float64{}
	for _, p := range prices {
		result = append(result, p)
	}
	return result
}

func assertInt(t *testing.T, expected int, actual int, msg string) {
	if expected != actual {
		t.Error(msg, fmt.Sprintf("expected : %v, got : %v", expected, actual))
	}
}

func assertFloat(t *testing.T, expected float64, actual float64, msg string) {
	if expected != actual {
		t.Error(msg, fmt.Sprintf("expected : %v, got : %v", expected, actual))
	}
}

func assertErr(t *testing.T, err error) {
	if err == nil {
		t.Error("expected error to be not nil")
	}
}

func assertFloats(t *testing.T, expected []float64, actual []float64, msg string) {
	if len(expected) != len(actual) {
		t.Error(msg, fmt.Sprintf("expected : %v, got : %v", expected, actual))
		return
	}
	sort.Float64s(expected)
	sort.Float64s(actual)
	for i, expectedValue := range expected {
		if expectedValue != actual[i] {
			t.Error(msg, fmt.Sprintf("expected : %v, got : %v", expected, actual))
			return
		}
	}
}

// Check that we are caching results (we should not call the external service for all calls)
func TestGetPriceFor_CachesResults(t *testing.T) {
	mockStorage := &mockStorage{
		mockResults: map[string]mockResult{
			"p1": {price: 5, err: nil},
		},
	}
	mockCache := &mockCache{
		maxAge: time.Millisecond * 200,
	}
	service := NewService(mockStorage, mockCache)

	assertFloat(t, 5, getPriceWithNoErr(t, service, "p1"), "wrong price returned")
	assertFloat(t, 5, getPriceWithNoErr(t, service, "p1"), "wrong price returned")
	assertFloat(t, 5, getPriceWithNoErr(t, service, "p1"), "wrong price returned")
	assertInt(t, 1, mockStorage.getNumCalls(), "wrong number of service calls")
}

// Check that cache returns an error if external service returns an error
func TestGetPriceFor_ReturnsErrorOnServiceError(t *testing.T) {
	mockStorage := &mockStorage{
		mockResults: map[string]mockResult{
			"p1": {price: 0, err: fmt.Errorf("some error")},
		},
	}
	mockCache := &mockCache{}
	service := NewService(mockStorage, mockCache)
	_, err := service.GetPricesFor("p1")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

// Check that cache can return more than one price at once, caching appropriately
func TestGetPricesFor_GetsSeveralPricesAtOnceAndCachesThem(t *testing.T) {
	mockStorage := &mockStorage{
		mockResults: map[string]mockResult{
			"p1": {price: 5, err: nil},
			"p2": {price: 7, err: nil},
		},
	}
	mockCache := &mockCache{
		maxAge: time.Millisecond * 200,
	}
	service := NewService(mockStorage, mockCache)

	assertFloat(t, 5, getPriceWithNoErr(t, service, "p1"), "wrong price returned")
	assertFloats(t, []float64{5, 7}, getPricesWithNoErr(t, service, "p1", "p2"), "wrong price returned")
	assertFloats(t, []float64{5, 7}, getPricesWithNoErr(t, service, "p1", "p2"), "wrong price returned")
	assertInt(t, 2, mockStorage.getNumCalls(), "wrong number of service calls")
}

// Check that we are expiring results when they exceed the max age
func TestGetPriceFor_DoesNotReturnOldResults(t *testing.T) {
	mockStorage := &mockStorage{
		mockResults: map[string]mockResult{
			"p1": {price: 5, err: nil},
			"p2": {price: 7, err: nil},
		},
	}
	mockCache := &mockCache{
		maxAge: time.Millisecond * 200,
	}
	maxAge70Pct := time.Millisecond * 140
	service := NewService(mockStorage, mockCache)

	// get price for "p1" twice (one external service call)
	assertFloat(t, 5, getPriceWithNoErr(t, service, "p1"), "wrong price returned")
	assertFloat(t, 5, getPriceWithNoErr(t, service, "p1"), "wrong price returned")
	assertInt(t, 1, mockStorage.getNumCalls(), "wrong number of service calls")
	// sleep 0.7 the maxAge
	time.Sleep(maxAge70Pct)
	// get price for "p1" and "p2", only "p2" should be retrieved from the external service (one more external call)
	assertFloat(t, 5, getPriceWithNoErr(t, service, "p1"), "wrong price returned")
	assertFloat(t, 5, getPriceWithNoErr(t, service, "p1"), "wrong price returned")
	assertFloat(t, 7, getPriceWithNoErr(t, service, "p2"), "wrong price returned")
	assertFloat(t, 7, getPriceWithNoErr(t, service, "p2"), "wrong price returned")
	assertInt(t, 2, mockStorage.getNumCalls(), "wrong number of service calls")
	// sleep 0.7 the maxAge
	time.Sleep(maxAge70Pct)
	// get price for "p1" and "p2", only "p1" should be retrieved from the cache ("p2" is still valid)
	assertFloat(t, 5, getPriceWithNoErr(t, service, "p1"), "wrong price returned")
	assertFloat(t, 5, getPriceWithNoErr(t, service, "p1"), "wrong price returned")
	assertFloat(t, 7, getPriceWithNoErr(t, service, "p2"), "wrong price returned")
	assertInt(t, 3, mockStorage.getNumCalls(), "wrong number of service calls")
}

// Check that cache parallelize service calls when getting several values at once
func TestGetPricesFor_ParallelizeCalls(t *testing.T) {
	mockStorage := &mockStorage{
		callDelay: time.Second, // each call to external service takes one full second
		mockResults: map[string]mockResult{
			"p1": {price: 5, err: nil},
			"p2": {price: 7, err: nil},
		},
	}
	mockCache := &mockCache{}
	cache := NewService(mockStorage, mockCache)

	start := time.Now()
	assertFloats(t, []float64{5, 7}, getPricesWithNoErr(t, cache, "p1", "p2"), "wrong price returned")
	elapsedTime := time.Since(start)
	if elapsedTime > (1200 * time.Millisecond) {
		t.Error("calls took too long, expected them to take a bit over one second")
	}
}

// If an item returns an error, the entire operation must return an error
func TestGetPricesFor_ParallelizeCallsWithError(t *testing.T) {
	mockService := &mockStorage{
		callDelay: time.Second,
		mockResults: map[string]mockResult{
			"p1": {price: 5, err: errors.New("Test error")},
			"p2": {price: 7, err: nil},
		},
	}
	mockCache := &mockCache{}
	cache := NewService(mockService, mockCache)
	start := time.Now()
	_, err := cache.GetPricesFor("p1", "p2")
	assertErr(t, err)
	elapsedTime := time.Since(start)
	if elapsedTime > (1200 * time.Millisecond) {
		t.Error("calls took too long, expected them to take a bit over one second")
	}
}

func TestGetPricesFor_NotFound(t *testing.T) {
	mockService := &mockStorage{}
	mockCache := &mockCache{}
	service := NewService(mockService, mockCache)
	_, err := service.GetPricesFor("p1", "p2")
	assert.Equal(t, "Items not found: p1,p2.", err.Message)
}

func TestSetPricesFor_InsertPrice(t *testing.T) {
	mockService := &mockStorage{
		mockResults: map[string]mockResult{
			"p2": {price: 7, err: nil},
		},
	}
	mockCache := &mockCache{}
	service := NewService(mockService, mockCache)
	err := service.SetPriceFor("p1", 10)
	assert.Nil(t, err)
}

func TestSetPricesFor_InsertPriceErr(t *testing.T) {
	mockService := &mockStorage{
		mockResults: map[string]mockResult{
			"p1": {price: 10, err: errors.New("Insert err")},
		},
	}
	mockCache := &mockCache{}
	service := NewService(mockService, mockCache)
	err := service.SetPriceFor("p1", 10)
	assert.Equal(t, "Internal server error.", err.Error())
}
