package images

import (
	"context"
)

type ImagesUsecase struct {
	imageRepository Repository
}

func NewImageUsecase(ir Repository) Usecase {
	return &ImagesUsecase{
		imageRepository: ir,
	}
}

func (usecase *ImagesUsecase) Insert(ctx context.Context, path string, types string) (int, error) {
	imageID, err := usecase.imageRepository.Insert(ctx, &Domain{
		Path: &path,
		Type: &types,
	})

	if err != nil {
		return 0, err
	}

	return imageID, nil
}


func (usecase *ImagesUsecase) Delete(ctx context.Context, userId int) error {
	err := usecase.imageRepository.Delete(ctx, userId)
	
	if err != nil {
		return err
	}

	return nil
}
