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

func NewMentionRepository(Db *sql.DB) MentionRepository {
	return &MentionRepositoryImpl{Db: Db}
}

// Delete implements MentionsRepository
func (b *MentionRepositoryImpl) Delete(ctx context.Context, mentionId int) {
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
		err := result.Scan(&mention.Id, &mention.Id)
		helper.PanicIfError(err)

		mentions = append(mentions, mention)
	}

	return mentions
}

// FindById implements MentionsRepository
func (b *MentionRepositoryImpl) FindById(ctx context.Context, mentionId int) (model.Mention, error) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL := "select id,name from mention where id=$1"
	result, errQuery := tx.QueryContext(ctx, SQL, mentionId)
	helper.PanicIfError(errQuery)
	defer result.Close()

	mention := model.Mention{}

	if result.Next() {
		err := result.Scan(&mention.Id, &mention.Id)
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
	_, err = tx.ExecContext(ctx, SQL, mention.Id)
	helper.PanicIfError(err)
}

// Update implements MentionsRepository
func (b *MentionRepositoryImpl) Update(ctx context.Context, mention model.Mention) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL := "update mention set name=$1 where id=$2"
	_, err = tx.ExecContext(ctx, SQL, mention.Tweet, mention.Id)
	helper.PanicIfError(err)
}