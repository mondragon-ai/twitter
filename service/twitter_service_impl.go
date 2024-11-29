package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/twitter/auth"
	"github.com/twitter/data/request"
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
	consumerKey := os.Getenv("CONSUMER_KEY")
    consumerSecret := os.Getenv("CONSUMER_SECRET")
    accessTokenKey := os.Getenv("ACCESS_TOKEN_KEY")
    accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	bearer := os.Getenv("ANGEL_BEARER")

	token := fmt.Sprintf("Bearer %s", bearer)
    if consumerKey == "" || consumerSecret == "" || accessTokenKey == "" || accessTokenSecret == "" {
        log.Fatal("One or more required environment variables are missing")
    }

    signature := auth.PrepareOAuthSignature(consumerKey, accessTokenKey, consumerSecret, accessTokenSecret)
    if signature == "" {
		log.Fatal("OAuth signature is empty; please verify your keys and signature generation logic.")
    }

    // Prepare request body if applicable
    var requestBody *bytes.Buffer
    if method != http.MethodGet && body != nil {
        bodyJSON, err := json.Marshal(body)
        if err != nil {
            log.Printf("Failed to marshal request body: %v", err)
            return nil, fmt.Errorf("failed to marshal request body: %w", err)
        }
        requestBody = bytes.NewBuffer(bodyJSON)
    }

	if requestBody == nil {
		log.Printf("Request body is nil as expected for %s requests", method)
	} else {
		log.Printf("Request body: %v", requestBody)
	}

	log.Printf("Request URL: %s", url)
	log.Printf("Authorization Header: %s", token)
	log.Printf("HTTP Method: %s", method)


    // Create HTTP request
    req, err := http.NewRequest(method, url, requestBody)
    log.Printf("Request: %v", req)
    if err != nil {
		log.Printf("Failed to create HTTP request: %v", err)
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	if req == nil {
		log.Fatal("http.NewRequest returned a nil request")
	}

    if requestBody != nil {
        req.Header.Set("Content-Type", "application/json")
    }
    req.Header.Set("Authorization", token)

    log.Printf("Making %s request to URL: %s", method, url)
    log.Printf("Request Headers: %v", req.Header)

    client := &http.Client{}
    resp, err := func() (*http.Response, error) {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Recovered from panic in client.Do: %v", r)
            }
        }()
        return client.Do(req)
    }()
    if err != nil {
		log.Printf("Error sending request: %v", err)
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	if resp == nil {
		log.Fatal("Received nil response from client.Do")
	}

    log.Printf("Received response with status: %d", resp.StatusCode)
    return resp, nil
}

func (t *TwitterServiceImpl) GetTwitterRequest(ctx context.Context, url string, body interface{}) (*http.Response, error) {
	bearer := os.Getenv("ANGEL_BEARER")
	token := fmt.Sprintf("Bearer %s", bearer)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Authorization", token)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    fmt.Println("Status:", resp.Status)

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(data))
    return resp, nil
}

func (t *TwitterServiceImpl) FetchMentions(ctx context.Context) ([]string, error) {
	// now := time.Now().UTC()
	// startTime := now.Add(-time.Duration(60*24*7) * time.Minute).Format(time.RFC3339)

	// Twitter API endpoint for mentions timeline
	user_id :=  "771833622286503936"
	url := fmt.Sprintf("https://api.x.com/2/users/%s/mentions", user_id)

    resp, err := t.GetTwitterRequest(ctx, url, nil)
    if err != nil {
        log.Printf("Failed to make Twitter request: %v", err)
        return nil, fmt.Errorf("failed to fetch mentions: %w", err)
    }

    if resp == nil || resp.Body == nil {
        log.Fatal("Received nil response or response body from Twitter API")
    }

    if resp.StatusCode != http.StatusOK {
        bodyBytes, _ := io.ReadAll(resp.Body)
        log.Printf("Twitter API error: %s", string(bodyBytes))
        return nil, fmt.Errorf("twitter API error: %d", resp.StatusCode)
    }

    var twitterResponse struct {
        Data []struct {
            Text string `json:"text"`
            ID   string `json:"id"`
        } `json:"data"`
    }

    err = json.NewDecoder(resp.Body).Decode(&twitterResponse)
    if err != nil {
        log.Printf("Failed to decode Twitter response: %v", err)
        return nil, fmt.Errorf("failed to parse mentions response: %w", err)
    }


	// Extract mentions from the tweet text
	var mentions []string
	for _, tweet := range twitterResponse.Data {
		extracted := extractMentions(tweet.Text)
		mentions = append(mentions, extracted...)
	}

	// Return the mentions or an empty array
	return mentions, nil
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

// Helper function to extract mentions from tweet text
func extractMentions(tweetText string) []string {
	var mentions []string
	words := strings.Split(tweetText, " ")
	for _, word := range words {
		if strings.HasPrefix(word, "@") {
			mentions = append(mentions, word)
		}
	}
	return mentions
}

func (t *TwitterServiceImpl) ReplyMention(ctx context.Context, mentionId string) {
}

func (t *TwitterServiceImpl) ReplyDM(ctx context.Context) {
}