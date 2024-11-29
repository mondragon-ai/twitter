package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/twitter/data/request"
	"github.com/twitter/data/response"
)

// Function to call OpenAI API for chat completion
func OpenAIChatCompletion(prompt string) (string, error) {

	messages := []request.OpenAIMessage{
		{Role: "system", Content: "You are a helpful assistant."},
		{Role: "user", Content: "Hello!"},
	}

	// Prepare the request body
	requestBody, err := json.Marshal(request.OpenAIRequest{
		Model:    "gpt-4o",
		Messages: messages,
	})
	if err != nil {
		return "", err
	}

	// Retrieve the API key from the environment variable
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiKey)

	// Send the request using an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the response body
	var response response.OpenAIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	// Check if there are any choices and return the first message content
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no choices found in the response")
}