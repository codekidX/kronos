package internal

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const DEBUG_VALUE = "1"

func CreateLogger() *zap.Logger {
	var config zapcore.EncoderConfig
	if os.Getenv("CN_DEBUG") == DEBUG_VALUE {
		config = zap.NewDevelopmentEncoderConfig()
		consoleJSONEncoder := zapcore.NewJSONEncoder(config)
		core := zapcore.NewTee(
			zapcore.NewCore(consoleJSONEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		)
		return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	} else {
		logFilePath := os.Getenv("CN_LOG_PATH")
		if logFilePath == "" {
			logFilePath = "nutlogs.json"
		}
		config = zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		fileEncoder := zapcore.NewJSONEncoder(config)
		// FIXME: we need to take log file from env variable
		logFile, _ := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		writer := zapcore.AddSync(logFile)
		defaultLogLevel := zapcore.DebugLevel
		core := zapcore.NewTee(
			zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		)
		return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	}
}
