package service

import (
	"boilerplate-backend-go/dto/response"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type OrderService interface {
	SearchOrder(ctx context.Context, soNo, orderNo string) (*response.SearchOrderResponse, error)
}

func (srv service) SearchOrder(ctx context.Context, soNo, orderNo string) (*response.SearchOrderResponse, error) {
	// 📝 เริ่มต้นการ Log การเรียก API โดยใช้ logFinish
	logFinish := srv.logger.With(
		zap.String("apiName", "SearchOrder"),
		zap.String("SoNo", soNo),
		zap.String("OrderNo", orderNo),
	)

	defer func() {
		if r := recover(); r != nil {
			srv.logger.Error("🔥 Panic occurred in SearchOrder", zap.Any("panic", r))
			logFinish.Error("Panic", zap.Any("error", fmt.Errorf("unexpected panic: %v", r)))
		}
	}()

	// 📌 Log การเริ่มต้นค้นหาคำสั่งขาย
	srv.logger.Info("🔍 Searching for Sale Order",
		zap.String("SoNo", soNo),
		zap.String("OrderNo", orderNo),
	)

	// ✅ ตรวจสอบเงื่อนไข: ต้องมีค่า SoNo หรือ OrderNo อย่างน้อยหนึ่งค่า
	if soNo == "" && orderNo == "" {
		err := errors.New("either SoNo or OrderNo must be provided")
		srv.logger.Warn("⚠️ Invalid request - Missing parameters", zap.Error(err))
		logFinish.Warn("Invalid Request", zap.Error(err))
		return nil, err
	}

	// 🔍 ค้นหาข้อมูลคำสั่งขายจาก Repository Layer
	order, err := srv.orderRepo.SearchOrder(ctx, soNo, orderNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ✅ ไม่พบข้อมูลคำสั่งขาย
			srv.logger.Warn("⚠️ No Sale Order found", zap.String("SoNo", soNo), zap.String("OrderNo", orderNo))
			logFinish.Warn("Not Found", zap.String("error", "Sale order not found"))
			return nil, errors.New("Sale order not found")
		}

		// ❌ กรณีเกิดปัญหาอื่น ๆ เช่น Database ล่ม
		srv.logger.Error("❌ Failed to search Sale Order",
			zap.String("SoNo", soNo),
			zap.String("OrderNo", orderNo),
			zap.Error(err),
		)
		logFinish.Error("Failed", zap.String("error", "Failed to retrieve sale order"), zap.Error(err))
		return nil, fmt.Errorf("Failed to retrieve sale order: %w", err)
	}

	// ✅ ถ้าพบข้อมูลคำสั่งขาย
	logFinish.Info("✅ Sale Order found", zap.Int("TotalItems", len(order.Items)))

	return order, nil
}
