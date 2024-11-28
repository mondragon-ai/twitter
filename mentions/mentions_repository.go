package mentions

import (
	"context"

	"github.com/twitter/model"
)

type MentionRepository interface {
	Save(ctx context.Context, mention model.Mention)
	Delete(ctx context.Context, mentiondID string)
	FindById(ctx context.Context, mentiondID string) (*model.Mention, error)
	FindAll(ctx context.Context) ([]model.Mention, error)
}