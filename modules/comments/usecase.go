package comments

import (
	"context"
)

type CommentUsecase struct {
	commentRepository Repository
}

func NewCommentUsecase(cr Repository) Usecase {
	return &CommentUsecase{
		commentRepository: cr,
	}
}

// Delete implements Usecase
func (usecase *CommentUsecase) Delete(ctx context.Context, articleId, commentId, userId int) error {
	err := usecase.commentRepository.Delete(ctx, articleId, commentId, userId)

	return err
}

// GetByArticleID implements Usecase
func (usecase *CommentUsecase) GetByArticleID(ctx context.Context, articlesId int) ([]Domain, error) {
	data, err := usecase.commentRepository.GetByArticleID(ctx, articlesId)
	if err != nil {
		return []Domain{}, err
	}

	return data, nil
}

// Insert implements Usecase
func (usecase *CommentUsecase) Insert(ctx context.Context, data *Domain, articlesId int) error {
	err := usecase.commentRepository.Insert(ctx, data, articlesId)

	return err
}

// Update implements Usecase
func (usecase *CommentUsecase) Update(ctx context.Context, data *Domain, commentId int) error {
	err := usecase.commentRepository.Update(ctx, data, commentId)

	return err
}
