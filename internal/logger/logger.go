package logger

import (
	"go.uber.org/zap"
)

type log struct {
	*zap.Logger
}

func (l log) Error(msg string, err error) {
	l.Logger.Error(msg, zap.Error(err))
}

func (l log) ErrorMsg(msg string) {
	l.Logger.Error(msg)
}

var Log log

func InitLogger(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log.Logger = zl
	return nil
}
