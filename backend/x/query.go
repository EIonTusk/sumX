// Package query provides functions to interact with the Twitter API.
package x

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Tweet struct {
	Text string `json:"text"`
}

type tweetsResponse struct {
	Data []Tweet `json:"data"`
}

// FetchTweetsByUsernameTimeframe fetches tweets between start and end time
func FetchTweetsByUsernameTimeframe(userID string, from string, to string, limit int, bearerToken string) ([]Tweet, string, error) {
	url := fmt.Sprintf(
		"https://api.twitter.com/2/users/%s/tweets?tweet.fields=text&exclude=retweets,replies",
		userID,
	)

	if from != "" {
		_, err := time.Parse(time.RFC3339, from)
		if err != nil {
			return nil, "", fmt.Errorf("invalid from date format: %v", err)
		}

		url += fmt.Sprintf("&start_time=%s", from)
	}

	if to != "" {
		_, err := time.Parse(time.RFC3339, to)
		if err != nil {
			return nil, "", fmt.Errorf("invalid to date format: %v", err)
		}

		url += fmt.Sprintf("&end_time=%s", to)
	}

	if limit != -1 {
		url += fmt.Sprintf("&max_results=%d", limit)
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("User-Agent", "PostmanRuntime/7.43.0")
	req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	nextReset := resp.Header.Get("x-rate-limit-reset")

	if resp.StatusCode == 429 {
		return nil, nextReset, errors.New("request limit reached for X")
	}
	if resp.StatusCode != 200 {
		return nil, "", errors.New("failed to fetch tweets")
	}

	var res tweetsResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, "", err
	}

	return res.Data, nextReset, nil
}
func FetchUserID(username string, bearerToken string) (string, error) {
	url := fmt.Sprintf("https://api.twitter.com/2/users/by/username/%s", username)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("User-Agent", "PostmanRuntime/7.43.0")
	req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 429 {
		return "", errors.New("request limit reached for X")
	}
	if resp.StatusCode != 200 {
		return "", errors.New("failed to fetch userID")
	}

	var data struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	return data.Data.ID, nil
}
