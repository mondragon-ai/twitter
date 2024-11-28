package service

import (
	"context"
	"net/http"

	"github.com/twitter/data/request"
	"github.com/twitter/data/response"
)

type TwitterService interface {
	PostTweet(ctx context.Context, request request.TweetCreateRequest)  (*http.Response, error)
	FetchMentions(ctx context.Context) (*[]response.MentionResponse, error)
	ReplyMention(ctx context.Context, mentionId string) 
	ReplyDM(ctx context.Context)
}