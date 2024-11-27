package mentions

import (
	"context"
	"database/sql"
	"errors"

	"github.com/twitter/helper"
	"github.com/twitter/model"
)

type MentionRepositoryImpl struct {
	Db *sql.DB
}

func MentionCrud(Db *sql.DB) MentionRepository {
	return &MentionRepositoryImpl{Db: Db}
}

// Delete implements MentionsRepository
func (b *MentionRepositoryImpl) Delete(ctx context.Context, mentionId string) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL := "delete from mention where id =$1"
	_, errExec := tx.ExecContext(ctx, SQL, mentionId)
	helper.PanicIfError(errExec)
}

// FindAll implements MentionsRepository
func (b *MentionRepositoryImpl) FindAll(ctx context.Context) []model.Mention {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL := "select id,name from mention"
	result, errQuery := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(errQuery)
	defer result.Close()

	var mentions []model.Mention

	for result.Next() {
		mention := model.Mention{}
		err := result.Scan(&mention.ID, &mention.ID)
		helper.PanicIfError(err)

		mentions = append(mentions, mention)
	}

	return mentions
}

// FindById implements MentionsRepository
func (b *MentionRepositoryImpl) FindById(ctx context.Context, mentionId string) (model.Mention, error) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL := "select id,name from mention where id=$1"
	result, errQuery := tx.QueryContext(ctx, SQL, mentionId)
	helper.PanicIfError(errQuery)
	defer result.Close()

	mention := model.Mention{}

	if result.Next() {
		err := result.Scan(&mention.ID, &mention.ID)
		helper.PanicIfError(err)
		return mention, nil
	} else {
		return mention, errors.New("mention id not found")
	}
}

// Save implements MentionsRepository
func (b *MentionRepositoryImpl) Save(ctx context.Context, mention model.Mention) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL := "insert into mention(name) values ($1)"
	_, err = tx.ExecContext(ctx, SQL, mention.ID)
	helper.PanicIfError(err)
}