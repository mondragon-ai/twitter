package mentions

import (
	"context"

	"github.com/twitter/model"
)

type MentionRepository interface {
	SaveMention(ctx context.Context, mention model.Mention)
	DeleteMention(ctx context.Context, mentiondID string)
	FindMentionById(ctx context.Context, mentiondID string) (*model.Mention, error)
	FindAllMentions(ctx context.Context) ([]model.Mention, error)

	SaveTweetIdea(ctx context.Context, ideas model.TweetIdea)
	DeleteTweetIdea(ctx context.Context, ideaID string)
	FindAllTweetIdeas(ctx context.Context) ([]model.TweetIdea, error)

	SaveThreadIdea(ctx context.Context, thread model.ThreadIdea)
	DeleteThreadIdea(ctx context.Context, threadID string)
	FindAllThreadIdeas(ctx context.Context) ([]model.ThreadIdea, error)

	SaveTweetClone(ctx context.Context, clone model.TweetClone)
	DeleteTweetClone(ctx context.Context, cloneID string)
	FindAllTweetClones(ctx context.Context) ([]model.TweetClone, error)

	SaveArticleUrl(ctx context.Context, article model.ArticleUrl)
	DeleteArticleUrl(ctx context.Context, articleID string)
	FindAllArticleUrls(ctx context.Context) ([]model.ArticleUrl, error)
}