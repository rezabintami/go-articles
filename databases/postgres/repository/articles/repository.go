package articles

import (
	"context"
	"database/sql"
	"go-articles/constants"
	"go-articles/modules/articles"
	"log"
	"time"
)

type postgreArticlesRepository struct {
	db *sql.DB
}

func NewPostgreArticlesRepository(db *sql.DB) articles.Repository {
	return &postgreArticlesRepository{
		db: db,
	}
}

func (repository *postgreArticlesRepository) GetByID(ctx context.Context, id int) (articles.Domain, error) {
	article := Articles{}
	query := `SELECT articles.id, users.name, im.path, articles.title, articles.description, articles.created_at  FROM "articles" 
	LEFT JOIN "users" ON users.id = articles.user_id
	LEFT JOIN "images" im ON im.id = articles.image_id
	WHERE articles.id = $1 AND "deleted_at" is NULL`

	err := repository.db.QueryRowContext(ctx, query, id).Scan(
		&article.ID,
		&article.Role.Name,
		&article.Image.Path,
		&article.Title,
		&article.Description,
		&article.CreatedAt,
	)

	if err != nil {
		log.Println("[error] articles.repository.GetByID : failed to execute get data article query", err)
		if err == sql.ErrNoRows {
			err = constants.ErrRecordNotFound
		}
		return articles.Domain{}, err
	}

	return *article.ToDomain(), nil
}

func (repository *postgreArticlesRepository) Insert(ctx context.Context, data *articles.Domain) error {
	record := fromDomain(*data)

	query := `INSERT INTO "articles" ("user_id", "title", "description","created_at") VALUES ($1,$2,$3,$4)`
	now := time.Now().Format("2006-01-02T15:04:05")
	_, err := repository.db.ExecContext(ctx, query, record.UserID, record.Title, record.Description, now)
	if err != nil {
		log.Println("[error] articles.repository.Insert : failed to execute article query", err)
		return err
	}

	return nil
}

func (repository *postgreArticlesRepository) Update(ctx context.Context, data *articles.Domain, id int) error {
	record := fromDomain(*data)

	query := `UPDATE "articles" SET "user_id" = $1, "image_id" = $2, "title" = $3, "description" = $4, "updated_at" = $5 WHERE "id" = $6 AND "deleted_at" is NULL`
	now := time.Now().Format("2006-01-02T15:04:05")
	_, err := repository.db.ExecContext(ctx, query, record.UserID, record.ImageID, record.Title, record.Description, now, id)
	if err != nil {
		log.Println("[error] articles.repository.Update : failed to execute update article query", err)
		return err
	}

	return nil
}

func (repository *postgreArticlesRepository) Fetch(ctx context.Context, page, perpage, userId int, by, search, sort string) ([]articles.Domain, int, error) {
	activeString := ``
	offset := (page - 1) * perpage

	if search != "" {
		activeString = ` AND LOWER(articles.title) LIKE '%` + search + `%'`
	}

	query := `SELECT articles.id, articles.title, articles.user_id, articles.description, im.path, articles.created_at, articles.updated_at, articles.deleted_at FROM "articles" 
	LEFT JOIN "images" im ON im.id = articles.image_id
	WHERE "user_id" = $1 AND "deleted_at" is NULL` + activeString + `
	ORDER BY ` + by + ` ` + sort + `
	OFFSET $2 LIMIT $3`

	rows, err := repository.db.QueryContext(ctx, query, userId, offset, perpage)
	if err != nil {
		log.Println("[error] articles.repository.Fetch : failed to execute fetch articles query", err)
		return []articles.Domain{}, 0, err
	}

	var total int
	query = `SELECT COUNT(*) FROM "articles" WHERE "user_id" = $1 AND "deleted_at" is NULL` + activeString
	err = repository.db.QueryRowContext(ctx, query, userId).Scan(&total)
	if err != nil {
		log.Println("[error] articles.repository.Fetch : failed to execute count articles query", err)
		return []articles.Domain{}, 0, err
	}

	var result []articles.Domain
	defer rows.Close()
	for rows.Next() {
		temporary := Articles{}

		if err := rows.Scan(
			&temporary.ID,
			&temporary.Title,
			&temporary.UserID,
			&temporary.Description,
			&temporary.Image.Path,
			&temporary.CreatedAt,
			&temporary.UpdatedAt,
			&temporary.DeletedAt,
		); err != nil {
			log.Println("[error] articles.repository.Fetch : failed to execute row fetch articles query", err)
			return []articles.Domain{}, 0, err
		}
		result = append(result, *temporary.ToDomain())
	}

	return result, total, nil
}

func (repository *postgreArticlesRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE "articles" SET "deleted_at" = $1 WHERE "id" = $2`
	now := time.Now().Format("2006-01-02T15:04:05")
	_, err := repository.db.ExecContext(ctx, query, now, id)

	if err != nil {
		log.Println("[error] articles.repository.Delete : failed to execute delete article query", err)
		return err
	}

	return nil
}

func (repository *postgreArticlesRepository) Search(ctx context.Context, search string) ([]articles.Domain, error) {
	query := `SELECT articles.id, articles.title, us.name, articles.description, im.path, articles.created_at, articles.updated_at FROM "articles" 
	LEFT JOIN "images" im ON im.id = articles.image_id
	LEFT JOIN "users" us ON us.id = articles.user_id
	WHERE "deleted_at" is NULL AND (LOWER(articles.title) LIKE '%` + search + `%' OR LOWER(us.name) LIKE '%` + search + `%')`

	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		log.Println("[error] articles.repository.Search : failed to execute search articles query", err)
		return []articles.Domain{}, err
	}

	var result []articles.Domain
	defer rows.Close()
	for rows.Next() {
		temporary := Articles{}

		if err := rows.Scan(
			&temporary.ID,
			&temporary.Title,
			&temporary.User.Name,
			&temporary.Description,
			&temporary.Image.Path,
			&temporary.CreatedAt,
			&temporary.UpdatedAt,
		); err != nil {
			log.Println("[error] articles.repository.Search : failed to execute row search articles query", err)
			return []articles.Domain{}, err
		}
		result = append(result, *temporary.ToDomain())
	}

	return result, nil
}
