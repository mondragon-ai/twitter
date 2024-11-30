package mentions

import (
	"context"
	"database/sql"
	"fmt"

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
func (b *MentionRepositoryImpl) DeleteMention(ctx context.Context, mentionId string) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL := "delete from mention where id =$1"
	_, errExec := tx.ExecContext(ctx, SQL, mentionId)
	helper.PanicIfError(errExec)
}

// FindAll implements MentionsRepository
func (b *MentionRepositoryImpl) FindAllMentions(ctx context.Context) ([]model.Mention, error) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL := "SELECT parent_id, author_id, tweet_id, content, author_name, created_at FROM mention"
	result, errQuery := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(errQuery)
	defer result.Close()

	var mentions []model.Mention

	for result.Next() {
		mention := model.Mention{}
		err := result.Scan(&mention.ParentID, &mention.TweetID,&mention.AuthorID, &mention.AuthorName,&mention.Content, &mention.CreatedAt)
		helper.PanicIfError(err)

		mentions = append(mentions, mention)
	}

	return mentions, nil
}

// FindById implements MentionsRepository
func (b *MentionRepositoryImpl) FindMentionById(ctx context.Context, mentionId string) (*model.Mention, error) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL := "SELECT parent_id, author_id, tweet_id, content, author_name, created_at FROM mention where id=$1"
	result, errQuery := tx.QueryContext(ctx, SQL, mentionId)
	helper.PanicIfError(errQuery)
	defer result.Close()

	mention := model.Mention{}

	if result.Next() {
		err := result.Scan(&mention.ParentID, &mention.TweetID,&mention.AuthorID, &mention.AuthorName,&mention.Content, &mention.CreatedAt)
		helper.PanicIfError(err)
		return &mention, nil
	} else {
		return &mention, fmt.Errorf("mention with id %s not found", mentionId)

	}
}

// Save implements MentionsRepository
func (b *MentionRepositoryImpl) SaveMention(ctx context.Context, mention model.Mention) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// "INSERT INTO mention (id, content, author, created) VALUES ($1, $2, $3, $4)"
	// SQL query for inserting into the `mentions` table
	SQL :=  `
		INSERT INTO mention (parent_id, author_id, tweet_id, content, author_name, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (tweet_id) DO NOTHING;
	`

	// Execute the query with the mention data
	_, err = tx.ExecContext(ctx, SQL, mention.ParentID, mention.AuthorID, mention.TweetID, mention.Content, mention.AuthorName, mention.CreatedAt)
	helper.PanicIfError(err)
}