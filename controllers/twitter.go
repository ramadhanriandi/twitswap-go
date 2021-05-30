package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"twitswap-go/request"
	"twitswap-go/response"
)

type TwitterController struct{}

var (
	streamURL string = "https://api.twitter.com/2/tweets/search/stream"
	rulesURL  string = "https://api.twitter.com/2/tweets/search/stream/rules"
)

/* Get stream for tweets */
func (t *TwitterController) GetStream() {
	authorizationBearer := "Bearer " + os.Getenv("TWITTER_AUTH_BEARER")
	client := &http.Client{}

	req, _ := http.NewRequest("GET", streamURL, nil)
	req.Header.Add("Authorization", authorizationBearer)

	resp, _ := client.Do(req)

	reader := bufio.NewReader(resp.Body)

	for {
		line, _ := reader.ReadBytes('\n')
		fmt.Println(string(line))
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
