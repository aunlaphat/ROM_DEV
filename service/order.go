package service

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// เพิ่ม constant สำหรับ timeout
const (
	defaultTimeout = 10 * time.Second
	txTimeout      = 30 * time.Second
)

// ReturnOrderService คือ interface ที่กำหนดความสามารถของ service
type ReturnOrderService interface {
	CreateOrderWithTransaction(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
}

// Implementation สำหรับ CreateOrderWithTransaction
func (srv service) CreateOrderWithTransaction(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	// สร้าง request ID สำหรับการ tracking
	requestID := uuid.New().String()

	// เริ่มต้น logging
	srv.logger.Info("Starting create order with transaction process",
		zap.String("requestID", requestID),
		zap.String("orderNo", req.OrderNo))

	// เพิ่ม context timeout
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()

	// เริ่มต้น transaction
	tx, err := srv.returnOrderRepo.db.BeginTxx(ctx, nil)
	if err != nil {
		srv.logger.Error("Failed to begin transaction",
			zap.String("requestID", requestID),
			zap.Error(err))
		return nil, err
	}

	// Handle panic และ rollback transaction ถ้ามี panic เกิดขึ้น
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			srv.logger.Error("Panic occurred! Transaction rolled back",
				zap.String("requestID", requestID))
			panic(p)
		}
	}()

	// 1. ตรวจสอบข้อมูลพื้นฐาน
	if err := srv.validateOrder(req); err != nil {
		srv.logger.Error("Validation failed",
			zap.String("requestID", requestID),
			zap.Error(err))
		tx.Rollback()
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// 2. สร้าง order header
	if err := srv.returnOrderRepo.CreateBeforeReturnOrder(ctx, req); err != nil {
		srv.logger.Error("Failed to create order",
			zap.String("requestID", requestID),
			zap.Error(err))
		tx.Rollback()
		return nil, err
	}

	// 3. สร้าง order lines
	if err := srv.returnOrderRepo.CreateBeforeReturnOrderLine(ctx, req.OrderNo, req.BeforeReturnOrderLines); err != nil {
		srv.logger.Error("Failed to create order lines",
			zap.String("requestID", requestID),
			zap.Error(err))
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		srv.logger.Error("Failed to commit transaction",
			zap.String("requestID", requestID),
			zap.Error(err))
		return nil, err
	}

	// 4. ดึงข้อมูลที่สร้างเพื่อส่งกลับ
	order, err := srv.returnOrderRepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("Failed to get created order",
			zap.String("requestID", requestID),
			zap.Error(err))
		return nil, err
	}

	srv.logger.Info("Successfully created order with transaction",
		zap.String("requestID", requestID),
		zap.String("orderNo", req.OrderNo))

	// 5. สร้าง response
	return &response.BeforeReturnOrderResponse{
		Success: true,
		OrderNo: req.OrderNo,
		Message: "Return order created successfully",
		Data:    order,
	}, nil
}

// ...existing code...
