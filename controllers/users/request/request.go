package request

import (
	"go-articles/modules/users"
)

type Users struct {
	Name     string `json:"name" validate:"required" validName:"name"`
	Password string `json:"password,omitempty" validate:"required" validName:"password"`
	Email    string `json:"email" validate:"required,email,max=100" validName:"email"`
}

type UpdateUsers struct {
	Name     string `json:"name" validate:"required" validName:"name"`
	Password string `json:"password,omitempty"`
}

type RequestPassword struct {
	Email string `json:"email" validate:"required,email,max=100" validName:"email"`
}

type RequestNewPassword struct {
	Password        string `json:"password" validate:"required" validName:"password"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password" validName:"confirm_password"`
}

func (request *Users) ToDomain() *users.Domain {
	return &users.Domain{
		Name:     request.Name,
		Password: request.Password,
		Email:    request.Email,
	}
}

func (request *UpdateUsers) ToDomain() *users.Domain {
	return &users.Domain{
		Name:     request.Name,
		Password: request.Password,
	}
}