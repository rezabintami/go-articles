package comments

import (
	"context"
	"go-articles/modules/users"
	"time"
)

type Domain struct {
	ID        int
	UserID    int
	User      users.Domain
	Comment   string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type Usecase interface {
	GetByArticleID(ctx context.Context, articlesId int) ([]Domain, error)
	Insert(ctx context.Context, data *Domain, articlesId int) error
	Update(ctx context.Context, data *Domain, commentId int) error
	Delete(ctx context.Context, articleId, commentId, userId int) error
}

type Repository interface {
	GetByArticleID(ctx context.Context, articlesId int) ([]Domain, error)
	Insert(ctx context.Context, data *Domain, articlesId int) error
	Update(ctx context.Context, data *Domain, commentId int) error
	Delete(ctx context.Context, articleId, commentId, userId int) error
}
