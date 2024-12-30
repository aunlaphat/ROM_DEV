package service

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"context"
	"fmt"
	"database/sql"

	"go.uber.org/zap"
)

type BefROService interface {
	CreateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error)
	GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error)
	ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error)
	GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error)
	UpdateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
<<<<<<< HEAD

	GetAllOrderDetail() ([]response.OrderDetail, error)
	GetOrderDetailBySO(soNo string) (*response.OrderDetail, error)
	DeleteBeforeReturnOrderLine(recID string) error
=======
	SearchSaleOrder(ctx context.Context, soNo string) ([]response.SaleOrderResponse, error)
>>>>>>> rom
}

func (srv service) CreateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Debug("🚀 Starting CreateOrderWithLines", zap.String("OrderNo", req.OrderNo))
	err := srv.befRORepo.CreateReturnOrderWithTransaction(ctx, req)
	if err != nil {
		srv.logger.Error("❌ Failed to create order with lines", zap.Error(err))
		return nil, err
	}

	// Fetch the created order to ensure all fields are correctly populated
	createdOrder, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to fetch created order", zap.Error(err))
		return nil, err
	}

	srv.logger.Debug("✅ Successfully created order with lines", zap.String("OrderNo", req.OrderNo))
	return createdOrder, nil
}

func (srv service) UpdateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Debug("🚀 Starting UpdateBeforeReturnOrderWithLines", zap.String("OrderNo", req.OrderNo))
	err := srv.befRORepo.UpdateBeforeReturnOrderWithTransaction(ctx, req)
	if err != nil {
		srv.logger.Error("❌ Failed to update order with lines", zap.Error(err))
		return nil, err
	}

	// Fetch the updated order to ensure all fields are correctly populated
	updatedOrder, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to fetch updated order", zap.Error(err))
		return nil, err
	}

	srv.logger.Debug("✅ Successfully updated order with lines", zap.String("OrderNo", req.OrderNo))
	return updatedOrder, nil
}

func (srv service) ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error) {
	srv.logger.Debug("🚀 Starting ListBeforeReturnOrders")
	orders, err := srv.befRORepo.ListBeforeReturnOrders(ctx)
	if err != nil {
		srv.logger.Error("❌ Failed to list return orders", zap.Error(err))
		return nil, err
	}
	srv.logger.Debug("✅ Successfully listed return orders", zap.Int("Count", len(orders)))
	return orders, nil
}

func (srv service) GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Debug("🚀 Starting GetBeforeReturnOrderByOrderNo", zap.String("OrderNo", orderNo))
	order, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to get return order by order number", zap.Error(err))
		return nil, err
	}
	srv.logger.Debug("✅ Successfully fetched return order", zap.String("OrderNo", orderNo))
	return order, nil
}

func (srv service) ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error) {
	srv.logger.Debug("🚀 Starting ListBeforeReturnOrderLines")
	lines, err := srv.befRORepo.ListBeforeReturnOrderLines(ctx)
	if err != nil {
		srv.logger.Error("❌ Failed to list return order lines", zap.Error(err))
		return nil, err
	}
	srv.logger.Debug("✅ Successfully listed return order lines", zap.Int("Count", len(lines)))
	return lines, nil
}

func (srv service) GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error) {
	srv.logger.Debug("🚀 Starting GetBeforeReturnOrderLineByOrderNo", zap.String("OrderNo", orderNo))
	lines, err := srv.befRORepo.GetBeforeReturnOrderLineByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to get return order lines by order number", zap.Error(err))
		return nil, err
	}
	srv.logger.Debug("✅ Successfully fetched return order lines", zap.String("OrderNo", orderNo))
	return lines, nil
}

<<<<<<< HEAD
// service เชื่อมกับ repo ต่อเพื่อดึงข้อมูลออกมา แต่ต้องมีการ validation ก่อนดึง
func (srv service) 	GetAllOrderDetail() ([]response.OrderDetail, error) {
	allorder, err := srv.befRORepo.GetAllOrderDetail()
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			srv.logger.Error(err)
			return nil, fmt.Errorf("no order data: %w", err)
		default:
			srv.logger.Error(err)
			return nil, fmt.Errorf("get order error: %w", err)
		}
	}
	return allorder, nil
}

func (srv service) 	GetOrderDetailBySO(soNo string) (*response.OrderDetail, error) {
	soOrder, err := srv.befRORepo.GetOrderDetailBySO(soNo)
	if err != nil {
		return nil, err
	}
	return soOrder, nil
}

func (srv service) DeleteBeforeReturnOrderLine(recID string) error {
	if recID == "" {
		return fmt.Errorf("RecID is required")
	}

	// ส่งไปยัง Repository Layer
	err := srv.befRORepo.DeleteBeforeReturnOrderLine(recID)
	if err != nil {
		return fmt.Errorf("failed to delete before return order line: %w", err)
	}

	return nil
}
=======
// Implementation สำหรับ SearchSaleOrder
func (srv service) SearchSaleOrder(ctx context.Context, soNo string) ([]response.SaleOrderResponse, error) {
	srv.logger.Debug("🚀 Starting SearchSaleOrder", zap.String("SoNo", soNo))
	orders, err := srv.befRORepo.SearchSaleOrder(ctx, soNo)
	if err != nil {
		srv.logger.Error("❌ Failed to search sale orders", zap.Error(err))
		return nil, err
	}
	if orders == nil {
		srv.logger.Debug("❗ No sale order found", zap.String("SoNo", soNo))
		return nil, nil
	}
	srv.logger.Debug("✅ Successfully searched sale orders", zap.String("SoNo", soNo))
	return orders, nil
}
>>>>>>> rom
