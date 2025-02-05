package service

// BefROService interface กำหนด method สำหรับการทำงานกับ Before Return Order
type BeforeReturnService interface {
	/* // Method สำหรับสร้าง Before Return Order พร้อมกับ Lines
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

	// Create Return Order MKP 🚨//
	//SearchOrder(ctx context.Context, soNo, orderNo string) (*response.SearchOrderResponse, error)
	CreateSaleReturn(ctx context.Context, req request.CreateSaleReturnOrder, userID string) (*response.BeforeReturnOrderResponse, error)
	UpdateSaleReturn(ctx context.Context, req request.UpdateSaleReturn, userID string) (*response.UpdateSaleReturnResponse, error)
	ConfirmSaleReturn(ctx context.Context, orderNo string, roleID int, userID string) (*response.ConfirmSaleReturnResponse, error)
	CancelSaleReturn(ctx context.Context, req request.CancelSaleReturn, userID string) (*response.CancelSaleReturnResponse, error)

	// Draft & Confirm MKP 🚨//
	ListDraftOrders(ctx context.Context, startDate, endDate string) ([]response.ListDraftConfirmOrdersResponse, error)
	ListConfirmOrders(ctx context.Context, startDate, endDate string) ([]response.ListDraftConfirmOrdersResponse, error)
	GetDraftConfirmOrderByOrderNo(ctx context.Context, orderNo string) (*response.DraftHeadResponse, error)
	ListCodeR(ctx context.Context) ([]response.ListCodeRResponse, error)
	AddCodeR(ctx context.Context, req request.AddCodeR, userID string) ([]response.AddCodeRResponse, error)
	DeleteCodeR(ctx context.Context, orderNo string, sku string, userID string) error
	UpdateDraftOrder(ctx context.Context, orderNo string, userID string) (*response.UpdateOrderStatusResponse, error)

	// Method ดึงข้อมูลรายละเอียดคำสั่งซื้อทั้งหมด
	//GetAllOrderDetail(ctx context.Context) ([]response.OrderDetail, error)
	// Method ดึงข้อมูลรายละเอียดคำสั่งซื้อทั้งหมดพร้อมการแบ่งหน้า
	//GetAllOrderDetails(ctx context.Context, page, limit int) ([]response.OrderDetail, error)
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
	ValidateCreate(req request.BeforeReturnOrder) error */
}

// Create Return Order MKP
/* func (srv service) SearchOrder(ctx context.Context, soNo, orderNo string) (*response.SearchOrderResponse, error) {
	// 📝 เริ่มต้นการ Log การเรียก API โดยใช้ logFinish
	logFinish := srv.logger.LogAPICall(ctx, "SearchOrder",
		zap.String("SoNo", soNo),
		zap.String("OrderNo", orderNo),
	)
	defer func() {
		// 🚀 ใช้ defer เพื่อจับ panic และ log ข้อมูล
		if r := recover(); r != nil {
			srv.logger.Error("🔥 Panic occurred in SearchOrder", zap.Any("panic", r))
			logFinish("Panic", fmt.Errorf("unexpected panic: %v", r))
		}
	}()

	// 📌 Log การเริ่มต้นค้นหาคำสั่งขาย
	srv.logger.Info("🔍 Searching for Sale Order",
		zap.String("SoNo", soNo),
		zap.String("OrderNo", orderNo),
	)

	// ✅ ตรวจสอบเงื่อนไข: ต้องมีค่า SoNo หรือ OrderNo อย่างน้อยหนึ่งค่า
	if soNo == "" && orderNo == "" {
		err := errors.New("either SoNo or OrderNo must be provided")
		srv.logger.Warn("⚠️ Invalid request - Missing parameters", zap.Error(err))
		logFinish("Invalid Request", err)
		return nil, err
	}

	// 🔍 ค้นหาข้อมูลคำสั่งขายจาก Repository Layer
	order, err := srv.beforeReturnRepo.SearchOrder(ctx, soNo, orderNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ✅ ไม่พบข้อมูลคำสั่งขาย
			errMsg := "Sale order not found"
			srv.logger.Warn("⚠️ No Sale Order found", zap.String("SoNo", soNo), zap.String("OrderNo", orderNo))
			logFinish("Not Found", errors.New(errMsg))
			return nil, errors.New(errMsg)
		}

		// ❌ กรณีเกิดปัญหาอื่น ๆ เช่น Database ล่ม
		errMsg := "Failed to retrieve sale order"
		srv.logger.Error("❌ Failed to search Sale Order",
			zap.String("SoNo", soNo),
			zap.String("OrderNo", orderNo),
			zap.Error(err),
		)
		logFinish("Failed", err)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	// ✅ ถ้าพบข้อมูลคำสั่งขาย
	srv.logger.Info("✅ Sale Order found",
		zap.String("SoNo", order.SoNo),
		zap.String("OrderNo", order.OrderNo),
		zap.Int("TotalItems", len(order.Items)), // ✅ แสดงจำนวนรายการสินค้าในออเดอร์
	)

	logFinish("Success", nil)
	return order, nil
} */
/*
func (srv service) CreateSaleReturn(ctx context.Context, req request.CreateSaleReturnOrder, userID string) (*response.BeforeReturnOrderResponse, error) {
	logFinish := srv.logger.LogAPICall(ctx, "CreateSaleReturn",
		zap.String("OrderNo", req.OrderNo),
		zap.String("SoNo", req.SoNo),
		zap.String("UserID", userID),
		zap.Int("TotalItems", len(req.OrderLines)),
	)
	defer func() {
		logFinish("Completed", nil)
	}()

	srv.logger.Info("📝 Processing Sale Return Order",
		zap.String("OrderNo", req.OrderNo),
		zap.String("SoNo", req.SoNo),
		zap.Int("TotalItems", len(req.OrderLines)),
	)

	// ✅ Validate Input
	if err := utils.ValidateCreateSaleReturn(req); err != nil {
		errMsg := "validation failed"
		srv.logger.Warn("⚠️ Validation failed", zap.String("Error", errMsg), zap.Error(err))
		logFinish("Validation Failed", err)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	// ✅ Assign CreateBy & OrderNo to all OrderLines
	req.CreateBy = userID
	for i := range req.OrderLines {
		req.OrderLines[i].CreateBy = userID
		req.OrderLines[i].OrderNo = req.OrderNo
	}

	// 🔄 Call Repository Layer
	createdOrder, err := srv.beforeReturnRepo.CreateSaleReturn(ctx, req)
	if err != nil {
		errMsg := "failed to create Sale Return Order"
		srv.logger.Error("❌ Failed to create Sale Return Order",
			zap.String("OrderNo", req.OrderNo),
			zap.String("Error", errMsg),
			zap.Error(err),
		)
		logFinish("Failed", err)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	// ✅ Successfully created Sale Return Order
	srv.logger.Info("✅ Sale Return Order created",
		zap.String("OrderNo", createdOrder.OrderNo),
		zap.String("SoNo", createdOrder.SoNo),
		zap.Int("TotalItems", len(createdOrder.BeforeReturnOrderLines)),
	)

	logFinish("Success", nil)
	return createdOrder, nil
}

func (srv service) UpdateSaleReturn(ctx context.Context, req request.UpdateSaleReturn, userID string) (*response.UpdateSaleReturnResponse, error) {
	logFinish := srv.logger.LogAPICall(ctx, "UpdateSaleReturn", zap.String("OrderNo", req.OrderNo), zap.String("UserID", userID))
	defer logFinish("Completed", nil)

	srv.logger.Info("🔄 Updating Sale Return Order",
		zap.String("OrderNo", req.OrderNo),
		zap.String("SrNo", req.SrNo),
	)

	// ✅ Validate Input
	if req.OrderNo == "" || req.SrNo == "" {
		srv.logger.Warn("⚠️ Invalid request: OrderNo or SrNo is missing")
		logFinish("Invalid Request", nil)
		return nil, errors.New("orderNo and srNo are required")
	}

	// ✅ Call Repository Layer with userID
	updatedOrder, err := srv.beforeReturnRepo.UpdateSaleReturn(ctx, req, userID)
	if err != nil {
		srv.logger.Error("❌ Failed to update Sale Return Order",
			zap.String("OrderNo", req.OrderNo),
			zap.Error(err),
		)
		logFinish("Failed", err)
		return nil, fmt.Errorf("failed to update Sale Return Order: %w", err)
	}

	srv.logger.Info("✅ Sale Return Order updated successfully",
		zap.String("OrderNo", updatedOrder.OrderNo),
		zap.String("SrNo", updatedOrder.SrNo),
		zap.Int("StatusReturnID", updatedOrder.StatusReturnID),
		zap.Int("StatusConfID", updatedOrder.StatusConfID),
		zap.String("UpdateBy", updatedOrder.UpdateBy),
		zap.Time("UpdateDate", updatedOrder.UpdateDate),
	)

	logFinish("Success", nil)
	return updatedOrder, nil
}

func (srv service) ConfirmSaleReturn(ctx context.Context, orderNo string, roleID int, userID string) (*response.ConfirmSaleReturnResponse, error) {
	// 🪄 Start Logging
	logFinish := srv.logger.LogAPICall(ctx, "ConfirmSaleReturn", zap.String("OrderNo", orderNo), zap.Int("RoleID", roleID))
	defer func() { logFinish("Completed", nil) }()

	// ✅ 1. Retrieve Order Details
	order, err := srv.beforeReturnRepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to get order", zap.String("OrderNo", orderNo), zap.Error(err))
		logFinish("Failed", err)
		return nil, fmt.Errorf("failed to get order: %w", err)
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
			statusReturnID = 1 // Pending
			statusConfID = 1   // Draft
		} else {
			statusReturnID = 3 // Booking
			statusConfID = 2   // Confirm
		}
	case 3: // WAREHOUSE
		if order.IsEdited != nil && !*order.IsEdited {
			statusReturnID = 3 // Booking
			statusConfID = 2   // Confirm
		} else {
			statusReturnID = 1 // Pending
			statusConfID = 1   // Draft
		}
	default:
		srv.logger.Warn("⚠️ Role has limited confirmation permissions - Defaulting to Pending/Draft",
			zap.Int("RoleID", roleID),
			zap.String("OrderNo", orderNo),
		)

		statusReturnID = 1 // Pending
		statusConfID = 1   // Draft
	}

	// ✅ 4. Call Repository Layer to Update Status
	confirmedOrder, err := srv.beforeReturnRepo.ConfirmSaleReturn(ctx, orderNo, statusReturnID, statusConfID, userID)
	if err != nil {
		srv.logger.Error("❌ Failed to confirm Sale Return Order",
			zap.String("OrderNo", orderNo),
			zap.Error(err),
		)
		logFinish("Failed", err)
		return nil, fmt.Errorf("failed to confirm Sale Return Order: %w", err)
	}

	// ✅ 5. Construct Response
	confirmedOrder.ConfirmBy = userID

	srv.logger.Info("✅ Sale return order confirmed successfully",
		zap.String("OrderNo", confirmedOrder.RefID),
		zap.Int("StatusReturnID", confirmedOrder.StatusReturnID),
		zap.Int("StatusConfID", confirmedOrder.StatusConfID),
		zap.String("ConfirmBy", confirmedOrder.ConfirmBy),
		zap.Time("ConfirmDate", confirmedOrder.ConfirmDate),
	)
	logFinish("Success", nil)

	return confirmedOrder, nil
}

func (srv service) CancelSaleReturn(ctx context.Context, req request.CancelSaleReturn, userID string) (*response.CancelSaleReturnResponse, error) {
	logFinish := srv.logger.LogAPICall(ctx, "CancelSaleReturn", zap.String("OrderNo", req.OrderNo), zap.String("UpdateBy", userID))
	defer func() { logFinish("Completed", nil) }()

	// ✅ ตรวจสอบ Input
	if strings.TrimSpace(req.OrderNo) == "" || strings.TrimSpace(req.Remark) == "" || strings.TrimSpace(userID) == "" {
		err := errors.New("orderNo, updateBy, and remark are required")
		srv.logger.Error("❌ Invalid input", zap.Error(err))
		logFinish("Failed", err)
		return nil, err
	}

	// ✅ ตรวจสอบว่าสามารถยกเลิกได้หรือไม่
	order, err := srv.beforeReturnRepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		err = errors.Wrap(err, "failed to get order")
		srv.logger.Error("❌ Failed to get order", zap.String("OrderNo", req.OrderNo), zap.Error(err))
		logFinish("Failed", err)
		return nil, err
	}
	if order == nil {
		err := fmt.Errorf("order not found: %s", req.OrderNo)
		srv.logger.Warn("⚠️ Order not found", zap.String("OrderNo", req.OrderNo))
		logFinish("Not Found", err)
		return nil, err
	}

	// ✅ เรียกใช้ Repository Layer (แต่ไม่ต้องรับ `CancelID`)
	err = srv.beforeReturnRepo.CancelSaleReturn(ctx, req, userID)
	if err != nil {
		err = errors.Wrap(err, "failed to cancel order")
		srv.logger.Error("❌ Failed to cancel order", zap.String("OrderNo", req.OrderNo), zap.Error(err))
		logFinish("Failed", err)
		return nil, err
	}

	// ✅ สร้าง Response (ไม่ต้องมี CancelID)
	response := &response.CancelSaleReturnResponse{
		RefID:        req.OrderNo,
		CancelStatus: true,
		CancelBy:     userID,
		Remark:       req.Remark,
		CancelDate:   time.Now(),
	}

	// 🪄 Logging Success
	srv.logger.Info("✅ Order canceled successfully", zap.String("OrderNo", req.OrderNo), zap.String("CanceledBy", userID))
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
	}

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

/* func (srv service) GetAllOrderDetail(ctx context.Context) ([]response.OrderDetail, error) {
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
} *

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

// Draft & Confirm MKP 🚨//
// ListDraftOrders ดึงรายการ Draft Status Orders 🚗
func (srv service) ListDraftOrders(ctx context.Context, startDate, endDate string) ([]response.ListDraftConfirmOrdersResponse, error) {
	logFinish := srv.logger.LogAPICall(ctx, "ListDraftOrders")
	defer logFinish("Completed", nil)

	srv.logger.Info("🔎 Fetching all draft orders...",
		zap.String("method", "ListDraftOrders"),
		zap.String("startDate", startDate),
		zap.String("endDate", endDate),
	)

	// 📌 ตรวจสอบว่า startDate < endDate หรือไม่
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		srv.logger.Warn("⚠️ Invalid startDate format ⚠️", zap.String("startDate", startDate))
		logFinish("Failed", err)
		return nil, fmt.Errorf("invalid startDate format (expected YYYY-MM-DD): %w", err)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		srv.logger.Warn("⚠️ Invalid endDate format ⚠️", zap.String("endDate", endDate))
		logFinish("Failed", err)
		return nil, fmt.Errorf("invalid endDate format (expected YYYY-MM-DD): %w", err)
	}

	if start.After(end) {
		srv.logger.Warn("⚠️ startDate cannot be after endDate ⚠️",
			zap.String("startDate", startDate),
			zap.String("endDate", endDate),
		)
		logFinish("Failed", fmt.Errorf("startDate cannot be after endDate"))
		return nil, fmt.Errorf("startDate cannot be after endDate")
	}

	// 📌 ดึงข้อมูลจาก Repository Layer
	orders, err := srv.beforeReturnRepo.ListDraftOrders(ctx, startDate, endDate)
	if err != nil {
		srv.logger.Error("❌ Failed to list draft orders",
			zap.Error(err),
			zap.String("startDate", startDate),
			zap.String("endDate", endDate),
		)
		logFinish("Failed", err)
		return nil, fmt.Errorf("ListDraftOrders failed: %w", err)
	}

	srv.logger.Info("✅ Successfully retrieved draft orders",
		zap.Int("count", len(orders)),
		zap.String("startDate", startDate),
		zap.String("endDate", endDate),
	)
	logFinish(fmt.Sprintf("Success - %d orders", len(orders)), nil)

	return orders, nil
}

// ListConfirmOrders ดึงรายการ Confirm Satus Orders 🚗
func (srv service) ListConfirmOrders(ctx context.Context, startDate, endDate string) ([]response.ListDraftConfirmOrdersResponse, error) {
	logFinish := srv.logger.LogAPICall(ctx, "ListConfirmOrders")
	defer logFinish("Completed", nil)

	srv.logger.Info("🔎 Fetching all confirm orders...",
		zap.String("startDate", startDate),
		zap.String("endDate", endDate),
	)

	// 📌 ตรวจสอบว่า startDate < endDate หรือไม่
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		srv.logger.Warn("⚠️ Invalid startDate format ⚠️", zap.String("startDate", startDate))
		logFinish("Failed", err)
		return nil, fmt.Errorf("invalid startDate format (expected YYYY-MM-DD): %w", err)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		srv.logger.Warn("⚠️ Invalid endDate format ⚠️", zap.String("endDate", endDate))
		logFinish("Failed", err)
		return nil, fmt.Errorf("invalid endDate format (expected YYYY-MM-DD): %w", err)
	}

	if start.After(end) {
		srv.logger.Warn("⚠️ startDate cannot be after endDate ⚠️",
			zap.String("startDate", startDate),
			zap.String("endDate", endDate),
		)
		logFinish("Failed", fmt.Errorf("startDate cannot be after endDate"))
		return nil, fmt.Errorf("startDate cannot be after endDate")
	}

	// 📌 ดึงข้อมูลจาก Repository Layer
	orders, err := srv.beforeReturnRepo.ListConfirmOrders(ctx, startDate, endDate)
	if err != nil {
		srv.logger.Error("❌ Failed to list confirm orders",
			zap.Error(err),
			zap.String("startDate", startDate),
			zap.String("endDate", endDate),
		)
		logFinish("Failed", err)
		return nil, fmt.Errorf("ListConfirmOrders failed: %w", err)
	}

	srv.logger.Info("✅ Successfully retrieved confirm orders",
		zap.Int("count", len(orders)),
		zap.String("startDate", startDate),
		zap.String("endDate", endDate),
	)
	logFinish(fmt.Sprintf("Success - %d orders", len(orders)), nil)

	return orders, nil
}

// GetDraftConfirmOrderByOrderNo ดึงข้อมูล Order และทำ Logging 🚗
func (srv service) GetDraftConfirmOrderByOrderNo(ctx context.Context, orderNo string) (*response.DraftHeadResponse, error) {
	logFinish := srv.logger.LogAPICall(ctx, "GetDraftConfirmOrderByOrderNo")
	defer logFinish("Completed", nil)

	srv.logger.Info("🔎 Fetching Draft Confirm Order...", zap.String("orderNo", orderNo))

	// 📌 เรียกใช้ Repository Layer
	order, err := srv.beforeReturnRepo.GetDraftConfirmOrderByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("❌ Failed to get Draft Confirm Order", zap.String("orderNo", orderNo), zap.Error(err))
		logFinish("Failed", err)
		return nil, err
	}

	srv.logger.Info("✅ Successfully retrieved Draft Confirm Order", zap.String("orderNo", orderNo), zap.Int("lineCount", len(order.OrderLines)))
	logFinish("Success", nil)

	return order, nil
}

// ListCodeR ดึงรายการ CodeR ที่ขึ้นต้นด้วย 'R' 🚗
func (srv service) ListCodeR(ctx context.Context) ([]response.ListCodeRResponse, error) {
	logFinish := srv.logger.LogAPICall(ctx, "ListCodeR")
	defer logFinish("Completed", nil)

	srv.logger.Info("🔎 Fetching all CodeR from ROM_V_ProductAll (WHERE SKU LIKE 'R%')...")

	codeRList, err := srv.beforeReturnRepo.ListCodeR(ctx)
	if err != nil {
		srv.logger.Error("❌ Failed to list CodeR", zap.Error(err))
		logFinish("Failed", err)
		return nil, fmt.Errorf("ListCodeR failed: %w", err)
	}

	srv.logger.Info("✅ Successfully retrieved CodeR list", zap.Int("count", len(codeRList)))
	logFinish(fmt.Sprintf("Success - %d CodeR", len(codeRList)), nil)

	return codeRList, nil
}

func (srv service) AddCodeR(ctx context.Context, req request.AddCodeR, userID string) ([]response.AddCodeRResponse, error) {
	logFinish := srv.logger.LogAPICall(ctx, "AddCodeR")
	defer logFinish("Completed", nil)

	// ✅ ตรวจสอบค่า `QTY` และ `Price` (ต้องเป็นค่าบวก)
	if req.QTY <= 0 || req.Price <= 0 {
		srv.logger.Warn("⚠️ Invalid QTY or Price",
			zap.Int("qty", req.QTY),
			zap.Float64("price", req.Price),
		)
		logFinish("Failed - Invalid QTY or Price", nil)
		return nil, fmt.Errorf("invalid QTY (%d) or Price (%.2f)", req.QTY, req.Price)
	}

	// ✅ ตั้งค่า `ReturnQTY = QTY`
	req.ReturnQTY = req.QTY

	srv.logger.Info("➕ Adding new CodeR...",
		zap.String("orderNo", req.OrderNo),
		zap.String("sku", req.SKU),
		zap.String("itemName", req.ItemName),
		zap.Int("qty", req.QTY),
		zap.Float64("price", req.Price),
		zap.String("createBy", userID),
	)

	// ✅ เรียกใช้งาน Repository Layer
	results, err := srv.beforeReturnRepo.AddCodeR(ctx, req)
	if err != nil {
		srv.logger.Error("❌ Failed to add CodeR", zap.Error(err))
		logFinish("Failed", err)
		return nil, fmt.Errorf("AddCodeR failed: %w", err)
	}

	srv.logger.Info("✅ Successfully added CodeR", zap.Int("count", len(results)))
	logFinish(fmt.Sprintf("Success - %d records", len(results)), nil)

	return results, nil
}

func (srv service) DeleteCodeR(ctx context.Context, orderNo string, sku string, userID string) error {
	logFinish := srv.logger.LogAPICall(ctx, "DeleteCodeR")
	defer logFinish("Completed", nil)

	srv.logger.Info("🗑️ Deleting CodeR...",
		zap.String("orderNo", orderNo),
		zap.String("sku", sku),
		zap.String("deletedBy", userID),
	)

	// ✅ เรียกใช้งาน Repository Layer
	rowsAffected, err := srv.beforeReturnRepo.DeleteCodeR(ctx, orderNo, sku)
	if err != nil {
		srv.logger.Error("❌ Failed to delete CodeR", zap.Error(err))
		logFinish("Failed", err)
		return fmt.Errorf("DeleteCodeR failed: %w", err)
	}

	// ✅ ตรวจสอบว่ามีข้อมูลถูกลบหรือไม่
	if rowsAffected == 0 {
		srv.logger.Warn("⚠️ CodeR not found", zap.String("orderNo", orderNo), zap.String("sku", sku))
		return fmt.Errorf("no CodeR found with OrderNo: %s and SKU: %s", orderNo, sku)
	}

	srv.logger.Info("✅ Successfully deleted CodeR",
		zap.String("orderNo", orderNo),
		zap.String("sku", sku),
		zap.Int64("rowsAffected", rowsAffected),
	)

	logFinish(fmt.Sprintf("Success - Deleted %d rows", rowsAffected), nil)
	return nil
}

// Method สำหรับอัพเดท Draft Order
func (srv service) UpdateDraftOrder(ctx context.Context, orderNo string, userID string) (*response.UpdateOrderStatusResponse, error) {
	logFinish := srv.logger.LogAPICall(ctx, "UpdateDraftOrder", zap.String("OrderNo", orderNo), zap.String("UserID", userID))
	defer logFinish("Completed", nil)

	srv.logger.Info("🔎 Starting draft order update process 🔎", zap.String("OrderNo", orderNo))

	// ✅ Update order status
	updatedOrder, err := srv.beforeReturnRepo.UpdateOrderStatus(ctx, orderNo, 2, 3, userID) // StatusConfID = 2 (Confirm), StatusReturnID = 3 (Booking)
	if err != nil {
		logFinish("Failed", err)
		srv.logger.Error("❌ Failed to update order status", zap.Error(err))
		return nil, err
	}

	srv.logger.Info("✅ Successfully updated draft order",
		zap.String("OrderNo", updatedOrder.OrderNo),
		zap.Int("StatusConfID", updatedOrder.StatusConfID),
		zap.Int("StatusReturnID", updatedOrder.StatusReturnID),
		zap.String("UpdateBy", updatedOrder.UpdateBy),
		zap.Time("UpdateDate", updatedOrder.UpdateDate),
	)

	logFinish("Success", nil)
	return updatedOrder, nil
}
*/
