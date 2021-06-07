package util

import "go.uber.org/zap"

var Logger *zap.SugaredLogger

func SetupLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	Logger = logger.Sugar()
}
