package logs

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	logger *zap.Logger
}

func NewLogger(serviceName, logPath string, maxSize, maxBackups, maxAge int) (*Logger, func(), error) {
	hook := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   true,
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.RFC3339TimeEncoder
	jsonEncoder := zapcore.NewJSONEncoder(config)

	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(hook), zapcore.ErrorLevel),
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)

	logger := zap.New(
		core,
		zap.Fields(zap.String("service", serviceName)),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
	)

	close := func() {
		_ = logger.Sync()
	}

	return &Logger{logger: logger}, close, nil
}

// 📌 Logging Methods
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}

// 📌 เพิ่ม `With()` ให้ Logger รองรับการเพิ่ม Context Fields
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{logger: l.logger.With(fields...)}
}

// 📌 Sync เพื่อปิด Logger อย่างปลอดภัย
func (l *Logger) Sync() {
	_ = l.logger.Sync()
}
