package zaplog

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

var globalLogger *Logger

func StartLogger() error {
	cfg := zapcore.EncoderConfig{
		TimeKey:          "dt",
		LevelKey:         "level",
		MessageKey:       "message",
		CallerKey:        "caller",
		StacktraceKey:    "stack_trace",
		ConsoleSeparator: " ",
		EncodeLevel:      zapcore.LowercaseLevelEncoder,
		EncodeTime:       zapcore.RFC3339NanoTimeEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.AddSync(os.Stdout),
		zap.InfoLevel,
	)

	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
	)

	globalLogger = &Logger{Logger: logger}
	return nil
}

func GetGlobalLogger() *Logger {
	if globalLogger == nil {
		_ = StartLogger()
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

func Infof(ctx context.Context, format string, args ...any) {
	GetGlobalLogger().WithContext(ctx).Info(fmt.Sprintf(format, args...))
}

func Debug(ctx context.Context, message string, fields ...zap.Field) {
	GetGlobalLogger().WithContext(ctx).Debug(message, fields...)
}

func Debugf(ctx context.Context, format string, args ...any) {
	GetGlobalLogger().WithContext(ctx).Debug(fmt.Sprintf(format, args...))
}

func Error(ctx context.Context, message string, fields ...zap.Field) {
	GetGlobalLogger().WithContext(ctx).Error(message, fields...)
}

func Errorf(ctx context.Context, format string, args ...any) {
	GetGlobalLogger().WithContext(ctx).Error(fmt.Sprintf(format, args...))
}

func Fatal(ctx context.Context, message string, fields ...zap.Field) {
	GetGlobalLogger().WithContext(ctx).Fatal(message, fields...)
}

func Fatalf(ctx context.Context, format string, args ...any) {
	GetGlobalLogger().WithContext(ctx).Fatal(fmt.Sprintf(format, args...))
}
