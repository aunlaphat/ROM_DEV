package logs

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	logger *zap.Logger
}

// NewLogger สร้าง instance ใหม่ของ Logger พร้อมการหมุนเวียนไฟล์และการตั้งค่าที่กำหนดเอง
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

	filePriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(hook), filePriority),
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
		logger.Sync()
	}

	return &Logger{logger: logger}, close, nil
}

// ฟังก์ชันทั่วไปสำหรับการบันทึก log
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

type LogConfig struct {
	ServiceName   string
	LogPath       string
	MaxSize       int
	MaxBackups    int
	MaxAge        int
	SlowThreshold time.Duration
}

func (l *Logger) LogAPICall(ctx context.Context, apiName string, fields ...zap.Field) func(status string, err error, additionalFields ...zap.Field) {
	start := time.Now()            // บันทึกเวลาที่เริ่มต้นการเรียก API
	traceID := uuid.New().String() // สร้าง traceID ที่ไม่ซ้ำกันสำหรับการติดตาม

	// เพิ่มข้อมูลพื้นฐานสำหรับการบันทึก log
	baseFields := append(fields,
		zap.String("traceID", traceID), // เพิ่ม traceID
		zap.String("apiName", apiName), // เพิ่มชื่อ API
		zap.Time("startTime", start))   // เพิ่มเวลาที่เริ่มต้น)

	// ดึงข้อมูลจาก context และเพิ่มลงใน log fields
	for _, key := range []string{"RequestID", "UserID", "ClientIP", "UserAgent"} {
		if val, ok := ctx.Value(key).(string); ok {
			baseFields = append(baseFields, zap.String(key, val)) // เพิ่มข้อมูลจาก context
		}
	}

	l.Info("⏳ Starting API Call", baseFields...) // บันทึก log ว่าเริ่มต้นการเรียก API

	return func(status string, err error, additionalFields ...zap.Field) {
		duration := time.Since(start) // คำนวณระยะเวลาที่ใช้ในการเรียก API
		logFields := append(baseFields,
			zap.Duration("duration", duration), // เพิ่มระยะเวลา
			zap.String("status", status))       // เพิ่มสถานะของการเรียก API

		// ตรวจสอบว่าการเรียก API ใช้เวลานานเกินไปหรือไม่
		if duration > 5*time.Second {
			logFields = append(logFields,
				zap.Bool("slowExecution", true),                    // ระบุว่าการเรียก API ใช้เวลานาน
				zap.Float64("durationSeconds", duration.Seconds())) // เพิ่มระยะเวลาในหน่วยวินาที
		}

		// ตรวจสอบว่ามีข้อผิดพลาดเกิดขึ้นหรือไม่
		if err != nil {
			logFields = append(logFields,
				zap.String("error", err.Error()),                // เพิ่มข้อความข้อผิดพลาด
				zap.String("errorType", fmt.Sprintf("%T", err)), // เพิ่มประเภทของข้อผิดพลาด
				zap.Stack("stackTrace"))                         // เพิ่ม stack trace ของข้อผิดพลาด
		}

		logFields = append(logFields, additionalFields...) // เพิ่มข้อมูลเพิ่มเติมลงใน log fields

		// บันทึก log ตามสถานะของการเรียก API
		switch status {
		case "Success":
			l.Info("✅ API Call Success ✅", logFields...) // บันทึก log ว่าสำเร็จ
		case "Failed":
			l.Error("❌ API Call Failed ❌", logFields...) // บันทึก log ว่าล้มเหลว
		case "Not Found":
			l.Warn("⚠️ API Call Not Found ⚠️", logFields...) // บันทึก log ว่าไม่พบข้อมูล
		default:
			l.Info("🔚 API Call End 🔚", logFields...) // บันทึก log ว่าการเรียก API เสร็จสิ้น
		}
	}
}
