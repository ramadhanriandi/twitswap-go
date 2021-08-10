package routes

import (
	"github.com/gin-gonic/gin"
)

func ruleRouter(r *gin.RouterGroup) {
	rule := r.Group("/rule")
	{
		rule.GET("/:id/visualization", ruleController.GetVisualizationByRuleID)
	}
}
