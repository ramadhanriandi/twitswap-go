package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nanmu42/gzip"

	"twitswap-go/routes"
	"twitswap-go/services"
)

func main() {
	fmt.Println("twitswap-go (Raw Tweet Kafka Topic Producer)")

	gin.SetMode(os.Getenv("GIN_MODE"))

	router := gin.New()
	router.Use(gzip.DefaultHandler().Gin)
	router.Use(services.GetCORSMiddleware())

	api := router.Group("/api")

	routes.SetupRouter(router, api)

	router.Run()
}
