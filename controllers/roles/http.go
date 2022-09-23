package roles

import (
	"go-articles/controllers/roles/request"
	"go-articles/controllers/roles/response"
	"go-articles/helpers"
	"go-articles/modules/roles"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RoleController struct {
	roleUsecase roles.Usecase
}

func NewRoleController(ru roles.Usecase) *RoleController {
	return &RoleController{
		roleUsecase: ru,
	}
}

func (controller *RoleController) GetByID(c echo.Context) error {
	ctx := c.Request().Context()

	id, _ := strconv.Atoi(c.Param("roleId"))
	role, err := controller.roleUsecase.GetByID(ctx, id)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, response.FromDomain(role))
}

func (controller *RoleController) Insert(c echo.Context) error {
	ctx := c.Request().Context()

	req := request.Roles{}
	if err := c.Bind(&req); err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	validateMessage, validate, err := helpers.Validate(&req)

	if validate {
		return helpers.ErrorValidateResponse(c, http.StatusBadRequest, err, validateMessage)
	}

	err = controller.roleUsecase.Insert(ctx, req.ToDomain())
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}
	return helpers.SuccessResponse(c, http.StatusCreated, nil)
}

func (controller *RoleController) Update(c echo.Context) error {
	ctx := c.Request().Context()

	id, _ := strconv.Atoi(c.Param("roleId"))

	req := request.Roles{}
	if err := c.Bind(&req); err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	validateMessage, validate, err := helpers.Validate(&req)

	if validate {
		return helpers.ErrorValidateResponse(c, http.StatusBadRequest, err, validateMessage)
	}

	err = controller.roleUsecase.Update(ctx, req.ToDomain(), id)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}
	return helpers.SuccessResponse(c, http.StatusOK, nil)
}

func (controller *RoleController) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	id, _ := strconv.Atoi(c.Param("roleId"))

	err := controller.roleUsecase.Delete(ctx, id)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, nil)
}

func (controller *RoleController) Fetch(c echo.Context) error {
	ctx := c.Request().Context()

	page, _ := strconv.Atoi(c.QueryParam("page"))
	perpage, _ := strconv.Atoi(c.QueryParam("per_page"))

	roles, count, err := controller.roleUsecase.Fetch(ctx, page, perpage)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, response.FromListDomain(roles, count))
}
