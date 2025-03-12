package logger

import (
	"os"
	"path/filepath"

	"github.com/cryptoPickle/go-ddd-example/internal/common/errors"
)

type LogLevel string

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
	Fatal LogLevel = "fatal"
)

type Logger interface {
	Debug(...any)
	Debugf(format string, v ...any)
	Info(...any)
	Infof(format string, v ...any)
	Warn(...any)
	Warnf(format string, v ...any)
	Error(...any)
	Errorf(format string, v ...any)
	Fatal(...any)
	Fatalf(format string, v ...any)
	SetLevel(lvl LogLevel) error
}

func CreateLogFile(filename string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(filename), 0744)

	if err != nil && err != os.ErrExist {
		return nil, errors.NewContextualError("can not create log file", "logger").Wrap(err)
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, errors.NewContextualError("can not open log file", "logger").Wrap(err)
	}

	return file, nil
}
