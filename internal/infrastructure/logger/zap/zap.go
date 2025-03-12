package zap

import (
	"slices"

	"go.uber.org/zap"
	"github.com/cryptoPickle/go-ddd-example/internal/infrastructure/logger"
)

type ZapLogger struct {
	logger  *zap.Logger
	sugared *zap.SugaredLogger
	config  *ZapConfig
}

type ZapConfig struct {
	LogLevel logger.LogLevel
	Env      string
	LogFile  string
}

func NewZapLogger(config *ZapConfig) (*ZapLogger, error) {
	if config.LogLevel == "" {
		config.LogLevel = logger.Info
	}
	logger := ZapLogger{}
	logger.config = config

	if err := logger.setup(config); err != nil {
		return nil, err
	}

	return &logger, nil
}

func (l *ZapLogger) Debug(args ...any) {
	l.sugared.Debug(args...)
}

func (l *ZapLogger) Debugf(format string, args ...any) {
	defer l.sugared.Sync()
	l.sugared.Debugf(format, args...)
}

func (l *ZapLogger) Info(args ...any) {
	l.sugared.Info(args...)
}

func (l *ZapLogger) Infof(format string, args ...any) {
	defer l.sugared.Sync()
	l.sugared.Infof(format, args...)
}

func (l *ZapLogger) Warn(args ...any) {
	l.sugared.Warn(args...)
}

func (l *ZapLogger) Warnf(format string, args ...any) {
	defer l.sugared.Sync()
	l.sugared.Warnf(format, args...)
}

func (l *ZapLogger) Error(args ...any) {
	l.sugared.Error(args...)
}

func (l *ZapLogger) Errorf(format string, args ...any) {
	defer l.sugared.Sync()
	l.sugared.Errorf(format, args...)
}

func (l *ZapLogger) Fatal(args ...any) {
	l.sugared.Fatal(args...)
}

func (l *ZapLogger) Fatalf(format string, args ...any) {
	defer l.sugared.Sync()
	l.sugared.Fatalf(format, args...)
}

func (l *ZapLogger) SetLevel(level logger.LogLevel) error {
	levels := []logger.LogLevel{logger.Debug, logger.Info, logger.Warn, logger.Error, logger.Fatal}
	if !slices.Contains(levels, level) {
		level = logger.Info
	}

	if level != l.config.LogLevel {
		l.config.LogLevel = level
		if err := l.setup(l.config); err != nil {
			return err
		}
	}
	return nil
}

func getLevel(level logger.LogLevel) zap.AtomicLevel {
	var lvl zap.AtomicLevel

	switch level {
	case logger.Debug:
		lvl = zap.NewAtomicLevelAt(zap.DebugLevel)
	case logger.Info:
		lvl = zap.NewAtomicLevelAt(zap.InfoLevel)
	case logger.Warn:
		lvl = zap.NewAtomicLevelAt(zap.WarnLevel)
	case logger.Error:
		lvl = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case logger.Fatal:
		lvl = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		lvl = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	return lvl
}

func (l *ZapLogger) setup(config *ZapConfig) error {
	var (
		encoding      = "console"
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	)

	if config.Env == "prod" {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoding = "json"
	}

	paths := []string{"stderr"}
	if config.LogFile != "" {
		if _, err := logger.CreateLogFile(config.LogFile); err != nil {
			return err
		}
		paths = append(paths, config.LogFile)
	}

	zapConfig := zap.Config{
		Level:            getLevel(config.LogLevel),
		Encoding:         encoding,
		EncoderConfig:    encoderConfig,
		OutputPaths:      paths,
		ErrorOutputPaths: paths,
	}

	z, _ := zapConfig.Build()
	l.logger = z
	l.sugared = z.Sugar()
	return nil
}
