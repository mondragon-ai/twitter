package service

import (
	"context"
	"net/http"

	"github.com/twitter/data/request"
	"github.com/twitter/model"
)

type TwitterService interface {
	PostTweet(ctx context.Context, request request.TweetCreateRequest)  (*http.Response, error)
	FetchMentions(ctx context.Context) ([]model.Mention, error) 
	ReplyMention(ctx context.Context, mentionId string) 
	ReplyDM(ctx context.Context)
}