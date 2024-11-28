package service

import (
	"context"

	"github.com/twitter/data/request"
	"github.com/twitter/data/response"
	"github.com/twitter/helper"
	"github.com/twitter/mentions"
	"github.com/twitter/model"
)

type MentionsServiceImpl struct {
	MentionRepository mentions.MentionRepository
}

func NewMentionServiceImpl(MentionRepository mentions.MentionRepository) MentionService {
	return &MentionsServiceImpl{MentionRepository: MentionRepository}
}

// Create implements BookService
func (b *MentionsServiceImpl) Create(ctx context.Context, request request.MentionCreateRequest) {
	book := model.Mention{
		ID: request.ID,
		Content: request.Content,
	}
	b.MentionRepository.Save(ctx, book)
}

// Delete implements BookService
func (b *MentionsServiceImpl) Delete(ctx context.Context, mentionId string) {
	mention, err := b.MentionRepository.FindById(ctx, mentionId)
	helper.PanicIfError(err)
	b.MentionRepository.Delete(ctx, mention.ID)
}

// FindAll implements BookService
func (b *MentionsServiceImpl) FindAll(ctx context.Context) []response.MentionResponse {
	books, err := b.MentionRepository.FindAll(ctx)
	helper.PanicIfError(err)

	var bookResp []response.MentionResponse

	for _, value := range books {
		book := response.MentionResponse{ID: value.ID, Content: value.Content}
		bookResp = append(bookResp, book)
	}
	return bookResp

}

// FindById implements BookService
func (b *MentionsServiceImpl) FindById(ctx context.Context, mentionId string) response.MentionResponse {
	book, err := b.MentionRepository.FindById(ctx, mentionId)
	helper.PanicIfError(err)
	return response.MentionResponse(book)
}