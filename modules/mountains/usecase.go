package mountains

import (
	"context"
	"strings"
)

type MountainsUsecase struct {
	mountainRepository Repository
}

func NewMountainUsecase(mr Repository) Usecase {
	return &MountainsUsecase{
		mountainRepository: mr,
	}
}

// Delete implements Usecase
func (usecase *MountainsUsecase) Delete(ctx context.Context, mountainId int) error {
	err := usecase.mountainRepository.Delete(ctx, mountainId)
	if err != nil {
		return err
	}

	return nil
}

// Fetch implements Usecase
func (usecase *MountainsUsecase) Fetch(ctx context.Context, page, perpage int, by, search, sort string) ([]Domain, int, error) {
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
	case "name":
		by = `mountains."name"`
	case "province":
		by = `mountains."province"`
	case "country":
		by = `mountains."country"`
	case "type":
		by = `mountains."type"`
	case "height":
		by = `mountains."height"`
	default:
		by = `mountains."created_at"`
	}

	res, total, err := usecase.mountainRepository.Fetch(ctx, page, perpage, by, strings.ToLower(search), sort)
	if err != nil {
		return []Domain{}, 0, err
	}

	return res, total, nil
}

// GetByID implements Usecase
func (usecase *MountainsUsecase) GetByID(ctx context.Context, mountainId int) (Domain, error) {
	mountain, err := usecase.mountainRepository.GetByID(ctx, mountainId)
	if err != nil {
		return Domain{}, err
	}

	return mountain, nil
}

// Insert implements Usecase
func (usecase *MountainsUsecase) Insert(ctx context.Context, data *Domain) error {
	err := usecase.mountainRepository.Insert(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// Search implements Usecase
func (usecase *MountainsUsecase) Search(ctx context.Context, search string) ([]Domain, error) {
	if search == "" {
		return []Domain{}, nil
	}

	res, err := usecase.mountainRepository.Search(ctx, strings.ToLower(search))
	if err != nil {
		return []Domain{}, err
	}

	return res, nil
}

// Update implements Usecase
func (usecase *MountainsUsecase) Update(ctx context.Context, data *Domain, mountainId int) error {
	err := usecase.mountainRepository.Update(ctx, data, mountainId)
	if err != nil {
		return err
	}

	return nil
}
