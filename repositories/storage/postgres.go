package storage

import (
	"database/sql"
	"fmt"

	"github.com/labstack/gommon/log"

	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres driver
	"github.com/ldegaetano/go-ddd-example/settings"
)

const initQuery = `CREATE TABLE IF NOT EXISTS items (
	item_code  VARCHAR NOT NULL,
	item_price NUMERIC(10,2) NOT NULL,
			
	CONSTRAINT items_pk PRIMARY KEY (item_code)
);`

type storageRepository struct {
	db *sql.DB
}

var storage *storageRepository

func New() storageRepository {
	if storage != nil {
		return *storage
	}

	c := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		settings.Postgres.UserName,
		settings.Postgres.Password,
		settings.Postgres.Host,
		settings.Postgres.Port,
		settings.Postgres.DBName,
		"disable",
	)
	db, err := sql.Open("postgres", c)
	if err != nil {
		log.Errorf("[build_db_err:%s]", err.Error())
		return *storage
	}

	db.Exec(initQuery)

	storage = &storageRepository{db}
	return *storage
}
