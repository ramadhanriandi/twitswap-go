package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	dbConn "twitswap-go/db"
	"twitswap-go/model"
	"twitswap-go/request"
	"twitswap-go/response"
)

type StreamingController struct{}

var (
	// Errors
	errInvalidRequest          = errors.New("invalid request body")
	errFailedGetRules          = errors.New("failed to get rules")
	errFailedDeleteRule        = errors.New("failed to delete rule")
	errFailedPostRule          = errors.New("failed to post rule")
	errFailedInsertRuleDB      = errors.New("failed to insert rule into DB")
	errFailedInsertStreamingDB = errors.New("failed to insert streaming into DB")
	errFailedStartStreaming    = errors.New("failed to start streaming")
	errFailedStopStreaming     = errors.New("failed to stop streaming")

	// Twitter controller
	twitterController = new(TwitterController)

	// Wait group
	wg = &sync.WaitGroup{}

	// Channel
	quit chan struct{}
)

/* Start the streaming */
func (s *StreamingController) StartStreaming(c *gin.Context) {
	var err error

	// Request body validation
	var requestBody request.StartStreaming
	err = c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errInvalidRequest.Error(),
			"error":   err.Error(),
		})
		return
	}

	// Get all existed rules
	getRulesResp, err := twitterController.GetRules()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errFailedGetRules.Error(),
			"error":   err.Error(),
		})
		return
	}
	fmt.Println("Get all existed rules")

	// Delete all existed rules if any
	if len(getRulesResp.Data) > 0 {
		var ruleIds []string
		for _, data := range getRulesResp.Data {
			ruleIds = append(ruleIds, data.Id)
		}

		_, err := twitterController.DeleteRules(ruleIds)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": errFailedDeleteRule.Error(),
				"error":   err.Error(),
			})
			return
		}
	}
	fmt.Println("Delete all existed rules if any")

	// Post the received rule
	postRulesResp, err := twitterController.PostRules([]string{requestBody.Rule})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errFailedPostRule.Error(),
			"error":   err.Error(),
		})
		return
	}
	fmt.Println("Post the received rule")

	// Store rule and streaming data into DB
	db := dbConn.OpenConnection()
	defer db.Close()

	rule := model.Rule{
		ID:    postRulesResp.Data[0].Id,
		Value: requestBody.Rule,
	}
	ruleQuery := "INSERT INTO rules (id, value) VALUES ($1, $2)"
	_, err = db.Exec(ruleQuery, rule.ID, rule.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errFailedInsertRuleDB.Error(),
			"error":   err.Error(),
		})
		return
	}

	streamingID := 0
	streaming := model.Streaming{
		Name:      requestBody.Name,
		StartTime: time.Now(),
		RuleID:    rule.ID,
	}
	streamingQuery := "INSERT INTO streamings (name, start_time, rule_id) VALUES ($1, $2, $3) RETURNING id"
	err = db.QueryRow(streamingQuery, streaming.Name, streaming.StartTime, streaming.RuleID).Scan(&streamingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errFailedInsertRuleDB.Error(),
			"error":   err.Error(),
		})
		return
	}

	fmt.Println("Store rule and streaming data into DB")

	// Start streaming
	wg.Add(1)
	quit = make(chan struct{})
	go twitterController.GetStream()
	fmt.Println("Start streaming (goroutine)")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": response.StartStreaming{
			ID:        int64(streamingID),
			Name:      streaming.Name,
			StartTime: streaming.StartTime,
			RuleID:    rule.ID,
			Rule:      rule.Value,
		},
	})
}

/* Stop the streaming */
func (s *StreamingController) StopStreaming(c *gin.Context) {
	close(quit)
	wg.Wait()
	c.JSON(http.StatusOK, gin.H{"success": true})
}
