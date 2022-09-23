package articles

import (
	"context"
	"strings"
)

type ArticlesUsecase struct {
	articleRepository Repository
}

func NewArticleUsecase(ar Repository) Usecase {
	return &ArticlesUsecase{
		articleRepository: ar,
	}
}

func (usecase *ArticlesUsecase) GetByID(ctx context.Context, id int) (Domain, error) {
	article, err := usecase.articleRepository.GetByID(ctx, id)
	if err != nil {
		return Domain{}, err
	}

	return article, nil
}

func (usecase *ArticlesUsecase) Insert(ctx context.Context, data *Domain, userId int) error {
	data.UserID = userId

	err := usecase.articleRepository.Insert(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *ArticlesUsecase) Update(ctx context.Context, data *Domain, id, userId int) error {
	data.UserID = userId

	err := usecase.articleRepository.Update(ctx, data, id)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *ArticlesUsecase) Fetch(ctx context.Context, page, perpage, userId int, by, search, sort string) ([]Domain, int, error) {
	if page <= 0 {
		page = 1
	}

	if perpage <= 0 {
		perpage = 25
	}

	if strings.ToUpper(sort) != "ASC" && strings.ToUpper(sort) != "DESC" {
		sort = "ASC"
	}

	switch by {
	case "title":
		by = `articles."title"`
	default:
		by = `articles."created_at"`
	}

	res, total, err := usecase.articleRepository.Fetch(ctx, page, perpage, userId, by, strings.ToLower(search), sort)
	if err != nil {
		return []Domain{}, 0, err
	}

	return res, total, nil
}

func (usecase *ArticlesUsecase) Delete(ctx context.Context, id int) error {
	err := usecase.articleRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *ArticlesUsecase) Search(ctx context.Context, search string) ([]Domain, error) {
	if search == "" {
		return []Domain{}, nil
	}

	res, err := usecase.articleRepository.Search(ctx, strings.ToLower(search))
	if err != nil {
		return []Domain{}, err
	}

	return res, nil
}
