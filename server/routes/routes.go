package routes

import (
	"go-articles/controllers/articles"
	"go-articles/controllers/mountains"
	"go-articles/controllers/roles"
	"go-articles/controllers/users"
	_middleware "go-articles/server/middleware"

	_config "go-articles/server/config"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type ControllerList struct {
	JWTMiddleware      middleware.JWTConfig
	JWTResetPassword   middleware.JWTConfig
	UserController     users.UserController
	ArticleController  articles.ArticleController
	RoleController     roles.RoleController
	MountainController mountains.MountainController
}

func (cl *ControllerList) RouteRegister(e *echo.Echo) {
	e.Use(_middleware.MiddlewareLogging)
	e.HTTPErrorHandler = _middleware.CustomHTTPErrorHandler

	e.Static("/static", "static")

	// showing swagger files
	if _config.GetConfiguration("app.env") != "PROD" {
		e.Static("/files", "files")
		url := echoSwagger.URL("/files/swagger.yaml")
		e.GET("/swagger/*", echoSwagger.EchoWrapHandler(url))
	}
	apiV1 := e.Group("/api/v1")

	//! GENERAL-ARTICLES
	article := apiV1.Group("/article")
	article.GET("/:articleId", cl.ArticleController.GetByID)
	article.GET("", cl.ArticleController.SearchArticles)
	article.POST("/:articleId/comment", cl.ArticleController.AddComment, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("USER", "ADMIN", "SUPERUSER"))
	article.DELETE("/:articleId/comment/:commentId", cl.ArticleController.DeleteComment, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("USER", "ADMIN", "SUPERUSER"))
	article.PUT("/:articleId/comment/:commentId", cl.ArticleController.UpdateComment, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("USER", "ADMIN", "SUPERUSER"))

	//! GENERAL-MOUNTAINS
	mountain := apiV1.Group("/mountain")
	mountain.GET("/:mountainId", cl.MountainController.GetByID)
	mountain.GET("", cl.MountainController.SearchMountains)

	//! AUTH
	auth := apiV1.Group("/auth")
	auth.POST("/register", cl.UserController.Register)
	auth.POST("/login", cl.UserController.Login)

	forgotPassword := auth.Group("/forgotPassword")
	forgotPassword.POST("", cl.UserController.ForgotPassword)
	forgotPassword.GET("/token/key/:key", cl.UserController.VerifyTokenForgotPassword)
	forgotPassword.POST("/newPassword", cl.UserController.NewPasswordSubmit, middleware.JWTWithConfig(cl.JWTResetPassword))

	// auth.POST("/logout", cl.UserController.Logout)
	// auth.POST("/refresh", cl.UserController.Refresh)

	//! ADMIN
	admin := apiV1.Group("/admin")

	//! ADMIN-ARTICLES
	adminArticle := admin.Group("/article")
	adminArticle.POST("", cl.ArticleController.Insert, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("ADMIN", "SUPERUSER"))
	adminArticle.POST("/:articleId/upload", cl.ArticleController.Upload, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("ADMIN", "SUPERUSER"))
	adminArticle.PUT("/:articleId", cl.ArticleController.Update, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("ADMIN", "SUPERUSER"))
	adminArticle.GET("", cl.ArticleController.Fetch, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("ADMIN", "SUPERUSER"))
	adminArticle.GET("/:articleId", cl.ArticleController.GetByID, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("ADMIN", "SUPERUSER"))
	adminArticle.DELETE("/:articleId", cl.ArticleController.Delete, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("ADMIN", "SUPERUSER"))

	//! ADMIN-USER
	adminUser := admin.Group("/user")
	adminUser.GET("", cl.UserController.AdminUserFetch, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("ADMIN", "SUPERUSER"))
	adminUser.GET("/:userId", cl.UserController.AdminGetUserByID, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("ADMIN", "SUPERUSER"))
	adminUser.DELETE("/:userId", cl.UserController.AdminDeleteUser, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("ADMIN", "SUPERUSER"))

	//! ADMIN-MOUNTAINS
	adminMountain := admin.Group("/mountain")
	adminMountain.POST("", cl.MountainController.Insert, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("ADMIN", "SUPERUSER"))
	adminMountain.GET("", cl.MountainController.Fetch, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("ADMIN", "SUPERUSER"))
	adminMountain.GET("/:mountainId", cl.MountainController.GetByID, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("ADMIN", "SUPERUSER"))
	adminMountain.PUT("/:mountainId", cl.MountainController.Update, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("ADMIN", "SUPERUSER"))
	adminMountain.DELETE("/:mountainId", cl.MountainController.Delete, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("ADMIN", "SUPERUSER"))

	//! ADMIN-ROLES
	adminRoles := admin.Group("/roles")
	adminRoles.POST("", cl.RoleController.Insert, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("SUPERUSER"))
	adminRoles.GET("/:roleId", cl.RoleController.GetByID, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("SUPERUSER"))
	adminRoles.GET("", cl.RoleController.Fetch, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("SUPERUSER"))
	adminRoles.PUT("/:roleId", cl.RoleController.Update, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("SUPERUSER"))
	adminRoles.DELETE("/:roleId", cl.RoleController.Delete, middleware.JWTWithConfig(cl.JWTMiddleware), _middleware.RoleValidation("SUPERUSER"))

	//! USERS
	user := apiV1.Group("/user")

	//! PROFILE
	userProfile := user.Group("/profile")
	userProfile.GET("", cl.UserController.GetByID, middleware.JWTWithConfig(cl.JWTMiddleware))
	userProfile.PUT("", cl.UserController.Update, middleware.JWTWithConfig(cl.JWTMiddleware))
	userProfile.POST("/upload-photo", cl.UserController.UploadPhoto, middleware.JWTWithConfig(cl.JWTMiddleware))
	userProfile.DELETE("/delete-photo", cl.UserController.DeletePhoto, middleware.JWTWithConfig(cl.JWTMiddleware))
}
