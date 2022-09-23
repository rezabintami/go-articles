package roles

import (
	"context"
	"time"
)

type Domain struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type Usecase interface {
	Fetch(ctx context.Context, start, last int) ([]Domain, int, error)
	GetByID(ctx context.Context, id int) (Domain, error)
	Insert(ctx context.Context, data *Domain) error
	Update(ctx context.Context, data *Domain, id int) error
	Delete(ctx context.Context, id int) error
}

type Repository interface {
	Fetch(ctx context.Context, start, last int) ([]Domain, int, error)
	GetByID(ctx context.Context, id int) (Domain, error)
	Insert(ctx context.Context, data *Domain) error
	Update(ctx context.Context, data *Domain, id int) error
	Delete(ctx context.Context, id int) error
}
