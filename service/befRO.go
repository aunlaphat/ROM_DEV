package service

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

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

	CreateBeforeReturn(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)

	UpdateSaleReturn(ctx context.Context, orderNo string, srNo string, updateBy string) error
	ConfirmSaleReturn(ctx context.Context, orderNo string, confirmBy string) error
	CancelSaleReturn(ctx context.Context, orderNo string, updateBy string, remark string) error

	//  CreateTrade Head + Line
	CreateTradeReturn(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	// CreateTrade Line
	CreateTradeReturnLine(ctx context.Context, orderNo string, lines request.TradeReturnLine) error

	ConfirmReceipt(ctx context.Context, req request.ConfirmTradeReturnRequest, updateBy string) error
	ConfirmReturn(ctx context.Context, req request.ConfirmToReturnRequest, updateBy string) error

	ValidateCreate(req request.BeforeReturnOrder) error
}

// ใช้ตรวจสอบ Create of BeforeReturnOrder
func (srv service) ValidateCreate(req request.BeforeReturnOrder) error {
	// 1. ตรวจสอบข้อมูลพื้นฐาน
	if req.OrderNo == "" {
		return fmt.Errorf("order number is required")
	}
	if req.SoNo == "" {
		return fmt.Errorf("SO number is required")
	}
	if req.CustomerID == "" {
		return fmt.Errorf("customer ID is required")
	}

	// 2. ตรวจสอบค่าที่ต้องมากกว่า 0
	if req.ChannelID <= 0 {
		return fmt.Errorf("invalid channel ID")
	}
	if req.WarehouseID <= 0 {
		return fmt.Errorf("invalid warehouse ID")
	}

	// 3. ตรวจสอบ ReturnType
	/* validReturnTypes := map[string]bool{
		"NORMAL": true,
		"DAMAGE": true,
		// เพิ่ม type อื่นๆ ตามต้องการ
	}
	if !validReturnTypes[req.ReturnType] {
		return fmt.Errorf("invalid return type: %s", req.ReturnType)
	} */

	// 4. ตรวจสอบ order lines
	if len(req.BeforeReturnOrderLines) == 0 {
		return fmt.Errorf("at least one order line is required")
	}

	for i, line := range req.BeforeReturnOrderLines {
		if line.SKU == "" {
			return fmt.Errorf("SKU is required for line %d", i+1)
		}
		if line.ItemName == "" {
			return fmt.Errorf("ItemName is required for line %d", i+1)
		}
		if line.QTY <= 0 {
			return fmt.Errorf("quantity must be greater than 0 for line %d", i+1)
		}
		if line.ReturnQTY < 0 {
			return fmt.Errorf("return quantity cannot be negative for line %d", i+1)
		}
		if line.ReturnQTY > line.QTY {
			return fmt.Errorf("return quantity cannot be greater than quantity for line %d", i+1)
		}
		if line.Price < 0 {
			return fmt.Errorf("price cannot be negative for line %d", i+1)
		}
		// ตรวจสอบ AlterSKU ถ้ามี
		if line.AlterSKU != nil && *line.AlterSKU == "" {
			return fmt.Errorf("alter SKU cannot be empty if provided for line %d", i+1)
		}
	}

	return nil
}

// create trade , statusReturnID = 3 (booking)
func (srv service) CreateTradeReturn(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🏁 Starting order creation process", zap.String("OrderNo", req.OrderNo))
	srv.logger.Debug("Creating order head", zap.String("OrderNo", req.OrderNo), zap.String("SoNo", req.SoNo))

	// 2. Validate request
	if err := srv.ValidateCreate(req); err != nil {
		srv.logger.Error("Invalid request", zap.Error(err))
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 3. ตรวจสอบว่า order มีอยู่แล้วหรือไม่
	existingOrder, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("Failed to check existing order", zap.Error(err))
		return nil, err
	}
	if existingOrder != nil {
		return nil, fmt.Errorf("order already exists: %s", req.OrderNo)
	}

	// 4. สร้าง sale return order
	createdOrder, err := srv.befRORepo.CreateTradeReturn(ctx, req)
	if err != nil {
		srv.logger.Error("❌ Failed to create sale return order", zap.Error(err))
		return nil, err
	}

	srv.logger.Info("✅ Successfully created order with lines",
		zap.String("OrderNo", req.OrderNo))
	return createdOrder, nil
}

// create line for if want to add line on page
func (srv service) CreateTradeReturnLine(ctx context.Context, orderNo string, lines request.TradeReturnLine) error {
	// ตรวจสอบว่ามี OrderNo อยู่ใน BeforeReturnOrder หรือไม่
	exists, err := srv.befRORepo.CheckBefOrderNoExists(ctx, orderNo)
	if err != nil {
		return fmt.Errorf("failed to check order existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("order not found: %s", orderNo)
	}

	// สร้างข้อมูลใน BeforeReturnOrderLine
	err = srv.befRORepo.CreateTradeReturnLine(ctx, orderNo, lines.TradeReturnLine)
	if err != nil {
		return fmt.Errorf("failed to create trade return line: %w", err)
	}

	return nil
}
func (srv service) ConfirmReceipt(ctx context.Context, req request.ConfirmTradeReturnRequest, updateBy string) error {
	srv.logger.Info("🏁 Starting trade return confirmation process",
		zap.String("Identifier", req.Identifier),
		zap.String("UpdateBy", updateBy))

	// ตรวจสอบความถูกต้องของข้อมูล
	if req.Identifier == "" || updateBy == "" {
		return fmt.Errorf("identifier (OrderNo or TrackingNo) and updateBy are required")
	}

	// 1. อัปเดตสถานะใน BeforeReturnOrder
	if err := srv.befRORepo.UpdateBefToWaiting(ctx, req, updateBy); err != nil {
		return fmt.Errorf("failed to update BeforeReturnOrder: %w", err)
	}

	// 2. ดึงข้อมูลจาก BeforeReturnOrder
	returnOrderData, err := srv.befRORepo.GetBeforeReturnOrderData(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to fetch BeforeReturnOrder: %w", err)
	}

	// กำหนดค่าเริ่มต้นให้กับ StatusCheckID ให้เป็นสถานะ waiting ทันทีเมื่อกด
	returnOrderData.StatusCheckID = 1

	// 3. Insert ข้อมูลลงใน ReturnOrder
	if err := srv.befRORepo.InsertReturnOrder(ctx, returnOrderData); err != nil {
		return fmt.Errorf("failed to insert into ReturnOrder: %w", err)
	}

	// 4. Insert ข้อมูลจาก importLines ลงใน ReturnOrderLine
	if err := srv.befRORepo.InsertReturnOrderLine(ctx, returnOrderData, req); err != nil {
		return fmt.Errorf("failed to insert into ReturnOrderLine: %w", err)
	}

	srv.logger.Info("✅ Successfully confirmed trade return",
		zap.String("Identifier", req.Identifier),
		zap.String("UpdateBy", updateBy))

	return nil
}

// check trade line from scan => confirm => success (unsuccess in process future..)
func (srv service) ConfirmReturn(ctx context.Context, req request.ConfirmToReturnRequest, updateBy string) error {
	srv.logger.Info("🏁 Starting return confirmation process",
		zap.String("OrderNo", req.OrderNo),
		zap.String("UpdateBy", updateBy))

	// ตรวจสอบว่ามี OrderNo และ UpdateBy หรือไม่
	if req.OrderNo == "" || updateBy == "" {
		return fmt.Errorf("OrderNo and UpdateBy are required")
	}

	// ตรวจสอบ SKU
	for _, line := range req.ImportLinesActual {
		if line.SKU == "" {
			return fmt.Errorf("SKU is required")
		}
		exists, err := srv.befRORepo.CheckReLineSKUExists(ctx, line.SKU)
		if err != nil {
			return fmt.Errorf("failed to check SKU existence: %w", err)
		}
		if !exists {
			return fmt.Errorf("SKU %s does not exist in the database", line.SKU)
		}
	}

	// ตรวจสอบว่า OrderNo ตรงกับฐานข้อมูลใน BeforeReturn หรือไม่
	exists, err := srv.befRORepo.CheckBefOrderNoExists(ctx, req.OrderNo)
	if err != nil {
		return fmt.Errorf("failed to check order existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("OrderNo does not exist in BeforeReturnOrder")
	}

	// อัปเดต BeforeReturnOrder
	if err := srv.befRORepo.UpdateStatusToSuccess(ctx, req.OrderNo, updateBy); err != nil {
		return fmt.Errorf("failed to update BeforeReturnOrder: %w", err)
	}

	// ดึงข้อมูล BeforeReturnOrder
	beforeReturnOrder, err := srv.befRORepo.GetBeforeOrderDetails(ctx, req.OrderNo)
	if err != nil {
		return fmt.Errorf("failed to fetch BeforeReturnOrder details: %w", err)
	}

	// อัปเดต ReturnOrder และ ReturnOrderLine
	if err := srv.befRORepo.UpdateReturnOrderAndLines(ctx, req, beforeReturnOrder); err != nil {
		return fmt.Errorf("failed to update ReturnOrder and ReturnOrderLine: %w", err)
	}

	srv.logger.Info("✅ Successfully confirmed return",
		zap.String("OrderNo", req.OrderNo),
		zap.String("UpdateBy", updateBy))
	return nil
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

func (srv service) CreateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🏁 Starting order creation process", zap.String("OrderNo", req.OrderNo))
	srv.logger.Debug("Creating order head", zap.String("OrderNo", req.OrderNo), zap.String("SoNo", req.SoNo))

	// 2. Validate request
	if err := srv.ValidateCreate(req); err != nil {
		srv.logger.Error("Invalid request", zap.Error(err))
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 3. ตรวจสอบว่า order มีอยู่แล้วหรือไม่
	existingOrder, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("Failed to check existing order", zap.Error(err))
		return nil, err
	}
	if existingOrder != nil {
		return nil, fmt.Errorf("order already exists: %s", req.OrderNo)
	}

	// 4. สร้าง sale return order
	createdOrder, err := srv.befRORepo.CreateBeforeReturnOrderWithTransaction(ctx, req)
	if err != nil {
		srv.logger.Error("❌ Failed to create sale return order", zap.Error(err))
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
		srv.logger.Error("failed to delete before return order line: ", zap.Error(err))
		return order, errors.UnexpectedError()
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

func (srv service) CreateBeforeReturn(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	// 1. เริ่มต้น log การทำงาน
	srv.logger.Info("🏁 Starting sale return creation process", zap.String("OrderNo", req.OrderNo))
	srv.logger.Debug("Creating sale return order", zap.String("OrderNo", req.OrderNo), zap.String("SoNo", req.SoNo))

	// 2. Validate request
	if err := srv.validateCreateSaleReturn(req); err != nil {
		srv.logger.Error("Invalid request", zap.Error(err))
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 3. ตรวจสอบว่า order มีอยู่แล้วหรือไม่
	existingOrder, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("Failed to check existing order", zap.Error(err))
		return nil, err
	}
	if existingOrder != nil {
		return nil, fmt.Errorf("order already exists: %s", req.OrderNo)
	}

	// 4. สร้าง sale return order
	createdOrder, err := srv.befRORepo.CreateBeforeReturn(ctx, req)
	if err != nil {
		srv.logger.Error("❌ Failed to create sale return order", zap.Error(err))
		return nil, err
	}

	// 5. บันทึก log สำเร็จ
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
	if order.StatusConfID != nil && *order.StatusConfID == 3 { // ถ้าถูกยกเลิกแล้ว
		srv.logger.Error("❌ Cannot update canceled order", zap.String("OrderNo", orderNo))
		return fmt.Errorf("cannot update canceled order")
	}

	// เพิ่มการตรวจสอบสถานะเพิ่มเติม
	if order.StatusReturnID != nil && *order.StatusReturnID != 1 { // ถ้าไม่ใช่สถานะเริ่มต้น
		srv.logger.Error("❌ Cannot update SR number: invalid status", zap.String("OrderNo", orderNo))
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
	// 1. Input validation
	if orderNo == "" || updateBy == "" || remark == "" {
		return fmt.Errorf("orderNo, updateBy and remark are required")
	}

	// 2. ตรวจสอบสถานะปัจจุบันของ order
	order, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("Failed to get order status", zap.Error(err))
		return err
	}
	if order == nil {
		return fmt.Errorf("order not found: %s", orderNo)
	}

	// 3. ตรวจสอบว่าถูกยกเลิกไปแล้วหรือไม่
	if order.StatusConfID != nil && *order.StatusConfID == 3 {
		srv.logger.Error("Order already canceled",
			zap.String("OrderNo", orderNo))
		return fmt.Errorf("order already canceled")
	}

	// 4. เรียกใช้ repository layer เพื่อยกเลิก order และสร้าง cancel status
	err = srv.befRORepo.CancelSaleReturn(ctx, orderNo, updateBy, remark)
	if err != nil {
		srv.logger.Error("Failed to process cancellation", zap.Error(err))
		return err
	}

	// 5. Log สำเร็จ
	srv.logger.Info("Successfully canceled sale return",
		zap.String("OrderNo", orderNo),
		zap.String("UpdateBy", updateBy),
		zap.String("Remark", remark))

	return nil
}

// เพิ่มฟังก์ชัน validate สำหรับ CreateSaleReturn
func (srv service) validateCreateSaleReturn(req request.BeforeReturnOrder) error {
	// 1. ตรวจสอบข้อมูลพื้นฐาน
	if req.OrderNo == "" {
		return fmt.Errorf("order number is required")
	}
	if req.SoNo == "" {
		return fmt.Errorf("SO number is required")
	}
	if req.CustomerID == "" {
		return fmt.Errorf("customer ID is required")
	}

	// 2. ตรวจสอบค่าที่ต้องมากกว่า 0
	if req.ChannelID <= 0 {
		return fmt.Errorf("invalid channel ID")
	}
	if req.WarehouseID <= 0 {
		return fmt.Errorf("invalid warehouse ID")
	}

	// 3. ตรวจสอบ Reason
	/* validReturnTypes := map[string]bool{
		"NORMAL": true,
		"DAMAGE": true,
		// เพิ่ม type อื่นๆ ตามต้องการ
	}
	if !validReturnTypes[req.Reason] {
		return fmt.Errorf("invalid return type: %s", req.Reason)
	} */

	// 4. ตรวจสอบ order lines
	if len(req.BeforeReturnOrderLines) == 0 {
		return fmt.Errorf("at least one order line is required")
	}

	for i, line := range req.BeforeReturnOrderLines {
		if line.SKU == "" {
			return fmt.Errorf("SKU is required for line %d", i+1)
		}
		if line.QTY <= 0 {
			return fmt.Errorf("quantity must be greater than 0 for line %d", i+1)
		}
		if line.ReturnQTY < 0 {
			return fmt.Errorf("return quantity cannot be negative for line %d", i+1)
		}
		if line.ReturnQTY > line.QTY {
			return fmt.Errorf("return quantity cannot be greater than quantity for line %d", i+1)
		}
		if line.Price < 0 {
			return fmt.Errorf("price cannot be negative for line %d", i+1)
		}
		// ตรวจสอบ AlterSKU ถ้ามี
		if line.AlterSKU != nil && *line.AlterSKU == "" {
			return fmt.Errorf("alter SKU cannot be empty if provided for line %d", i+1)
		}
	}

	return nil
}
