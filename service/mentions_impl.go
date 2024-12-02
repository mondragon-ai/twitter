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




// Save Tweet Idea
func (b *MentionsServiceImpl) CreateTweetIdea(ctx context.Context, request request.TweetIdea) {
	idea := model.TweetIdea{
		ID: 1234,
		Idea: request.Idea,
		UsedCount:  request.UsedCount,
	}
	b.MentionRepository.SaveTweetIdea(ctx, idea)
}

// Delete Tweet Idea
func (b *MentionsServiceImpl) DeleteTweetIdea(ctx context.Context, ideaID string) {
	b.MentionRepository.DeleteMention(ctx, ideaID)
}

// Find All Tweet Ideas
func (b *MentionsServiceImpl) FindAllTweetIdea(ctx context.Context) []model.TweetIdea {
	mentions, err := b.MentionRepository.FindAllTweetIdeas(ctx)
	helper.PanicIfError(err)

	var mentionsResp []model.TweetIdea

	for _, value := range mentions {
		mention := model.TweetIdea{
			ID: value.ID,
			Idea: value.Idea,
			UsedCount: value.UsedCount,
		}
		mentionsResp = append(mentionsResp, mention)
	}
	return mentionsResp
}




// Save Tweet Clone Idea
func (b *MentionsServiceImpl) CreateTweetClone(ctx context.Context, request request.TweetClone) {
	clone := model.TweetClone{
		ID: 1234,
		AuthorName: request.AuthorName,
		Tweet:  request.Tweet,
	}
	b.MentionRepository.SaveTweetClone(ctx, clone)
}

// Delete Tweet Clone Idea
func (b *MentionsServiceImpl) DeleteTweetClone(ctx context.Context, cloneID string) {
	b.MentionRepository.DeleteMention(ctx, cloneID)
}

// Find All Tweet Clone Ideas
func (b *MentionsServiceImpl) FindAllTweetClone(ctx context.Context) []model.TweetClone {
	mentions, err := b.MentionRepository.FindAllTweetClones(ctx)
	helper.PanicIfError(err)

	var mentionsResp []model.TweetClone

	for _, value := range mentions {
		mention := model.TweetClone{
			ID: value.ID,
			AuthorName: value.AuthorName,
			Tweet:  value.Tweet,
		}
		mentionsResp = append(mentionsResp, mention)
	}
	return mentionsResp
}



// Save Tweet Thread Idea
func (b *MentionsServiceImpl) CreateThreadIdea(ctx context.Context, request request.ThreadIdea) {
	thread := model.ThreadIdea{
		ID: 1234,
		Idea: request.Idea,
		UsedCount:  request.UsedCount,
	}
	b.MentionRepository.SaveThreadIdea(ctx, thread)
}

// Delete Tweet Thread Idea
func (b *MentionsServiceImpl) DeleteThreadIdea(ctx context.Context, threadID string) {
	b.MentionRepository.DeleteMention(ctx, threadID)
}

// Find All Tweet Thread Ideas
func (b *MentionsServiceImpl) FindAllThreadIdea(ctx context.Context) []model.ThreadIdea {
	mentions, err := b.MentionRepository.FindAllThreadIdeas(ctx)
	helper.PanicIfError(err)

	var mentionsResp []model.ThreadIdea

	for _, value := range mentions {
		mention := model.ThreadIdea{
			ID: value.ID,
			Idea: value.Idea,
			UsedCount:  value.UsedCount,
		}
		mentionsResp = append(mentionsResp, mention)
	}
	return mentionsResp
}



// Save Article Url
func (b *MentionsServiceImpl) CreateArticleUrl(ctx context.Context, request request.UrlCreate) {
	article := model.ArticleUrl{
		ID: 1234,
		Url: request.Url,
		Title:  request.Title,
	}
	b.MentionRepository.SaveArticleUrl(ctx, article)
}

// Delete Article Url
func (b *MentionsServiceImpl) DeleteArticleUrl(ctx context.Context, articleID string) {
	b.MentionRepository.DeleteMention(ctx, articleID)
}

// Find All Article Urls
func (b *MentionsServiceImpl) FindAllArticleUrl(ctx context.Context) []model.ArticleUrl {
	mentions, err := b.MentionRepository.FindAllArticleUrls(ctx)
	helper.PanicIfError(err)

	var mentionsResp []model.ArticleUrl

	for _, value := range mentions {
		mention := model.ArticleUrl{
			ID: value.ID,
			Url: value.Url,
			Title:  value.Title,
		}
		mentionsResp = append(mentionsResp, mention)
	}
	return mentionsResp
}



