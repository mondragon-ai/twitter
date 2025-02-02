package service

import (
	"context"

	"github.com/twitter/data/request"
	"github.com/twitter/data/response"
	"github.com/twitter/model"
)

type MentionService interface {
	Create(ctx context.Context, request request.MentionCreateRequest)
	Delete(ctx context.Context, mentionId int)
	FindById(ctx context.Context, mentionId int) (*response.MentionResponse, error)
	FindAll(ctx context.Context) []response.MentionResponse


	CreateTweetIdea(ctx context.Context, request request.TweetIdea)
	DeleteTweetIdea(ctx context.Context, ideaID int)
	FindAllTweetIdea(ctx context.Context) []model.TweetIdea

	CreateThreadIdea(ctx context.Context, request request.ThreadIdea)
	DeleteThreadIdea(ctx context.Context, ideaID int)
	FindAllThreadIdea(ctx context.Context) []model.ThreadIdea

	CreateTweetClone(ctx context.Context, request request.TweetClone)
	DeleteTweetClone(ctx context.Context, cloneID int)
	FindAllTweetClone(ctx context.Context) []model.TweetClone

	CreateArticleUrl(ctx context.Context, request request.UrlCreate)
	DeleteArticleUrl(ctx context.Context, articleID int)
	FindAllArticleUrl(ctx context.Context) []model.ArticleUrl
}