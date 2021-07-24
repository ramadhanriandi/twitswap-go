package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
