package service

import (
	entity "boilerplate-backend-go/Entity"
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"context"
	"fmt"

	"go.uber.org/zap"
)

type DraftConfirmService interface {
	GetOrders(ctx context.Context, statusConfID int, startDate, endDate string) ([]response.OrderHeadResponse, error)
	GetOrderWithItems(ctx context.Context, orderNo string) (*response.DraftConfirmResponse, error)
	AddItemToDraftOrder(ctx context.Context, orderNo string, req request.AddItem, userID string) error
	RemoveItemFromDraftOrder(ctx context.Context, orderNo, sku string) error
	ConfirmDraftOrder(ctx context.Context, orderNo string, userID string) error
}

// ✅ ดึงรายการ Draft หรือ Confirm Orders ตาม `StatusConfID`
func (srv service) GetOrders(ctx context.Context, statusConfID int, startDate, endDate string) ([]response.OrderHeadResponse, error) {
	srv.logger.Info("📄 Fetching Orders (HEAD)",
		zap.Int("StatusConfID", statusConfID),
		zap.String("StartDate", startDate),
		zap.String("EndDate", endDate),
	)

	orders, err := srv.draftConfirmRepo.GetOrders(ctx, statusConfID, startDate, endDate)
	if err != nil {
		srv.logger.Error("❌ Failed to fetch orders", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}

	srv.logger.Info("✅ Orders fetched successfully", zap.Int("TotalOrders", len(orders)))
	return orders, nil
}

func (srv service) GetOrderWithItems(ctx context.Context, orderNo string) (*response.DraftConfirmResponse, error) {
	srv.logger.Info("📦 Fetching Order with Items", zap.String("OrderNo", orderNo))

	order, err := srv.draftConfirmRepo.GetOrderWithItems(ctx, orderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to fetch order", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}

	srv.logger.Info("✅ Order with Items fetched successfully", zap.String("OrderNo", order.OrderNo))
	return order, nil
}

// ✅ เพิ่มสินค้าเข้า Draft Order
func (srv service) AddItemToDraftOrder(ctx context.Context, orderNo string, req request.AddItem, userID string) error {
	srv.logger.Info("➕ Adding Item to Draft Order",
		zap.String("OrderNo", orderNo),
		zap.String("SKU", req.SKU),
		zap.Int("QTY", req.QTY),
		zap.Float64("Price", req.Price),
		zap.String("CreatedBy", userID),
	)

	// สร้าง `entity.BeforeReturnOrderLine`
	item := entity.BeforeReturnOrderLine{
		OrderNo:   orderNo,
		SKU:       req.SKU,
		ItemName:  req.ItemName,
		QTY:       req.QTY,
		ReturnQTY: req.QTY, // ค่า Default = จำนวนที่สั่ง
		Price:     req.Price,
		CreateBy:  userID,
	}

	// บันทึกลง DB
	err := srv.draftConfirmRepo.AddItemToDraftOrder(ctx, orderNo, item)
	if err != nil {
		srv.logger.Error("❌ Failed to add item to draft order", zap.String("OrderNo", orderNo), zap.Error(err))
		return fmt.Errorf("failed to add item: %w", err)
	}

	srv.logger.Info("✅ Item added successfully", zap.String("OrderNo", orderNo), zap.String("SKU", req.SKU))
	return nil
}

// ✅ ลบสินค้าออกจาก Draft Order
func (srv service) RemoveItemFromDraftOrder(ctx context.Context, orderNo, sku string) error {
	srv.logger.Info("❌ Removing Item from Draft Order",
		zap.String("OrderNo", orderNo),
		zap.String("SKU", sku),
	)

	// ลบรายการสินค้า
	err := srv.draftConfirmRepo.RemoveItemFromDraftOrder(ctx, orderNo, sku)
	if err != nil {
		srv.logger.Error("❌ Failed to remove item from draft order", zap.String("OrderNo", orderNo), zap.Error(err))
		return fmt.Errorf("failed to remove item: %w", err)
	}

	srv.logger.Info("✅ Item removed successfully", zap.String("OrderNo", orderNo), zap.String("SKU", sku))
	return nil
}

// ✅ อัปเดต Draft Order เป็น Confirm
func (srv service) ConfirmDraftOrder(ctx context.Context, orderNo string, userID string) error {
	srv.logger.Info("🔄 Confirming Draft Order",
		zap.String("OrderNo", orderNo),
		zap.String("UpdatedBy", userID),
	)

	// ✅ ใช้ `UpdateOrderStatus` จาก `repository/order.go`
	err := srv.orderRepo.UpdateOrderStatus(ctx, orderNo, 3, 2, userID) // (StatusReturnID = 3, StatusConfID = 2)
	if err != nil {
		srv.logger.Error("❌ Failed to confirm draft order", zap.String("OrderNo", orderNo), zap.Error(err))
		return fmt.Errorf("failed to confirm draft order: %w", err)
	}

	srv.logger.Info("✅ Draft Order confirmed", zap.String("OrderNo", orderNo))
	return nil
}
