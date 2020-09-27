package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func clearDB(storage storageRepository) {
	storage.db.Exec(`TRUNCATE TABLE items;`)
}

func TestStorage_GetPricesFor(t *testing.T) {
	storage := New()
	defer clearDB(storage)

	storage.SetPriceFor("p1", 10)
	storage.SetPriceFor("p2", 3)
	storage.SetPriceFor("p3", 4)

	itemsPrice, err := storage.GetPricesFor([]string{"p1", "p2", "p3"})

	assert.Nil(t, err)

	assert.Equal(t, float64(10), itemsPrice["p1"])
	assert.Equal(t, float64(3), itemsPrice["p2"])
	assert.Equal(t, float64(4), itemsPrice["p3"])
}

func TestStorage_GetPricesForErr(t *testing.T) {
	storageRepo := New()
	defer clearDB(storageRepo)
	defer func() {
		storage = nil
	}()

	storage.db.Close()

	_, err := storageRepo.GetPricesFor([]string{"p1", "p2", "p3"})

	assert.Equal(t, "Price query error", err.Error())
}

func TestStorage_SetPricesForErr(t *testing.T) {
	storageRepo := New()
	defer clearDB(storageRepo)
	defer func() {
		storage = nil
	}()

	storage.db.Close()

	err := storageRepo.SetPriceFor("p1", 10)

	assert.Equal(t, "Price insert error", err.Error())
}
