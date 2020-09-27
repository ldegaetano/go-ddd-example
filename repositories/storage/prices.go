package storage

import (
	"errors"

	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
)

const (
	priceQuery  = "SELECT item_code, item_price FROM items WHERE item_code = ANY ($1);"
	insertQuery = "INSERT INTO items (item_code, item_price) VALUES ($1, $2::decimal) ON CONFLICT (item_code) DO UPDATE SET item_price = EXCLUDED.item_price;"
)

func (sr storageRepository) GetPricesFor(itemsCode []string) (map[string]float64, error) {
	res := map[string]float64{}

	rows, err := sr.db.Query(priceQuery, pq.Array(itemsCode))
	if err != nil {
		log.Errorf("[price_query_err:%s]", err.Error())
		return res, errors.New("Price query error")
	}
	defer rows.Close()

	for rows.Next() {
		var itemCode string
		var itemPrice float64
		rows.Scan(&itemCode, &itemPrice)
		res[itemCode] = itemPrice
	}
	return res, nil
}

func (sr storageRepository) SetPriceFor(itemCode string, price float64) error {
	_, err := sr.db.Exec(insertQuery, itemCode, price)
	if err != nil {
		log.Errorf("[price_insert_err:%s]", err.Error())
		return errors.New("Price insert error")
	}
	return nil
}
