package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

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

	// Wait group
	wg = &sync.WaitGroup{}

	// Channel
	quit chan struct{}
)

/* Start the streaming */
func (s *StreamingController) StartStreaming(c *gin.Context) {
	// Request body validation
	var requestBody request.StartStreaming
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
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
	postRulesResp, _ := twitterController.PostRules([]string{requestBody.Rule})
	fmt.Println("Post the received rule")

	// Start streaming
	wg.Add(1)
	quit = make(chan struct{})
	go twitterController.GetStream()
	fmt.Println("Start streaming (goroutine)")

	c.JSON(http.StatusOK, gin.H{"success": true, "rule_id": postRulesResp.Data[0].Id})
}

/* Stop the streaming */
func (s *StreamingController) StopStreaming(c *gin.Context) {
	// cancel()
	close(quit)
	wg.Wait()
	c.JSON(http.StatusOK, gin.H{"success": true})
}
