package mountains

import (
	"context"
	"database/sql"
	"go-articles/constants"
	"go-articles/modules/mountains"
	"log"
	"time"
)

type postgreMountainsRepository struct {
	db *sql.DB
}

func NewPostgreMountainsRepository(db *sql.DB) mountains.Repository {
	return &postgreMountainsRepository{
		db: db,
	}
}

// Delete implements mountains.Repository
func (repository *postgreMountainsRepository) Delete(ctx context.Context, mountainId int) error {
	query := `UPDATE "mountains" SET "deleted_at" = $1 WHERE "id" = $2`
	now := time.Now().Format("2006-01-02T15:04:05")
	_, err := repository.db.ExecContext(ctx, query, now, mountainId)

	if err != nil {
		log.Println("[error] mountains.repository.Delete : failed to execute delete mountain query", err)
		return err
	}

	return nil
}

// Fetch implements mountains.Repository
func (repository *postgreMountainsRepository) Fetch(ctx context.Context, page, perpage int, by, search, sort string) ([]mountains.Domain, int, error) {
	activeString := ``
	offset := (page - 1) * perpage

	if search != "" {
		activeString = ` AND (LOWER(mountains.name) LIKE '%` + search + `%' OR LOWER(mountains.province) LIKE '%` + search + `%' 
		OR LOWER(mountains.country) LIKE '%` + search + `%')`
	}

	query := `SELECT mountains.id, mountains.name, mountains.description, mountains.status, mountains.province, mountains.country, mountains.type, mountains.height,
	mountains.difficult, mountains.last_eruption, mountains.temperature_min, mountains.temperature_max, mountains.updated_at FROM "mountains"
	WHERE "deleted_at" is NULL` + activeString + `
	ORDER BY ` + by + ` ` + sort + `
	OFFSET $1 LIMIT $2`

	rows, err := repository.db.QueryContext(ctx, query, offset, perpage)
	if err != nil {
		log.Println("[error] mountain.repository.Fetch : failed to execute fetch mountain query", err)
		return []mountains.Domain{}, 0, err
	}

	var total int
	query = `SELECT COUNT(*) FROM "mountains" WHERE "deleted_at" is NULL` + activeString
	err = repository.db.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		log.Println("[error] mountains.repository.Fetch : failed to execute count mountain query", err)
		return []mountains.Domain{}, 0, err
	}

	var result []mountains.Domain
	defer rows.Close()
	for rows.Next() {
		temporary := Mountains{}

		if err := rows.Scan(
			&temporary.ID, &temporary.Name, &temporary.Description, &temporary.Status, &temporary.Province, &temporary.Country, &temporary.Type, &temporary.Height,
			&temporary.Difficult, &temporary.LastEruption, &temporary.TemperatureMin, &temporary.TemperatureMax, &temporary.UpdatedAt,
		); err != nil {
			log.Println("[error] mountains.repository.Fetch : failed to execute row fetch mountain query", err)
			return []mountains.Domain{}, 0, err
		}
		result = append(result, *temporary.ToDomain())
	}

	return result, total, nil
}

// GetByID implements mountains.Repository
func (repository *postgreMountainsRepository) GetByID(ctx context.Context, mountainId int) (mountains.Domain, error) {
	mountain := Mountains{}
	query := `SELECT mountains.id, mountains.name, mountains.description, mountains.about, mountains.status,
	mountains.latitude, mountains.longitude, mountains.province, mountains.country, mountains.type, mountains.height,
	mountains.difficult, mountains.last_eruption, mountains.temperature_min, mountains.temperature_max, mountains.created_at,
	mountains.updated_at FROM "mountains" WHERE mountains.id = $1 AND "deleted_at" is NULL`

	err := repository.db.QueryRowContext(ctx, query, mountainId).Scan(
		&mountain.ID, &mountain.Name, &mountain.Description, &mountain.About, &mountain.Status,
		&mountain.Latitude, &mountain.Longitude, &mountain.Province, &mountain.Country, &mountain.Type, &mountain.Height,
		&mountain.Difficult, &mountain.LastEruption, &mountain.TemperatureMin, &mountain.TemperatureMax, &mountain.CreatedAt, &mountain.UpdatedAt,
	)

	if err != nil {
		log.Println("[error] mountains.repository.GetByID : failed to execute get data mountain query", err)
		if err == sql.ErrNoRows {
			err = constants.ErrRecordNotFound
		}
		return mountains.Domain{}, err
	}

	return *mountain.ToDomain(), nil
}

// Insert implements mountains.Repository
func (repository *postgreMountainsRepository) Insert(ctx context.Context, data *mountains.Domain) error {
	record := fromDomain(*data)

	query := `INSERT INTO "mountains" ("name", "description", "about", "status", "latitude", "longitude", "province", 
	"country", "type", "height", "difficult", "last_eruption", "temperature_min", "temperature_max", "created_at", "updated_at")
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`
	now := time.Now().Format("2006-01-02T15:04:05")
	_, err := repository.db.ExecContext(ctx, query, record.Name, record.Description, record.About, record.Status, record.Latitude,
		record.Longitude, record.Province, record.Country, record.Type, record.Height, record.Difficult, record.LastEruption,
		record.TemperatureMin, record.TemperatureMax, now, now)
	if err != nil {
		log.Println("[error] mountains.repository.Insert : failed to execute mountain query", err)
		return err
	}

	return nil
}

// Search implements mountains.Repository
func (repository *postgreMountainsRepository) Search(ctx context.Context, search string) ([]mountains.Domain, error) {
	query := `SELECT mountains.id, mountains.name, mountains.description, mountains.about, mountains.status,
	mountains.latitude, mountains.longitude, mountains.province, mountains.country, mountains.type, mountains.height,
	mountains.difficult, mountains.last_eruption, mountains.temperature_min, mountains.temperature_max, mountains.created_at,
	mountains.updated_at FROM "mountains" WHERE "deleted_at" is NULL AND LOWER(mountains.name) LIKE '%` + search + `%'`

	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		log.Println("[error] mountains.repository.Search : failed to execute search mountain query", err)
		return []mountains.Domain{}, err
	}

	var result []mountains.Domain
	defer rows.Close()
	for rows.Next() {
		temporary := Mountains{}
		if err := rows.Scan(
			&temporary.ID, &temporary.Name, &temporary.Description, &temporary.About, &temporary.Status,
			&temporary.Latitude, &temporary.Longitude, &temporary.Province, &temporary.Country, &temporary.Type, &temporary.Height,
			&temporary.Difficult, &temporary.LastEruption, &temporary.TemperatureMin, &temporary.TemperatureMax, &temporary.CreatedAt, &temporary.UpdatedAt,
		); err != nil {
			log.Println("[error] mountains.repository.Search : failed to execute row search mountain query", err)
			return []mountains.Domain{}, err
		}
		result = append(result, *temporary.ToDomain())
	}

	return result, nil

}

// Update implements mountains.Repository
func (repository *postgreMountainsRepository) Update(ctx context.Context, data *mountains.Domain, mountainId int) error {
	record := fromDomain(*data)

	query := `UPDATE "mountains" SET "name" = $1, "description" = $2, "about" = $3, "status" = $4, "latitude" = $5,
	"longitude" = $6, "province" = $7, "country" = $8, "type" = $9, "height" = $10, "difficult" = $11, "last_eruption" = $12,
	"temperature_min" =  $13, "temperature_max" = $14, "updated_at" = $15 WHERE "id" = $16 AND "deleted_at" is NULL`
	now := time.Now().Format("2006-01-02T15:04:05")
	_, err := repository.db.ExecContext(ctx, query, record.Name, record.Description, record.About, record.Status, record.Latitude,
		record.Longitude, record.Province, record.Country, record.Type, record.Height, record.Difficult, record.LastEruption,
		record.TemperatureMin, record.TemperatureMax, now, mountainId)
	if err != nil {
		log.Println("[error] mountains.repository.Update : failed to execute update mountain query", err)
		return err
	}

	return nil
}
