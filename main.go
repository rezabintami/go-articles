package main

import (
	"log"
	"os"

	_mongoDriver "go-articles/databases/mongodb"
	_postgreDriver "go-articles/databases/postgres"
	_redisDriver "go-articles/databases/redis"
	"go-articles/helpers"
	"go-articles/logger"
	_config "go-articles/server/config"
	_plugins "go-articles/server/plugins"

	echo "github.com/labstack/echo/v4"
)

func main() {
	log.Println("Starting application version :", _config.GetConfiguration("app.version"))
	log.Println("Environment :", _config.GetConfiguration("app.env"))
	log.Println("Server started at port :" + _config.GetConfiguration("server.port"))

	log.Println("User :", _config.GetConfiguration("postgres.user"))
	log.Println("Host :", _config.GetConfiguration("postgres.host"))
	log.Println("Port :", _config.GetConfiguration("postgres.port"))
	log.Println("Name :", _config.GetConfiguration("postgres.name"))

	// * Init Postgresql
	postgreDriver := _postgreDriver.Init()

	// * Init Mongodb
	mongoDriver := _mongoDriver.Init()

	// * Init Redis
	redisDriver := _redisDriver.InitialRedis(_config.GetConfiguration("redis.host"), _config.GetConfiguration("redis.pass"))

	// * Init Logger Mongodb
	logger.Init(mongoDriver)
	
	// * Init Validation
	helpers.InitValidation()

	// * Init Mail
	mailConnection := helpers.MailConnection{
		Host:     _config.GetConfiguration("mail.host"),
		Port:     helpers.ConvertStringtoInt(_config.GetConfiguration("mail.port")),
		Username: _config.GetConfiguration("mail.username"),
		Password: _config.GetConfiguration("mail.password"),
	}

	plugins := _plugins.ConfigurationPlugins{
		PostgreDriver:  postgreDriver,
		RedisDriver:    redisDriver,
		MailConnection: mailConnection,
	}

	e := echo.New()

	// * Init Dependencies Plugin
	route := plugins.RoutePlugins()
	route.RouteRegister(e)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	e.Start(":" + port)
}
