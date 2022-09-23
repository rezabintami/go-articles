package logger

import (
	// log "go-articles/logger"

	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoggerHook struct {
	DB  *mongo.Database
	Log logger.Logger
}

type Repository interface {
	Levels() []logger.Level
	Fire(entry *logger.Entry) error
}

func (repository *LoggerHook) Levels() []logger.Level {
	return []logger.Level{
		logger.ErrorLevel,
		logger.WarnLevel,
		logger.InfoLevel,
		logger.DebugLevel,
	}
}

func (repository *LoggerHook) Fire(entry *logger.Entry) error {
	_, err := repository.DB.Collection("logger").InsertOne(entry.Context, Logger{
		Level:   entry.Level.String(),
		Message: entry.Message,
		Time:    entry.Time,
	})
	if err != nil {
		repository.Log.Error("Can't insert data :%s", err.Error())
	}
	return nil
}
