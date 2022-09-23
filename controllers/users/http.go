package users

import (
	"errors"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"strconv"

	"go-articles/constants"
	"go-articles/controllers/users/request"
	"go-articles/controllers/users/response"
	helpers "go-articles/helpers"
	"go-articles/logger"
	"go-articles/modules/images"
	"go-articles/modules/users"
	"go-articles/server/middleware"

	"github.com/JoinVerse/xid"
	echo "github.com/labstack/echo/v4"
)

type UserController struct {
	userUsecase  users.Usecase
	imageUsecase images.Usecase
}

func NewUserController(uc users.Usecase, iu images.Usecase) *UserController {
	return &UserController{
		userUsecase:  uc,
		imageUsecase: iu,
	}
}

func (controller *UserController) Register(c echo.Context) error {
	ctx := c.Request().Context()

	req := request.Users{}
	if err := c.Bind(&req); err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	validateMessage, validate, err := helpers.Validate(&req)

	if validate {
		return helpers.ErrorValidateResponse(c, http.StatusBadRequest, err, validateMessage)
	}

	err = controller.userUsecase.Register(ctx, req.ToDomain())
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}
	return helpers.SuccessResponse(c, http.StatusCreated, nil)
}

func (controller *UserController) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var userLogin request.Users
	if err := c.Bind(&userLogin); err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	accessToken, refreshToken, err := controller.userUsecase.Login(ctx, userLogin.Email, userLogin.Password)

	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, response.TokenFromDomain(accessToken, refreshToken))
}

func (controller *UserController) GetByID(c echo.Context) error {
	ctx := c.Request().Context()

	id := middleware.GetUser(c).ID

	user, err := controller.userUsecase.GetByID(ctx, id)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	logger.Info("Get Data User successfully, id: %d", id)
	return helpers.SuccessResponse(c, http.StatusOK, response.FromDomain(user))
}

func (controller *UserController) ForgotPassword(c echo.Context) error {
	ctx := c.Request().Context()

	req := request.RequestPassword{}
	if err := c.Bind(&req); err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	validateMessage, validate, err := helpers.Validate(&req)
	if validate {
		return helpers.ErrorValidateResponse(c, http.StatusBadRequest, err, validateMessage)
	}

	err = controller.userUsecase.ForgotPassword(ctx, req.Email)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, nil)
}

func (controller *UserController) VerifyTokenForgotPassword(c echo.Context) error {
	ctx := c.Request().Context()

	key := c.Param("key")
	if key == "" {
		return helpers.ErrorResponse(c, http.StatusBadRequest, errors.New(constants.PARAMETER_MUST_BE_FILLED))
	}

	token, tokenExpiredAt, err := controller.userUsecase.VerifyTokenForgotPassword(ctx, key)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, response.VerifyTokenFromDomain(token, tokenExpiredAt))
}

func (controller *UserController) NewPasswordSubmit(c echo.Context) error {
	ctx := c.Request().Context()

	id := middleware.GetUserIdResetPassword(c).ID

	req := request.RequestNewPassword{}
	if err := c.Bind(&req); err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	validateMessage, validate, err := helpers.Validate(&req)
	if validate {
		return helpers.ErrorValidateResponse(c, http.StatusBadRequest, err, validateMessage)
	}

	err = controller.userUsecase.SetForgotPassword(ctx, id, req.Password)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, nil)
}

func (controller *UserController) Update(c echo.Context) error {
	ctx := c.Request().Context()

	id := middleware.GetUser(c).ID

	req := request.UpdateUsers{}
	if err := c.Bind(&req); err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	validateMessage, validate, err := helpers.Validate(&req)
	if validate {
		return helpers.ErrorValidateResponse(c, http.StatusBadRequest, err, validateMessage)
	}

	err = controller.userUsecase.Update(ctx, req.ToDomain(), id)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, nil)
}

func (controller *UserController) UploadPhoto(c echo.Context) error {
	ctx := c.Request().Context()

	id := middleware.GetUser(c).ID
	name := middleware.GetUser(c).Name

	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, 5<<20)
	err := c.Request().ParseMultipartForm(5 << 20)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusInternalServerError, errors.New("file is too large"))
	}

	file, err := c.FormFile("file")
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	src, err := file.Open()
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusInternalServerError, errors.New("can't open file"))
	}

	defer src.Close()

	fileBytes, err := ioutil.ReadAll(src)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusInternalServerError, errors.New("invalid file"))
	}

	filetype := http.DetectContentType(fileBytes)
	if filetype != "image/jpeg" && filetype != "image/jpg" &&
		filetype != "image/gif" && filetype != "image/png" {
		return helpers.ErrorResponse(c, http.StatusInternalServerError, errors.New("invalid file type"))
	}

	fileType := "profile"
	fileName := fileType + "-" + xid.New().String()
	fileEndings, err := mime.ExtensionsByType(filetype)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusInternalServerError, errors.New("can't read File Type"))
	}

	fullFileName := "/" + fileName + fileEndings[len(fileEndings)-1]
	uploadURL := "./static" + fullFileName

	err = os.WriteFile(uploadURL, fileBytes, 0644)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	imageID, err := controller.imageUsecase.Insert(ctx, "/static"+fullFileName, fileType)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	err = controller.userUsecase.Update(ctx, &users.Domain{
		Name:    name,
		ImageID: &imageID,
	}, id)

	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusCreated, nil)
}

func (controller *UserController) DeletePhoto(c echo.Context) error {
	ctx := c.Request().Context()

	id := middleware.GetUser(c).ID

	user, err := controller.userUsecase.GetByID(ctx, id)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	controller.userUsecase.Update(ctx, &users.Domain{
		Name:    user.Name,
		ImageID: nil,
	}, id)

	err = controller.imageUsecase.Delete(ctx, *user.ImageID)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, nil)
}

func (controller *UserController) AdminGetUserByID(c echo.Context) error {
	ctx := c.Request().Context()

	id, _ := strconv.Atoi(c.Param("userId"))

	user, err := controller.userUsecase.GetByID(ctx, id)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, response.AdminUserFromDomain(user))
}

func (controller *UserController) AdminUserFetch(c echo.Context) error {
	ctx := c.Request().Context()

	page, _ := strconv.Atoi(c.QueryParam("page"))
	perpage, _ := strconv.Atoi(c.QueryParam("per_page"))
	search := c.QueryParam("search")
	sort := c.QueryParam("sort")
	by := c.QueryParam("by")

	users, count, err := controller.userUsecase.Fetch(ctx, page, perpage, by, search, sort)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, response.AdminUserFromListDomain(users, count))
}

func (controller *UserController) AdminDeleteUser(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("userId"))

	err := controller.userUsecase.Delete(ctx, id)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, nil)
}
