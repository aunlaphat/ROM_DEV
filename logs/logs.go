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

	// ⚡ ปรับแต่ง Console Encoder (แบบข้อความที่อ่านง่าย)
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(hook), zapcore.ErrorLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
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

/*
func (l *Logger) LogAPICall(ctx context.Context, apiName string, fields ...zap.Field) func(status string, err error, additionalFields ...zap.Field) {
	start := time.Now()
	logger := l.logger.With(zap.String("apiName", apiName))

	// 📌 Log ว่า API Call เริ่มต้นแล้ว
	logger.Info("⏳ Starting API Call ⏳", fields...)

	// 📌 ฟังก์ชันปิด Log เมื่อ API Call จบ
	return func(status string, err error, additionalFields ...zap.Field) {
		duration := time.Since(start)
		logFields := append(fields,
			zap.Duration("duration", duration),
			zap.String("status", status),
		)

		// 🔥 ตรวจสอบว่าใช้เวลานานเกินไปหรือไม่
		if duration > 5*time.Second {
			logFields = append(logFields, zap.Bool("slowExecution", true))
		}

		// ❌ ถ้ามีข้อผิดพลาด
		if err != nil {
			logFields = append(logFields,
				zap.String("error", err.Error()),
				zap.String("errorType", fmt.Sprintf("%T", err)),
				zap.Stack("stackTrace"),
			)
		}

		// ✅ เพิ่ม additionalFields (ข้อมูลเพิ่มเติมที่ส่งมาตอนจบ API)
		logFields = append(logFields, additionalFields...)

		// 🔥 Log ตามสถานะของ API Call
		switch status {
		case "Success":
			logger.Info("✅ API Call Success ✅", logFields...)
		case "Failed":
			logger.Error("❌ API Call Failed ❌", logFields...)
		case "Not Found":
			logger.Warn("⚠️ API Call Not Found ⚠️", logFields...)
		default:
			logger.Info("🔚 API Call End 🔚", logFields...)
		}
	}
}
*/
