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
	srv.logger.Info("🔎 Searching for Order",
		zap.String("SoNo", req.SoNo),
		zap.String("OrderNo", req.OrderNo),
	)

	if req.SoNo == "" && req.OrderNo == "" {
		err := errors.New("either SoNo or OrderNo must be provided")
		srv.logger.Warn("⚠️ Invalid request - Missing parameters", zap.Error(err))
		return nil, err
	}

	order, err := srv.orderRepo.SearchOrder(ctx, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			srv.logger.Warn("⚠️ No Sale Order found",
				zap.String("SoNo", req.SoNo),
				zap.String("OrderNo", req.OrderNo),
			)
			return nil, sql.ErrNoRows
		}

		srv.logger.Error("❌ Failed to search Order",
			zap.String("SoNo", req.SoNo),
			zap.String("OrderNo", req.OrderNo),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to retrieve order: %w", err)
	}

	srv.logger.Info("✅ Order found",
		zap.String("SoNo", order.SoNo),
		zap.String("OrderNo", order.OrderNo),
		zap.Int("TotalItems", len(order.Items)),
	)

	return order, nil
}
