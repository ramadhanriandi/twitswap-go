package routes

import (
	"twitswap-go/controllers"

	"github.com/gin-gonic/gin"
)

var (
	streamingController = new(controllers.StreamingController)
)

func SetupRouter(router *gin.Engine, routerGroup *gin.RouterGroup) {
	streamingRouter(routerGroup)
}
