package service

import (
	"context"

	"github.com/twitter/data/request"
	"github.com/twitter/data/response"
)

type MentionService interface {
	Create(ctx context.Context, request request.MentionCreateRequest)
	Delete(ctx context.Context, mentionId string)
	FindById(ctx context.Context, mentionId string) (*response.MentionResponse, error)
	FindAll(ctx context.Context) []response.MentionResponse
}