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
}