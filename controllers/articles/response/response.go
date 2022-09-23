package response

import (
	"go-articles/modules/articles"
	"go-articles/modules/comments"
	"time"
)

type Articles struct {
	ID          int       `json:"id"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	ImagePath   *string   `json:"image_path"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
type ArticlesList struct {
	Articles *[]Articles `json:"articles"`
	Total    int         `json:"total"`
}

type ArticlesComment struct {
	ID          int         `json:"id"`
	Author      string      `json:"author"`
	Title       string      `json:"title"`
	ImagePath   *string     `json:"image_path"`
	Description string      `json:"description"`
	Comments    *[]Comments `json:"comments"`
	CreatedAt   time.Time   `json:"created_at"`
}

type Comments struct {
	ID        int       `json:"id"`
	Author    string    `json:"author"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

func CommentFromDomain(domain *[]comments.Domain) (result *[]Comments) {
	if domain != nil {
		result = &[]Comments{}
		for _, value := range *domain {
			comment := Comments{
				ID:        value.ID,
				Author:    value.User.Name,
				Comment:   value.Comment,
				CreatedAt: value.CreatedAt,
			}
			*result = append(*result, comment)
		}
	}

	return result
}

func FromDomain(domain articles.Domain) ArticlesComment {
	return ArticlesComment{
		ID:          domain.ID,
		Author:      domain.Role.Name,
		Title:       domain.Title,
		ImagePath:   domain.Image.Path,
		Description: domain.Description,
		Comments:    CommentFromDomain(&domain.Comments),
		CreatedAt:   domain.CreatedAt,
	}
}

func FetchFromListDomain(domain []articles.Domain, count int, author string) *ArticlesList {
	articlesList := []Articles{}
	for _, value := range domain {
		article := Articles{
			ID:          value.ID,
			Title:       value.Title,
			Author:      author,
			ImagePath:   value.Image.Path,
			Description: value.Description,
			CreatedAt:   value.CreatedAt,
		}
		articlesList = append(articlesList, article)
	}

	result := ArticlesList{}
	result.Articles = &articlesList
	result.Total = count
	return &result
}

func FromListDomain(domain []articles.Domain) *[]Articles {
	result := []Articles{}
	for _, value := range domain {
		article := Articles{
			ID:          value.ID,
			Title:       value.Title,
			Author:      value.User.Name,
			ImagePath:   value.Image.Path,
			Description: value.Description,
			CreatedAt:   value.CreatedAt,
		}
		result = append(result, article)
	}

	return &result
}
