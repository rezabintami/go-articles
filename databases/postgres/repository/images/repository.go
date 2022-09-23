package images

import (
	"context"
	"database/sql"
	"go-articles/constants"
	"go-articles/modules/images"
	"log"
	"time"
)

type postgreImagesRepository struct {
	db *sql.DB
}

func NewPostgreImagesRepository(db *sql.DB) images.Repository {
	return &postgreImagesRepository{
		db: db,
	}
}

func (repository *postgreImagesRepository) Insert(ctx context.Context, data *images.Domain) (int, error) {
	record := fromDomain(*data)

	query := `INSERT INTO "images" ("path", "type", "created_at") VALUES ($1,$2,$3) RETURNING id`
	now := time.Now().Format("2006-01-02T15:04:05")
	err := repository.db.QueryRowContext(ctx, query, record.Path, record.Type, now).Scan(&data.ID)
	if err != nil {
		log.Println("[error] images.repository.Insert : failed to execute insert image query", err)
		return 0, err
	}
	return *data.ID, nil
}

func (repository *postgreImagesRepository) Delete(ctx context.Context, userId int) error {
	query := `DELETE FROM "images" WHERE "id" = $1`

	_, err := repository.db.ExecContext(ctx, query, userId)

	if err != nil {
		log.Println("[error] images.repository.Delete : failed to execute delete image query", err)
		if err == sql.ErrNoRows {
			err = constants.ErrRecordNotFound
		}
		return err
	}

	return nil
}
