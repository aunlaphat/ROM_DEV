package service

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type OrderService interface {
	SearchOrder(ctx context.Context, req request.SearchOrder) (*response.SearchOrderResponse, error)
}

func (srv service) SearchOrder(ctx context.Context, req request.SearchOrder) (*response.SearchOrderResponse, error) {
	// ✅ Log Request
	srv.logger.Info("🔎 Searching for Order",
		zap.String("SoNo", req.SoNo),
		zap.String("OrderNo", req.OrderNo),
	)

	// ✅ ตรวจสอบว่ามีค่า SoNo หรือ OrderNo
	if req.SoNo == "" && req.OrderNo == "" {
		err := errors.New("either SoNo or OrderNo must be provided")
		srv.logger.Warn("⚠️ Invalid request - Missing parameters", zap.Error(err))
		return nil, err
	}

	// 🛠 Call Repository Layer (No Logging in Repo)
	order, err := srv.orderRepo.SearchOrder(ctx, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ✅ Log เมื่อไม่พบข้อมูล
			srv.logger.Warn("⚠️ No Sale Order found",
				zap.String("SoNo", req.SoNo),
				zap.String("OrderNo", req.OrderNo),
			)
			return nil, sql.ErrNoRows
		}

		// ❌ Log เมื่อ Query ล้มเหลว
		srv.logger.Error("❌ Failed to search Order",
			zap.String("SoNo", req.SoNo),
			zap.String("OrderNo", req.OrderNo),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to retrieve order: %w", err)
	}

	// ✅ Log สำเร็จ
	srv.logger.Info("✅ Order found",
		zap.String("SoNo", order.SoNo),
		zap.String("OrderNo", order.OrderNo),
		zap.Int("TotalItems", len(order.Items)),
	)

	return order, nil
}
