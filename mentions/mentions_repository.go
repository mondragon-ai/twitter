package mentions

import (
	"context"

	"github.com/twitter/model"
)

type MentionRepository interface {
	Save(ctx context.Context, book model.Mention)
	Update(ctx context.Context, book model.Mention)
	Delete(ctx context.Context, bookId int)
	FindById(ctx context.Context, bookId int) (model.Mention, error)
	FindAll(ctx context.Context) []model.Mention
}