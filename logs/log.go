package logs

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"message/config"
)

var LogInfo, LogError, LogAccess *zap.SugaredLogger
var logInfo, logError, logAccess *zap.Logger

func InitLog() {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths: []string{"stdout"},
		InitialFields: map[string]interface{}{
			"version": "v1.0.0",
		},
	}

	// 信息日志
	logInfo = createLog(cfg, config.AppConfig.App.Log.Info)
	LogInfo = logInfo.Sugar()

	// 错误日志
	logError = createLog(cfg, config.AppConfig.App.Log.Error)
	LogError = logError.Sugar()

	// 访问日志
	logAccess = createLog(cfg, config.AppConfig.App.Log.Access)
	LogAccess = logAccess.Sugar()
}

func createLog(cfg zap.Config, outputPath string) *zap.Logger {
	cfg.OutputPaths = append(cfg.OutputPaths, outputPath)
	logger, err := cfg.Build()
	if err != nil {
		panic(fmt.Sprintln("日志功能出现错误！", "\n", err))
	}
	return logger
}

func Sync() {
	logInfo.Sync()
	logError.Sync()
	logAccess.Sync()
}
