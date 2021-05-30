package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nanmu42/gzip"

	"twitswap-go/routes"
	"twitswap-go/services"
)

func main() {
	fmt.Println("twitswap-go (Raw Tweet Kafka Topic Producer)")

	// Set environment variables
	godotenv.Load()

	// Set gin mode
	gin.SetMode(os.Getenv("GIN_MODE"))

	// Routers
	router := gin.New()
	router.Use(gzip.DefaultHandler().Gin)
	router.Use(services.GetCORSMiddleware())

	api := router.Group("/api")

	routes.SetupRouter(router, api)

	// Run the application
	router.Run()
}
