// Package llm provides a client for interacting with the Ollama LLM API.
package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type HFClient struct {
	BaseURL string
	Model   string
}

type message struct {
	Role   string `json:"role"`
	Content string `json:"content"`
}


type generateRequest struct {
	Messages []message `json:"messages"`
	Model string `json:"model"`
	Stream bool  `json:"stream"`
}

type choice struct {
	Message struct{
		Role string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}

type generateResponse struct {
	Choices []choice `json:"choices"`
}

func NewHFClient(model string) *HFClient {
	return &HFClient{
		BaseURL: "https://router.huggingface.co/v1/chat/completions",
		Model:   model,
	}
}

// Summarize sends text to Ollama and returns a summary
func (c *HFClient) Summarize(input string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %v", err)
	}

	task := `### Instruction:
You are given multiple tweets, each tweet is seperated by '\n---\n'. Summarize them into a single tweet that captures the overall information. Give each block of information a headline.

### Output format:
	[
		{
			"heading": "Headline1",
			"text": "Summary of topic 1"
		},
		{
			"heading": "Headline2",
			"text": "Summary of topic 2"
		},
		...
	]

### Input:


	`

	bearerToken := os.Getenv("HF_BEARER_TOKEN")
	payload := generateRequest{
		Model: c.Model,
		Stream: false,
		Messages: []message{
			{	
				Role:    "user",
				Content: task + input,
			},
		},
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal error: %v", err)
	}

	req, _ := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var res generateResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return "", fmt.Errorf("unmarshal error of %s: %v", body, err)
	}

	return res.Choices[0].Message.Content, nil
}
