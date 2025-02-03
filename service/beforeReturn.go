package service

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/utils"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// BefROService interface กำหนด method สำหรับการทำงานกับ Before Return Order
type BeforeReturnService interface {
	// Method สำหรับสร้าง Before Return Order พร้อมกับ Lines
	CreateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	// Method สำหรับดึงรายการ Before Return Orders ทั้งหมด
	ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error)
	// Method สำหรับดึง Before Return Order โดยใช้ OrderNo
	GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error)
	// Method สำหรับดึงรายการ Before Return Order Lines ทั้งหมด
	ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error)
	// Method สำหรับดึง Before Return Order Lines โดยใช้ OrderNo
	GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error)
	// Method สำหรับอัพเดท Before Return Order พร้อมกับ Lines
	UpdateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)

	// ************************ Create Sale Return ************************ //
	// Method สำหรับค้นหา Order โดยใช้ SoNo และ OrderNo
	SearchOrder(ctx context.Context, soNo, orderNo string) ([]response.SaleOrderResponse, error)
	// Method สำหรับสร้าง Sale Return
	CreateSaleReturn(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	// Method สำหรับอัพเดท Sale Return
	UpdateSaleReturn(ctx context.Context, req request.UpdateSaleReturn) error
	// Method สำหรับยืนยัน Sale Return
	ConfirmSaleReturn(ctx context.Context, orderNo string, roleID int, userID string) (*response.ConfirmSaleReturnResponse, error)
	// Method สำหรับยกเลิก Sale Return
	CancelSaleReturn(ctx context.Context, orderNo, updateBy, remark string) (*response.CancelSaleReturnResponse, error)

	// Method สำหรับดึงรายการ Draft Orders ทั้งหมด
	ListDraftOrders(ctx context.Context) ([]response.ListDraftConfirmOrdersResponse, error)
	// Method สำหรับดึงรายการ Confirm Orders ทั้งหมด
	ListConfirmOrders(ctx context.Context) ([]response.ListDraftConfirmOrdersResponse, error)
	// Method สำหรับดึง Draft Confirm Order โดยใช้ OrderNo
	GetDraftConfirmOrderByOrderNo(ctx context.Context, orderNo string) (*response.DraftHeadResponse, error)
	// Method สำหรับดึง CodeR ทั้งหมด
	ListCodeR(ctx context.Context) ([]response.CodeRResponse, error)
	// Method สำหรับเพิ่ม CodeR
	AddCodeR(ctx context.Context, req request.CodeR) (*response.DraftLineResponse, error)
	// Method สำหรับลบ CodeR
	DeleteCodeR(ctx context.Context, orderNo string, sku string) error
	// Method สำหรับอัพเดท Draft Order
	UpdateDraftOrder(ctx context.Context, orderNo string, userID string) error

	// Method ดึงข้อมูลรายละเอียดคำสั่งซื้อทั้งหมด
	GetAllOrderDetail(ctx context.Context) ([]response.OrderDetail, error)
	// Method ดึงข้อมูลรายละเอียดคำสั่งซื้อทั้งหมดพร้อมการแบ่งหน้า
	GetAllOrderDetails(ctx context.Context, page, limit int) ([]response.OrderDetail, error)
	// Method ดึงข้อมูลรายละเอียดคำสั่งซื้อโดยใช้หมายเลข SO
	GetOrderDetailBySO(ctx context.Context, soNo string) (*response.OrderDetail, error)
	// Method ลบรายการ BeforeReturnOrderLine โดยใช้ RecID
	DeleteBeforeReturnOrderLine(ctx context.Context, recID string) error
	// Method สร้างคำสั่งซื้อคืนสินค้า
	CreateTradeReturn(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	// Method สร้างรายการคืนสินค้า
	CreateTradeReturnLine(ctx context.Context, orderNo string, lines request.TradeReturnLine) ([]response.BeforeReturnOrderLineResponse, error)
	// Method ยืนยันการรับสินค้าคืนจากหน้าคลัง
	ConfirmReceipt(ctx context.Context, req request.ConfirmTradeReturnRequest, updateBy string) error
	// Method ยืนยันการคืนสินค้าโดยสมบูรณ์
	ConfirmReturn(ctx context.Context, req request.ConfirmToReturnRequest, updateBy string) error
	// Method ตรวจสอบความถูกต้องของข้อมูลก่อนสร้างคำสั่งซื้อคืนสินค้า
	ValidateCreate(req request.BeforeReturnOrder) error
}

// Create Return Order MKP
func (srv service) SearchOrder(ctx context.Context, soNo, orderNo string) ([]response.SaleOrderResponse, error) {
	logFinish := srv.logger.LogAPICall(ctx, "SearchOrder", zap.String("SoNo", soNo), zap.String("OrderNo", orderNo))
	defer logFinish("Completed", nil)

	srv.logger.Info("🔎 Searching Sale Order 🔎",
		zap.String("SoNo", soNo),
		zap.String("OrderNo", orderNo),
	)

	order, err := srv.beforeReturnRepo.SearchOrder(ctx, soNo, orderNo)
	if err != nil {
		logFinish("Failed", err)
		srv.logger.Error("❌ Search Order Failed ❌",
			zap.Error(err),
			zap.String("SoNo", soNo),
			zap.String("OrderNo", orderNo),
		)
		return nil, fmt.Errorf("search order failed: %w", err)
	}

	if order == nil {
		logFinish("Not Found", nil)
		srv.logger.Warn("⚠️ No Sale Order Found ⚠️",
			zap.String("SoNo", soNo),
			zap.String("OrderNo", orderNo),
		)
		return nil, nil
	}

	srv.logger.Info("✅ Sale Order Found ✅",
		zap.String("SoNo", soNo),
		zap.String("OrderNo", orderNo),
		zap.Any("OrderDetails", order),
	)

	logFinish("Success", nil)
	return []response.SaleOrderResponse{*order}, nil
}

// Method สำหรับสร้าง Sale Return
func (srv service) CreateSaleReturn(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "CreateSaleReturn", zap.String("OrderNo", req.OrderNo))
	defer logFinish("Completed", nil) // สร้าง closure สำหรับบันทึกสถานะเมื่อฟังก์ชันจบ

	// Logging ว่าเริ่มการทำงาน
	srv.logger.Info("🔎 Starting sale return creation process",
		zap.String("OrderNo", req.OrderNo),
		zap.String("CustomerID", req.CustomerID),
		zap.Time("CreateDate", time.Now()),
	)

	// Validate request
	if err := utils.ValidateCreateSaleReturn(req); err != nil {
		logFinish("Failed", err)
		srv.logger.Error("❌ Validation failed", zap.Error(err))
		return nil, errors.Wrap(err, "validation failed")
	}

	// ตรวจสอบว่า Order มีอยู่แล้วหรือไม่
	existingOrder, err := srv.beforeReturnRepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logFinish("Failed", err)
		srv.logger.Error("❌ Failed to fetch order", zap.Error(err))
		return nil, errors.Wrap(err, "failed to check existing order")
	}
	if existingOrder != nil {
		err := fmt.Errorf("⚠️ Order already exists: %s", req.OrderNo)
		logFinish("Failed", err)
		srv.logger.Warn("⚠️ Duplicate order found", zap.String("OrderNo", req.OrderNo))
		return nil, err
	}

	// สร้าง Sale Return Order
	createdOrder, err := srv.beforeReturnRepo.CreateSaleReturn(ctx, req)
	if err != nil {
		logFinish("Failed", err)
		srv.logger.Error("❌ Failed to create order", zap.Error(err))
		return nil, errors.Wrap(err, "failed to create sale return order")
	}

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	srv.logger.Info("✅ Sale return order created successfully", zap.String("OrderNo", createdOrder.OrderNo))

	return createdOrder, nil
}

// Method สำหรับอัพเดท Sale Return
func (srv service) UpdateSaleReturn(ctx context.Context, req request.UpdateSaleReturn) error {
	logFinish := srv.logger.LogAPICall(ctx, "UpdateSaleReturn", zap.String("OrderNo", req.OrderNo), zap.String("SrNo", req.SrNo))
	defer logFinish("Completed", nil)

	srv.logger.Info("🔎 Updating Sale Return Order", zap.String("OrderNo", req.OrderNo), zap.String("SrNo", req.SrNo))

	// Validation
	if req.OrderNo == "" || req.SrNo == "" {
		err := errors.New("OrderNo and SrNo are required")
		srv.logger.Warn("⚠️ Validation failed", zap.String("OrderNo", req.OrderNo), zap.Error(err))
		logFinish("Failed", err)
		return err
	}

	// Fetch existing order
	order, err := srv.beforeReturnRepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("❌ Error fetching order", zap.String("OrderNo", req.OrderNo), zap.Error(err))
		logFinish("Failed", err)
		return errors.Wrap(err, "failed to fetch order")
	}
	if order == nil {
		err := fmt.Errorf("⚠️ Order not found: %s", req.OrderNo)
		srv.logger.Warn("⚠️ Order not found", zap.String("OrderNo", req.OrderNo))
		logFinish("Failed", err)
		return err
	}

	// Update Sale Return
	if err := srv.beforeReturnRepo.UpdateSaleReturn(ctx, req); err != nil {
		srv.logger.Error("❌ Failed to update SR number", zap.String("OrderNo", req.OrderNo), zap.Error(err))
		logFinish("Failed", err)
		return errors.Wrap(err, "failed to update SR number")
	}

	srv.logger.Info("✅ Sale Return Updated Successfully", zap.String("OrderNo", req.OrderNo))
	logFinish("Success", nil)
	return nil
}

func (srv service) ConfirmSaleReturn(ctx context.Context, orderNo string, roleID int, userID string) (*response.ConfirmSaleReturnResponse, error) {
	// 🪄 Start Logging
	logFinish := srv.logger.LogAPICall(ctx, "ConfirmSaleReturn", zap.String("OrderNo", orderNo), zap.Int("RoleID", roleID))
	defer func() { logFinish("Completed", nil) }()

	// ✅ 1. Retrieve Order Details
	order, err := srv.beforeReturnRepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo)
	if err != nil {
		err = errors.Wrap(err, "failed to get order")
		srv.logger.Error("❌ Failed to get order", zap.String("OrderNo", orderNo), zap.Error(err))
		logFinish("Failed", err)
		return nil, err
	}
	if order == nil {
		err := fmt.Errorf("⚠️ Order not found: %s", orderNo)
		srv.logger.Warn("⚠️ Order not found", zap.String("OrderNo", orderNo))
		logFinish("Not Found", err)
		return nil, err
	}

	// ✅ 2. Ensure required fields are not nil
	if order.IsCNCreated == nil || order.IsEdited == nil {
		err := fmt.Errorf("❌ Missing required fields in BeforeReturnOrder (IsCNCreated or IsEdited is nil)")
		srv.logger.Error("❌ Missing fields in BeforeReturnOrder", zap.String("OrderNo", orderNo), zap.Error(err))
		logFinish("Failed", err)
		return nil, err
	}

	// ✅ 3. Validate RoleID and Determine Status Updates
	var statusReturnID, statusConfID int

	switch roleID {
	case 2: // ACCOUNTING
		if order.IsCNCreated != nil && !*order.IsCNCreated {
			// If CN is not created
			statusReturnID = 1 // Pending
			statusConfID = 1   // Draft
		} else {
			// CN already created, confirmation is allowed
			statusReturnID = 3 // Booking
			statusConfID = 2   // Confirm
		}
	case 3: // WAREHOUSE
		if order.IsEdited != nil && !*order.IsEdited {
			// No edits, direct confirmation
			statusReturnID = 3 // Booking
			statusConfID = 2   // Confirm
		} else {
			// Edits made, confirmation is not allowed
			statusReturnID = 1 // Pending
			statusConfID = 1   // Draft
		}
	default:
		// ✅ ถ้า Role อื่น ๆ ที่ไม่ใช่ Accounting หรือ Warehouse ให้ตั้งค่าตามที่กำหนด
		srv.logger.Warn("⚠️ Role has limited confirmation permissions - Defaulting to Pending/Draft",
			zap.Int("RoleID", roleID),
			zap.String("OrderNo", orderNo),
		)

		statusReturnID = 1 // Pending
		statusConfID = 1   // Draft
	}

	// ✅ 4. Log Determined Status Before Updating
	srv.logger.Info("📝 Determined Status",
		zap.String("OrderNo", orderNo),
		zap.Int("RoleID", roleID),
		zap.Int("StatusReturnID", statusReturnID),
		zap.Int("StatusConfID", statusConfID),
		zap.String("ConfirmBy", userID),
	)

	// ✅ 5. Call Repository Layer to Update Status
	err = srv.beforeReturnRepo.ConfirmSaleReturn(ctx, orderNo, statusReturnID, statusConfID, userID)
	if err != nil {
		err = errors.Wrap(err, "failed to update return order status")
		srv.logger.Error("❌ Failed to update return order status", zap.String("OrderNo", orderNo), zap.Error(err))
		logFinish("Failed", err)
		return nil, err
	}

	// ✅ 6. Construct Response
	response := &response.ConfirmSaleReturnResponse{
		RefID:          orderNo,
		StatusReturnID: statusReturnID,
		StatusConfID:   statusConfID,
		ConfirmBy:      userID,
		ConfirmDate:    time.Now(),
	}

	// 🪄 Logging Success
	srv.logger.Info("✅ Sale return order confirmed successfully",
		zap.String("OrderNo", orderNo),
		zap.Int("RoleID", roleID),
		zap.String("ConfirmedBy", userID),
		zap.Time("ConfirmedDate", response.ConfirmDate),
	)
	logFinish("Success", nil)

	return response, nil
}

func (srv service) CancelSaleReturn(ctx context.Context, orderNo, updateBy, remark string) (*response.CancelSaleReturnResponse, error) {
	logFinish := srv.logger.LogAPICall(ctx, "CancelSaleReturn", zap.String("OrderNo", orderNo), zap.String("UpdateBy", updateBy))
	defer func() { logFinish("Completed", nil) }()

	// ✅ ตรวจสอบ Input
	if strings.TrimSpace(orderNo) == "" || strings.TrimSpace(updateBy) == "" || strings.TrimSpace(remark) == "" {
		err := errors.New("orderNo, updateBy, and remark are required")
		srv.logger.Error("❌ Invalid input", zap.Error(err))
		logFinish("Failed", err)
		return nil, err
	}

	// ✅ ตรวจสอบว่าสามารถยกเลิกได้หรือไม่
	order, err := srv.beforeReturnRepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo)
	if err != nil {
		err = errors.Wrap(err, "failed to get order")
		srv.logger.Error("❌ Failed to get order", zap.String("OrderNo", orderNo), zap.Error(err))
		logFinish("Failed", err)
		return nil, err
	}
	if order == nil {
		err := fmt.Errorf("order not found: %s", orderNo)
		srv.logger.Warn("⚠️ Order not found", zap.String("OrderNo", orderNo))
		logFinish("Not Found", err)
		return nil, err
	}

	// ✅ เรียกใช้ Repository Layer (แต่ไม่ต้องรับ `CancelID`)
	err = srv.beforeReturnRepo.CancelSaleReturn(ctx, orderNo, updateBy, remark)
	if err != nil {
		err = errors.Wrap(err, "failed to cancel order")
		srv.logger.Error("❌ Failed to cancel order", zap.String("OrderNo", orderNo), zap.Error(err))
		logFinish("Failed", err)
		return nil, err
	}

	// ✅ สร้าง Response (ไม่ต้องมี CancelID)
	response := &response.CancelSaleReturnResponse{
		RefID:        orderNo,
		CancelStatus: true,
		CancelBy:     updateBy,
		Remark:       remark,
		CancelDate:   time.Now(),
	}

	// 🪄 Logging Success
	srv.logger.Info("✅ Order canceled successfully", zap.String("OrderNo", orderNo), zap.String("CanceledBy", updateBy))
	logFinish("Success", nil)

	return response, nil
}

func (srv service) DeleteBeforeReturnOrderLine(ctx context.Context, recID string) error {
	if recID == "" {
		return fmt.Errorf("RecID is required")
	}

	// ส่งไปยัง Repository Layer
	err := srv.beforeReturnRepo.DeleteBeforeReturnOrderLine(ctx, recID)
	if err != nil {
		return fmt.Errorf("failed to delete before return order line: %w", err)
	}

	return nil
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

// create trade , set statusReturnID = 3 (booking)
func (srv service) CreateTradeReturn(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🏁 Starting order creation process", zap.String("OrderNo", req.OrderNo))
	srv.logger.Debug("Creating order head", zap.String("OrderNo", req.OrderNo), zap.String("SoNo", req.SoNo))

	// Validate request
	if err := srv.ValidateCreate(req); err != nil {
		srv.logger.Error("Invalid request", zap.Error(err))
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// ตรวจสอบว่า order มีอยู่แล้วหรือไม่
	existingOrder, err := srv.beforeReturnRepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("Failed to check existing order", zap.Error(err))
		return nil, err
	}
	if existingOrder != nil {
		return nil, fmt.Errorf("order already exists: %s", req.OrderNo)
	}

	// สร้าง trade return order
	createdOrder, err := srv.beforeReturnRepo.CreateTradeReturn(ctx, req)
	if err != nil {
		srv.logger.Error("❌ Failed to create trade return order", zap.Error(err))
		return nil, err
	}

	srv.logger.Info("✅ Successfully created order with lines",
		zap.String("OrderNo", req.OrderNo))
	return createdOrder, nil
}

// add line create trade
func (srv service) CreateTradeReturnLine(ctx context.Context, orderNo string, lines request.TradeReturnLine) ([]response.BeforeReturnOrderLineResponse, error) {

	// ตรวจสอบ OrderNo ที่สร้างว่าซ้ำกับตัวที่มีหรือไม่
	exists, err := srv.beforeReturnRepo.CheckBefOrderNoExists(ctx, orderNo)
	if err != nil {
		return nil, fmt.Errorf("failed to check order existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("order not found: %s", orderNo)
	}

	// สร้างข้อมูลใน BeforeReturnOrderLine
	err = srv.beforeReturnRepo.CreateTradeReturnLine(ctx, orderNo, lines.TradeReturnLine)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to create trade return line: %w", err)
	}

	// สร้าง trade return order
	createdOrderLines, err := srv.beforeReturnRepo.GetBeforeReturnOrderLineByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to create trade return order", zap.Error(err))
		return nil, err
	}

	srv.logger.Info("✅ Successfully created order lines",
		zap.String("OrderNo", orderNo))
	return createdOrderLines, nil
}

func (srv service) ConfirmReceipt(ctx context.Context, req request.ConfirmTradeReturnRequest, updateBy string) error {
	srv.logger.Info("🏁 Starting trade return confirmation process",
		zap.String("Identifier", req.Identifier),
		zap.String("UpdateBy", updateBy))

	// ตรวจสอบค่าว่าง
	if req.Identifier == "" || updateBy == "" {
		return fmt.Errorf("identifier (OrderNo or TrackingNo) and updateBy are required")
	}

	// ตรวจสอบว่า orderNo or trackingNo มีอยู่ในฐานข้อมูล BeforeReturnOrder หรือไม่
	exists, err := srv.beforeReturnRepo.CheckBefOrderOrTrackingExists(ctx, req.Identifier)
	if err != nil {
		return fmt.Errorf("failed to check orderNo or trackingNo existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("orderNo or trackingNo not found: %s", req.Identifier)
	}

	// ตรวจสอบว่ามี sku ที่ Identifier เดียวกันหรือไม่ หากมีสามารถเพิ่มได้ เพราะของหน้าคลังต้องตรงกับข้อมูลที่กรอกเข้าระบบ
	for _, line := range req.ImportLines {
		exists, err := srv.beforeReturnRepo.CheckBefLineSKUExists(ctx, req.Identifier, line.SKU)
		if err != nil {
			return fmt.Errorf("failed to check SKU existence: %w", err)
		}
		if !exists {
			return fmt.Errorf("SKU %s does not exist in BeforeReturnOrderLine for Identifier %s", line.SKU, req.Identifier)
		}
	}

	// 1. อัปเดตสถานะใน BeforeReturnOrder
	if err := srv.beforeReturnRepo.UpdateBefToWaiting(ctx, req, updateBy); err != nil {
		return fmt.Errorf("failed to update BeforeReturnOrder: %w", err)
	}

	// 2. ดึงข้อมูลจาก BeforeReturnOrder
	returnOrderData, err := srv.beforeReturnRepo.GetBeforeReturnOrderData(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to fetch BeforeReturnOrder: %w", err)
	}

	// กำหนดค่าเริ่มต้นให้กับ StatusCheckID ให้เป็นสถานะ waiting ทันทีเมื่อกด
	returnOrderData.StatusCheckID = 1

	// 3. Insert ข้อมูลลงใน ReturnOrder
	if err := srv.beforeReturnRepo.InsertReturnOrder(ctx, returnOrderData); err != nil {
		return fmt.Errorf("failed to insert into ReturnOrder: %w", err)
	}

	// 4. Insert ข้อมูลจาก importLines ลงใน ReturnOrderLine + Check ว่า SKU ตรงกับใน BeforeOD ก่อนถึงเพิ่มได้
	if err := srv.beforeReturnRepo.InsertReturnOrderLine(ctx, returnOrderData, req); err != nil {
		return fmt.Errorf("failed to insert into ReturnOrderLine: %w", err)
	}

	// 5. Insert ข้อมูลภาพลงใน Images (ไฟล์ภาพ)
	if err := srv.beforeReturnRepo.InsertImages(ctx, returnOrderData, req); err != nil {
		return fmt.Errorf("failed to insert images: %w", err)
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

	// ตรวจสอบว่า OrderNo ตรงกับฐานข้อมูลใน BeforeReturn หรือไม่
	exists, err := srv.beforeReturnRepo.CheckBefOrderNoExists(ctx, req.OrderNo)
	if err != nil {
		return fmt.Errorf("failed to check order existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("OrderNo does not exist in BeforeReturnOrder")
	}

	// ตรวจสอบ SKU
	for _, line := range req.ImportLinesActual {
		if line.SKU == "" {
			return fmt.Errorf("SKU is required")
		}
		exists, err := srv.beforeReturnRepo.CheckReLineSKUExists(ctx, req.OrderNo, line.SKU)
		if err != nil {
			return fmt.Errorf("failed to check SKU existence: %w", err)
		}
		if !exists {
			return fmt.Errorf("SKU %s does not exist in ReturnOrderLine for OrderNo %s", line.SKU, req.OrderNo)
		}
	}

	// อัปเดต BeforeReturnOrder
	if err := srv.beforeReturnRepo.UpdateStatusToSuccess(ctx, req.OrderNo, updateBy); err != nil {
		return fmt.Errorf("failed to update BeforeReturnOrder: %w", err)
	}

	// ดึงข้อมูล BeforeReturnOrder
	beforeReturnOrder, err := srv.beforeReturnRepo.GetBeforeOrderDetails(ctx, req.OrderNo)
	if err != nil {
		return fmt.Errorf("failed to fetch BeforeReturnOrder details: %w", err)
	}

	// อัปเดต ReturnOrder และ ReturnOrderLine
	if err := srv.beforeReturnRepo.UpdateReturnOrderAndLines(ctx, req, beforeReturnOrder); err != nil {
		return fmt.Errorf("failed to update ReturnOrder and ReturnOrderLine: %w", err)
	}

	srv.logger.Info("✅ Successfully confirmed return",
		zap.String("OrderNo", req.OrderNo),
		zap.String("UpdateBy", updateBy))
	return nil
}

func (srv service) GetAllOrderDetail(ctx context.Context) ([]response.OrderDetail, error) {
	allorder, err := srv.beforeReturnRepo.GetAllOrderDetail(ctx)
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

	allorder, err := srv.beforeReturnRepo.GetAllOrderDetails(ctx, offset, limit)
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
	soOrder, err := srv.beforeReturnRepo.GetOrderDetailBySO(ctx, soNo)
	if err != nil {
		return nil, err
	}
	return soOrder, nil
}

// Method สำหรับสร้าง Before Return Order พร้อมกับ Lines
func (srv service) CreateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🔎 Starting order creation process", zap.String("OrderNo", req.OrderNo))                  // Logging ว่าเริ่มการสร้าง order
	srv.logger.Debug("Creating order head", zap.String("OrderNo", req.OrderNo), zap.String("SoNo", req.SoNo)) // Logging ว่ากำลังสร้าง order head

	err := srv.beforeReturnRepo.CreateBeforeReturnOrderWithTransaction(ctx, req) // เรียก repository เพื่อสร้าง order พร้อมกับ transaction
	if err != nil {
		srv.logger.Error("❌ Failed to create order with lines", zap.Error(err)) // Logging ว่าการสร้าง order ล้มเหลว
		return nil, err
	}

	createdOrder, err := srv.beforeReturnRepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo) // ดึงข้อมูล order ที่สร้างเสร็จแล้ว
	if err != nil {
		srv.logger.Error("❌ Failed to fetch created order", zap.Error(err)) // Logging ว่าการดึงข้อมูล order ล้มเหลว
		return nil, err
	}

	srv.logger.Info("✅ Successfully created order with lines", zap.String("OrderNo", req.OrderNo)) // Logging ว่าการสร้าง order สำเร็จ
	return createdOrder, nil
}

// Method สำหรับอัพเดท Before Return Order พร้อมกับ Lines
func (srv service) UpdateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🔎 Starting order update process", zap.String("OrderNo", req.OrderNo))                    // Logging ว่าเริ่มการอัพเดท order
	srv.logger.Debug("Updating order head", zap.String("OrderNo", req.OrderNo), zap.String("SoNo", req.SoNo)) // Logging ว่ากำลังอัพเดท order head

	err := srv.beforeReturnRepo.UpdateBeforeReturnOrderWithTransaction(ctx, req) // เรียก repository เพื่ออัพเดท order พร้อมกับ transaction
	if err != nil {
		srv.logger.Error("❌ Failed to update order with lines", zap.Error(err)) // Logging ว่าการอัพเดท order ล้มเหลว
		return nil, err
	}

	updatedOrder, err := srv.beforeReturnRepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo) // ดึงข้อมูล order ที่อัพเดทเสร็จแล้ว
	if err != nil {
		srv.logger.Error("❌ Failed to fetch updated order", zap.Error(err)) // Logging ว่าการดึงข้อมูล order ล้มเหลว
		return nil, err
	}

	srv.logger.Info("✅ Successfully updated order with lines", zap.String("OrderNo", req.OrderNo)) // Logging ว่าการอัพเดท order สำเร็จ
	return updatedOrder, nil
}

// Method สำหรับดึงรายการ Before Return Orders ทั้งหมด
func (srv service) ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🔎 Starting to list all return orders")         // Logging ว่าเริ่มการดึงรายการ return orders ทั้งหมด
	orders, err := srv.beforeReturnRepo.ListBeforeReturnOrders(ctx) // เรียก repository เพื่อดึงรายการ return orders ทั้งหมด
	if err != nil {
		srv.logger.Error("❌ Failed to list return orders", zap.Error(err)) // Logging ว่าการดึงรายการ return orders ล้มเหลว
		return nil, err
	}
	srv.logger.Info("✅ Successfully listed return orders") // Logging ว่าการดึงรายการ return orders สำเร็จ
	return orders, nil
}

// Method สำหรับดึง Before Return Order โดยใช้ OrderNo
func (srv service) GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("🔎 Starting to get return order by order number", zap.String("OrderNo", orderNo)) // Logging ว่าเริ่มการดึง return order โดยใช้ order number
	order, err := srv.beforeReturnRepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo)                    // เรียก repository เพื่อดึง return order โดยใช้ order number
	if err != nil {
		srv.logger.Error("❌ Failed to get return order by order number", zap.Error(err)) // Logging ว่าการดึง return order ล้มเหลว
		return nil, err
	}
	return order, nil
}

// Method สำหรับดึงรายการ Before Return Order Lines ทั้งหมด
func (srv service) ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error) {
	srv.logger.Info("🔎 Starting to list all return order lines")       // Logging ว่าเริ่มการดึงรายการ return order lines ทั้งหมด
	lines, err := srv.beforeReturnRepo.ListBeforeReturnOrderLines(ctx) // เรียก repository เพื่อดึงรายการ return order lines ทั้งหมด
	if err != nil {
		srv.logger.Error("❌ Failed to list return order lines", zap.Error(err)) // Logging ว่าการดึงรายการ return order lines ล้มเหลว
		return nil, err
	}
	srv.logger.Info("✅ Successfully listed return order lines") // Logging ว่าการดึงรายการ return order lines สำเร็จ
	return lines, nil
}

// Method สำหรับดึง Before Return Order Lines โดยใช้ OrderNo
func (srv service) GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error) {
	srv.logger.Info("🔎 Starting to get return order lines by order number", zap.String("OrderNo", orderNo)) // Logging ว่าเริ่มการดึง return order lines โดยใช้ order number
	lines, err := srv.beforeReturnRepo.GetBeforeReturnOrderLineByOrderNo(ctx, orderNo)                      // เรียก repository เพื่อดึง return order lines โดยใช้ order number
	if err != nil {
		srv.logger.Error("❌ Failed to get return order lines by order number", zap.Error(err)) // Logging ว่าการดึง return order lines ล้มเหลว
		return nil, err
	}
	srv.logger.Info("✅ Successfully fetched return order lines") // Logging ว่าการดึง return order lines สำเร็จ
	return lines, nil
}

// Method สำหรับดึงรายการ Draft Orders ทั้งหมด
func (srv service) ListDraftOrders(ctx context.Context) ([]response.ListDraftConfirmOrdersResponse, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "ListDraftOrders")
	defer logFinish("Completed", nil)

	// Logging ว่าเริ่มการทำงาน
	srv.logger.Info("🔎 Starting to list all draft orders 🔎")

	// เรียก Repository เพื่อค้นหา Order ทั้งหมดที่ Status เป็น Draft
	orders, err := srv.beforeReturnRepo.ListDraftOrders(ctx)
	if err != nil {
		// หากเกิดข้อผิดพลาด อัปเดต Log ที่ Error
		logFinish("Failed", fmt.Errorf("❌ Failed to list draft orders : %v", err))
		srv.logger.Error("❌ Failed to list draft orders", zap.Error(err))
		return nil, err
	}

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	return orders, nil
}

// Method สำหรับดึงรายการ Confirm Orders ทั้งหมด
func (srv service) ListConfirmOrders(ctx context.Context) ([]response.ListDraftConfirmOrdersResponse, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "ListConfirmOrders")
	defer logFinish("Completed", nil)

	srv.logger.Info("🔎 Starting to list all confirm orders 🔎")

	// เรียก Repository เพื่อค้นหา Order ทั้งหมดที่ Status เป็น Confirm
	orders, err := srv.beforeReturnRepo.ListConfirmOrders(ctx)
	if err != nil {
		// หากเกิดข้อผิดพลาด อัปเดต Log ที่ Error
		logFinish("Failed", fmt.Errorf("❌ Failed to list confirm orders : %v", err))
		srv.logger.Error("❌ Failed to list confirm orders", zap.Error(err))
		return nil, err
	}

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	return orders, nil
}

// Method สำหรับดึง Draft Confirm Order โดยใช้ OrderNo
func (srv service) GetDraftConfirmOrderByOrderNo(ctx context.Context, orderNo string) (*response.DraftHeadResponse, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "GetDraftConfirmOrderByOrderNo", zap.String("OrderNo", orderNo))
	defer logFinish("Completed", nil)

	// Logging ว่าเริ่มการทำงาน
	srv.logger.Info("🔎 Starting to get draft order by order number 🔎", zap.String("OrderNo", orderNo))

	head, lines, err := srv.beforeReturnRepo.GetDraftConfirmOrderByOrderNo(ctx, orderNo)
	if err != nil {
		// อัปเดต Log ว่าไม่สามารถดึงข้อมูลได้
		logFinish("Failed", fmt.Errorf("❌ Failed to get draft order : %v", err))
		srv.logger.Error("❌ Failed to get draft order", zap.Error(err))
		return nil, err
	}

	head.OrderLines = lines

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	return head, nil
}

// Method สำหรับดึง CodeR ทั้งหมด
func (srv service) ListCodeR(ctx context.Context) ([]response.CodeRResponse, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "GetCodeR")
	defer logFinish("Completed", nil)

	// Logging ว่าเริ่มการทำงาน
	srv.logger.Info("🔎 Starting to get CodeR 🔎")

	// เรียก Repository เพื่อค้นหา CodeR ทั้งหมด
	codeR, err := srv.beforeReturnRepo.ListCodeR(ctx)
	if err != nil {
		// หากเกิดข้อผิดพลาด อัปเดต Log ที่ Error
		logFinish("Failed", fmt.Errorf("❌ Failed to get CodeR : %v", err))
		srv.logger.Error("❌ Failed to get CodeR", zap.Error(err))
		return nil, err
	}

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	return codeR, nil
}

// Method สำหรับเพิ่ม CodeR
func (srv service) AddCodeR(ctx context.Context, req request.CodeR) (*response.DraftLineResponse, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "AddCodeR")
	defer logFinish("Completed", nil)

	// Logging ว่าเริ่มการทำงาน
	srv.logger.Info("🔎 Starting to add CodeR 🔎")

	// ตรวจสอบว่า SKU มีอยู่แล้วหรือไม่
	existingLines, err := srv.beforeReturnRepo.GetBeforeReturnOrderLineByOrderNo(ctx, req.OrderNo)
	if err != nil {
		logFinish("Failed", fmt.Errorf("failed to check existing SKUs: %v", err))
		srv.logger.Error("❌ Failed to check existing SKUs", zap.Error(err))
		return nil, err
	}

	for _, line := range existingLines {
		if line.SKU == req.SKU {
			err := fmt.Errorf("SKU already exists for OrderNo: %s", req.OrderNo)
			logFinish("Failed", err)
			srv.logger.Warn("⚠️ Duplicate SKU found", zap.String("OrderNo", req.OrderNo), zap.String("SKU", req.SKU))
			return nil, err
		}
	}

	// เรียกใช้ repository layer
	result, err := srv.beforeReturnRepo.AddCodeR(ctx, req)
	if err != nil {
		logFinish("Failed", fmt.Errorf("failed to add CodeR: %v", err))
		srv.logger.Error("❌ Failed to add CodeR", zap.Error(err))
		return nil, err
	}

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	return result, nil
}

// Method สำหรับลบ CodeR
func (srv service) DeleteCodeR(ctx context.Context, orderNo string, sku string) error {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "DeleteCodeR", zap.String("OrderNo", orderNo), zap.String("SKU", sku))
	defer logFinish("Completed", nil)

	// เรียกใช้ repository layer
	if err := srv.beforeReturnRepo.DeleteCodeR(ctx, orderNo, sku); err != nil {
		// อัปเดต Log ว่าไม่สามารถลบ CodeR ได้ เนื่องจากเกิดข้อผิดพลาด
		logFinish("Failed", fmt.Errorf("failed to delete CodeR: %v", err))
		srv.logger.Error("❌ Failed to delete CodeR", zap.Error(err))
		return err
	}

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	return nil
}

// Method สำหรับอัพเดท Draft Order
func (srv service) UpdateDraftOrder(ctx context.Context, orderNo string, userID string) error {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "UpdateDraftOrder", zap.String("OrderNo", orderNo), zap.String("UserID", userID))
	defer logFinish("Completed", nil)

	// Logging ว่าเริ่มการทำงาน
	srv.logger.Info("🔎 Starting draft order update process 🔎", zap.String("OrderNo", orderNo))

	// Update order status
	err := srv.beforeReturnRepo.UpdateOrderStatus(ctx, orderNo, 2, 3, userID) // StatusConfID = 2 (Confirm), StatusReturnID = 3 (Booking)
	if err != nil {
		logFinish("Failed", err)
		srv.logger.Error("❌ Failed to update order status", zap.Error(err))
		return err
	}

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	return nil
}
