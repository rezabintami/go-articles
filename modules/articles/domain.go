package articles

import (
	"context"
	"go-articles/modules/comments"
	"go-articles/modules/images"
	"go-articles/modules/roles"
	"go-articles/modules/users"
	"time"
)

type Domain struct {
	ID          int
	UserID      int
	ImageID     *int
	User		users.Domain
	Role        roles.Domain
	Image       images.Domain
	Title       string
	Description string
	Comments    []comments.Domain
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type Usecase interface {
	Insert(ctx context.Context, data *Domain, userId int) error
	GetByID(ctx context.Context, id int) (Domain, error)
	Update(ctx context.Context, data *Domain, id, userId int) error
	Fetch(ctx context.Context, start, last, userId int, by, search, sort string) ([]Domain, int, error)
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, search string) ([]Domain, error)
}

type Repository interface {
	GetByID(ctx context.Context, id int) (Domain, error)
	Insert(ctx context.Context, data *Domain) error
	Update(ctx context.Context, data *Domain, id int) error
	Fetch(ctx context.Context, start, last, userId int, by, search, sort string) ([]Domain, int, error)
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, search string) ([]Domain, error)
}
