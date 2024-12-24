package service

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"context"

	"go.uber.org/zap"
)

type ReturnOrderService interface {
	CreateOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error)
	GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error)
	ListBeforeReturnOrderLines(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error)
	GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderLineResponse, error)
}

func (srv service) CreateOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Debug("🚀 Starting CreateOrderWithLines", zap.String("OrderNo", req.OrderNo))
	err := srv.returnOrderRepo.CreateReturnOrderWithTransaction(ctx, req)
	if err != nil {
		srv.logger.Error("❌ Failed to create order with lines", zap.Error(err))
		return nil, err
	}

	returnOrderResponse := &response.BeforeReturnOrderResponse{
		OrderNo:                req.OrderNo,
		SaleOrder:              req.SaleOrder,
		SaleReturn:             req.SaleReturn,
		ChannelID:              req.ChannelID,
		ReturnType:             req.ReturnType,
		CustomerID:             req.CustomerID,
		TrackingNo:             req.TrackingNo,
		Logistic:               req.Logistic,
		WarehouseID:            req.WarehouseID,
		StatusReturnID:         req.StatusReturnID,
		StatusConfID:           req.StatusConfID,
		CreateDate:             *req.CreateDate,
		CreateBy:               req.CreateBy,
		BeforeReturnOrderLines: make([]response.BeforeReturnOrderLineResponse, len(req.BeforeReturnOrderLines)),
	}

	for i, line := range req.BeforeReturnOrderLines {
		returnOrderResponse.BeforeReturnOrderLines[i] = response.BeforeReturnOrderLineResponse{
			OrderNo:    line.OrderNo,
			SKU:        line.SKU,
			QTY:        line.QTY,
			ReturnQTY:  line.ReturnQTY,
			Price:      line.Price,
			TrackingNo: line.TrackingNo,
			CreateDate: *line.CreateDate,
		}
	}

	srv.logger.Debug("✅ Successfully created order with lines", zap.String("OrderNo", req.OrderNo))
	return returnOrderResponse, nil
}

func (srv service) ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error) {
	srv.logger.Debug("🚀 Starting ListBeforeReturnOrders")
	orders, err := srv.returnOrderRepo.ListBeforeReturnOrders(ctx)
	if err != nil {
		srv.logger.Error("❌ Failed to list return orders", zap.Error(err))
		return nil, err
	}
	srv.logger.Debug("✅ Successfully listed return orders", zap.Int("Count", len(orders)))
	return orders, nil
}

func (srv service) GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Debug("🚀 Starting GetBeforeReturnOrderByOrderNo", zap.String("OrderNo", orderNo))
	order, err := srv.returnOrderRepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to get return order by order number", zap.Error(err))
		return nil, err
	}
	srv.logger.Debug("✅ Successfully fetched return order", zap.String("OrderNo", orderNo))
	return order, nil
}

func (srv service) ListBeforeReturnOrderLines(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error) {
	srv.logger.Debug("🚀 Starting ListBeforeReturnOrderLines", zap.String("OrderNo", orderNo))
	lines, err := srv.returnOrderRepo.ListBeforeReturnOrderLines(ctx, orderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to list return order lines", zap.Error(err))
		return nil, err
	}
	srv.logger.Debug("✅ Successfully listed return order lines", zap.Int("Count", len(lines)))
	return lines, nil
}

func (srv service) GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderLineResponse, error) {
	srv.logger.Debug("🚀 Starting GetBeforeReturnOrderLineByOrderNo", zap.String("OrderNo", orderNo))
	line, err := srv.returnOrderRepo.GetBeforeReturnOrderLineByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to get return order line by order number", zap.Error(err))
		return nil, err
	}
	srv.logger.Debug("✅ Successfully fetched return order line", zap.String("OrderNo", orderNo))
	return line, nil
}
