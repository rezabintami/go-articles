package images

import (
	"context"
	"time"
)

type Domain struct {
	ID        *int
	Path      *string
	Type      *string
	CreatedAt time.Time
}

type Usecase interface {
	Insert(ctx context.Context, path, types string) (int, error)
	Delete(ctx context.Context, userId int) error
}

type Repository interface {
	Insert(ctx context.Context, data *Domain) (int, error)
	Delete(ctx context.Context, userId int) error
}
