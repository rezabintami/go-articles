package users

import (
	"context"
	"database/sql"
	"go-articles/constants"
	"go-articles/modules/users"
	"log"
	"time"
)

type postgreUsersRepository struct {
	db *sql.DB
}

func NewPostgreUsersRepository(db *sql.DB) users.Repository {
	return &postgreUsersRepository{
		db: db,
	}
}

func (repository *postgreUsersRepository) GetByID(ctx context.Context, id int) (users.Domain, error) {
	user := Users{}
	query := `SELECT users.id, users.name, users.email, im.id, im.path, roles.name, users.created_at FROM "users" 
	LEFT JOIN "images" im ON im.id = users.image_id
	LEFT JOIN "roles" ON roles.id = users.role_id
	WHERE users.id = $1`

	err := repository.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.ImageID,
		&user.Image.Path,
		&user.Role.Name,
		&user.CreatedAt,
	)
	if err != nil {
		log.Println("[error] users.repository.GetByID : failed to execute get data user query", err)
		if err == sql.ErrNoRows {
			err = constants.ErrRecordNotFound
		}
		return users.Domain{}, err
	}

	return *user.ToDomain(), nil
}

func (repository *postgreUsersRepository) GetByEmail(ctx context.Context, email string) (users.Domain, error) {
	user := Users{}
	query := `SELECT users.id, users.name, users.password, users.email, roles.id, roles.name FROM "users"
	LEFT JOIN "roles" ON roles.id = users.role_id
	WHERE "email" = $1`

	err := repository.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.Email,
		&user.RoleID,
		&user.Role.Name,
	)
	if err != nil {
		log.Println("[error] users.repository.GetByEmail : failed to execute get data user query", err)
		if err == sql.ErrNoRows {
			err = constants.ErrRecordNotFound
		}
		return users.Domain{}, err
	}

	return *user.ToDomain(), nil
}

func (repository *postgreUsersRepository) Register(ctx context.Context, data *users.Domain) error {
	record := fromDomain(*data)

	query := `INSERT INTO "users" ("role_id", "name", "email", "password","created_at") VALUES ($1,$2,$3,$4,$5)`
	now := time.Now().Format("2006-01-02T15:04:05")
	_, err := repository.db.ExecContext(ctx, query, 1, record.Name, record.Email, record.Password, now)
	if err != nil {
		log.Println("[error] users.repository.Register : failed to execute register user query", err)
		return err
	}
	return nil
}

func (repository *postgreUsersRepository) Update(ctx context.Context, data *users.Domain, hasPassword bool, id int) error {
	record := fromDomain(*data)

	if hasPassword {
		query := `UPDATE "users" SET "name" = $1, "image_id" = $2, "password" = $3, "updated_at" = $4 WHERE "id" = $5`
		now := time.Now().Format("2006-01-02T15:04:05")
		_, err := repository.db.ExecContext(ctx, query, record.Name, record.ImageID, record.Password, now, id)
		if err != nil {
			log.Println("[error] users.repository.Update : failed to execute update user query", err)
			return err
		}
		return nil
	}

	query := `UPDATE "users" SET "name" = $1, "image_id" = $2, "updated_at" = $3 WHERE "id" = $4`
	now := time.Now().Format("2006-01-02T15:04:05")
	_, err := repository.db.ExecContext(ctx, query, record.Name, record.ImageID, now, id)
	if err != nil {
		log.Println("[error] users.repository.Update : failed to execute update user query", err)
		return err
	}
	return nil
}

func (repository *postgreUsersRepository) SetPassword(ctx context.Context, id int, password string) error {
	query := `UPDATE "users" SET "password" = $1 WHERE "id" = $2`
	_, err := repository.db.ExecContext(ctx, query, password, id)
	if err != nil {
		log.Println("[error] users.repository.SetPassword : failed to execute update password user query", err)
		return err
	}

	return nil
}

func (repository *postgreUsersRepository) Fetch(ctx context.Context, page, perpage int, by, search, sort string) ([]users.Domain, int, error) {
	activeString := ``
	offset := (page - 1) * perpage

	if search != "" {
		activeString = ` WHERE LOWER(users.name) LIKE '%` + search + `%' OR LOWER(users.email) LIKE '%` + search + `%'`
	}

	query := `SELECT users.id, users.name, users.email, im.id, im.path, roles.name, users.created_at FROM "users" 
	LEFT JOIN "images" im ON im.id = users.image_id
	LEFT JOIN "roles" ON roles.id = users.role_id
	` + activeString + `
	ORDER BY ` + by + ` ` + sort + `
	OFFSET $1 LIMIT $2`

	rows, err := repository.db.QueryContext(ctx, query, offset, perpage)
	if err != nil {
		log.Println("[error] users.repository.Fetch : failed to execute fetch users query", err)
		return []users.Domain{}, 0, err
	}

	var total int
	query = `SELECT COUNT(*) FROM "users"` + activeString
	err = repository.db.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		log.Println("[error] users.repository.Fetch : failed to execute count users query", err)
		return []users.Domain{}, 0, err
	}

	var result []users.Domain
	defer rows.Close()
	for rows.Next() {
		temporary := Users{}

		if err := rows.Scan(
			&temporary.ID,
			&temporary.Name,
			&temporary.Email,
			&temporary.ImageID,
			&temporary.Image.Path,
			&temporary.Role.Name,
			&temporary.CreatedAt,
		); err != nil {
			log.Println("[error] users.repository.Fetch : failed to execute row fetch users query", err)
			return []users.Domain{}, 0, err
		}
		result = append(result, *temporary.ToDomain())
	}

	return result, total, nil
}

func (repository *postgreUsersRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM "users" 
	USING "users" AS us
	LEFT OUTER JOIN "roles" AS rl ON us.role_id = rl.id
	WHERE us.id = $1 AND rl.name = 'USER' AND users.id = us.id`

	result, _ := repository.db.ExecContext(ctx, query, id)
	rowAffected, _ := result.RowsAffected()

	if rowAffected == 0 {
		log.Println("[error] users.repository.Delete : failed to execute delete users query", constants.ErrDoNotHavePermission)
		return constants.ErrDoNotHavePermission
	}

	return nil
}
