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

// Create implements MentionService
func (b *MentionsServiceImpl) Create(ctx context.Context, request request.MentionCreateRequest) {
	mention := model.Mention{
		ID: request.ID,
		Content: request.Content,
		Author:  request.Author,
		Created: request.Created,
	}
	b.MentionRepository.Save(ctx, mention)
}

// Delete implements MentionService
func (b *MentionsServiceImpl) Delete(ctx context.Context, mentionId string) {
	mention, err := b.MentionRepository.FindById(ctx, mentionId)
	helper.PanicIfError(err)
	b.MentionRepository.Delete(ctx, mention.ID)
}

// FindAll implements MentionService
func (b *MentionsServiceImpl) FindAll(ctx context.Context) []response.MentionResponse {
	mentions, err := b.MentionRepository.FindAll(ctx)
	helper.PanicIfError(err)

	var mentionsResp []response.MentionResponse

	for _, value := range mentions {
		mention := response.MentionResponse{ID: value.ID, Content: value.Content}
		mentionsResp = append(mentionsResp, mention)
	}
	return mentionsResp
}

// FindById implements MentionService
func (b *MentionsServiceImpl) FindById(ctx context.Context, mentionId string) (*response.MentionResponse, error) {
	mention, err := b.MentionRepository.FindById(ctx, mentionId)
	if err != nil {
		return nil, err
	}

	return &response.MentionResponse{
        ID:   mention.ID,
		Content: mention.Content,
		Author: mention.Author,
		Created: mention.Created,
    }, nil
}