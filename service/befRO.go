package service

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"context"

	"go.uber.org/zap"
)

type BefROService interface {
	CreateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error)
	GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error)
	ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error)
	GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error)
	UpdateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	SearchSaleOrder(ctx context.Context, soNo string) ([]response.SaleOrderResponse, error)
	ConfirmSaleReturn(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	UpdateSrNo(ctx context.Context, orderNo string, srNo string) error
	UpdateDynamicFields(ctx context.Context, orderNo string, fields map[string]interface{}) error
	ListDrafts(ctx context.Context) ([]response.BeforeReturnOrderResponse, error)
	EditOrder(ctx context.Context, req request.EditOrderRequest) (*response.BeforeReturnOrderResponse, error)
	ConfirmOrder(ctx context.Context, req request.ConfirmOrderRequest) (*response.BeforeReturnOrderResponse, error)
}

func (srv service) CreateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🏁 Starting order creation process", zap.String("OrderNo", req.OrderNo))
	srv.logger.Debug("Creating order head", zap.String("OrderNo", req.OrderNo), zap.String("SoNo", req.SoNo))

	err := srv.befRORepo.CreateReturnOrderWithTransaction(ctx, req)
	if err != nil {
		srv.logger.Error("❌ Failed to create order with lines", zap.Error(err))
		return nil, err
	}

	createdOrder, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to fetch created order", zap.Error(err))
		return nil, err
	}

	srv.logger.Info("✅ Successfully created order with lines",
		zap.String("OrderNo", req.OrderNo),
		zap.Any("CreatedOrder", createdOrder))
	return createdOrder, nil
}

func (srv service) UpdateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🏁 Starting order update process", zap.String("OrderNo", req.OrderNo))
	srv.logger.Debug("Updating order head", zap.String("OrderNo", req.OrderNo), zap.String("SoNo", req.SoNo))

	err := srv.befRORepo.UpdateBeforeReturnOrderWithTransaction(ctx, req)
	if err != nil {
		srv.logger.Error("❌ Failed to update order with lines", zap.Error(err))
		return nil, err
	}

	updatedOrder, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to fetch updated order", zap.Error(err))
		return nil, err
	}

	srv.logger.Info("✅ Successfully updated order with lines",
		zap.String("OrderNo", req.OrderNo),
		zap.Any("UpdatedOrder", updatedOrder))
	return updatedOrder, nil
}

func (srv service) ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🏁 Starting to list all return orders")
	orders, err := srv.befRORepo.ListBeforeReturnOrders(ctx)
	if err != nil {
		srv.logger.Error("❌ Failed to list return orders", zap.Error(err))
		return nil, err
	}
	srv.logger.Info("✅ Successfully listed return orders",
		zap.Int("Count", len(orders)),
		zap.Any("Orders", orders))
	return orders, nil
}

func (srv service) GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🏁 Starting to get return order by order number", zap.String("OrderNo", orderNo))
	order, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to get return order by order number", zap.Error(err))
		return nil, err
	}
	srv.logger.Info("✅ Successfully fetched return order",
		zap.String("OrderNo", orderNo),
		zap.Any("Order", order))
	return order, nil
}

func (srv service) ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error) {
	srv.logger.Info("🏁 Starting to list all return order lines")
	lines, err := srv.befRORepo.ListBeforeReturnOrderLines(ctx)
	if err != nil {
		srv.logger.Error("❌ Failed to list return order lines", zap.Error(err))
		return nil, err
	}
	srv.logger.Info("✅ Successfully listed return order lines",
		zap.Int("Count", len(lines)),
		zap.Any("Lines", lines))
	return lines, nil
}

func (srv service) GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error) {
	srv.logger.Info("🏁 Starting to get return order lines by order number", zap.String("OrderNo", orderNo))
	lines, err := srv.befRORepo.GetBeforeReturnOrderLineByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to get return order lines by order number", zap.Error(err))
		return nil, err
	}
	srv.logger.Info("✅ Successfully fetched return order lines",
		zap.String("OrderNo", orderNo),
		zap.Any("Lines", lines))
	return lines, nil
}

func (srv service) SearchSaleOrder(ctx context.Context, soNo string) ([]response.SaleOrderResponse, error) {
	srv.logger.Info("🏁 Starting to search sale order", zap.String("SoNo", soNo))
	order, err := srv.befRORepo.SearchSaleOrder(ctx, soNo)
	if err != nil {
		srv.logger.Error("❌ Failed to search sale orders", zap.Error(err))
		return nil, err
	}
	if order == nil {
		srv.logger.Info("❗ No sale order found", zap.String("SoNo", soNo))
		return nil, nil
	}
	srv.logger.Info("✅ Successfully searched sale orders",
		zap.String("SoNo", soNo),
		zap.Any("Order", order))
	return []response.SaleOrderResponse{*order}, nil
}

func (srv service) UpdateSrNo(ctx context.Context, orderNo string, srNo string) error {
	srv.logger.Info("🏁 Starting to update SR number", zap.String("OrderNo", orderNo), zap.String("SrNo", srNo))

	err := srv.befRORepo.UpdateSrNo(ctx, orderNo, srNo)
	if err != nil {
		srv.logger.Error("❌ Failed to update SR number", zap.Error(err))
		return err
	}

	srv.logger.Info("✅ Successfully updated SR number", zap.String("OrderNo", orderNo), zap.String("SrNo", srNo))
	return nil
}

func (srv service) ConfirmSaleReturn(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🏁 Starting to confirm sale return", zap.String("OrderNo", req.OrderNo))

	fields := map[string]interface{}{
		"StatusConfID": req.StatusConfID,
		"ConfirmBy":    req.ConfirmBy,
		"UpdateBy":     req.UpdateBy,
	}

	err := srv.befRORepo.UpdateDynamicFields(ctx, req.OrderNo, fields)
	if err != nil {
		srv.logger.Error("❌ Failed to confirm sale return", zap.Error(err))
		return nil, err
	}

	confirmedOrder, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to fetch confirmed order", zap.Error(err))
		return nil, err
	}

	srv.logger.Info("✅ Successfully confirmed sale return",
		zap.String("OrderNo", req.OrderNo),
		zap.Any("ConfirmedOrder", confirmedOrder))
	return confirmedOrder, nil
}

func (srv service) UpdateDynamicFields(ctx context.Context, orderNo string, fields map[string]interface{}) error {
	srv.logger.Info("🏁 Starting to update dynamic fields", zap.String("OrderNo", orderNo), zap.Any("Fields", fields))

	err := srv.befRORepo.UpdateDynamicFields(ctx, orderNo, fields)
	if err != nil {
		srv.logger.Error("❌ Failed to update dynamic fields", zap.Error(err))
		return err
	}

	srv.logger.Info("✅ Successfully updated dynamic fields", zap.String("OrderNo", orderNo), zap.Any("Fields", fields))
	return nil
}

func (srv service) ListDrafts(ctx context.Context) ([]response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🏁 Starting to list all drafts")
	drafts, err := srv.befRORepo.ListDrafts(ctx)
	if err != nil {
		srv.logger.Error("❌ Failed to list drafts", zap.Error(err))
		return nil, err
	}
	srv.logger.Info("✅ Successfully listed drafts",
		zap.Int("Count", len(drafts)),
		zap.Any("Drafts", drafts))
	return drafts, nil
}

func (srv service) EditOrder(ctx context.Context, req request.EditOrderRequest) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🏁 Starting order edit process", zap.String("OrderNo", req.OrderNo))
	srv.logger.Debug("Editing order", zap.String("OrderNo", req.OrderNo), zap.String("SoNo", req.SoNo))

	err := srv.befRORepo.EditOrder(ctx, req)
	if err != nil {
		srv.logger.Error("❌ Failed to edit order", zap.Error(err))
		return nil, err
	}

	editedOrder, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to fetch edited order", zap.Error(err))
		return nil, err
	}

	srv.logger.Info("✅ Successfully edited order",
		zap.String("OrderNo", req.OrderNo),
		zap.Any("EditedOrder", editedOrder))
	return editedOrder, nil
}

func (srv service) ConfirmOrder(ctx context.Context, req request.ConfirmOrderRequest) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🏁 Starting to confirm order", zap.String("OrderNo", req.OrderNo))

	fields := map[string]interface{}{
		"StatusConfID": req.StatusConfID,
		"ConfirmBy":    req.ConfirmBy,
	}

	err := srv.befRORepo.UpdateDynamicFields(ctx, req.OrderNo, fields)
	if err != nil {
		srv.logger.Error("❌ Failed to confirm order", zap.Error(err))
		return nil, err
	}

	confirmedOrder, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to fetch confirmed order", zap.Error(err))
		return nil, err
	}

	srv.logger.Info("✅ Successfully confirmed order",
		zap.String("OrderNo", req.OrderNo),
		zap.Any("ConfirmedOrder", confirmedOrder))
	return confirmedOrder, nil
}
