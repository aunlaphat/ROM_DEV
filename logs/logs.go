// Log package
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

// New Logger instance with save to error.json for view error log when service is on server production
func NewLogger(logPath string, maxSize, maxBackups, maxAge int) (*Logger, func(), error) {

	hook := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   true,
	}

	config := zap.NewProductionEncoderConfig()

	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	filePiority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(hook), filePiority),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)

	logger := zap.New(core, zap.AddCallerSkip(1), zap.AddStacktrace(zap.WarnLevel))

	close := func() {
		logger.Sync()
	}

	return &Logger{logger: logger}, close, nil

}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info("ℹ️ " + msg, fields...)
}

func (l *Logger) Error(msg interface{}, fields ...zap.Field) {
	switch v := msg.(type) {
	case error:
		l.logger.Error("❌ " + v.Error(), fields...)
	case string:
		l.logger.Error("❌ " + v, fields...)
	}
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn("⚠️ " + msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug("🐛 " + msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal("💀 " + msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic("😱 " + msg, fields...)
}

func (l *Logger) Sync() error {
	return l.logger.Sync()
}
