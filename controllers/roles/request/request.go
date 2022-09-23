package request

import "go-articles/modules/roles"

type Roles struct {
	Name string `json:"name" validate:"required" validName:"name"`
}

func (request *Roles) ToDomain() *roles.Domain {
	return &roles.Domain{
		Name: request.Name,
	}
}
