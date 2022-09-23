package comments

import (
	"context"
	"database/sql"
	"go-articles/constants"
	"go-articles/modules/comments"
	"log"
	"time"
)

type postgreCommentRepository struct {
	db *sql.DB
}

func NewPostgreCommentRepository(db *sql.DB) comments.Repository {
	return &postgreCommentRepository{
		db: db,
	}
}

// Delete implements comments.Repository
func (repository *postgreCommentRepository) Delete(ctx context.Context, articleId, commentId, userId int) error {
	query := `DELETE FROM "comments" WHERE "id" = $1 AND "user_id" = $2`

	_, err := repository.db.ExecContext(ctx, query, commentId, userId)
	if err != nil {
		log.Println("[error] images.repository.Delete : failed to execute delete comment query", err)
		if err == sql.ErrNoRows {
			err = constants.ErrRecordNotFound
		}
		return err
	}

	query = `DELETE FROM "articles_comments" WHERE "article_id" = $1 AND "comment_id" = $2`
	_, err = repository.db.ExecContext(ctx, query, articleId, commentId)
	if err != nil {
		log.Println("[error] images.repository.Delete : failed to execute delete article_comment query", err)
		if err == sql.ErrNoRows {
			err = constants.ErrRecordNotFound
		}
		return err
	}
	
	return nil
}

// GetByArticleID implements comments.Repository
func (repository *postgreCommentRepository) GetByArticleID(ctx context.Context, articlesId int) (data []comments.Domain, err error) {
	query := `SELECT comments.id, users.name, comments.comment, comments.created_at  FROM "articles_comments" 
	LEFT JOIN "comments" ON comments.id = articles_comments.comment_id
	LEFT JOIN "users" ON users.id = comments.user_id
	WHERE articles_comments.article_id = $1 AND comments.deleted_at is NULL`

	rows, err := repository.db.QueryContext(ctx, query, articlesId)
	if err != nil {
		log.Println("[error] comments.repository.GetByArticleID : failed to execute get comment query", err)
		return []comments.Domain{}, err
	}

	defer rows.Close()
	for rows.Next() {
		temporary := Comments{}

		if err := rows.Scan(
			&temporary.ID,
			&temporary.User.Name,
			&temporary.Comment,
			&temporary.CreatedAt,
		); err != nil {
			log.Println("[error] comments.repository.GetByArticleID : failed to execute row get comment query", err)
			return []comments.Domain{}, err
		}

		data = append(data, *temporary.ToDomain())
	}

	return data, nil
}

// Insert implements comments.Repository
func (repository *postgreCommentRepository) Insert(ctx context.Context, data *comments.Domain, articlesId int) error {
	query := `INSERT INTO "comments" ("user_id", "comment","created_at") VALUES ($1,$2,$3) RETURNING id`
	now := time.Now().Format("2006-01-02T15:04:05")
	err := repository.db.QueryRowContext(ctx, query, data.UserID, data.Comment, now).Scan(&data.ID)
	if err != nil {
		log.Println("[error] comments.repository.Insert : failed to execute insert comment query", err)
		return err
	}

	query = `INSERT INTO "articles_comments" ("article_id","comment_id") VALUES ($1,$2)`
	_, err = repository.db.ExecContext(ctx, query, articlesId, data.ID)
	if err != nil {
		log.Println("[error] articles_comments.repository.Insert : failed to execute insert article_comment query", err)
		return err
	}

	return nil
}

// Update implements comments.Repository
func (repository *postgreCommentRepository) Update(ctx context.Context, data *comments.Domain, commentId int) error {
	query := `UPDATE "comments" SET "comment" = $1, "updated_at" = $2 WHERE "id" = $3 AND "user_id" = $4 AND "deleted_at" is NULL`
	now := time.Now().Format("2006-01-02T15:04:05")
	_, err := repository.db.ExecContext(ctx, query, data.Comment, now, commentId, data.UserID)
	if err != nil {
		log.Println("[error] comments.repository.Update : failed to execute update comment query", err)
		return err
	}

	return nil
}
