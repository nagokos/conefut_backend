package logger

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/nagokos/connefut_backend/config"
	sqldblogger "github.com/simukti/sqldb-logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct{}

func NewLogger() *zap.Logger {
	encoderConcoleConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoderLogConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	path := config.Config.LogFile
	errPath := config.Config.LogErrorFile

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		NewLogger().Error(err.Error())
	}

	errFile, err := os.OpenFile(errPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		NewLogger().Error(err.Error())
	}

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConcoleConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	logCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderLogConfig),
		zapcore.AddSync(file),
		zapcore.DebugLevel,
	)

	errLogCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderLogConfig),
		zapcore.AddSync(errFile),
		zapcore.ErrorLevel,
	)

	logger := zap.New(
		zapcore.NewTee(
			consoleCore,
			logCore,
			errLogCore,
		),
		zap.AddCaller(),
	)

	return logger
}

func (l *Logger) Log(ctx context.Context, level sqldblogger.Level, msg string, data map[string]interface{}) {
	if string(level) == "error" {
		w := bytes.NewBufferString(fmt.Sprintf("%s:%v", level, msg))

		for k, v := range data {
			fmt.Fprintf(w, "\t%s:%v", k, v)
		}

		NewLogger().Error(w.String())
	} else {
		w := bytes.NewBufferString(fmt.Sprintf("%s:%v", level, msg))

		for k, v := range data {
			fmt.Fprintf(w, "\t%s:%v", k, v)
		}

		NewLogger().Debug(w.String())
	}
}
