package service

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type DraftConfirmService interface {
	GetOrders(ctx context.Context, statusConfID int, startDate, endDate string) ([]response.OrderHeadResponse, error)
	GetOrderWithItems(ctx context.Context, orderNo string, statusConfID int) (*response.DraftConfirmResponse, error)
	ListCodeR(ctx context.Context) ([]response.ListCodeRResponse, error)
	AddItemToDraftOrder(ctx context.Context, req request.AddItem, userID string) ([]response.AddItemResponse, error)
	RemoveItemFromDraftOrder(ctx context.Context, orderNo, sku string) error
	ConfirmDraftOrder(ctx context.Context, orderNo string, userID string) (*response.UpdateOrderStatusResponse, error)
}

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

func (srv service) GetOrderWithItems(ctx context.Context, orderNo string, statusConfID int) (*response.DraftConfirmResponse, error) {
	srv.logger.Info("📦 Fetching Order with Items",
		zap.String("OrderNo", orderNo),
		zap.Int("StatusConfID", statusConfID),
	)

	order, err := srv.draftConfirmRepo.GetOrderWithItems(ctx, orderNo, statusConfID)
	if err != nil {
		srv.logger.Error("❌ Failed to fetch order", zap.String("OrderNo", orderNo), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}

	if order.Items == nil {
		order.Items = []response.DraftConfirmItem{}
	}

	srv.logger.Info("✅ Order with Items fetched successfully",
		zap.String("OrderNo", order.OrderNo),
		zap.Int("StatusConfID", statusConfID),
	)
	return order, nil
}

func (srv service) ListCodeR(ctx context.Context) ([]response.ListCodeRResponse, error) {
	srv.logger.Info("📦 Fetching List of CodeR")

	codeRList, err := srv.draftConfirmRepo.ListCodeR(ctx)
	if err != nil {
		srv.logger.Error("❌ Failed to fetch CodeR list", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch CodeR list: %w", err)
	}

	srv.logger.Info("✅ CodeR list fetched successfully", zap.Int("TotalItems", len(codeRList)))
	return codeRList, nil
}

func (srv service) AddItemToDraftOrder(ctx context.Context, req request.AddItem, userID string) ([]response.AddItemResponse, error) {
	srv.logger.Info("➕ Adding Item to Draft Order",
		zap.String("OrderNo", req.OrderNo),
		zap.String("SKU", req.SKU),
		zap.Int("QTY", req.QTY),
		zap.Float64("Price", req.Price),
		zap.String("CreatedBy", userID),
	)

	if req.ReturnQTY == 0 {
		req.ReturnQTY = req.QTY
	}

	results, err := srv.draftConfirmRepo.AddItemToDraftOrder(ctx, req, userID)
	if err != nil {
		srv.logger.Error("❌ Failed to add item to draft order", zap.String("OrderNo", req.OrderNo), zap.Error(err))
		return nil, fmt.Errorf("failed to add item to draft order: %w", err)
	}

	srv.logger.Info("✅ Item added successfully", zap.String("OrderNo", req.OrderNo), zap.String("SKU", req.SKU))
	return results, nil
}

func (srv service) RemoveItemFromDraftOrder(ctx context.Context, orderNo, sku string) error {
	srv.logger.Info("❌ Removing Item from Draft Order",
		zap.String("OrderNo", orderNo),
		zap.String("SKU", sku),
	)

	rowsAffected, err := srv.draftConfirmRepo.RemoveItemFromDraftOrder(ctx, orderNo, sku)
	if err != nil {
		srv.logger.Error("❌ Failed to remove item", zap.String("OrderNo", orderNo), zap.Error(err))
		return fmt.Errorf("failed to remove item: %w", err)
	}

	if rowsAffected == 0 {
		srv.logger.Warn("⚠️ No item found to delete", zap.String("OrderNo", orderNo), zap.String("SKU", sku))
		return fmt.Errorf("no item found to delete")
	}

	srv.logger.Info("✅ Item removed successfully", zap.String("OrderNo", orderNo), zap.String("SKU", sku))
	return nil
}

func (srv service) ConfirmDraftOrder(ctx context.Context, orderNo string, userID string) (*response.UpdateOrderStatusResponse, error) {
	srv.logger.Info("🔄 Confirming Draft Order",
		zap.String("OrderNo", orderNo),
		zap.String("UpdatedBy", userID),
	)

	order, err := srv.draftConfirmRepo.GetOrderWithItems(ctx, orderNo, 1) // 1 = Draft
	if err != nil {
		srv.logger.Error("❌ Failed to fetch order", zap.String("OrderNo", orderNo), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}
	if len(order.Items) == 0 {
		srv.logger.Warn("⚠️ Cannot confirm order without items", zap.String("OrderNo", orderNo))
		return nil, fmt.Errorf("cannot confirm order without items")
	}

	err = srv.orderRepo.UpdateOrderStatus(ctx, orderNo, 3, 2, userID) // 3 = Booking, 2 = Confirm
	if err != nil {
		srv.logger.Error("❌ Failed to update order status", zap.String("OrderNo", orderNo), zap.Error(err))
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	srv.logger.Info("✅ Draft Order confirmed", zap.String("OrderNo", orderNo))

	return &response.UpdateOrderStatusResponse{
		OrderNo:        orderNo,
		StatusReturnID: 3, // Booking
		StatusConfID:   2, // Confirm
		ConfirmBy:      userID,
		ConfirmDate:    time.Now(),
	}, nil
}
