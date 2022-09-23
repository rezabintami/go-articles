package users

import (
	"go-articles/databases/postgres/repository/images"
	"go-articles/databases/postgres/repository/roles"
	"go-articles/modules/users"
	"time"
)

type Users struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Password  string `db:"password"`
	Email     string `db:"email"`
	RoleID    int    `db:"role_id"`
	ImageID   *int   `db:"image_id"`
	Role      roles.Roles
	Image     images.Images
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func (record *Users) ToDomain() *users.Domain {
	return &users.Domain{
		ID:        record.ID,
		Name:      record.Name,
		Password:  record.Password,
		Email:     record.Email,
		ImageID:   record.ImageID,
		Role:      *record.Role.ToDomain(),
		Image:     *record.Image.ToDomain(),
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
	}
}

func fromDomain(domain users.Domain) *Users {
	return &Users{
		ID:        domain.ID,
		Name:      domain.Name,
		ImageID:   domain.ImageID,
		Password:  domain.Password,
		Email:     domain.Email,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}
