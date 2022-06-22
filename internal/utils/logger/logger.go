package logger

import (
	"github.com/sirupsen/logrus"
	"log"
)

type Logger struct {
	Logger *logrus.Entry
}

type LoggerInterface interface {
	LogInfo(data interface{})
	LogError(data interface{})
}

func (l *Logger) LogInfo(data interface{}) {
	log.Println(data)
	//l.Logger.Info(data)
}

func (l *Logger) LogError(data interface{}) {
	log.Println(data)
	//l.Logger.Error(data)
}
