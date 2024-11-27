package service

import (
	"context"

	"github.com/twitter/data/request"
	"github.com/twitter/data/response"
)

type MentionService interface {
	Create(ctx context.Context, request request.MentionCreateRequest)
	Delete(ctx context.Context, bookId int)
	FindById(ctx context.Context, bookId int) response.MentionResponse
	FindAll(ctx context.Context) []response.MentionResponse
}