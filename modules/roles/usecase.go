package roles

import (
	"context"
	"go-articles/constants"
	"strings"
)

type RolesUsecase struct {
	roleRepository Repository
}

func NewRoleUsecase(ru Repository) Usecase {
	return &RolesUsecase{
		roleRepository: ru,
	}
}

func (usecase *RolesUsecase) GetByID(ctx context.Context, id int) (Domain, error) {
	role, err := usecase.roleRepository.GetByID(ctx, id)
	if err != nil {
		return Domain{}, err
	}

	return role, nil
}

func (usecase *RolesUsecase) Insert(ctx context.Context, data *Domain) error {
	err := usecase.roleRepository.Insert(ctx, data)

	if strings.Contains(err.Error(), constants.DUPLICATE_DATA_UNIQUE) {
		return constants.ErrDuplicateData
	}

	if err != nil {
		return err
	}

	return nil
}

func (usecase *RolesUsecase) Update(ctx context.Context, data *Domain, id int) error {
	err := usecase.roleRepository.Update(ctx, data, id)

	if strings.Contains(err.Error(), constants.DUPLICATE_DATA_UNIQUE) {
		return constants.ErrDuplicateData
	}

	if err != nil {
		return err
	}

	return nil
}

func (usecase *RolesUsecase) Delete(ctx context.Context, id int) error {
	err := usecase.roleRepository.Delete(ctx, id)

	if err != nil {
		return err
	}

	return nil
}

func (usecase *RolesUsecase) Fetch(ctx context.Context, page, perpage int) (result []Domain, count int, err error) {
	if page <= 0 {
		page = 1
	}
	if perpage <= 0 {
		perpage = 25
	}

	result, count, err = usecase.roleRepository.Fetch(ctx, page, perpage)
	if err != nil {
		return []Domain{}, 0, err
	}

	return result, count, nil
}
