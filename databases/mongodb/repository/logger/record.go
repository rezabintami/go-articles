package logger

import "time"

type Logger struct {
	Level   string    `bson:"level"`
	Message string    `bson:"message"`
	Time    time.Time `bson:"time"`
}
