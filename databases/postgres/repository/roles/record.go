package roles

import (
	"go-articles/modules/roles"
	"time"
)

type Roles struct {
	ID        int        `db:"id"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func (record *Roles) ToDomain() *roles.Domain {
	return &roles.Domain{
		ID:        record.ID,
		Name:      record.Name,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
	}
}

func fromDomain(domain roles.Domain) *Roles {
	return &Roles{
		ID:        domain.ID,
		Name:      domain.Name,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}
