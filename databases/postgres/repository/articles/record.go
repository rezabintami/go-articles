package articles

import (
	"go-articles/databases/postgres/repository/images"
	"go-articles/databases/postgres/repository/roles"
	"go-articles/databases/postgres/repository/users"
	"go-articles/modules/articles"
	"time"
)

type Articles struct {
	ID          int  `db:"id"`
	UserID      int  `db:"user_id"`
	ImageID     *int `db:"image_id"`
	User        users.Users
	Role        roles.Roles
	Image       images.Images
	Title       string     `db:"title"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

func (record *Articles) ToDomain() *articles.Domain {
	return &articles.Domain{
		ID:          record.ID,
		UserID:      record.UserID,
		Title:       record.Title,
		ImageID:     record.ImageID,
		User:        *record.User.ToDomain(),
		Role:        *record.Role.ToDomain(),
		Image:       *record.Image.ToDomain(),
		Description: record.Description,
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   record.UpdatedAt,
	}
}

func fromDomain(domain articles.Domain) *Articles {
	return &Articles{
		ID:          domain.ID,
		UserID:      domain.UserID,
		ImageID:     domain.ImageID,
		Title:       domain.Title,
		Description: domain.Description,
		CreatedAt:   domain.CreatedAt,
		UpdatedAt:   domain.UpdatedAt,
	}
}
