package logger

import (
	loggerHook "go-articles/databases/mongodb/repository/logger"

	"go-articles/server/config"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	HttpRequest = "send http request to [%s] %s %s"
	UserId      = "User-ID"
)

//! Logger output file
// func NewLogger() Logger {
// 	logger := log.New()
// 	f, err := os.OpenFile("api.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
// 	if err != nil {
// 		log.Errorf("cannot open 'testlogfile', (%s)", err.Error())
// 		flag.Usage()
// 		os.Exit(-1)
// 	}
// 	logger.SetFormatter(&log.JSONFormatter{})
// 	logger.SetOutput(f)
// 	return Logger{
// 		Log: logger,
// 	}
// }

func Init(db *mongo.Database) {
	logger.SetFormatter(&logger.JSONFormatter{})
	switch strings.ToLower(config.GetConfiguration("log.level")) {
	case "info":
		logger.SetLevel(logger.InfoLevel)
	case "warn":
		logger.SetLevel(logger.WarnLevel)
	case "error":
		logger.SetLevel(logger.ErrorLevel)
	case "debug":
		logger.SetLevel(logger.DebugLevel)
	default:
		logger.SetLevel(logger.InfoLevel)
	}
	logger.AddHook(&loggerHook.LoggerHook{
		DB: db,
		Log: *logger.StandardLogger(),
	})
	logger.SetOutput(os.Stdout)
}

func Info(format string, values ...interface{}) {
	logger.WithFields(logger.Fields{}).Infof(format, values...)
}

func Warn(format string, values ...interface{}) {
	logger.WithFields(logger.Fields{}).Warnf(format, values...)
}

func Error(format string, values ...interface{}) {
	logger.WithFields(logger.Fields{}).Errorf(format, values...)
}

func Debug(format string, values ...interface{}) {
	logger.WithFields(logger.Fields{}).Debugf(format, values...)
}

func Fatal(format string, values ...interface{}) {
	logger.WithFields(logger.Fields{}).Fatalf(format, values...)
}

//! CAN USE ADDHOOK TO CREATE FIELD IN MONGODB
func Logging(c echo.Context) {
	logger.WithContext(c.Request().Context()).WithFields(logger.Fields{}).Debugf(
		HttpRequest,
		c.Request().Method,
		c.Request().URL,
		c.Request().RemoteAddr,
	)
}
