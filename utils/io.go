package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/labstack/gommon/log"
)

func ReadQuery(path string, query string) (string, error) {
	wd, _ := os.Getwd()
	queryPath := filepath.Join(wd, path, query+".sql")
	file, err := ioutil.ReadFile(queryPath)

	if err != nil {
		log.Info("Reading file ", query, err)
	}

	if err != nil {
		return "", err
	}

	return string(file), nil
}
