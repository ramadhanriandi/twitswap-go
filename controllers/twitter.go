package controllers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/segmentio/kafka-go"

	"twitswap-go/request"
	"twitswap-go/response"
)

type TwitterController struct{}

var (
	// Twitter API v2 URLs
	streamURL string = "https://api.twitter.com/2/tweets/search/stream"
	rulesURL  string = "https://api.twitter.com/2/tweets/search/stream/rules"

	// Kafka configurations
	topic     = "raw-tweet-topic"
	partition = 0
)

/* Get stream for tweets */
func (t *TwitterController) GetStream() {
	defer wg.Done()

	// Kafka producer setup
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	// HTTP request setup
	authorizationBearer := "Bearer " + os.Getenv("TWITTER_AUTH_BEARER")
	client := &http.Client{}

	req, _ := http.NewRequest("GET", streamURL, nil)
	req.Header.Add("Authorization", authorizationBearer)

	resp, _ := client.Do(req)
	reader := bufio.NewReader(resp.Body)

	for {
		select {
		// When receive quit channel, stop streaming
		case <-quit:
			// Close HTTP connection
			defer resp.Body.Close()

			// Close Kafka connection
			if err := conn.Close(); err != nil {
				log.Fatal("failed to close writer:", err)
			}

			fmt.Print("Stop streaming...")
			return

		default:
			// Read the response
			tweet, _ := reader.ReadBytes('\n')
			line := string(tweet)

			fmt.Print(line)

			// Send raw tweet to Kafka topic
			_, err = conn.WriteMessages(
				kafka.Message{Value: []byte(line)},
			)

			if err != nil {
				log.Fatal("Failed to write messages:", err)
			}
		}
	}
}

/* Get all rules in streaming */
func (t *TwitterController) GetRules() (response.GetRules, error) {
	authorizationBearer := "Bearer " + os.Getenv("TWITTER_AUTH_BEARER")
	client := &http.Client{}

	req, _ := http.NewRequest("GET", rulesURL, nil)
	req.Header.Add("Authorization", authorizationBearer)

	resp, err := client.Do(req)
	if err != nil {
		return response.GetRules{}, err
	}

	defer resp.Body.Close()
	jsonResp, _ := ioutil.ReadAll(resp.Body)

	var getRulesResp response.GetRules
	json.Unmarshal([]byte(jsonResp), &getRulesResp)

	return getRulesResp, nil
}

/* Post rules in streaming */
func (t *TwitterController) PostRules(rules []string) (response.PostRules, error) {
	authorizationBearer := "Bearer " + os.Getenv("TWITTER_AUTH_BEARER")
	client := &http.Client{}

	var postRequestBody request.PostRules

	for _, rule := range rules {
		postRequestBody.Add = append(postRequestBody.Add, request.PostRulesValue{
			Value: rule,
		})
	}

	requestBody, _ := json.Marshal(postRequestBody)

	req, _ := http.NewRequest("POST", rulesURL, strings.NewReader(string(requestBody)))
	req.Header.Add("Authorization", authorizationBearer)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return response.PostRules{}, err
	}

	defer resp.Body.Close()
	jsonResp, _ := ioutil.ReadAll(resp.Body)

	var postRulesResp response.PostRules
	json.Unmarshal([]byte(jsonResp), &postRulesResp)

	return postRulesResp, nil
}

/* Delete rules in streaming */
func (t *TwitterController) DeleteRules(ruleIds []string) (response.DeleteRules, error) {
	authorizationBearer := "Bearer " + os.Getenv("TWITTER_AUTH_BEARER")
	client := &http.Client{}

	var deleteRequestBody request.DeleteRules

	deleteRequestBody.Delete.Ids = make([]string, len(ruleIds))
	copy(deleteRequestBody.Delete.Ids, ruleIds)

	requestBody, _ := json.Marshal(deleteRequestBody)

	req, _ := http.NewRequest("POST", rulesURL, strings.NewReader(string(requestBody)))
	req.Header.Add("Authorization", authorizationBearer)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return response.DeleteRules{}, err
	}

	defer resp.Body.Close()
	jsonResp, _ := ioutil.ReadAll(resp.Body)

	var deleteRulesResp response.DeleteRules
	json.Unmarshal([]byte(jsonResp), &deleteRulesResp)

	return deleteRulesResp, nil
}
