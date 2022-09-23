package articles

import (
	"errors"
	"go-articles/constants"
	"go-articles/controllers/articles/request"
	"go-articles/controllers/articles/response"
	"go-articles/helpers"
	"go-articles/modules/articles"
	"go-articles/modules/comments"
	"go-articles/modules/images"
	"go-articles/server/middleware"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"strconv"

	"github.com/JoinVerse/xid"
	"github.com/labstack/echo/v4"
)

type ArticleController struct {
	articleUsecase articles.Usecase
	imageUsecase   images.Usecase
	commentUsecase comments.Usecase
}

func NewArticleController(uc articles.Usecase, iu images.Usecase, cu comments.Usecase) *ArticleController {
	return &ArticleController{
		articleUsecase: uc,
		imageUsecase:   iu,
		commentUsecase: cu,
	}
}

func (controller *ArticleController) GetByID(c echo.Context) error {
	ctx := c.Request().Context()

	id, _ := strconv.Atoi(c.Param("articleId"))
	article, err := controller.articleUsecase.GetByID(ctx, id)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	article.Comments, err = controller.commentUsecase.GetByArticleID(ctx, article.ID)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, response.FromDomain(article))
}

func (controller *ArticleController) Insert(c echo.Context) error {
	ctx := c.Request().Context()

	request := request.Articles{}
	userId := middleware.GetUser(c).ID
	if err := c.Bind(&request); err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	validateMessage, validate, err := helpers.Validate(&request)

	if validate {
		return helpers.ErrorValidateResponse(c, http.StatusBadRequest, err, validateMessage)
	}

	err = controller.articleUsecase.Insert(ctx, request.ToDomain(), userId)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusCreated, nil)
}

func (controller *ArticleController) Upload(c echo.Context) error {
	ctx := c.Request().Context()

	id, _ := strconv.Atoi(c.Param("articleId"))
	userId := middleware.GetUser(c).ID

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

	fileType := "articles"
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

	article, err := controller.articleUsecase.GetByID(ctx, id)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	err = controller.articleUsecase.Update(ctx, &articles.Domain{
		Title:       article.Title,
		Description: article.Description,
		ImageID:     &imageID,
	}, id, userId)

	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusCreated, nil)
}

func (controller *ArticleController) Update(c echo.Context) error {
	ctx := c.Request().Context()

	request := request.Articles{}
	id, _ := strconv.Atoi(c.Param("articleId"))
	userId := middleware.GetUser(c).ID

	if err := c.Bind(&request); err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	validateMessage, validate, err := helpers.Validate(&request)

	if validate {
		return helpers.ErrorValidateResponse(c, http.StatusBadRequest, err, validateMessage)
	}

	article, _ := controller.articleUsecase.GetByID(ctx, id)
	if article.UserID != userId {
		return helpers.ErrorResponse(c, http.StatusBadRequest, constants.ErrDoNotHavePermission)
	}

	err = controller.articleUsecase.Update(ctx, request.ToDomain(), id, userId)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, nil)
}

func (controller *ArticleController) Fetch(c echo.Context) error {
	ctx := c.Request().Context()

	userId := middleware.GetUser(c).ID
	userName := middleware.GetUser(c).Name
	page, _ := strconv.Atoi(c.QueryParam("page"))
	perpage, _ := strconv.Atoi(c.QueryParam("per_page"))
	search := c.QueryParam("search")
	sort := c.QueryParam("sort")
	by := c.QueryParam("by")

	articles, count, err := controller.articleUsecase.Fetch(ctx, page, perpage, userId, by, search, sort)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, response.FetchFromListDomain(articles, count, userName))
}

func (controller *ArticleController) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	id, _ := strconv.Atoi(c.Param("articleId"))
	userId := middleware.GetUser(c).ID

	article, _ := controller.articleUsecase.GetByID(ctx, id)
	if article.UserID != userId {
		return helpers.ErrorResponse(c, http.StatusBadRequest, constants.ErrDoNotHavePermission)
	}

	err := controller.articleUsecase.Delete(ctx, id)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, nil)
}

func (controller *ArticleController) AddComment(c echo.Context) error {
	ctx := c.Request().Context()

	id, _ := strconv.Atoi(c.Param("articleId"))
	userId := middleware.GetUser(c).ID

	request := request.Comments{}
	if err := c.Bind(&request); err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	validateMessage, validate, err := helpers.Validate(&request)
	if validate {
		return helpers.ErrorValidateResponse(c, http.StatusBadRequest, err, validateMessage)
	}

	err = controller.commentUsecase.Insert(ctx, request.ToDomain(userId), id)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusCreated, nil)
}

func (controller *ArticleController) DeleteComment(c echo.Context) error {
	ctx := c.Request().Context()

	articleId, _ := strconv.Atoi(c.Param("articleId"))
	commentId, _ := strconv.Atoi(c.Param("commentId"))
	userId := middleware.GetUser(c).ID

	err := controller.commentUsecase.Delete(ctx, articleId, commentId, userId)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, nil)
}

func (controller *ArticleController) UpdateComment(c echo.Context) error {
	ctx := c.Request().Context()

	commentId, _ := strconv.Atoi(c.Param("commentId"))
	userId := middleware.GetUser(c).ID

	request := request.Comments{}
	if err := c.Bind(&request); err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	validateMessage, validate, err := helpers.Validate(&request)
	if validate {
		return helpers.ErrorValidateResponse(c, http.StatusBadRequest, err, validateMessage)
	}

	controller.commentUsecase.Update(ctx, request.ToDomain(userId), commentId)
	return helpers.SuccessResponse(c, http.StatusOK, nil)
}

func (controller *ArticleController) SearchArticles(c echo.Context) error {
	ctx := c.Request().Context()

	search := c.QueryParam("search")

	articles, err := controller.articleUsecase.Search(ctx, search)
	if err != nil {
		return helpers.ErrorResponse(c, http.StatusBadRequest, err)
	}

	return helpers.SuccessResponse(c, http.StatusOK, response.FromListDomain(articles))
}
