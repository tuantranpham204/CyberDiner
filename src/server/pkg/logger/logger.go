package logger

import (
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func Init(mode string) error {
	var z *zap.Logger
	var err error
	if mode == "debug" {
		z, err = zap.NewDevelopment()
	} else {
		z, err = zap.NewProduction()
	}
	if err != nil {
		return err
	}
	log = z.Sugar()
	return nil
}

func L() *zap.SugaredLogger {
	if log == nil {
		z, _ := zap.NewDevelopment()
		log = z.Sugar()
	}
	return log
}

func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
