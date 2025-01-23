package service

/*
type BefROService interface {
ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error)
	ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error)

	GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error)
	GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error)
	GetAllOrderDetail(ctx context.Context) ([]response.OrderDetail, error)
	GetAllOrderDetails(ctx context.Context, page, limit int) ([]response.OrderDetail, error)

	GetOrderDetailBySO(ctx context.Context, soNo string) (*response.OrderDetail, error)
	SearchSaleOrder(ctx context.Context, soNo string) ([]response.SaleOrderResponse, error)

	CreateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	UpdateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	DeleteBeforeReturnOrderLine(ctx context.Context, recID string) error

	CreateTradeReturnLine(ctx context.Context, orderNo string, line request.TradeReturnLineRequest) error
	ConfirmReturn(ctx context.Context, orderNo string, confirmBy string) error
	CancelReturn(ctx context.Context, orderNo string, cancelBy string, remark string) error
}

func (srv service) CreateTradeReturnLine(ctx context.Context, orderNo string, line request.TradeReturnLineRequest) error {
	// ตรวจสอบว่ามี OrderNo อยู่ใน BeforeReturnOrder หรือไม่
	exists, err := srv.befRORepo.CheckOrderNoExists(ctx, orderNo)
	if err != nil {
		return fmt.Errorf("failed to check order existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("order not found: %s", orderNo)
	}

	// สร้างข้อมูลใน BeforeReturnOrderLine
	err = srv.befRORepo.CreateTradeReturnLine(ctx, orderNo, line)
	if err != nil {
		return fmt.Errorf("failed to create trade return line: %w", err)
	}

	return nil
}

/* func (srv service) ConfirmReturn(ctx context.Context, orderNo string, confirmBy string) error {
	// 1. เริ่มต้น log การทำงาน
	srv.logger.Info("🏁 Starting sale return confirmation process",
		zap.String("OrderNo", orderNo),
		zap.String("ConfirmBy", confirmBy))

	// 2. เรียกใช้ repository layer เพื่อ update ข้อมูลในฐานข้อมูล
	err := srv.befRORepo.ConfirmOrderNo(ctx, orderNo, confirmBy)
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

func (srv service) CancelReturn(ctx context.Context, orderNo string, cancelBy string, remark string) error {
	// 1. บันทึก log เริ่มต้น
	srv.logger.Info("🏁 Starting sale return cancellation process",
		zap.String("OrderNo", orderNo),
		zap.String("CancelBy", cancelBy))

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
	err = srv.befRORepo.CancelOrderNo(ctx, orderNo, cancelBy, remark)
	if err != nil {
		srv.logger.Error("❌ Failed to cancel sale return", zap.Error(err))
		return err
	}

	// 4. บันทึก log เมื่อสำเร็จ
	srv.logger.Info("✅ Successfully canceled sale return",
		zap.String("OrderNo", orderNo),
		zap.String("CancelBy", cancelBy))
	return nil
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

func (srv service) GetAllOrderDetail(ctx context.Context) ([]response.OrderDetail, error) {
	allorder, err := srv.befRORepo.GetAllOrderDetail(ctx)
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

func (srv service) GetAllOrderDetails(ctx context.Context, page, limit int) ([]response.OrderDetail, error) {
	offset := (page - 1) * limit

	allorder, err := srv.befRORepo.GetAllOrderDetails(ctx, offset, limit)
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

func (srv service) GetOrderDetailBySO(ctx context.Context, soNo string) (*response.OrderDetail, error) {
	soOrder, err := srv.befRORepo.GetOrderDetailBySO(ctx, soNo)
	if err != nil {
		return nil, err
	}
	return soOrder, nil
}

// Implementation สำหรับ SearchSaleOrder
func (srv service) SearchSaleOrder(ctx context.Context, soNo string) ([]response.SaleOrderResponse, error) {
	srv.logger.Debug("🚀 Starting SearchSaleOrder", zap.String("SoNo", soNo))
	order, err := srv.befRORepo.SearchSaleOrder(ctx, soNo)
	if err != nil {
		srv.logger.Error("❌ Failed to search sale orders", zap.Error(err))
		return nil, err
	}
	if order == nil {
		srv.logger.Debug("❗ No sale order found", zap.String("SoNo", soNo))
		return nil, nil
	}
	srv.logger.Debug("✅ Successfully searched sale orders", zap.String("SoNo", soNo))
	return []response.SaleOrderResponse{*order}, nil
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

func (srv service) DeleteBeforeReturnOrderLine(ctx context.Context, recID string) error {
	if recID == "" {
		return errors.ValidationError("ReturnID is required")
	}

	err := srv.befRORepo.DeleteBeforeReturnOrderLine(ctx, recID)
	if err != nil {
		srv.logger.Error("failed to delete before return order line: ", zap.Error(err))
		return errors.UnexpectedError()
	}

	return nil
}
*/
