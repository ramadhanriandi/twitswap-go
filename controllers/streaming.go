package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"twitswap-go/request"
)

type StreamingController struct{}

var (
	// Errors
	errInvalidRequest       = errors.New("invalid request body")
	errFailedStartStreaming = errors.New("failed to start streaming")
	errFailedStopStreaming  = errors.New("failed to stop streaming")

	// Twitter controller
	twitterController = new(TwitterController)
)

/* Start the streaming */
func (s *StreamingController) StartStreaming(c *gin.Context) {
	// Request body validation
	var requestBody request.StartStreaming
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": errInvalidRequest.Error(),
			"error":   err.Error(),
		})

		return
	}

	// Get all existed rules
	getRulesResp, _ := twitterController.GetRules()
	fmt.Println("Get all existed rules")

	// Delete all existed rules if any
	if len(getRulesResp.Data) > 0 {
		var ruleIds []string
		for _, data := range getRulesResp.Data {
			ruleIds = append(ruleIds, data.Id)
		}

		twitterController.DeleteRules(ruleIds)
	}
	fmt.Println("Delete all existed rules if any")

	// Post the received rule
	twitterController.PostRules([]string{requestBody.Rule})
	fmt.Println("Post the received rule")

	// Start streaming
	go twitterController.GetStream()

	c.JSON(http.StatusOK, gin.H{"message": "Successfully start streaming"})
}

/* Stop the streaming */
func (s *StreamingController) StopStreaming(c *gin.Context) {
	var data request.StopStreaming

	// Request body validation
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": errInvalidRequest.Error(),
			"error":   err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully stop streaming"})
}
