package request

import (
	"go-articles/modules/articles"
	"go-articles/modules/comments"
)

type Articles struct {
	Title       string `json:"title"  validate:"required" validName:"title"`
	Description string `json:"description"  validate:"required" validName:"description"`
}

type Comments struct {
	Comment string `json:"comment" validate:"required" validName:"comment"`
}

func (request *Articles) ToDomain() *articles.Domain {
	return &articles.Domain{
		Title:       request.Title,
		Description: request.Description,
	}
}

func (request *Comments) ToDomain(userId int) *comments.Domain {
	return &comments.Domain{
		UserID:  userId,
		Comment: request.Comment,
	}
}
