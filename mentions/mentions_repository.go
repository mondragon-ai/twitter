package mentions

import (
	"context"

	"github.com/twitter/model"
)

type MentionRepository interface {
	SaveMention(ctx context.Context, mention model.Mention)
	DeleteMention(ctx context.Context, mentiondID int)
	FindMentionById(ctx context.Context, mentiondID int) (*model.Mention, error)
	FindAllMentions(ctx context.Context) ([]model.Mention, error)

	SaveTweetIdea(ctx context.Context, ideas model.TweetIdea)
	DeleteTweetIdea(ctx context.Context, ideaID int)
	FindAllTweetIdeas(ctx context.Context) ([]model.TweetIdea, error)

	SaveThreadIdea(ctx context.Context, thread model.ThreadIdea)
	DeleteThreadIdea(ctx context.Context, threadID int)
	FindAllThreadIdeas(ctx context.Context) ([]model.ThreadIdea, error)

	SaveTweetClone(ctx context.Context, clone model.TweetClone)
	DeleteTweetClone(ctx context.Context, cloneID int)
	FindAllTweetClones(ctx context.Context) ([]model.TweetClone, error)

	SaveArticleUrl(ctx context.Context, article model.ArticleUrl)
	DeleteArticleUrl(ctx context.Context, articleID int)
	FindAllArticleUrls(ctx context.Context) ([]model.ArticleUrl, error)
}