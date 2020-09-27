package server

import (
	"github.com/gin-gonic/gin"
	"github.com/ldegaetano/go-ddd-example/handlers/prices"
)

//Start new server in 8080 port
func Start() {
	router := gin.Default()

	pricesHandler := prices.StartHandler()
	pricesBase := router.Group(pricesHandler.BasePath)
	{
		pricesBase.GET(pricesHandler.PricesPath, pricesHandler.GetPricesFor)
		pricesBase.POST(pricesHandler.PricesPath, pricesHandler.SetPricesFor)
	}

	router.Run()
}
