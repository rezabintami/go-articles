package response

import (
	"go-articles/modules/roles"
	"time"
)

type Roles struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type RolesList struct {
	Roles *[]Roles `json:"roles"`
	Total int      `json:"total"`
}

func FromDomain(domain roles.Domain) Roles {
	return Roles{
		ID:        domain.ID,
		Name:      domain.Name,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func FromListDomain(domain []roles.Domain, count int)  *RolesList {
	rolesList := []Roles{}
	for _, value := range domain {
		role := Roles{
			ID:        value.ID,
			Name:      value.Name,
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		}

		rolesList = append(rolesList, role)
	}
	result := RolesList{}
	result.Roles = &rolesList
	result.Total = count
	return &result
}
