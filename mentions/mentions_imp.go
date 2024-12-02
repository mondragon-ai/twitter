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



// Delete Tweet Idea
func (b *MentionRepositoryImpl) DeleteTweetIdea(ctx context.Context, ideaID string) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL := "delete from ideas where id =$1"
	_, errExec := tx.ExecContext(ctx, SQL, ideaID)
	helper.PanicIfError(errExec)
}

// FindAll Tweet Ideas
func (b *MentionRepositoryImpl) FindAllTweetIdeas(ctx context.Context) ([]model.TweetIdea, error) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL := "SELECT id, idea, used_count FROM ideas"
	result, errQuery := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(errQuery)
	defer result.Close()

	var ideas []model.TweetIdea

	for result.Next() {
		idea := model.TweetIdea{}
		err := result.Scan(&idea.ID, &idea.Idea, &idea.UsedCount)
		helper.PanicIfError(err)

		ideas = append(ideas, idea)
	}

	return ideas, nil
}

// Save Tweet Idea
func (b *MentionRepositoryImpl) SaveTweetIdea(ctx context.Context, idea model.TweetIdea) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL :=  `
		INSERT INTO ideas (idea, used_count)
		VALUES ($1, $2,)
		ON CONFLICT (id) DO NOTHING;
	`

	// Execute the query with the idea data
	_, err = tx.ExecContext(ctx, SQL, idea.Idea, idea.UsedCount)
	helper.PanicIfError(err)
}




// Delete Tweet Thread Idea
func (b *MentionRepositoryImpl) DeleteThreadIdea(ctx context.Context, ideaID string) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL := "delete from threads where id =$1"
	_, errExec := tx.ExecContext(ctx, SQL, ideaID)
	helper.PanicIfError(errExec)
}

// FindAll Tweet Thread Idea
func (b *MentionRepositoryImpl) FindAllThreadIdeas(ctx context.Context) ([]model.ThreadIdea, error) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL := "SELECT id, idea, used_count FROM threads"
	result, errQuery := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(errQuery)
	defer result.Close()

	var threads []model.ThreadIdea

	for result.Next() {
		thread := model.ThreadIdea{}
		err := result.Scan(&thread.ID, &thread.Idea, &thread.UsedCount)
		helper.PanicIfError(err)

		threads = append(threads, thread)
	}

	return threads, nil
}

// Save Tweet Thread Idea
func (b *MentionRepositoryImpl) SaveThreadIdea(ctx context.Context, thread model.ThreadIdea) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL :=  `
		INSERT INTO threads (idea, used_count)
		VALUES ($1, $2)
		ON CONFLICT (id) DO NOTHING;
	`

	// Execute the query with the idea data
	_, err = tx.ExecContext(ctx, SQL, thread.Idea, thread.UsedCount)
	helper.PanicIfError(err)
}




// Delete Tweet clone Idea
func (b *MentionRepositoryImpl) DeleteTweetClone(ctx context.Context, ideaID string) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL := "delete from clones where id =$1"
	_, errExec := tx.ExecContext(ctx, SQL, ideaID)
	helper.PanicIfError(errExec)
}

// FindAll Tweet Clone Idea
func (b *MentionRepositoryImpl) FindAllTweetClones(ctx context.Context) ([]model.TweetClone, error) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL := "SELECT id, author_name, tweet FROM clones"
	result, errQuery := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(errQuery)
	defer result.Close()

	var clones []model.TweetClone

	for result.Next() {
		clone := model.TweetClone{}
		err := result.Scan(&clone.ID, &clone.AuthorName, &clone.Tweet)
		helper.PanicIfError(err)

		clones = append(clones, clone)
	}

	return clones, nil
}

// Save Tweet clone Idea
func (b *MentionRepositoryImpl) SaveTweetClone(ctx context.Context, clone model.TweetClone) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL :=  `
		INSERT INTO clones (author_name, tweet)
		VALUES ($1, $2)
		ON CONFLICT (id) DO NOTHING;
	`

	// Execute the query with the idea data
	_, err = tx.ExecContext(ctx, SQL, clone.AuthorName, clone.Tweet)
	helper.PanicIfError(err)
}




// Delete Tweet Article url
func (b *MentionRepositoryImpl) DeleteArticleUrl(ctx context.Context, articleID string) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL := "delete from articles where id =$1"
	_, errExec := tx.ExecContext(ctx, SQL, articleID)
	helper.PanicIfError(errExec)
}

// FindAll Tweet Article urls
func (b *MentionRepositoryImpl) FindAllArticleUrls(ctx context.Context) ([]model.ArticleUrl, error) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL := "SELECT id, url, title FROM articles"
	result, errQuery := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(errQuery)
	defer result.Close()

	var articles []model.ArticleUrl

	for result.Next() {
		article := model.ArticleUrl{}
		err := result.Scan(&article.ID, &article.Url, &article.Title)
		helper.PanicIfError(err)

		articles = append(articles, article)
	}

	return articles, nil
}

// Save Tweet article url
func (b *MentionRepositoryImpl) SaveArticleUrl(ctx context.Context, article model.ArticleUrl) {
	tx, err := b.Db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	SQL :=  `
		INSERT INTO articles (title, url)
		VALUES ($1, $2)
		ON CONFLICT (id) DO NOTHING;
	`

	// Execute the query with the url data
	_, err = tx.ExecContext(ctx, SQL, article.Title, article.Url)
	helper.PanicIfError(err)
}

