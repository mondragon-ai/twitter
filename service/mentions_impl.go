package service

import (
	"context"

	"github.com/twitter/data/request"
	"github.com/twitter/data/response"
	"github.com/twitter/helper"
	"github.com/twitter/mentions"
	"github.com/twitter/model"
)

type BookServiceImpl struct {
	BookRepository mentions.MentionRepository
}

func NewMentionServiceImpl(bookRepository mentions.MentionRepository) MentionService {
	return &BookServiceImpl{BookRepository: bookRepository}
}

// Create implements BookService
func (b *BookServiceImpl) Create(ctx context.Context, request request.MentionCreateRequest) {
	book := model.Mention{
		Id: request.ID,
		Tweet: request.Tweet,
	}
	b.BookRepository.Save(ctx, book)
}

// Delete implements BookService
func (b *BookServiceImpl) Delete(ctx context.Context, bookId int) {
	mention, err := b.BookRepository.FindById(ctx, bookId)
	helper.PanicIfError(err)
	b.BookRepository.Delete(ctx, mention.Id)
}

// FindAll implements BookService
func (b *BookServiceImpl) FindAll(ctx context.Context) []response.MentionResponse {
	books := b.BookRepository.FindAll(ctx)

	var bookResp []response.MentionResponse

	for _, value := range books {
		book := response.MentionResponse{Id: value.Id, Tweet: value.Tweet}
		bookResp = append(bookResp, book)
	}
	return bookResp

}

// FindById implements BookService
func (b *BookServiceImpl) FindById(ctx context.Context, bookId int) response.MentionResponse {
	book, err := b.BookRepository.FindById(ctx, bookId)
	helper.PanicIfError(err)
	return response.MentionResponse(book)
}

// Update implements BookService
func (b *BookServiceImpl) Update(ctx context.Context, request request.MentionCreateRequest) {
	mention, err := b.BookRepository.FindById(ctx, request.ID)
	helper.PanicIfError(err)

	mention.Tweet = request.Tweet
	b.BookRepository.Update(ctx, mention)
}