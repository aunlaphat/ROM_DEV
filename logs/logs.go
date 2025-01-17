// Log package บันทึกข้อมูลและ error ที่เกิดขั้นภายในระบบ
package logs

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	logger *zap.Logger
}

// New Logger instance with save to error.json for view error log when service is on server production
func NewLogger(serviceName, logPath string, maxSize, maxBackups, maxAge int) (*Logger, func(), error) {

	hook := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   true,
	}

	config := zap.NewProductionEncoderConfig()

	config.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEncoder := zapcore.NewJSONEncoder(config)

	filePiority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(hook), filePiority),
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)

	logger := zap.New(
		core,
		zap.Fields(zap.String("service", serviceName)),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel))

	close := func() {
		logger.Sync()
	}

	return &Logger{logger: logger}, close, nil
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Error(msg interface{}, fields ...zap.Field) {
	switch v := msg.(type) {
	case error:
		l.logger.Error(v.Error(), fields...)
	case string:
		l.logger.Error(v, fields...)
	}
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}
func (l *Logger) Sync() error {
	return l.logger.Sync()
}

// log api call
func (l *Logger) LogAPICall(apiName string, fields ...zap.Field) func(status string, err error) {
	start := time.Now()
	
	// Log จุดเริ่มต้น
	l.Info(fmt.Sprintf("⏳ Starting API Call: %s", apiName), fields...)

	return func(status string, err error) {
		duration := time.Since(start)
		logFields := append(fields, 
			zap.Duration("duration", duration),
			zap.String("status", status))

		// สร้าง message ตาม status
		var msg string
		switch status {
		case "Success":
			msg = fmt.Sprintf("✅ API Call Success: %s", apiName)
			l.Info(msg, logFields...)
		case "Failed":
			msg = fmt.Sprintf("❌ API Call Failed: %s", apiName)
			if err != nil {
				// ใช้เฉพาะ error message โดยไม่แสดง stacktrace
				logFields = append(logFields, zap.String("error", err.Error()))
				l.Error(msg, logFields...)
			}
		case "Not Found":
			msg = fmt.Sprintf("⚠️ API Call Not Found: %s", apiName)
			if err != nil {
				logFields = append(logFields, zap.String("error", err.Error()))
			}
			l.Warn(msg, logFields...)
		default:
			msg = fmt.Sprintf("🔄 API Call Completed: %s", apiName)
			l.Info(msg, logFields...)
		}
	}
}
