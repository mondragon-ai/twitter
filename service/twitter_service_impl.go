package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/twitter/auth"
	"github.com/twitter/data/request"
	"github.com/twitter/data/response"
	"github.com/twitter/mentions"
)


type TwitterServiceImpl struct {
	MentionRepository mentions.MentionRepository
}

func NewTwitterServiceImpl(MentionRepository mentions.MentionRepository) TwitterService {
	return &TwitterServiceImpl{
		MentionRepository: MentionRepository,
	}
}

func (t *TwitterServiceImpl) MakeTwitterRequest(ctx context.Context, method, url string, body interface{}) (*http.Response, error) {
	// Retrieve environment variables
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessTokenKey := os.Getenv("ACCESS_TOKEN_KEY")
	accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")

	// Generate OAuth signature
	signature := auth.PrepareOAuthSignature(consumerKey, accessTokenKey, consumerSecret, accessTokenSecret)

	// Prepare request body if applicable
	var requestBody *bytes.Buffer
	if method != http.MethodGet && body != nil {
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		requestBody = bytes.NewBuffer(bodyJSON)
	} else {
		requestBody = nil
	}

	// Create HTTP request
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", signature)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	// Return response
	return resp, nil
}

func (t *TwitterServiceImpl) PostTweet(ctx context.Context, request request.TweetCreateRequest) (*http.Response, error) {
	

	text := ""
	switch request.Type {
	case "create":

		messages := generateCreativePrompt()
		completion, err := OpenAIChatCompletion(messages)
		if err != nil {
			return nil, fmt.Errorf("could not compelte chat: %w", err)
		}

		filteredHash := filterWordsWithHash(completion);
		text = strings.Replace(filteredHash, "\"", "", -1)
	case "clone":
		messages := generateClonePrompt(t.MentionRepository, ctx)
		completion, err := OpenAIChatCompletion(messages)
		if err != nil {
			return nil, fmt.Errorf("could not compelte chat: %w", err)
		}

		filteredHash := filterWordsWithHash(completion);
		text = strings.Replace(filteredHash, "\"", "", -1)
	case "article":
		messages := generateArticlePrompt(t.MentionRepository, ctx)
		completion, err := OpenAIChatCompletion(messages)
		if err != nil {
			return nil, fmt.Errorf("could not compelte chat: %w", err)
		}

		filteredHash := filterWordsWithHash(completion);
		text = strings.Replace(filteredHash, "\"", "", -1)
	default:
		return nil, fmt.Errorf("default to post tweet: %v", "TEST")
	}
	// t.MentionRepository.


	log.Print(request)
    log.Printf("text: %s", text)
	// // Prepare the body for the POST request
	// body := map[string]string{
	// 	"text": text,
	// }

	// // Make the Twitter request
	// url := "https://api.twitter.com/2/tweets"
	// resp, err := t.MakeTwitterRequest(ctx, http.MethodPost, url, body)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to post tweet: %w", err)
	// }

	return nil, fmt.Errorf("failed to post tweet: %v", "TEST")
}

func (t *TwitterServiceImpl) FetchMentions(ctx context.Context) (*[]response.MentionResponse, error) {

	return nil, nil
}

func (t *TwitterServiceImpl) ReplyMention(ctx context.Context, mentionId string) {
}

func (t *TwitterServiceImpl) ReplyDM(ctx context.Context) {
}