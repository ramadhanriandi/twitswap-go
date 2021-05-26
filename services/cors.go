package services

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func GetCORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Options{
		AllowedOrigins: strings.Split(os.Getenv("CORS_ALLOWED_ORIGIN"), ","),
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodOptions,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-Requested-With", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})
}
