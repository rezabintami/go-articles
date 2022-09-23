package mountains

import (
	"context"
	"time"
)

type Domain struct {
	ID             int
	Name           string
	Description    string
	About          string
	Status         string
	Latitude       float64
	Longitude      float64
	Province       string
	Country        string
	Type           string
	Height         int
	Difficult      string
	LastEruption   time.Time
	TemperatureMin float64
	TemperatureMax float64
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}

type Usecase interface {
	GetByID(ctx context.Context, mountainId int) (Domain, error)
	Insert(ctx context.Context, data *Domain) error
	Update(ctx context.Context, data *Domain, mountainId int) error
	Fetch(ctx context.Context, start, last int, by, search, sort string) ([]Domain, int, error)
	Delete(ctx context.Context, mountainId int) error
	Search(ctx context.Context, search string) ([]Domain, error)
}

type Repository interface {
	GetByID(ctx context.Context, mountainId int) (Domain, error)
	Insert(ctx context.Context, data *Domain) error
	Update(ctx context.Context, data *Domain, mountainId int) error
	Fetch(ctx context.Context, start, last int, by, search, sort string) ([]Domain, int, error)
	Delete(ctx context.Context, mountainId int) error
	Search(ctx context.Context, search string) ([]Domain, error)
}
