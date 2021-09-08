package routes

import (
	"twitswap-go/controllers"

	"github.com/gin-gonic/gin"
)

var (
	ruleController      = new(controllers.RuleController)
	streamingController = new(controllers.StreamingController)
	testController      = new(controllers.TestController)
)

func SetupRouter(router *gin.Engine, routerGroup *gin.RouterGroup) {
	ruleRouter(routerGroup)
	streamingRouter(routerGroup)
	testRouter(routerGroup)
}
