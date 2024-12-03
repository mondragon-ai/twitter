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
	"github.com/twitter/model"
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
	// bearer := os.Getenv("ANGEL_BEARER")

	// token := fmt.Sprintf("Bearer %s", bearer)
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

    // Create HTTP request
    req, err := http.NewRequest(method, url, requestBody)
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
    req.Header.Set("Authorization", signature)

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

func (t *TwitterServiceImpl) GetTwitterRequest(ctx context.Context, url string, body interface{}) ([]byte, error) {
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

    if resp == nil || resp.Body == nil {
        log.Fatal("Received nil response or response body from Twitter API")
    }

    if resp.StatusCode != http.StatusOK {
        bodyBytes, _ := io.ReadAll(resp.Body)
        log.Printf("Twitter API error: %s", string(bodyBytes))
        return nil, fmt.Errorf("twitter API error: %d", resp.StatusCode)
    }

    data, err := io.ReadAll(resp.Body)
    defer resp.Body.Close()
    if err != nil {
        log.Fatal(err)
    }

    return data, nil
}

func (t *TwitterServiceImpl) MakeThreadTweet(ctx context.Context, text *string, id *string) (*http.Response, *string, error) {
	body := map[string]interface{}{
		"text": *text, // Dereference the pointer to get the string value
	}

	if id != nil {
		body["reply"] = map[string]string{
			"in_reply_to_tweet_id": *id, // Dereference the pointer
		}
	}

	// Make the Twitter request
	url := "https://api.twitter.com/2/tweets"
	resp, err := t.MakeTwitterRequest(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to post tweet: %w", err)
	}

	// Check for a non-2xx status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("twitter API error: %s", string(respBody))
	}

	// Parse the response body
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract the tweet ID
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("unexpected response format: %v", response)
	}

	tweetID, ok := data["id"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("tweet ID not found in response: %v", data)
	}

	return resp, &tweetID, nil
}

func (t *TwitterServiceImpl) FetchMentions(ctx context.Context) ([]model.Mention, error)  {
	// now := time.Now()
	// startTime := now.Add(-time.Duration(60*24*7) * time.Minute).Format(time.RFC3339)

	// Twitter API endpoint for mentions timeline
	user_id :=  "771833622286503936"
	url := fmt.Sprintf("https://api.x.com/2/users/%s/mentions?expansions=attachments.poll_ids,attachments.media_keys,author_id,geo.place_id,in_reply_to_user_id,referenced_tweets.id,entities.mentions.username,referenced_tweets.id.author_id,edit_history_tweet_ids&tweet.fields=created_at,author_id,text,attachments,conversation_id,referenced_tweets", user_id)

    data, err := t.GetTwitterRequest(ctx, url, nil)
    if err != nil {
        log.Printf("Failed to make Twitter request: %v", err)
        return nil, fmt.Errorf("failed to fetch mentions: %w", err)
    }

	twitterResponse := &request.RootTweetMentions{}
	err = json.Unmarshal(data, twitterResponse)
	if err != nil {
		log.Printf("Failed to decode Twitter response: %v", err)
		return nil, fmt.Errorf("failed to parse mentions response: %w", err)
	}

	mentions := []model.Mention{}
	for _, tweet := range twitterResponse.Data {
		if tweet.ConversationID == tweet.ID {
			for _, u := range twitterResponse.Includes.Users {
				if u.ID == tweet.AuthorID {
					mentions = append(mentions, model.Mention{
						ParentID: tweet.ConversationID,
						AuthorID: tweet.AuthorID,
						TweetID: tweet.ID,
						Content: tweet.Text,
						AuthorName: u.Name,
						ParentContent: "",
						CreatedAt: tweet.CreatedAt,
					})
				}
			}
		} else {
			for _, t := range twitterResponse.Includes.Tweets {
				if t.ID == tweet.ConversationID {
					for _, u := range twitterResponse.Includes.Users {
						if u.ID == tweet.AuthorID {
							mentions = append(mentions, model.Mention{
								ParentID: tweet.ConversationID,
								AuthorID: tweet.AuthorID,
								TweetID: tweet.ID,
								Content: tweet.Text,
								AuthorName: u.Name,
								ParentContent: t.Text,
								CreatedAt: tweet.CreatedAt,
							})
						}
					}
				}
			}
		}
	}

	// Return the mentions or an empty array
	return mentions, nil
}

func (t *TwitterServiceImpl) PostTweet(ctx context.Context, request request.TweetCreateRequest) (*http.Response, error) {
	

	text := ""
	var resp *http.Response
	var err error
	switch request.Type {
	case "create":

		messages := generateCreativePrompt(t.MentionRepository, ctx)
		completion, err := OpenAIChatCompletion(messages, 300)
		if err != nil {
			return nil, fmt.Errorf("could not compelte chat: %w", err)
		}

		filteredHash := filterWordsWithHash(completion);
		text = strings.Replace(filteredHash, "\"", "", -1)
	case "clone":
		messages := generateClonePrompt(t.MentionRepository, ctx)
		completion, err := OpenAIChatCompletion(messages, 300)
		if err != nil {
			return nil, fmt.Errorf("could not compelte chat: %w", err)
		}

		filteredHash := filterWordsWithHash(completion);
		text = strings.Replace(filteredHash, "\"", "", -1)
	case "article":
		messages := generateArticlePrompt(t.MentionRepository, ctx)
		completion, err := OpenAIChatCompletion(messages, 300)
		if err != nil {
			return nil, fmt.Errorf("could not compelte chat: %w", err)
		}

		filteredHash := filterWordsWithHash(completion);
		text = strings.Replace(filteredHash, "\"", "", -1)

	case "thread":
		messages := generateThreadPrompt(t.MentionRepository, ctx)
		completion, err := OpenAIChatCompletion(messages, 1000)
		if err != nil {
			return nil, fmt.Errorf("could not compelte chat: %w", err)
		}

		filteredHash := filterWordsWithHash(completion);
		resp := cleanThread(filteredHash, t, ctx)
		// Ensure the response is not nil
		if resp == nil {
			return nil, fmt.Errorf("failed to post thread: response is nil")
		}

		return resp, nil
	default:
		return nil, fmt.Errorf("default to post tweet: %v", "TEST")
	}

	// t.MentionRepository.

	// Prepare the body for the POST request
	body := map[string]string{
		"text": text,
	}

	// Make the Twitter request
	url := "https://api.twitter.com/2/tweets"
	resp, err = t.MakeTwitterRequest(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to post tweet: %w", err)
	}

	return resp, nil
}

func (t *TwitterServiceImpl) ReplyMention(ctx context.Context, mentionId string) {
}

func (t *TwitterServiceImpl) ReplyDM(ctx context.Context) {
}