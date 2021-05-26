package routes

import (
	"github.com/gin-gonic/gin"
)

func streamingRouter(r *gin.RouterGroup) {
	streaming := r.Group("/streaming")
	{
		streaming.POST("/start", streamingController.StartStreaming)
		streaming.POST("/stop", streamingController.StopStreaming)
	}
}
