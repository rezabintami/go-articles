package helpers

import (
	echo "github.com/labstack/echo/v4"
)

type Meta struct {
	Message string `json:"message"`
	Errors  string `json:"error,omitempty"`
}
type BaseResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

func SuccessResponse(c echo.Context, status int, param interface{}) error {
	response := BaseResponse{}
	response.Meta.Message = "success"
	response.Data = param

	return c.JSON(status, response)
}

func ErrorResponse(c echo.Context, status int, err error) error {
	response := BaseResponse{}
	response.Meta.Message = "something wrong"
	response.Meta.Errors = err.Error()

	return c.JSON(status, response)
}

func ErrorValidateResponse(c echo.Context, status int, err error, param interface{}) error {
	response := BaseResponse{}
	response.Meta.Message = "something wrong"
	response.Data = param

	return c.JSON(status, response)
}
