package mountains

import (
	"go-articles/modules/mountains"
	"time"
)

type Mountains struct {
	ID             int        `db:"id"`
	Name           string     `db:"name"`
	Description    string     `db:"description"`
	About          string     `db:"about"`
	Status         string     `db:"status"`
	Latitude       float64    `db:"latitude"`
	Longitude      float64    `db:"longitude"`
	Province       string     `db:"province"`
	Country        string     `db:"country"`
	Type           string     `db:"type"`
	Height         int        `db:"height"`
	Difficult      string     `db:"difficult"`
	LastEruption   time.Time  `db:"last_eruption"`
	TemperatureMin float64    `db:"temperature_min"`
	TemperatureMax float64    `db:"temperature_max"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
}

func (record *Mountains) ToDomain() *mountains.Domain {
	return &mountains.Domain{
		ID:             record.ID,
		Name:           record.Name,
		Description:    record.Description,
		About:          record.About,
		Status:         record.Status,
		Latitude:       record.Latitude,
		Longitude:      record.Longitude,
		Province:       record.Province,
		Country:        record.Country,
		Type:           record.Type,
		Height:         record.Height,
		Difficult:      record.Difficult,
		LastEruption:   record.LastEruption,
		TemperatureMin: record.TemperatureMin,
		TemperatureMax: record.TemperatureMax,
		CreatedAt:      record.CreatedAt,
		UpdatedAt:      record.UpdatedAt,
	}
}

func fromDomain(record mountains.Domain) *Mountains {
	return &Mountains{
		ID:             record.ID,
		Name:           record.Name,
		Description:    record.Description,
		About:          record.About,
		Status:         record.Status,
		Latitude:       record.Latitude,
		Longitude:      record.Longitude,
		Province:       record.Province,
		Country:        record.Country,
		Type:           record.Type,
		Height:         record.Height,
		Difficult:      record.Difficult,
		LastEruption:   record.LastEruption,
		TemperatureMin: record.TemperatureMin,
		TemperatureMax: record.TemperatureMax,
		CreatedAt:      record.CreatedAt,
		UpdatedAt:      record.UpdatedAt,
	}
}
