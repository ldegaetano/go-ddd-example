package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func ChangeWorkDir(path string) {
	old, errs := os.Getwd()
	if errs != nil {
		log.Fatalf("err: %s", errs)
	}
	newDir := filepath.Join(old, path)
	if errs := os.Chdir(newDir); errs != nil {
		log.Fatalf("err: %s", errs)
	}
}

func ServeTestRequest(method string, endPoint string, payload io.Reader, handler gin.HandlerFunc, queryParams string) *httptest.ResponseRecorder {
	router := gin.Default()
	switch method {
	case "GET":
		router.GET(endPoint, handler)
	case "POST":
		router.POST(endPoint, handler)
	default:
		return nil
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, fmt.Sprintf("%s?%s", endPoint, queryParams), payload)
	router.ServeHTTP(w, req)

	return w
}
