package comments

import (
	"go-articles/databases/postgres/repository/users"
	"go-articles/modules/comments"
	"time"
)

type Comments struct {
	ID        int `db:"id"`
	UserID    int `db:"user_id"`
	User      users.Users
	Comment   string     `db:"comment"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (record *Comments) ToDomain() *comments.Domain {
	return &comments.Domain{
		ID:        record.ID,
		UserID:    record.UserID,
		User:      *record.User.ToDomain(),
		Comment:   record.Comment,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
	}
}
