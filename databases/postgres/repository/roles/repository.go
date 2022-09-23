package roles

import (
	"context"
	"database/sql"
	"go-articles/constants"
	"go-articles/modules/roles"
	"log"
	"strings"
	"time"
)

type postgreRolesRepository struct {
	db *sql.DB
}

func NewPostgreRolesRepository(db *sql.DB) roles.Repository {
	return &postgreRolesRepository{
		db: db,
	}
}

func (repository *postgreRolesRepository) GetByID(ctx context.Context, id int) (roles.Domain, error) {
	role := Roles{}
	query := `SELECT * FROM "roles" WHERE "id" = $1`

	err := repository.db.QueryRowContext(ctx, query, id).Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt)

	if err != nil {
		log.Println("[error] roles.repository.GetByID : failed to execute get data role query", err)
		if err == sql.ErrNoRows {
			err = constants.ErrRecordNotFound
		}
		return roles.Domain{}, err
	}

	return *role.ToDomain(), nil
}

func (repository *postgreRolesRepository) Insert(ctx context.Context, data *roles.Domain) error {
	record := fromDomain(*data)

	query := `INSERT INTO "roles" ("name", "created_at") VALUES ($1,$2)`
	now := time.Now().Format("2006-01-02T15:04:05")
	_, err := repository.db.ExecContext(ctx, query, strings.ToUpper(record.Name), now)
	if err != nil {
		log.Println("[error] roles.repository.Insert : failed to execute insert role query", err)
		return err
	}
	return nil
}

func (repository *postgreRolesRepository) Update(ctx context.Context, data *roles.Domain, id int) error {
	record := fromDomain(*data)

	query := `UPDATE "roles" SET "name" = $1, "updated_at" = $2 WHERE "id" = $3`
	now := time.Now().Format("2006-01-02T15:04:05")
	_, err := repository.db.ExecContext(ctx, query, strings.ToUpper(record.Name), now, id)
	if err != nil {
		log.Println("[error] roles.repository.Update : failed to execute update role query", err)
		return err
	}
	return nil
}

func (repository *postgreRolesRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM "roles" WHERE "id" = $1`

	_, err := repository.db.ExecContext(ctx, query, id)

	if err != nil {
		log.Println("[error] roles.repository.Delete : failed to execute delete role query", err)
		if err == sql.ErrNoRows {
			err = constants.ErrRecordNotFound
		}
		return err
	}

	return nil
}

func (repository *postgreRolesRepository) Fetch(ctx context.Context, page, perpage int) (data []roles.Domain, count int, err error) {
	offset := (page - 1) * perpage

	query := `SELECT * FROM "roles"	OFFSET $1 LIMIT $2`

	rows, err := repository.db.QueryContext(ctx, query, offset, perpage)
	if err != nil {
		log.Println("[error] roles.repository.Fetch : failed to execute fetch roles query", err)
		return []roles.Domain{}, 0, err
	}

	query = `SELECT COUNT(*) FROM "roles"`
	err = repository.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		log.Println("[error] roles.repository.Fetch : failed to execute count roles query", err)
		return []roles.Domain{}, 0, err
	}

	defer rows.Close()
	for rows.Next() {
		role := Roles{}

		if err = rows.Scan(
			&role.ID,
			&role.Name,
			&role.CreatedAt,
			&role.UpdatedAt,
		); err != nil {
			log.Println("[error] roles.repository.Fetch : failed to execute row fetch roles query", err)
			return []roles.Domain{}, 0, err
		}

		data = append(data, *role.ToDomain())
	}

	return data, count, nil
}
