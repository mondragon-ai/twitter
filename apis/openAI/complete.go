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
	Top      float64  `json:"top_p"`
	Model            string   `json:"model"`
	FrequencyPenalty int      `json:"frequency_penalty"`
	PresencePenalty int      `json:"presence_penalty"`
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
		Model:            "gpt-3.5-turbo-instruct",
		FrequencyPenalty: 0,
		Prompt:           prompt,
		PresencePenalty:  0,
		MaxTokens:        70,
		Temperature:      1,
	})

	openaiKey := os.Getenv("OPENAI_API_KEY")
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	var response OpenAIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println(err)
		return "", err
	}


	if len(response.Choices) > 0 {
		fmt.Println(response.Choices[0].Text)
		return response.Choices[0].Text, nil
	}

	return "", nil
}