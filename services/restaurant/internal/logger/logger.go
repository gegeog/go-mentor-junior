package logger

import (
	"strings"

	"github.com/gegeog/go-mentor-junior/services/restaurant/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(loggerConfig config.LoggerConfig) *zap.Logger {
	loggerLvl := getLogLevelFromEnv(loggerConfig.Level)

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(loggerLvl),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build(
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	// TODO (review): здесь паникуем. но мне кажется это нормальное поведение:
	//  если у нас какие-то базовые вещи не инициализируются, хотя мы в дальнейшем ожидаем их использовать.
	// я правильно понимаю, что ты имел ввиду что слишком жестко было паниковать если не нашли логгер в контексте?
	if err != nil {
		panic(err)
	}

	return logger
}

func getLogLevelFromEnv(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}
