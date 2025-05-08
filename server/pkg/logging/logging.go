package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var l *zap.Logger

func GetLogger() *zap.Logger {
	return l
}

func init() {
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		CallerKey:        "caller",
		MessageKey:       "msg",
		EncodeLevel:      zapcore.CapitalColorLevelEncoder,
		EncodeTime:       zapcore.ISO8601TimeEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		ConsoleSeparator: " | ",
	})

	fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		CallerKey:    "caller",
		MessageKey:   "msg",
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	})

	logFile, err := os.OpenFile("./pkg/logging/logs/all.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	fileWriter := zapcore.AddSync(logFile)

	consoleWriter := zapcore.AddSync(os.Stdout)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleWriter, zapcore.DebugLevel),
		zapcore.NewCore(fileEncoder, fileWriter, zapcore.DebugLevel),
	)

	logger := zap.New(core, zap.AddCaller())

	l = logger
}
