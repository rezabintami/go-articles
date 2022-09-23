package middleware

import (
	"errors"
	"go-articles/helpers"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

func RoleValidation(roles ...string) echo.MiddlewareFunc {
	return func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := GetUser(c)

			for _, role := range roles {
				if claims.Role == role {
					return hf(c)
				}
			}
			return helpers.ErrorResponse(c, http.StatusForbidden, errors.New("forbidden roles"))
		}
	}
}
