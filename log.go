package zaplog

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

var globalLogger *Logger

func InitGlobalLogger() error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	logger, err := config.Build()
	if err != nil {
		return err
	}

	globalLogger = &Logger{Logger: logger}

	return nil
}

func GetGlobalLogger() *Logger {
	if globalLogger == nil {
		_ = InitGlobalLogger()
	}

	return globalLogger
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	traceID := TraceIDFromContext(ctx)
	return &Logger{Logger: l.With(zap.String("traceID", traceID))}
}

func TraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value("traceID").(string); ok {
		return traceID
	}
	return ""
}

func Info(ctx context.Context, message string, fields ...zap.Field) {
	GetGlobalLogger().WithContext(ctx).Info(message, fields...)
}

func Debug(ctx context.Context, message string, fields ...zap.Field) {
	GetGlobalLogger().WithContext(ctx).Debug(message, fields...)
}

func Error(ctx context.Context, message string, fields ...zap.Field) {
	GetGlobalLogger().WithContext(ctx).Error(message, fields...)
}
