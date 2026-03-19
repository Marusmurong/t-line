package logger

import (
	"os"
	"path/filepath"

	"github.com/t-line/backend/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var L *zap.SugaredLogger

func Init(cfg config.LogConfig) {
	level := parseLevel(cfg.Level)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	var cores []zapcore.Core

	// console output
	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)
	cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level))

	// file output
	if cfg.Filename != "" {
		dir := filepath.Dir(cfg.Filename)
		if err := os.MkdirAll(dir, 0o755); err == nil {
			fileWriter := &lumberjack.Logger{
				Filename:   cfg.Filename,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
				Compress:   cfg.Compress,
			}
			jsonEncoder := zapcore.NewJSONEncoder(encoderCfg)
			cores = append(cores, zapcore.NewCore(jsonEncoder, zapcore.AddSync(fileWriter), level))
		}
	}

	core := zapcore.NewTee(cores...)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))
	L = zapLogger.Sugar()
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
