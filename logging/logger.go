package logging

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init() {
	logger, _ := zap.NewProduction()
	Log = logger
}

func Error(msg string, err error) {
	Log.Error(msg,
		zap.Error(err),
	)
}
