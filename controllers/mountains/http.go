package mountains

import (
	"go-articles/controllers/mountains/request"
	"go-articles/controllers/mountains/response"
	"go-articles/helpers"
	"go-articles/modules/mountains"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type MountainController struct {
	mountainUsecase mountains.Usecase
}

func NewMountainController(mu mountains.Usecase) *MountainController {
	return &MountainController{
		mountainUsecase: mu,
	}
}

func (controller *MountainController) SearchMountains(c echo.Context) error {
	ctx := c.Request().Context()

	search := c.QueryParam("search")

	mountains, err := controller.mountainUsecase.Search(ctx, search)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, response.FromListDomain(mountains))
}

func (controller *MountainController) GetByID(c echo.Context) error {
	ctx := c.Request().Context()

	id, _ := strconv.Atoi(c.Param("mountainId"))
	mountain, err := controller.mountainUsecase.GetByID(ctx, id)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, response.FromDomain(mountain))
}

func (controller *MountainController) Insert(c echo.Context) error {
	ctx := c.Request().Context()

	request := request.Mountains{}
	if err := c.Bind(&request); err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	validateMessage, validate, err := helpers.Validate(&request)

	if validate {
		return helpers.ErrorValidateResponse(c, http.StatusBadRequest, err, validateMessage)
	}

	err = controller.mountainUsecase.Insert(ctx, request.ToDomain())
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusCreated, nil)
}

func (controller *MountainController) Update(c echo.Context) error {
	ctx := c.Request().Context()

	request := request.Mountains{}
	id, _ := strconv.Atoi(c.Param("mountainId"))

	if err := c.Bind(&request); err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	validateMessage, validate, err := helpers.Validate(&request)

	if validate {
		return helpers.ErrorValidateResponse(c, http.StatusBadRequest, err, validateMessage)
	}

	err = controller.mountainUsecase.Update(ctx, request.ToDomain(), id)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, nil)
}

func (controller *MountainController) Fetch(c echo.Context) error {
	ctx := c.Request().Context()

	page, _ := strconv.Atoi(c.QueryParam("page"))
	perpage, _ := strconv.Atoi(c.QueryParam("per_page"))
	search := c.QueryParam("search")
	sort := c.QueryParam("sort")
	by := c.QueryParam("by")

	articles, count, err := controller.mountainUsecase.Fetch(ctx, page, perpage, by, search, sort)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, response.FetchFromListDomain(articles, count))
}

func (controller *MountainController) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	id, _ := strconv.Atoi(c.Param("mountainId"))
	err := controller.mountainUsecase.Delete(ctx, id)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, nil)
}
