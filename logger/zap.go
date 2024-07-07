package logger

import (
	"context"

	"github.com/everywan/common/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	*zap.SugaredLogger
	logger    *zap.Logger
	atomLevel zap.AtomicLevel
}

var _ Logger = new(zapLogger)

func NewZapLogger(options ...zap.Option) (*zapLogger, error) {
	atomLevel := zap.NewAtomicLevel()
	var cfg zap.Config
	if utils.IsProdEnv() || utils.IsStagingEnv() {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}
	cfg.Level = atomLevel
	logger, err := cfg.Build(options...)
	if err != nil {
		return &zapLogger{}, err
	}
	return &zapLogger{
		SugaredLogger: logger.Sugar(),
		logger:        logger,
		atomLevel:     atomLevel,
	}, nil
}

func (l *zapLogger) SetLevel(level Level) {
	l.atomLevel.SetLevel(zapcore.Level(level))
}

// todo 添加 ctx 的处理, 比如获取 traceid 等.
func (l *zapLogger) Debug(ctx context.Context, args ...interface{}) {
	l.SugaredLogger.Debug(args...)
}

func (l *zapLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.SugaredLogger.Debugf(format, args...)
}

func (l *zapLogger) Info(ctx context.Context, args ...interface{}) {
	l.SugaredLogger.Info(args...)
}

func (l *zapLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.SugaredLogger.Infof(format, args...)
}

func (l *zapLogger) Warn(ctx context.Context, args ...interface{}) {
	l.SugaredLogger.Warn(args...)
}

func (l *zapLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.SugaredLogger.Warnf(format, args...)
}

func (l *zapLogger) Error(ctx context.Context, args ...interface{}) {
	l.SugaredLogger.Error(args...)
}

func (l *zapLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.SugaredLogger.Errorf(format, args...)
}

func (l *zapLogger) Fatal(ctx context.Context, args ...interface{}) {
	l.SugaredLogger.Fatal(args...)
}

func (l *zapLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.SugaredLogger.Fatalf(format, args...)
}
