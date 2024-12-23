package service

import (
	request "boilerplate-backend-go/dto/request"
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
	// CRUD operations พื้นฐาน
	CreateOrder(ctx context.Context, req request.BeforeReturnOrderRequest) (*response.BeforeReturnOrderResponse, error)
	CreateOrderLine(ctx context.Context, orderNo string, lines []request.BeforeReturnOrderLineRequest) error
	GetOrder(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error)
	GetOrderLines(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineData, error)
	ListOrders(ctx context.Context, page, limit int) (*response.BeforeReturnOrderListResponse, error)
	UpdateStatus(ctx context.Context, orderNo string, statusID int, updateBy string) error
	CancelOrder(ctx context.Context, orderNo string, cancelBy string) error
	GetOrderByReturnID(ctx context.Context, returnID string) (response.ReturnOrder, error)
}

// ลบ GetAllOrder ออก และใช้ ListOrders แทน
func (srv service) ListOrders(ctx context.Context, page, limit int) (*response.BeforeReturnOrderListResponse, error) {
	requestID := uuid.New().String()
	srv.logger.Info("Starting list orders",
		zap.String("requestID", requestID),
		zap.Int("page", page),
		zap.Int("limit", limit))

	// คำนวณ offset
	offset := (page - 1) * limit

	// ดึงข้อมูล
	orders, err := srv.returnOrderRepo.ListBeforeReturnOrders(ctx, limit, offset)
	if err != nil {
		srv.logger.Error("Failed to list orders",
			zap.String("requestID", requestID),
			zap.Error(err))
		return nil, err
	}

	srv.logger.Info("Successfully retrieved orders",
		zap.String("requestID", requestID),
		zap.Int("totalOrders", len(orders)))

	return &response.BeforeReturnOrderListResponse{
		Page:   page,
		Limit:  limit,
		Orders: orders,
		Total:  len(orders),
	}, nil
}

// Implementation ของ GetOrder
func (srv service) GetOrder(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error) {
	requestID := uuid.New().String()

	// ตรวจสอบ input
	if orderNo == "" {
		return nil, fmt.Errorf("order number is required")
	}

	// ดึงข้อมูล
	order, err := srv.returnOrderRepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("Failed to get order",
			zap.String("requestID", requestID),
			zap.Error(err))
		return nil, err
	}

	// ตรวจสอบว่าพบข้อมูลหรือไม่
	if order == nil {
		return &response.BeforeReturnOrderResponse{
			Success: false,
			OrderNo: orderNo,
			Message: "Order not found",
		}, nil
	}

	return &response.BeforeReturnOrderResponse{
		Success: true,
		OrderNo: orderNo,
		Data:    order,
	}, nil
}

// Implementation สำหรับ CreateOrder
func (srv service) CreateOrder(ctx context.Context, req request.BeforeReturnOrderRequest) (*response.BeforeReturnOrderResponse, error) {
	// สร้าง request ID สำหรับการ tracking
	requestID := uuid.New().String()

	// เริ่มต้น logging
	srv.logger.Info("Starting create order process",
		zap.String("requestID", requestID),
		zap.String("orderNo", req.OrderNo))

	// 1. ตรวจสอบข้อมูลพื้นฐาน
	if err := srv.validateOrder(req); err != nil {
		srv.logger.Error("Validation failed",
			zap.String("requestID", requestID),
			zap.Error(err))
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// 2. สร้าง order header
	if err := srv.returnOrderRepo.CreateBeforeReturnOrder(ctx, req); err != nil {
		srv.logger.Error("Failed to create order",
			zap.String("requestID", requestID),
			zap.Error(err))
		return nil, err
	}

	// 3. สร้าง order lines
	if err := srv.returnOrderRepo.CreateBeforeReturnOrderLine(ctx, req.OrderNo, req.ReturnLines); err != nil {
		srv.logger.Error("Failed to create order lines",
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

	srv.logger.Info("Successfully created order",
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

// Implementation สำหรับ CreateOrderLine
func (srv service) CreateOrderLine(ctx context.Context, orderNo string, lines []request.BeforeReturnOrderLineRequest) error {
	requestID := uuid.New().String()
	srv.logger.Info("Starting create order line process",
		zap.String("requestID", requestID),
		zap.String("orderNo", orderNo))

	// Call repository to create order lines
	err := srv.returnOrderRepo.CreateBeforeReturnOrderLine(ctx, orderNo, lines)
	if err != nil {
		srv.logger.Error("Failed to create order lines",
			zap.String("requestID", requestID),
			zap.Error(err))
		return fmt.Errorf("failed to create order lines: %w", err)
	}

	srv.logger.Info("Successfully created order lines",
		zap.String("requestID", requestID),
		zap.String("orderNo", orderNo))

	return nil
}

// Implementation ของ GetOrderLines
func (srv service) GetOrderLines(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineData, error) {
	requestID := uuid.New().String()
	srv.logger.Info("Starting get order lines process",
		zap.String("requestID", requestID),
		zap.String("orderNo", orderNo))

	// Call repository to get order lines
	lines, err := srv.returnOrderRepo.GetBeforeReturnOrderLines(ctx, orderNo)
	if err != nil {
		srv.logger.Error("Failed to get order lines",
			zap.String("requestID", requestID),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get order lines: %w", err)
	}

	srv.logger.Info("Successfully retrieved order lines",
		zap.String("requestID", requestID),
		zap.String("orderNo", orderNo))

	return lines, nil
}

// Implementation สำหรับ UpdateStatus
func (srv service) UpdateStatus(ctx context.Context, orderNo string, statusID int, updateBy string) error {
	requestID := uuid.New().String()
	srv.logger.Info("Starting update order status process",
		zap.String("requestID", requestID),
		zap.String("orderNo", orderNo),
		zap.Int("statusID", statusID),
		zap.String("updateBy", updateBy))

	// Validate input
	if orderNo == "" {
		return fmt.Errorf("order number is required")
	}
	if updateBy == "" {
		return fmt.Errorf("update by is required")
	}

	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// Call repository to update status
	err := srv.returnOrderRepo.UpdateBeforeReturnOrderStatus(ctx, orderNo, statusID, updateBy)
	if err != nil {
		srv.logger.Error("Failed to update order status",
			zap.String("requestID", requestID),
			zap.Error(err))
		return fmt.Errorf("failed to update order status: %w", err)
	}

	srv.logger.Info("Successfully updated order status",
		zap.String("requestID", requestID),
		zap.String("orderNo", orderNo))

	return nil
}

// Implementation สำหรับ CancelOrder
func (srv service) CancelOrder(ctx context.Context, orderNo string, cancelBy string) error {
	requestID := uuid.New().String()
	srv.logger.Info("Starting cancel order process",
		zap.String("requestID", requestID),
		zap.String("orderNo", orderNo),
		zap.String("cancelBy", cancelBy))

	// Validate input
	if orderNo == "" {
		return fmt.Errorf("order number is required")
	}
	if cancelBy == "" {
		return fmt.Errorf("cancel by is required")
	}

	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// Call repository to cancel order
	err := srv.returnOrderRepo.CancelBeforeReturnOrder(ctx, orderNo, cancelBy)
	if err != nil {
		srv.logger.Error("Failed to cancel order",
			zap.String("requestID", requestID),
			zap.Error(err))
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	srv.logger.Info("Successfully cancelled order",
		zap.String("requestID", requestID),
		zap.String("orderNo", orderNo))

	return nil
}

// Implementation สำหรับ GetOrderByReturnID
func (srv service) GetOrderByReturnID(ctx context.Context, returnID string) (response.ReturnOrder, error) {
	requestID := uuid.New().String()
	srv.logger.Info("Starting get order by return ID process",
		zap.String("requestID", requestID),
		zap.String("returnID", returnID))

	// Call repository to get order by return ID
	order, err := srv.returnOrderRepo.GetOrderByReturnID(ctx, returnID)
	if err != nil {
		srv.logger.Error("Failed to get order by return ID",
			zap.String("requestID", requestID),
			zap.Error(err))
		return response.ReturnOrder{}, fmt.Errorf("failed to get order by return ID: %w", err)
	}

	srv.logger.Info("Successfully retrieved order by return ID",
		zap.String("requestID", requestID),
		zap.String("returnID", returnID))

	return order, nil
}

// Validation helper
func (srv service) validateOrder(req request.BeforeReturnOrderRequest) error {
	// ตรวจสอบข้อมูลจำเป็นพื้นฐาน
	if req.OrderNo == "" {
		return fmt.Errorf("order number is required")
	}
	if req.CustomerID == "" {
		return fmt.Errorf("customer ID is required")
	}
	if len(req.ReturnLines) == 0 {
		return fmt.Errorf("at least one return line is required")
	}

	// ตรวจสอบรายการสินค้า
	for i, line := range req.ReturnLines {
		if line.SKU == "" {
			return fmt.Errorf("line %d: SKU is required", i+1)
		}
		if line.ReturnQTY <= 0 {
			return fmt.Errorf("line %d: return quantity must be greater than 0", i+1)
		}
		if line.ReturnQTY > line.QTY {
			return fmt.Errorf("line %d: return quantity cannot be greater than original quantity", i+1)
		}
	}

	return nil
}
