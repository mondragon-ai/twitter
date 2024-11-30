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
		ParentID: request.ParentID,
		AuthorID:  request.AuthorID,
		TweetID:  request.TweetID,
		AuthorName:  request.AuthorName,
		Content: request.Content,
	}
	b.MentionRepository.SaveMention(ctx, mention)
}

// Delete implements MentionService
func (b *MentionsServiceImpl) Delete(ctx context.Context, mentionId string) {
	mention, err := b.MentionRepository.FindMentionById(ctx, mentionId)
	helper.PanicIfError(err)
	b.MentionRepository.DeleteMention(ctx, mention.TweetID)
}

// FindAll implements MentionService
func (b *MentionsServiceImpl) FindAll(ctx context.Context) []response.MentionResponse {
	mentions, err := b.MentionRepository.FindAllMentions(ctx)
	helper.PanicIfError(err)

	var mentionsResp []response.MentionResponse

	for _, value := range mentions {
		mention := response.MentionResponse{
			ParentID: value.ParentID,
			AuthorID:  value.AuthorID,
			TweetID:  value.TweetID,
			AuthorName:  value.AuthorName,
			Content: value.Content,
		}
		mentionsResp = append(mentionsResp, mention)
	}
	return mentionsResp
}

// FindById implements MentionService
func (b *MentionsServiceImpl) FindById(ctx context.Context, mentionId string) (*response.MentionResponse, error) {
	mention, err := b.MentionRepository.FindMentionById(ctx, mentionId)
	if err != nil {
		return nil, err
	}

	return &response.MentionResponse{
		ParentID: mention.ParentID,
		AuthorID:  mention.AuthorID,
		TweetID:  mention.TweetID,
		AuthorName:  mention.AuthorName,
		Content: mention.Content,
    }, nil
}