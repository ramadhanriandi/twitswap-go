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
	errFailedGetTweetAnnotationsDB  = errors.New("failed to get tweet annotations from DB")
	errFailedGetTweetDomainsDB      = errors.New("failed to get tweet domains from DB")
	errFailedGetTweetGeolocationsDB = errors.New("failed to get tweet geolocations from DB")
	errFailedGetTweetHashtagsDB     = errors.New("failed to get tweet hashtags from DB")
	errFailedGetTweetLanguagesDB    = errors.New("failed to get tweet languages from DB")
	errFailedGetTweetMetricsDB      = errors.New("failed to get tweet metrics from DB")
	errFailedGetTweetPopularitiesDB = errors.New("failed to get tweet popularities from DB")
	errFailedGetTweetTypesDB        = errors.New("failed to get tweet types from DB")
	errFailedParseTime              = errors.New("failed to parse time from query parameter")

	// Limit
	commonRowsLimit     = 20
	popularityRowsLimit = 3
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
		"SELECT name, COALESCE(SUM(count), 0)  AS total FROM tweet_annotations WHERE rule_id = $1 AND created_at <= $2 GROUP BY name ORDER BY total DESC LIMIT $3",
		ruleID,
		latestTime,
		commonRowsLimit,
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
		"SELECT name, COALESCE(SUM(count), 0) AS total FROM tweet_domains WHERE rule_id = $1 AND created_at <= $2 GROUP BY name ORDER BY total DESC LIMIT $3",
		ruleID,
		latestTime,
		commonRowsLimit,
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

	// Get tweet geolocations
	tweetGeolocationRows, tweetGeolocationErr := db.Query(
		"SELECT DISTINCT lat, long FROM tweet_geolocations WHERE rule_id = $1 AND created_at <= $2",
		ruleID,
		latestTime,
	)
	if tweetGeolocationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errFailedGetTweetGeolocationsDB.Error(),
			"error":   tweetGeolocationErr.Error(),
		})
		return
	}
	defer tweetGeolocationRows.Close()

	for tweetGeolocationRows.Next() {
		var data response.TweetGeolocation

		err := tweetGeolocationRows.Scan(&data.Lat, &data.Long)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": errFailedGetTweetGeolocationsDB.Error(),
				"error":   err.Error(),
			})
			return
		}

		resp.TweetGeolocations = append(resp.TweetGeolocations, data)
	}

	// Get tweet hashtags
	tweetHashtagRows, tweetHashtagErr := db.Query(
		"SELECT name, COALESCE(SUM(count), 0) AS total FROM tweet_hashtags WHERE rule_id = $1 AND created_at <= $2 GROUP BY name ORDER BY total DESC LIMIT $3",
		ruleID,
		latestTime,
		commonRowsLimit,
	)
	if tweetHashtagErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errFailedGetTweetHashtagsDB.Error(),
			"error":   tweetHashtagErr.Error(),
		})
		return
	}
	defer tweetHashtagRows.Close()

	for tweetHashtagRows.Next() {
		var data response.TweetHashtag

		err := tweetHashtagRows.Scan(&data.Name, &data.Count)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": errFailedGetTweetHashtagsDB.Error(),
				"error":   err.Error(),
			})
			return
		}

		resp.TweetHashtags = append(resp.TweetHashtags, data)
	}

	// Get tweet languages
	tweetLanguageQuery := "SELECT COALESCE(SUM(en_count), 0), COALESCE(SUM(in_count), 0),COALESCE(SUM(other_count), 0) FROM tweet_languages WHERE rule_id = $1 AND created_at <= $2"
	tweetLanguageErr := db.QueryRow(tweetLanguageQuery, ruleID, latestTime).
		Scan(&resp.TweetLanguages.En, &resp.TweetLanguages.In, &resp.TweetLanguages.Other)
	if tweetLanguageErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errFailedGetTweetLanguagesDB.Error(),
			"error":   tweetLanguageErr.Error(),
		})
		return
	}

	// Get tweet metrics
	tweetMetricCumulativeQuery := "SELECT COALESCE(SUM(like_count), 0), COALESCE(SUM(reply_count), 0), COALESCE(SUM(retweet_count), 0), COALESCE(SUM(quote_count), 0) FROM tweet_metrics WHERE rule_id = $1 AND created_at <= $2"
	tweetMetricCumulativeErr := db.QueryRow(tweetMetricCumulativeQuery, ruleID, latestTime).
		Scan(
			&resp.TweetMetrics.Cumulative.Like,
			&resp.TweetMetrics.Cumulative.Reply,
			&resp.TweetMetrics.Cumulative.Retweet,
			&resp.TweetMetrics.Cumulative.Quote,
		)
	if tweetMetricCumulativeErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errFailedGetTweetMetricsDB.Error(),
			"error":   tweetMetricCumulativeErr.Error(),
		})
		return
	}

	tweetMetricIntervalRows, tweetMetricIntervalErr := db.Query(
		"SELECT like_count, reply_count, retweet_count, quote_count, created_at FROM tweet_metrics WHERE rule_id = $1 AND created_at <= $2 ORDER BY created_at DESC LIMIT $3",
		ruleID,
		latestTime,
		commonRowsLimit,
	)
	if tweetMetricIntervalErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errFailedGetTweetMetricsDB.Error(),
			"error":   tweetMetricIntervalErr.Error(),
		})
		return
	}
	defer tweetMetricIntervalRows.Close()

	for tweetMetricIntervalRows.Next() {
		var data response.TweetMetricInterval

		err := tweetMetricIntervalRows.Scan(&data.Like, &data.Reply, &data.Retweet, &data.Quote, &data.CreatedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": errFailedGetTweetMetricsDB.Error(),
				"error":   err.Error(),
			})
			return
		}

		resp.TweetMetrics.Interval = append(resp.TweetMetrics.Interval, data)
	}

	// Get tweet popularities
	tweetPopularityRows, tweetPopularityErr := db.Query(
		"SELECT tweet_id, COALESCE(SUM(popularity), 0) AS total FROM tweet_popularities WHERE rule_id = $1 AND created_at <= $2 GROUP BY tweet_id ORDER BY total DESC LIMIT $3",
		ruleID,
		latestTime,
		popularityRowsLimit,
	)
	if tweetPopularityErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errFailedGetTweetPopularitiesDB.Error(),
			"error":   tweetPopularityErr.Error(),
		})
		return
	}
	defer tweetPopularityRows.Close()

	for tweetPopularityRows.Next() {
		var data response.TweetPopularity

		err := tweetPopularityRows.Scan(&data.TweetID, &data.Count)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": errFailedGetTweetPopularitiesDB.Error(),
				"error":   err.Error(),
			})
			return
		}

		resp.TweetPopularities = append(resp.TweetPopularities, data)
	}

	// Get tweet types
	tweetTypeCumulativeQuery := "SELECT COALESCE(SUM(tweet_count), 0), COALESCE(SUM(retweet_count), 0), COALESCE(SUM(quote_count), 0), COALESCE(SUM(reply_count), 0) FROM tweet_types WHERE rule_id = $1 AND created_at <= $2"
	tweetTypeCumulativeErr := db.QueryRow(tweetTypeCumulativeQuery, ruleID, latestTime).
		Scan(
			&resp.TweetTypes.Cumulative.Tweet,
			&resp.TweetTypes.Cumulative.Retweet,
			&resp.TweetTypes.Cumulative.Quote,
			&resp.TweetTypes.Cumulative.Reply,
		)
	if tweetTypeCumulativeErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errFailedGetTweetTypesDB.Error(),
			"error":   tweetTypeCumulativeErr.Error(),
		})
		return
	}

	tweetTypeIntervalRows, tweetTypeIntervalErr := db.Query(
		"SELECT tweet_count, retweet_count, quote_count, reply_count, created_at FROM tweet_types WHERE rule_id = $1 AND created_at <= $2 ORDER BY created_at DESC LIMIT $3",
		ruleID,
		latestTime,
		commonRowsLimit,
	)
	if tweetTypeIntervalErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errFailedGetTweetTypesDB.Error(),
			"error":   tweetTypeIntervalErr.Error(),
		})
		return
	}
	defer tweetTypeIntervalRows.Close()

	for tweetTypeIntervalRows.Next() {
		var data response.TweetTypeInterval

		err := tweetTypeIntervalRows.Scan(&data.Tweet, &data.Retweet, &data.Quote, &data.Reply, &data.CreatedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": errFailedGetTweetTypesDB.Error(),
				"error":   err.Error(),
			})
			return
		}

		resp.TweetTypes.Interval = append(resp.TweetTypes.Interval, data)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}
