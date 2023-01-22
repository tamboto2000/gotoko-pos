package logger

import (
	"errors"
	"os"

	"go.uber.org/zap"
)

const (
	ModeStaging    = "staging"
	ModeProduction = "production"
	logFile        = "./gotoko-pos.log"
)

func NewLogger(mode string) (*zap.Logger, error) {
	var conf zap.Config

	if mode == ModeStaging {
		conf = zap.NewDevelopmentConfig()
	} else if mode == ModeProduction {
		conf = zap.NewProductionConfig()
	}

	if _, err := os.Stat(logFile); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(logFile)
		if err != nil {
			return nil, err
		}

		defer f.Close()
	}

	conf.OutputPaths = append(conf.OutputPaths, logFile)
	return conf.Build()
}
