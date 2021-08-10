package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	dbConn "twitswap-go/db"
	"twitswap-go/helpers"
	"twitswap-go/response"
)

type RuleController struct{}

var (
	// Errors
	errFailedGetTweetAnnotationsDB = errors.New("failed to get tweet annotations from DB")
	errFailedGetTweetDomainsDB     = errors.New("failed to get tweet domains from DB")
	errFailedParseTime             = errors.New("failed to parse time from query parameter")

	// Limit
	rowsLimit = 20
)

/* Get the visualization data by rule ID */
func (s *RuleController) GetVisualizationByRuleID(c *gin.Context) {
	// Get rule ID from the path variable
	ruleID := c.Param("id")

	// Get latest time from query parameter
	latestTime := time.Now()

	queryMap := c.Request.URL.Query()
	latestTimeVal, ok := queryMap["latest_time"]

	if ok && len(latestTimeVal) == 1 {
		latestTime = helpers.Convert2DateTime(latestTimeVal[0])
	}

	// Get the visualization data based on rule ID from DB
	db := dbConn.OpenConnection()
	defer db.Close()

	var resp response.GetVisualizationByRuleID

	// Get tweet annotations
	tweetAnnotationRows, tweetAnnotationErr := db.Query(
		"SELECT name, SUM(count) AS total FROM tweet_annotations WHERE rule_id = $1 AND created_at <= $2 GROUP BY name ORDER BY total DESC LIMIT $3",
		ruleID,
		latestTime,
		rowsLimit,
	)
	if tweetAnnotationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errFailedGetTweetAnnotationsDB.Error(),
			"error":   tweetAnnotationErr.Error(),
		})
		return
	}
	defer tweetAnnotationRows.Close()

	for tweetAnnotationRows.Next() {
		var data response.TweetAnnotation

		err := tweetAnnotationRows.Scan(&data.Name, &data.Count)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": errFailedGetTweetAnnotationsDB.Error(),
				"error":   err.Error(),
			})
			return
		}

		resp.TweetAnnotations = append(resp.TweetAnnotations, data)
	}

	// Get tweet domains
	tweetDomainRows, tweetDomainErr := db.Query(
		"SELECT name, SUM(count) AS total FROM tweet_domains WHERE rule_id = $1 AND created_at <= $2 GROUP BY name ORDER BY total DESC LIMIT $3",
		ruleID,
		latestTime,
		rowsLimit,
	)
	if tweetDomainErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errFailedGetTweetDomainsDB.Error(),
			"error":   tweetDomainErr.Error(),
		})
		return
	}
	defer tweetDomainRows.Close()

	for tweetDomainRows.Next() {
		var data response.TweetDomain

		err := tweetDomainRows.Scan(&data.Name, &data.Count)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": errFailedGetTweetDomainsDB.Error(),
				"error":   err.Error(),
			})
			return
		}

		resp.TweetDomains = append(resp.TweetDomains, data)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}
