package logger

import "go.uber.org/zap"

type Logger struct {
	*zap.Logger
}

func New(env string) (*Logger, error) {
	var log *zap.Logger
	var err error
	if env == "prod" {
		log, err = zap.NewProduction()
	} else {
		log, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, err
	}
	return &Logger{Logger: log}, nil
}
