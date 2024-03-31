package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// OpenAI completion request structure
type OpenAIRequest struct {
	Prompt           string   `json:"prompt"`
	MaxTokens        int      `json:"max_tokens"`
	Temperature      float64  `json:"temperature"`
	Model            string   `json:"model"`
	FrequencyPenalty int      `json:"frequency_penalty"`
	Stop             []string `json:"stop"` 
}

// OpenAI completion response structure
type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens    int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Function to call OpenAI API for chat completion
func OpenAIChatCompletion(prompt string) (string, error) {
	requestBody, _ := json.Marshal(OpenAIRequest{
		Model:            "davinci-002",
		FrequencyPenalty: 1,
		Prompt:           prompt,
		MaxTokens:        50,
		Temperature:      0.2,
	})

	openaiKey := os.Getenv("OPENAI_API_KEY")
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response OpenAIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}
	fmt.Println(req)
	fmt.Println(response.Choices[0].Text)

	if len(response.Choices) > 0 {
		return response.Choices[0].Text, nil
	}

	return "", nil
}