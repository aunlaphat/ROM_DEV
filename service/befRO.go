package service

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"context"
	"fmt"

	"go.uber.org/zap"
)

type BefROService interface {
	CreateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error)
	GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error)
	ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error)
	GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error)
	UpdateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)

	// ************************ Create Sale Return ************************ //
	SearchSaleOrder(ctx context.Context, soNo string) ([]response.SaleOrderResponse, error)
	CreateSaleReturn(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	UpdateSaleReturn(ctx context.Context, orderNo string, srNo string, updateBy string) error
	ConfirmSaleReturn(ctx context.Context, orderNo string, confirmBy string) error
	CancelSaleReturn(ctx context.Context, orderNo string, updateBy string, remark string) error
}

func (srv service) CreateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🏁 Starting order creation process", zap.String("OrderNo", req.OrderNo))
	srv.logger.Debug("Creating order head", zap.String("OrderNo", req.OrderNo), zap.String("SoNo", req.SoNo))

	err := srv.befRORepo.CreateBeforeReturnOrderWithTransaction(ctx, req)
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
		zap.String("OrderNo", req.OrderNo))
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
		zap.String("OrderNo", req.OrderNo))
	return updatedOrder, nil
}

func (srv service) ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🏁 Starting to list all return orders")
	orders, err := srv.befRORepo.ListBeforeReturnOrders(ctx)
	if err != nil {
		srv.logger.Error("❌ Failed to list return orders", zap.Error(err))
		return nil, err
	}
	srv.logger.Info("✅ Successfully listed return orders")
	return orders, nil
}

func (srv service) GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🏁 Starting to get return order by order number", zap.String("OrderNo", orderNo))
	order, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to get return order by order number", zap.Error(err))
		return nil, err
	}
	srv.logger.Info("✅ Successfully fetched return order")
	return order, nil
}

func (srv service) ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error) {
	srv.logger.Info("🏁 Starting to list all return order lines")
	lines, err := srv.befRORepo.ListBeforeReturnOrderLines(ctx)
	if err != nil {
		srv.logger.Error("❌ Failed to list return order lines", zap.Error(err))
		return nil, err
	}
	srv.logger.Info("✅ Successfully listed return order lines")
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
		zap.String("OrderNo", orderNo))
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
		zap.String("SoNo", soNo))
	return []response.SaleOrderResponse{*order}, nil
}

func (srv service) CreateSaleReturn(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🏁 Starting sale return creation process", zap.String("OrderNo", req.OrderNo))
	srv.logger.Debug("Creating sale return order", zap.String("OrderNo", req.OrderNo), zap.String("SoNo", req.SoNo))

	createdOrder, err := srv.befRORepo.CreateSaleReturn(ctx, req)
	if err != nil {
		srv.logger.Error("❌ Failed to create sale return order", zap.Error(err))
		return nil, err
	}

	srv.logger.Info("✅ Successfully created sale return order",
		zap.String("OrderNo", req.OrderNo))
	return createdOrder, nil
}

func (srv service) UpdateSaleReturn(ctx context.Context, orderNo string, srNo string, updateBy string) error {
	// 1. Validation
	if orderNo == "" {
		return fmt.Errorf("order number is required")
	}
	if srNo == "" {
		return fmt.Errorf("SR number is required")
	}
	if updateBy == "" {
		return fmt.Errorf("updater information is required")
	}

	// 2. ตรวจสอบสถานะปัจจุบันของ order
	order, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to get order", zap.Error(err))
		return err
	}
	if order == nil {
		return fmt.Errorf("order not found: %s", orderNo)
	}

	// 3. เพิ่มการตรวจสอบสถานะก่อนอัพเดท (ถ้าจำเป็น)
	if order.StatusConfID == 3 { // ถ้าถูกยกเลิกแล้ว
		return fmt.Errorf("cannot update canceled order")
	}

	// เพิ่มการตรวจสอบสถานะเพิ่มเติม
	if order.StatusReturnID != 1 { // ถ้าไม่ใช่สถานะเริ่มต้น
		return fmt.Errorf("cannot update SR number: invalid status")
	}

	// 4. อัพเดท SR number
	err = srv.befRORepo.UpdateSaleReturn(ctx, orderNo, srNo, updateBy)
	if err != nil {
		srv.logger.Error("❌ Failed to update SR number", zap.Error(err))
		return err
	}

	// 5. บันทึก log สำเร็จ
	srv.logger.Info("✅ Successfully updated SR number",
		zap.String("OrderNo", orderNo),
		zap.String("SrNo", srNo),
		zap.String("UpdateBy", updateBy))

	return nil
}

func (srv service) ConfirmSaleReturn(ctx context.Context, orderNo string, confirmBy string) error {
	// 1. เริ่มต้น log การทำงาน
	srv.logger.Info("🏁 Starting sale return confirmation process",
		zap.String("OrderNo", orderNo),
		zap.String("ConfirmBy", confirmBy))

	// 2. เรียกใช้ repository layer เพื่อ update ข้อมูลในฐานข้อมูล
	err := srv.befRORepo.ConfirmSaleReturn(ctx, orderNo, confirmBy)
	if err != nil {
		// 3. log error ถ้าเกิดข้อผิดพลาด
		srv.logger.Error("❌ Failed to confirm sale return", zap.Error(err))
		return err
	}

	// 4. บันทึก log เมื่อทำงานสำเร็จ
	srv.logger.Info("✅ Successfully confirmed sale return",
		zap.String("OrderNo", orderNo),
		zap.String("ConfirmBy", confirmBy))
	return nil
}

func (srv service) CancelSaleReturn(ctx context.Context, orderNo string, updateBy string, remark string) error {
	// 1. บันทึก log เริ่มต้น
	srv.logger.Info("🏁 Starting sale return cancellation process",
		zap.String("OrderNo", orderNo),
		zap.String("updateBy", updateBy))

	// 2. ตรวจสอบสถานะปัจจุบันของ order
	order, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to get order status", zap.Error(err))
		return err
	}
	if order.StatusConfID == 3 {
		srv.logger.Error("❌ Order already canceled",
			zap.String("OrderNo", orderNo))
		return fmt.Errorf("order already canceled")
	}

	// 3. เรียกใช้ repository layer
	err = srv.befRORepo.CancelSaleReturn(ctx, orderNo, updateBy, remark)
	if err != nil {
		srv.logger.Error("❌ Failed to cancel sale return", zap.Error(err))
		return err
	}

	// 4. บันทึก log เมื่อสำเร็จ
	srv.logger.Info("✅ Successfully canceled sale return",
		zap.String("OrderNo", orderNo),
		zap.String("UpdateBy", updateBy))
	return nil
}
