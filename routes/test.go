package routes

import (
	"github.com/gin-gonic/gin"
)

func testRouter(r *gin.RouterGroup) {
	test := r.Group("/test")
	{
		test.POST("/", testController.StartTesting)
	}
}
