package users

import (
	"context"
	"go-articles/modules/images"
	"go-articles/modules/roles"
	"time"
)

type Domain struct {
	ID        int
	Name      string
	Password  string
	Email     string
	ImageID   *int
	Role      roles.Domain
	Image     images.Domain
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type Usecase interface {
	Login(ctx context.Context, email, password string) (string, string, error)
	Register(ctx context.Context, data *Domain) error
	GetByID(ctx context.Context, id int) (Domain, error)
	Fetch(ctx context.Context, start, last int, by, search, sort string) ([]Domain, int, error)
	ForgotPassword(ctx context.Context, email string) error
	VerifyTokenForgotPassword(ctx context.Context, key string) (string, string, error)
	Update(ctx context.Context, data *Domain, id int) error
	SetForgotPasswordRedis(key string, value string, duration time.Duration) error
	GetForgotPasswordRedis(key string) (string, error)
	DeleteForgotPasswordRedis(key string) error
	SetSessionRedis(sessionName string, id int, time time.Time) error
	SetForgotPassword(ctx context.Context, id int, password string) error
	Delete(ctx context.Context, id int) error
}

type Repository interface {
	GetByID(ctx context.Context, id int) (Domain, error)
	Fetch(ctx context.Context, start, last int, by, search, sort string) ([]Domain, int, error)
	GetByEmail(ctx context.Context, email string) (Domain, error)
	Register(ctx context.Context, data *Domain) error
	Update(ctx context.Context, data *Domain, hasPassword bool, id int) error
	SetPassword(ctx context.Context, id int, password string) error
	Delete(ctx context.Context, id int) error
}
