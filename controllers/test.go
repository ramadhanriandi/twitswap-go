package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"

	"twitswap-go/request"
	"twitswap-go/response"
)

type TestController struct{}

var (
	errWrongNumberOfMessages = errors.New("the number of messages should be in defined range")

	bufferedChannelMaxNumber   = 100
	maxMessageNumberForTesting = int64(1000000)
	sampleMessage              = "{\"data\":{\"entities\":{\"annotations\":[{\"normalized_text\":\"Mini Football\"}],\"hashtags\":[{\"tag\":\"MiniFootball\"}]},\"geo\":{},\"id\":\"21412312412312\",\"lang\":\"en\",\"public_metrics\":{\"retweet_count\":0,\"reply_count\":0,\"like_count\":0,\"quote_count\":0},\"referenced_tweets\":[{\"type\":\"replied_to\"}],\"source\":\"Twitter for iPhone\",\"text\":\"@TTLAMBO1 @CryptoCatBTC Mini Football is the best community! Right now is the time to buy!!! @MiniFootballBsc #MiniFootball\"},\"matching_rules\":[{\"id\":1435295589592297473}]}"
)

/* Generate messages for testing */
func (t *TestController) StartTesting(c *gin.Context) {
	var err error

	// Request body validation
	var requestBody request.StartTesting
	err = c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errInvalidRequest.Error(),
			"error":   err.Error(),
		})
		return
	}

	if requestBody.Count < 1 || requestBody.Count > maxMessageNumberForTesting {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errInvalidRequest.Error(),
			"error":   err.Error(),
		})
		return
	}

	// Kafka producer setup
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	// Start the testing
	go func() {
		ticker := time.NewTicker(1 * time.Second)

		for range ticker.C {
			count := 0

			for i := int64(0); i < requestBody.Count; i++ {
				count++

				// Send raw tweet to Kafka topic
				_, err = conn.WriteMessages(
					kafka.Message{Value: []byte(sampleMessage)},
				)

				if err != nil {
					log.Fatal("Failed to write message:", err)
				}
			}

			fmt.Println(fmt.Sprintf("%d messages generated", count))
		}
	}()

	fmt.Println(fmt.Sprintf("Start testing with %d messages per second", requestBody.Count))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": response.StartTesting{
			Count: requestBody.Count,
		},
	})
}
