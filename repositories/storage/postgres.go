package storage

import (
	"database/sql"
	"fmt"

	"github.com/ldegaetano/go-ddd-example/utils"

	"github.com/labstack/gommon/log"

	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres driver
	"github.com/ldegaetano/go-ddd-example/settings"
)

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

	initeQuery, err := utils.ReadQuery("migrations/", "init")
	if err != nil {
		log.Errorf("[init_db_err:%s]", err.Error())
	}
	db.Exec(initeQuery)

	storage = &storageRepository{db}
	return *storage
}
