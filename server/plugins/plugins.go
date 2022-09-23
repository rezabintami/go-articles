package app_plugins

import (
	"database/sql"
	_helpers "go-articles/helpers"
	_config "go-articles/server/config"
	_middleware "go-articles/server/middleware"
	_routes "go-articles/server/routes"
	"time"

	_userController "go-articles/controllers/users"
	_userRepo "go-articles/databases/postgres/repository/users"
	_userUsecase "go-articles/modules/users"

	_roleController "go-articles/controllers/roles"
	_roleRepo "go-articles/databases/postgres/repository/roles"
	_roleUsecase "go-articles/modules/roles"

	_articleController "go-articles/controllers/articles"
	_articleRepo "go-articles/databases/postgres/repository/articles"
	_articleUsecase "go-articles/modules/articles"

	_imageRepo "go-articles/databases/postgres/repository/images"
	_imageUsecase "go-articles/modules/images"

	_commentRepo "go-articles/databases/postgres/repository/comments"
	_commentUsecase "go-articles/modules/comments"

	_mountainController "go-articles/controllers/mountains"
	_mountainRepo "go-articles/databases/postgres/repository/mountains"
	_mountainUsecase "go-articles/modules/mountains"

	"github.com/go-redis/redis"
)

type ConfigurationPlugins struct {
	PostgreDriver  *sql.DB
	RedisDriver    *redis.Client
	MailConnection _helpers.MailConnection
}

func (route *ConfigurationPlugins) RoutePlugins() _routes.ControllerList {
	configJWT := _middleware.ConfigJWT{
		SecretJWT:        _config.GetConfiguration("jwt.access_token"),
		RefreshSecretJWT: _config.GetConfiguration("jwt.refresh_token"),
		ExpiredDuration:  _helpers.ConvertStringtoInt(_config.GetConfiguration("jwt.expired")),
	}

	configForgotJWT := _middleware.ConfigForgotJWT{
		ResetSecretJWT:  _config.GetConfiguration("jwt.reset_token"),
		ExpiredDuration: _helpers.ConvertStringtoInt(_config.GetConfiguration("jwt.expired")),
	}

	timeoutContext := time.Duration(_helpers.ConvertStringtoInt(_config.GetConfiguration("server.timeout"))) * time.Second

	//! REPOSITORY
	userRepo := _userRepo.NewPostgreUsersRepository(route.PostgreDriver)
	roleRepo := _roleRepo.NewPostgreRolesRepository(route.PostgreDriver)
	articleRepo := _articleRepo.NewPostgreArticlesRepository(route.PostgreDriver)
	imageRepo := _imageRepo.NewPostgreImagesRepository(route.PostgreDriver)
	commentRepo := _commentRepo.NewPostgreCommentRepository(route.PostgreDriver)
	mountainRepo := _mountainRepo.NewPostgreMountainsRepository(route.PostgreDriver)

	//! USECASE
	userUsecase := _userUsecase.NewUserUsecase(userRepo, roleRepo, &configJWT, &configForgotJWT, timeoutContext, route.MailConnection, route.RedisDriver)
	articleUsecase := _articleUsecase.NewArticleUsecase(articleRepo)
	roleUsecase := _roleUsecase.NewRoleUsecase(roleRepo)
	imageUsecase := _imageUsecase.NewImageUsecase(imageRepo)
	commentUsecase := _commentUsecase.NewCommentUsecase(commentRepo)
	mountainUsecase := _mountainUsecase.NewMountainUsecase(mountainRepo)

	//! CONTROLLER
	userCtrl := _userController.NewUserController(userUsecase, imageUsecase)
	articleCtrl := _articleController.NewArticleController(articleUsecase, imageUsecase, commentUsecase)
	roleCtrl := _roleController.NewRoleController(roleUsecase)
	mountainCtrl := _mountainController.NewMountainController(mountainUsecase)

	return _routes.ControllerList{
		JWTMiddleware:      configJWT.Init(),
		JWTResetPassword:   configForgotJWT.VerifyTokenForgotPassword(),
		UserController:     *userCtrl,
		ArticleController:  *articleCtrl,
		RoleController:     *roleCtrl,
		MountainController: *mountainCtrl,
	}
}
