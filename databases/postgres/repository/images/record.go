package images

import (
	"go-articles/modules/images"
	"time"
)

type Images struct {
	ID        *int      `db:"id"`
	Path      *string   `db:"path"`
	Type      *string   `db:"type"`
	CreatedAt time.Time `db:"created_at"`
}

func (record *Images) ToDomain() *images.Domain {
	return &images.Domain{
		ID:        record.ID,
		Path:      record.Path,
		Type:      record.Type,
		CreatedAt: record.CreatedAt,
	}
}

func fromDomain(domain images.Domain) *Images {
	return &Images{
		ID:        domain.ID,
		Path:      domain.Path,
		Type:      domain.Type,
		CreatedAt: domain.CreatedAt,
	}
}
