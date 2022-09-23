package response

import (
	"go-articles/modules/users"
	"time"
)

type Users struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	ImagePath *string `json:"photo_path"`
}

type AdminUser struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	ImagePath *string   `json:"photo_path"`
	RoleName  string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminUserList struct {
	Users *[]AdminUser `json:"users"`
	Total int          `json:"total"`
}
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type VerifyToken struct {
	Token       string `json:"token"`
	ExpiredDate string `json:"expired_date"`
}

func FromDomain(domain users.Domain) Users {
	return Users{
		ID:        domain.ID,
		Name:      domain.Name,
		Email:     domain.Email,
		ImagePath: domain.Image.Path,
	}
}

func AdminUserFromDomain(domain users.Domain) AdminUser {
	return AdminUser{
		ID:        domain.ID,
		Name:      domain.Name,
		Email:     domain.Email,
		ImagePath: domain.Image.Path,
		RoleName:  domain.Role.Name,
		CreatedAt: domain.CreatedAt,
	}
}

func TokenFromDomain(accessToken, refreshToken string) Token {
	return Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func VerifyTokenFromDomain(token, expiredDate string) VerifyToken {
	return VerifyToken{
		Token:       token,
		ExpiredDate: expiredDate,
	}
}

func AdminUserFromListDomain(domain []users.Domain, count int) *AdminUserList {
	userList := []AdminUser{}
	for _, value := range domain {
		user := AdminUser{
			ID:        value.ID,
			Name:      value.Name,
			Email:     value.Email,
			ImagePath: value.Image.Path,
			RoleName:  value.Role.Name,
			CreatedAt: value.CreatedAt,
		}
		userList = append(userList, user)
	}

	result := AdminUserList{}
	result.Users = &userList
	result.Total = count
	return &result
}
