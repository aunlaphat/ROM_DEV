package service

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"boilerplate-backend-go/utils"
	"context"
	// "database/sql"
	"fmt"

	"go.uber.org/zap"
)

// BefROService interface กำหนด method สำหรับการทำงานกับ Before Return Order
type BeforeReturnService interface {
	// // Method สำหรับสร้าง Before Return Order พร้อมกับ Lines
	// CreateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	// // Method สำหรับดึงรายการ Before Return Orders ทั้งหมด
	// ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error)
	// // Method สำหรับดึง Before Return Order โดยใช้ OrderNo
	GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error)
	// // Method สำหรับดึงรายการ Before Return Order Lines ทั้งหมด
	// ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error)
	// // Method สำหรับดึง Before Return Order Lines โดยใช้ OrderNo
	GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderItem, error)
	// // Method สำหรับอัพเดท Before Return Order พร้อมกับ Lines
	// UpdateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)

	// // ************************ Create Sale Return ************************ //
	// // Method สำหรับค้นหา Order โดยใช้ SoNo และ OrderNo
	// SearchOrder(ctx context.Context, soNo, orderNo string) ([]response.SaleOrderResponse, error)
	// // Method สำหรับสร้าง Sale Return
	// CreateSaleReturn(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	// // Method สำหรับอัพเดท Sale Return
	// UpdateSaleReturn(ctx context.Context, orderNo string, srNo string, updateBy string) error
	// // Method สำหรับยืนยัน Sale Return
	// ConfirmSaleReturn(ctx context.Context, orderNo string, confirmBy string) error
	// // Method สำหรับยกเลิก Sale Return
	// CancelSaleReturn(ctx context.Context, orderNo string, updateBy string, remark string) error

	// // Method สำหรับดึงรายการ Draft Orders ทั้งหมด
	// ListDraftOrders(ctx context.Context) ([]response.ListDraftConfirmOrdersResponse, error)
	// // Method สำหรับดึงรายการ Confirm Orders ทั้งหมด
	// ListConfirmOrders(ctx context.Context) ([]response.ListDraftConfirmOrdersResponse, error)
	// // Method สำหรับดึง Draft Confirm Order โดยใช้ OrderNo
	// GetDraftConfirmOrderByOrderNo(ctx context.Context, orderNo string) (*response.DraftHeadResponse, error)
	// // Method สำหรับดึง CodeR ทั้งหมด
	// ListCodeR(ctx context.Context) ([]response.CodeRResponse, error)
	// // Method สำหรับเพิ่ม CodeR
	// AddCodeR(ctx context.Context, req request.CodeR) (*response.DraftLineResponse, error)
	// // Method สำหรับลบ CodeR
	// DeleteCodeR(ctx context.Context, orderNo string, sku string) error
	// // Method สำหรับอัพเดท Draft Order
	// UpdateDraftOrder(ctx context.Context, orderNo string, userID string) error

	// Method ดึงข้อมูลรายละเอียดคำสั่งซื้อทั้งหมดพร้อมการแบ่งหน้า
	GetAllOrderDetails(ctx context.Context, page, limit int) ([]response.OrderDetail, error)
	// Method ดึงข้อมูลรายละเอียดคำสั่งซื้อโดยใช้หมายเลข SO
	GetOrderDetailBySO(ctx context.Context, soNo string) (*response.OrderDetail, error)
	// Method ลบรายการ BeforeReturnOrderLine
	DeleteBeforeReturnOrderLine(ctx context.Context, orderNo string, sku string) error
	// Method สร้างคำสั่งซื้อคืนสินค้า
	CreateTradeReturn(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	// Method สร้างรายการคืนสินค้า
	CreateTradeReturnLine(ctx context.Context, orderNo string, lines request.TradeReturnLine) ([]response.BeforeReturnOrderItem, error)
	// Method ยืนยันการรับสินค้าคืนจากหน้าคลัง
	ConfirmReceipt(ctx context.Context, req request.ConfirmTradeReturnRequest, updateBy string) error
	// Method ยืนยันการคืนสินค้าโดยสมบูรณ์
	ConfirmReturn(ctx context.Context, req request.ConfirmToReturnRequest, updateBy string) error
}

// *️⃣ create trade , set statusReturnID = 3 (booking)
func (srv service) CreateTradeReturn(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("[ Starting trade return creation process ]", zap.String("OrderNo", req.OrderNo))

	// *️⃣ ตรวจสอบว่า ReturnOrderLine ต้องไม่เป็นค่าว่าง (มีอย่างน้อย 1 รายการ) จึงจะสามารถสร้างการคืนของได้
	if len(req.BeforeReturnOrderLines) == 0 {
		srv.logger.Warn("[ ReturnOrderLine can't empty must be > 0 line ]")
		return nil, errors.ValidationError("[ ReturnOrderLine can't empty must be > 0 line ]")
	}

	// *️⃣ Validate request ที่ส่งมา
	if err := utils.ValidateCreateTradeReturn(req); err != nil {
		srv.logger.Warn("[ Validation failed ]", zap.Error(err))
		return nil, errors.ValidationError("[ Validation failed: %v ]", err)
	}

	// *️⃣ ตรวจสอบว่า OrderNo สร้างซ้ำหรือไม่
	exists, err := srv.beforeReturnRepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("[ [ Error checking OrderNo existence ]", zap.Error(err)) // db มีปัญหา
		return nil, errors.InternalError("[ Error checking OrderNo existence: %v ]", err)
	}
	if exists != nil {
		srv.logger.Warn("[ Order already exists ]", zap.String("OrderNo", req.OrderNo))
		return nil, errors.ValidationError("[ OrderNo %s already exists: %v ]", req.OrderNo, err)
	}

	// *️⃣ สร้าง trade return order
	createdOrder, err := srv.beforeReturnRepo.CreateTradeReturn(ctx, req)
	if err != nil {
		srv.logger.Error("[ Failed to create trade return order ]", zap.Error(err))
		return nil, errors.InternalError("[ Failed to create trade return order: %v ]", err)
	}

	srv.logger.Info("[ Trade Return created successfully ]", zap.String("OrderNo", req.OrderNo), zap.String("CreateBy", req.CreateBy))
	return createdOrder, nil
}

// *️⃣ add line create trade
func (srv service) CreateTradeReturnLine(ctx context.Context, orderNo string, lines request.TradeReturnLine) ([]response.BeforeReturnOrderItem, error) {
	srv.logger.Info("[ Starting trade return line creation process ]", zap.String("OrderNo", orderNo))

	if orderNo == "" {
		srv.logger.Warn("[ OrderNo is required ]")
		return nil, errors.ValidationError("[ OrderNo is required ]")
	}

	// *️⃣ ตรวจสอบ OrderNo ว่ามีอยู่ใน BeforeReturnOrder
	exists, err := srv.beforeReturnRepo.CheckBefOrderNoExists(ctx, orderNo)
	if err != nil {
		srv.logger.Error("[ Error checking OrderNo existence ]", zap.Error(err)) // db มีปัญหา
		return nil, errors.InternalError("[ Error checking OrderNo existence: %v ]", err)
	}
	if !exists {
		srv.logger.Warn("[ OrderNo not found ]", zap.String("OrderNo", orderNo))
		return nil, errors.NotFoundError("[ This OrderNo not found: %s ]", orderNo)
	}

	// *️⃣ Validate request ที่ส่งมา
	if err := utils.ValidateCreateTradeReturnLine(lines.TradeReturnLine); err != nil {
		srv.logger.Warn("[ Validation failed ]", zap.Error(err))
		return nil, errors.ValidationError("[ Validation failed: %v ]", err)
	}

	// *️⃣ สร้างข้อมูลใน BeforeReturnOrderLine
	err = srv.beforeReturnRepo.CreateTradeReturnLine(ctx, orderNo, lines.TradeReturnLine)
	if err != nil {
		srv.logger.Error("[ Failed to create trade return line ]", zap.Error(err))
		return nil,  errors.InternalError("[ Failed to create trade return line: %v ]", err)
	}

	// *️⃣ ดึงข้อมูลของ order lines ที่เพิ่งสร้างขึ้นมาแสดง
	createdOrderLines, err := srv.beforeReturnRepo.GetBeforeReturnOrderLineByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("[ Failed to fetch created order lines ]", zap.Error(err))
		return nil,  errors.InternalError("[ Failed to fetch created order lines: %v ]", err)
	}

	return createdOrderLines, nil
}

func (srv service) DeleteBeforeReturnOrderLine(ctx context.Context, orderNo string, sku string) error {
	srv.logger.Info("[ Starting delete process ]", zap.String("OrderNo", orderNo), zap.String("SKU", sku))

	if orderNo == "" {
		srv.logger.Warn("[ OrderNo is required ]")
		return errors.ValidationError("[ OrderNo is required ]")
	}

	if sku == "" {
		srv.logger.Warn("[ SKU is required ]")
		return errors.ValidationError("[ SKU is required ]")
	}

	err := srv.beforeReturnRepo.DeleteBeforeReturnOrderLine(ctx, orderNo, sku)
	if err != nil {
		srv.logger.Error("[ Failed to delete order line ]", zap.Error(err))
		return errors.InternalError("[ Failed to delete order line: %v ]", err)
	}

	srv.logger.Info("[ Order Line deleted successfully ]", zap.String("OrderNo", orderNo), zap.String("SKU", sku))
	return nil
}

func (srv service) ConfirmReceipt(ctx context.Context, req request.ConfirmTradeReturnRequest, updateBy string) error {
	srv.logger.Info("[ Starting confirm receipt process ]", zap.String("Identifier", req.Identifier))

	// *️⃣ ตรวจสอบค่าว่าง
	if req.Identifier == "" {
		srv.logger.Warn("[ OrderNo or TrackingNo are required ]")
		return errors.ValidationError("[ OrderNo or TrackingNo are required ]")
	}

	// *️⃣ ตรวจสอบว่า orderNo or trackingNo มีอยู่ในฐานข้อมูล BeforeReturnOrder หรือไม่
	exists, err := srv.beforeReturnRepo.CheckBefOrderOrTrackingExists(ctx, req.Identifier)
	if err != nil {
		srv.logger.Error("[ Failed to check order existence", zap.Error(err))
		return fmt.Errorf("[ Failed to check orderNo or trackingNo existence: %w", err)
	}
	if !exists {
		srv.logger.Warn("[ Not found", zap.String("Identifier", req.Identifier), zap.Error(err))
		return fmt.Errorf("[ Not found: %s", req.Identifier)
	}

	// *️⃣ ตรวจสอบ sku ที่เพิ่มมาว่าตรงกับใน BeforeReturn ที่กรอกเข้ามาไหม หากมีจึงจะสามารถเพิ่มได้ เพราะของหน้าคลังต้องตรงกับข้อมูลที่กรอกเข้าระบบ
	for _, line := range req.ImportLines {
		exists, err := srv.beforeReturnRepo.CheckBefLineSKUExists(ctx, req.Identifier, line.SKU)
		if err != nil {
			srv.logger.Error("[ Failed to check SKU existence", zap.String("SKU", line.SKU), zap.Error(err))
			return  errors.InternalError("[ failed to check SKU existence: %v ]", err)
		}
		if !exists {
			srv.logger.Warn("[ SKU does not exist in BeforeReturnOrderLine from Identifier ]", zap.Error(err))
			return errors.ValidationError("[ SKU %s does not exist in BeforeReturnOrderLine from Identifier %s: %v ]", line.SKU, req.Identifier, err)
		}
	}

	// 1. *️⃣อัปเดตสถานะใน BeforeReturnOrder
	if err := srv.beforeReturnRepo.UpdateBefToWaiting(ctx, req, updateBy); err != nil {
		srv.logger.Error("[ Failed to update BeforeReturnOrder ]", zap.Error(err))
		return  errors.InternalError("[ Failed to update BeforeReturnOrder: %v ]", err)
	}

	// 2. *️⃣ดึงข้อมูลจาก BeforeReturnOrder
	returnOrderData, err := srv.beforeReturnRepo.GetBeforeReturnOrderData(ctx, req)
	if err != nil {
		srv.logger.Error("[ Failed to fetch BeforeReturnOrder ]", zap.Error(err))
		return  errors.InternalError("[ Failed to fetch BeforeReturnOrder: %v ]", err)
	}

	// *️⃣ กำหนดค่าเริ่มต้นให้กับ StatusCheckID ให้เป็นสถานะ waiting
	returnOrderData.StatusCheckID = 1

	// 3. *️⃣Insert ข้อมูลลงใน ReturnOrder
	if err := srv.beforeReturnRepo.InsertReturnOrder(ctx, returnOrderData); err != nil {
		srv.logger.Error("[ Failed to insert into ReturnOrder ]", zap.Error(err))
		return  errors.InternalError("[ Failed to insert into ReturnOrder: %v ]", err)
	}

	// 4. *️⃣Insert ข้อมูลจาก importLines ลงใน ReturnOrderLine + Check ว่า SKU ตรงกับใน BeforeOD ก่อนถึงเพิ่มได้
	if err := srv.beforeReturnRepo.InsertReturnOrderLine(ctx, returnOrderData, req); err != nil {
		srv.logger.Error("[ Failed to insert into ReturnOrderLine ]", zap.Error(err))
		return  errors.InternalError("[ Failed to insert into ReturnOrderLine: %v ]", err)
	}

	// 5. *️⃣Insert ข้อมูลภาพลงใน Images (ไฟล์ภาพ)
	if err := srv.beforeReturnRepo.InsertImages(ctx, returnOrderData, req); err != nil {
		srv.logger.Error("[ Failed to insert images ]", zap.Error(err))
		return  errors.InternalError("[ Failed to insert images: %v ]", err)
	}

	srv.logger.Info("[ Confirm Receipt successfully ]", zap.String("UpdateBy", *req.UpdateBy))
	return nil
}

// *️⃣ check trade line from scan => confirm => success (unsuccess in process future..)
func (srv service) ConfirmReturn(ctx context.Context, req request.ConfirmToReturnRequest, updateBy string) error {
	srv.logger.Info("[ Starting confirm return process ]", zap.String("OrderNo", req.OrderNo))

	if req.OrderNo == "" {
		srv.logger.Warn("[ OrderNo is required ]")
		return errors.ValidationError("[ OrderNo is required ]")
	}

	// *️⃣ ตรวจสอบว่า OrderNo ตรงกับฐานข้อมูลใน BeforeReturn หรือไม่
	exists, err := srv.beforeReturnRepo.CheckBefOrderNoExists(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("[ Failed to check order existence ]", zap.Error(err))
		return  errors.InternalError("[ Failed to check order existence: %v ]", err)
	}
	if !exists {
		srv.logger.Warn("[ OrderNo does not exist in BeforeReturnOrder ]", zap.Error(err))
		return errors.ValidationError("[ OrderNo does not exist in BeforeReturnOrder: %v ]", err)
	}

	// *️⃣ ตรวจสอบ SKU
	for _, line := range req.ImportLinesActual {
		if line.SKU == "" {
			srv.logger.Warn("[ SKU is required ]")
			return errors.ValidationError("[ SKU is required ]")
		}

		exists, err := srv.beforeReturnRepo.CheckReLineSKUExists(ctx, req.OrderNo, line.SKU)
		if err != nil {
			srv.logger.Error("[ failed to check SKU existence ]", zap.Error(err))
			return  errors.InternalError("[ failed to check SKU existence: %v", err)
		}
		if !exists {
			srv.logger.Warn("[ SKU does not exist in ReturnOrderLine from OrderNo ]", zap.Error(err))
			return errors.ValidationError("[ SKU %s does not exist in ReturnOrderLine from OrderNo %s: %v ]", line.SKU, req.OrderNo, err)
		}
	}

	// *️⃣ อัปเดต BeforeReturnOrder
	if err := srv.beforeReturnRepo.UpdateStatusToSuccess(ctx, req.OrderNo, updateBy); err != nil {
		srv.logger.Error("[ Failed to update BeforeReturnOrder ]", zap.Error(err))
		return  errors.InternalError("[ Failed to update BeforeReturnOrder: %v ]", err)
	}

	// *️⃣ ดึงข้อมูล BeforeReturnOrder
	beforeReturnOrder, err := srv.beforeReturnRepo.GetBeforeOrderDetails(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("[ Failed to fetch BeforeReturnOrder details ]", zap.Error(err))
		return  errors.InternalError("[ Failed to fetch BeforeReturnOrder details: %v ]", err)
	}

	// *️⃣ อัปเดต ReturnOrder และ ReturnOrderLine
	if err := srv.beforeReturnRepo.UpdateReturnOrderAndLines(ctx, req, beforeReturnOrder); err != nil {
		srv.logger.Error("[ Failed to fetch updated ReturnOrder and ReturnOrderLine ]", zap.Error(err))
		return  errors.InternalError("[ Failed to fetch updated ReturnOrder and ReturnOrderLine: %v ]", err)
	}

	srv.logger.Info("[ Confirm Return successfully ]", zap.String("OrderNo", req.OrderNo), zap.String("UpdateBy", *req.UpdateBy))
	return nil
}

func (srv service) GetAllOrderDetails(ctx context.Context, page, limit int) ([]response.OrderDetail, error) {
	srv.logger.Info("[ Starting get all order detail process ]")

	offset := (page - 1) * limit

	allorder, err := srv.beforeReturnRepo.GetAllOrderDetails(ctx, offset, limit)
	if err != nil {
		srv.logger.Error("[ Failed to fetch order ]", zap.Error(err))
		return nil, errors.InternalError("[ Failed to fetch order: %v ]", err)
	}
	// if err != nil {
	// 	switch err {
	// 	case sql.ErrNoRows:

	// 		srv.logger.Error("[ no order data: %w", zap.Error(err))
	// 		return nil, fmt.Errorf("[ no order data: %w", err)
	// 	default:

	// 		srv.logger.Error("[ get order error: %w", zap.Error(err))
	// 		return nil, fmt.Errorf("[ get order error: %w", err)
	// 	}
	// }

	srv.logger.Info("[ Successfully fetched Order Details ]", zap.Int("Total amount of data", len(allorder)))
	return allorder, nil
}

func (srv service) GetOrderDetailBySO(ctx context.Context, soNo string) (*response.OrderDetail, error) {
	srv.logger.Info("[ Starting get order detail by SO process ]", zap.String("SoNo", soNo))

	orders, err := srv.beforeReturnRepo.GetOrderDetailBySO(ctx, soNo)
	if err != nil {
		srv.logger.Error("[ Failed to fetch order ]", zap.Error(err))
		return nil, errors.InternalError("[ Failed to fetch order: %v ]", err)
	}

	srv.logger.Info("[ Successfully fetched Order Details ]", zap.String("SoNo", soNo))
	return orders, nil
}

// // Method สำหรับสร้าง Before Return Order พร้อมกับ Lines
// func (srv service) CreateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
// 	logFinish := srv.logger.LogAPICall(ctx, "CreateBeforeReturnOrderWithLines", zap.String("OrderNo", req.OrderNo))
// 	defer logFinish("Completed", nil)
// 	srv.logger.Info("🔎 Starting order creation process", zap.String("OrderNo", req.OrderNo))

// 	err := srv.beforeReturnRepo.CreateBeforeReturnOrderWithTransaction(ctx, req) // เรียก repository เพื่อสร้าง order พร้อมกับ transaction
// 	if err != nil {

// 		srv.logger.Error("❌ Failed to create order with lines", zap.Error(err)) // Logging ว่าการสร้าง order ล้มเหลว
// 		return nil, err
// 	}

// 	createdOrder, err := srv.beforeReturnRepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo) // ดึงข้อมูล order ที่สร้างเสร็จแล้ว
// 	if err != nil {

// 		srv.logger.Error("❌ Failed to fetch created order", zap.Error(err)) // Logging ว่าการดึงข้อมูล order ล้มเหลว
// 		return nil, err
// 	}

// 	return createdOrder, nil
// }

// // Method สำหรับอัพเดท Before Return Order พร้อมกับ Lines
// func (srv service) UpdateBeforeReturnOrderWithLines(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
// 	srv.logger.Info("🔎 Starting order update process", zap.String("OrderNo", req.OrderNo))                    // Logging ว่าเริ่มการอัพเดท order
// 	srv.logger.Debug("Updating order head", zap.String("OrderNo", req.OrderNo), zap.String("SoNo", req.SoNo)) // Logging ว่ากำลังอัพเดท order head

// 	err := srv.beforeReturnRepo.UpdateBeforeReturnOrderWithTransaction(ctx, req) // เรียก repository เพื่ออัพเดท order พร้อมกับ transaction
// 	if err != nil {
// 		srv.logger.Error("❌ Failed to update order with lines", zap.Error(err)) // Logging ว่าการอัพเดท order ล้มเหลว
// 		return nil, err
// 	}

// 	updatedOrder, err := srv.beforeReturnRepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo) // ดึงข้อมูล order ที่อัพเดทเสร็จแล้ว
// 	if err != nil {
// 		srv.logger.Error("❌ Failed to fetch updated order", zap.Error(err)) // Logging ว่าการดึงข้อมูล order ล้มเหลว
// 		return nil, err
// 	}

// 	srv.logger.Info("✅ Successfully updated order with lines", zap.String("OrderNo", req.OrderNo)) // Logging ว่าการอัพเดท order สำเร็จ
// 	return updatedOrder, nil
// }

// // Method สำหรับดึงรายการ Before Return Orders ทั้งหมด
// func (srv service) ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error) {
// 	srv.logger.Info("🔎 Starting to list all return orders")         // Logging ว่าเริ่มการดึงรายการ return orders ทั้งหมด
// 	orders, err := srv.beforeReturnRepo.ListBeforeReturnOrders(ctx) // เรียก repository เพื่อดึงรายการ return orders ทั้งหมด
// 	if err != nil {
// 		srv.logger.Error("❌ Failed to list return orders", zap.Error(err)) // Logging ว่าการดึงรายการ return orders ล้มเหลว
// 		return nil, err
// 	}
// 	srv.logger.Info("✅ Successfully listed return orders") // Logging ว่าการดึงรายการ return orders สำเร็จ
// 	return orders, nil
// }

// Method สำหรับดึง Before Return Order โดยใช้ OrderNo
func (srv service) GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error) {
	srv.logger.Info("[ Starting to get return order by order number ]", zap.String("OrderNo", orderNo))

	order, err := srv.beforeReturnRepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo) 
	if err != nil {
		srv.logger.Error("[ Failed to get return order by order number ]", zap.Error(err)) 
		return nil, errors.InternalError("[ Failed to get return order by order number: %v ]", err)
	}

	return order, nil
}

// // Method สำหรับดึงรายการ Before Return Order Lines ทั้งหมด
// func (srv service) ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error) {
// 	srv.logger.Info("🔎 Starting to list all return order lines")       // Logging ว่าเริ่มการดึงรายการ return order lines ทั้งหมด
// 	lines, err := srv.beforeReturnRepo.ListBeforeReturnOrderLines(ctx) // เรียก repository เพื่อดึงรายการ return order lines ทั้งหมด
// 	if err != nil {
// 		srv.logger.Error("❌ Failed to list return order lines", zap.Error(err)) // Logging ว่าการดึงรายการ return order lines ล้มเหลว
// 		return nil, err
// 	}
// 	srv.logger.Info("✅ Successfully listed return order lines") // Logging ว่าการดึงรายการ return order lines สำเร็จ
// 	return lines, nil
// }

// Method สำหรับดึง Before Return Order Lines โดยใช้ OrderNo
func (srv service) GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderItem, error) {
	srv.logger.Info("[ Starting to get return order lines by order number ]", zap.String("OrderNo", orderNo)) 
	
	lines, err := srv.beforeReturnRepo.GetBeforeReturnOrderLineByOrderNo(ctx, orderNo)                  
	if err != nil {
		srv.logger.Error("[ Failed to get return order lines by order number ]", zap.Error(err))
		return nil, errors.InternalError("[ Failed to get return order order lines by order number: %v ]", err)
	}

	srv.logger.Info("[ Successfully fetched return order lines ]")
	return lines, nil
}

// // ************************ Create Sale Return ************************ //

// // Method สำหรับค้นหา Order โดยใช้ SoNo และ OrderNo
// func (srv service) SearchOrder(ctx context.Context, soNo, orderNo string) ([]response.SaleOrderResponse, error) {
// 	// เริ่มต้น Logging ของ API Call
// 	logFinish := srv.logger.LogAPICall(ctx, "SearchOrder", zap.String("SoNo", soNo), zap.String("OrderNo", orderNo))
// 	defer logFinish("Completed", nil)

// 	// Logging ว่าเริ่มการทำงาน
// 	srv.logger.Info("🔎 Starting to search sale order 🔎", zap.String("SoNo", soNo), zap.String("OrderNo", orderNo))

// 	// เรียก Repository เพื่อค้นหา Order ด้วย SoNo และ OrderNo
// 	order, err := srv.beforeReturnRepo.SearchOrder(ctx, soNo, orderNo)
// 	if err != nil {
// 		// หากเกิดข้อผิดพลาด อัปเดต Log ที่ Error

// 		srv.logger.Error("❌ Failed to search sale orders", zap.Error(err))
// 		return nil, err
// 	}

// 	// กรณีไม่พบข้อมูล
// 	if order == nil {
// 		// หากเกิดข้อผิดพลาด อัปเดต Log ว่าไม่พบข้อมูล
// 		logFinish("Not Found", nil)
// 		srv.logger.Warn("⚠️ No sale order found", zap.String("SoNo", soNo), zap.String("OrderNo", orderNo))
// 		return nil, nil
// 	}

// 	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ

// 	return []response.SaleOrderResponse{*order}, nil
// }

// // Method สำหรับสร้าง Sale Return
// func (srv service) CreateSaleReturn(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
// 	// เริ่มต้น Logging ของ API Call
// 	logFinish := srv.logger.LogAPICall(ctx, "CreateSaleReturn", zap.String("OrderNo", req.OrderNo))
// 	defer logFinish("Completed", nil) // สร้าง closure สำหรับบันทึกสถานะเมื่อฟังก์ชันจบ

// 	// Logging ว่าเริ่มการทำงาน
// 	srv.logger.Info("🔎 Starting sale return creation process 🔎", zap.String("OrderNo", req.OrderNo))

// 	// Validate request
// 	if err := utils.ValidateCreateBeforeReturn(req); err != nil {

// 		srv.logger.Error("❌ Validation failed", zap.Error(err))
// 		return nil, fmt.Errorf("validation failed: %w", err)
// 	}

// 	// ตรวจสอบว่า Order มีอยู่แล้วหรือไม่
// 	existingOrder, err := srv.beforeReturnRepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo)
// 	if err != nil {

// 		srv.logger.Error("❌ Failed to fetch order", zap.Error(err))
// 		return nil, err
// 	}
// 	if existingOrder != nil {
// 		err := fmt.Errorf("order already exists: %s", req.OrderNo)

// 		srv.logger.Warn("⚠️ Duplicate order found", zap.String("OrderNo", req.OrderNo))
// 		return nil, err
// 	}

// 	// สร้าง Sale Return Order
// 	createdOrder, err := srv.beforeReturnRepo.CreateSaleReturn(ctx, req)
// 	if err != nil {

// 		srv.logger.Error("❌ Failed to create order", zap.Error(err))
// 		return nil, err
// 	}

// 	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ

// 	return createdOrder, nil
// }

// // Method สำหรับอัพเดท Sale Return
// func (srv service) UpdateSaleReturn(ctx context.Context, orderNo string, srNo string, updateBy string) error {
// 	// เริ่มต้น Logging ของ API Call
// 	logFinish := srv.logger.LogAPICall(ctx, "UpdateSaleReturn", zap.String("OrderNo", orderNo), zap.String("SrNo", srNo), zap.String("UpdateBy", updateBy))
// 	defer logFinish("Completed", nil)

// 	// Logging ว่าเริ่มการทำงาน
// 	srv.logger.Info("🔎 Starting sale return update process 🔎",
// 		zap.String("OrderNo", orderNo),
// 		zap.String("SrNo", srNo),
// 		zap.String("UpdateBy", updateBy))

// 	// Validation ของ request
// 	if err := utils.ValidateUpdateSaleReturn(orderNo, srNo, updateBy); err != nil {
// 		// หากเกิดข้อผิดพลาด อัปเดต Log ว่าไม่สามารถอัพเดท order ได้

// 		srv.logger.Error("❌ Invalid request", zap.Error(err))
// 		return err
// 	}

// 	// ตรวจสอบสถานะปัจจุบันของ order
// 	order, err := srv.beforeReturnRepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo)
// 	if err != nil {
// 		// อัปเดต Log ว่าไม่สามารถดึงข้อมูล order ได้

// 		srv.logger.Error("❌ Failed to get order", zap.Error(err))
// 		return err
// 	}
// 	if order == nil {
// 		// อัปเดต Log ว่าไม่พบ order
// 		logFinish("Not Found", nil)
// 		srv.logger.Warn("⚠️ Order not found", zap.String("OrderNo", orderNo))
// 		return fmt.Errorf("order not found")
// 	}

// 	// เพิ่มการตรวจสอบสถานะก่อนอัพเดท (ถ้าจำเป็น)
// 	if order.StatusConfID != nil && *order.StatusConfID == 3 { // ถ้าถูกยกเลิกแล้ว
// 		// อัปเดต Log ว่าไม่สามารถอัพเดท order ได้
// 		logFinish("Failed", fmt.Errorf("order is canceled"))
// 		srv.logger.Error("❌ Cannot update canceled order", zap.String("OrderNo", orderNo))
// 		return fmt.Errorf("order is canceled")
// 	}

// 	// เพิ่มการตรวจสอบสถานะเพิ่มเติม
// 	if order.StatusReturnID != nil && *order.StatusReturnID != 1 { // ถ้าไม่ใช่สถานะเริ่มต้น
// 		// อัปเดต Log ว่าไม่สามารถอัพเดท order ได้
// 		logFinish("Failed", fmt.Errorf("invalid status"))
// 		srv.logger.Error("❌ Cannot update SR number: invalid status", zap.String("OrderNo", orderNo))
// 		return fmt.Errorf("invalid status")
// 	}

// 	// อัพเดท SR number
// 	err = srv.beforeReturnRepo.UpdateSaleReturn(ctx, orderNo, srNo, updateBy)
// 	if err != nil {
// 		// อัปเดต Log ว่าไม่สามารถอัพเดท SR number ได้

// 		srv.logger.Error("❌ Failed to update SR number", zap.Error(err))
// 		return err
// 	}

// 	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ

// 	return nil
// }

// // Method สำหรับยืนยัน Sale Return
// func (srv service) ConfirmSaleReturn(ctx context.Context, orderNo string, confirmBy string) error {
// 	// เริ่มต้น Logging ของ API Call
// 	logFinish := srv.logger.LogAPICall(ctx, "ConfirmSaleReturn", zap.String("OrderNo", orderNo), zap.String("ConfirmBy", confirmBy))
// 	defer logFinish("Completed", nil)

// 	// Logging ว่าเริ่มการทำงาน
// 	srv.logger.Info("🔎 Starting sale return confirm process 🔎", zap.String("OrderNo", orderNo))

// 	// ตรวจสอบสถานะปัจจุบันของ order
// 	order, err := srv.beforeReturnRepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo)
// 	if err != nil {
// 		// อัปเดต Log ว่าไม่สามารถยืนยัน order ได้ เนื่องจากเกิดข้อผิดพลาด
// 		logFinish("Failed", fmt.Errorf("failed to get order: %v", err))
// 		srv.logger.Error("❌ Failed to get order", zap.Error(err))
// 		return err
// 	}
// 	if order == nil {
// 		// อัปเดต Log ว่าไม่สามารถยืนยัน order ได้ เนื่องจากไม่พบ order
// 		err = fmt.Errorf("order not found: %s", orderNo)
// 		logFinish("Not Found", err)
// 		srv.logger.Warn("⚠️ Order not found", zap.String("OrderNo", orderNo))
// 		return err
// 	}

// 	// ตรวจสอบว่า order ถูก confirm ไปแล้วหรือไม่
// 	if order.StatusReturnID != nil && *order.StatusReturnID != 1 {
// 		// อัปเดต Log ว่าไม่สามารถยืนยัน order ได้ เนื่องจาก order ไม่ได้เริ่มต้น
// 		err = fmt.Errorf("order %s is not in pending status", orderNo)

// 		srv.logger.Error("❌ Order is not in pending status", zap.String("OrderNo", orderNo))
// 		return err
// 	}
// 	if order.StatusConfID != nil && *order.StatusConfID == 1 {
// 		// อัปเดต Log ว่าไม่สามารถยืนยัน order ได้ เนื่องจาก order ถูกยืนยันข้อมูลไปแล้ว
// 		err = fmt.Errorf("order %s is already confirmed", orderNo)

// 		srv.logger.Error("❌ Order is already confirmed", zap.String("OrderNo", orderNo))
// 		return err
// 	}

// 	// เรียกใช้ repository layer
// 	if err := srv.beforeReturnRepo.ConfirmSaleReturn(ctx, orderNo, confirmBy); err != nil {
// 		// อัปเดต Log ว่าไม่สามารถยืนยัน order ได้ เนื่องจากเกิดข้อผิดพลาด
// 		logFinish("Failed", fmt.Errorf("failed to confirm order: %v", err))
// 		srv.logger.Error("❌ Failed to confirm order", zap.Error(err))
// 		return err
// 	}

// 	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ

// 	return nil
// }

// func (srv service) CancelSaleReturn(ctx context.Context, orderNo string, updateBy string, remark string) error {
// 	// เริ่มต้น Logging ของ API Call
// 	logFinish := srv.logger.LogAPICall(ctx, "CancelSaleReturn", zap.String("OrderNo", orderNo), zap.String("UpdateBy", updateBy), zap.String("Remark", remark))
// 	defer logFinish("Completed", nil)

// 	// Logging ว่าเริ่มการทำงาน
// 	srv.logger.Info("🔎 Starting sale return cancel process 🔎", zap.String("OrderNo", orderNo))

// 	// Input validation
// 	if orderNo == "" || updateBy == "" || remark == "" {
// 		// อัปเดต Log ว่าไม่สามารถยกเลิก order ได้ เนื่องจากข้อมูลไม่ครบ
// 		err := fmt.Errorf("orderNo, updateBy and remark are required")

// 		srv.logger.Error("❌ Invalid input", zap.Error(err))
// 		return err
// 	}

// 	// ตรวจสอบสถานะปัจจุบันของ order
// 	order, err := srv.beforeReturnRepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo)
// 	if err != nil {
// 		// อัปเดต Log ว่าไม่สามารถยกเลิก order ได้ เนื่องจากเกิดข้อผิดพลาด
// 		logFinish("Failed", fmt.Errorf("failed to get order: %v", err))
// 		srv.logger.Error("❌ Failed to get order", zap.Error(err))
// 		return err
// 	}
// 	if order == nil {
// 		// อัปเดต Log ว่าไม่สามารถยกเลิก order ได้ เนื่องจากไม่พบ order
// 		err = fmt.Errorf("order not found: %s", orderNo)
// 		logFinish("Not Found", err)
// 		srv.logger.Warn("⚠️ Order not found", zap.String("OrderNo", orderNo))
// 		return err
// 	}

// 	// ตรวจสอบว่าถูกยกเลิกไปแล้วหรือไม่
// 	if order.StatusConfID != nil && *order.StatusConfID == 3 {
// 		// อัปเดต Log ว่าไม่สามารถยกเลิก order ได้ เนื่องจาก order ถูกยกเลิกไปแล้ว
// 		err = fmt.Errorf("order %s is already canceled", orderNo)

// 		srv.logger.Warn("⚠️ Order is already canceled", zap.String("OrderNo", orderNo))
// 		return err
// 	}
// 	if order.StatusReturnID != nil && *order.StatusReturnID == 2 {
// 		// อัปเดต Log ว่าไม่สามารถยกเลิก order ได้ เนื่องจาก order ถูกยกเลิกไปแล้ว
// 		err = fmt.Errorf("order %s is already canceled", orderNo)

// 		srv.logger.Warn("⚠️ Order is already canceled", zap.String("OrderNo", orderNo))
// 		return err
// 	}

// 	// เรียกใช้ repository layer เพื่อยกเลิก order
// 	if err = srv.beforeReturnRepo.CancelSaleReturn(ctx, orderNo, updateBy, remark); err != nil {
// 		logFinish("Failed", fmt.Errorf("failed to cancel order: %v", err))
// 		srv.logger.Error("❌ Failed to cancel order", zap.Error(err))
// 		return err
// 	}

// 	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ

// 	return nil
// }

// // Method สำหรับดึงรายการ Draft Orders ทั้งหมด
// func (srv service) ListDraftOrders(ctx context.Context) ([]response.ListDraftConfirmOrdersResponse, error) {
// 	// เริ่มต้น Logging ของ API Call
// 	logFinish := srv.logger.LogAPICall(ctx, "ListDraftOrders")
// 	defer logFinish("Completed", nil)

// 	// Logging ว่าเริ่มการทำงาน
// 	srv.logger.Info("🔎 Starting to list all draft orders 🔎")

// 	// เรียก Repository เพื่อค้นหา Order ทั้งหมดที่ Status เป็น Draft
// 	orders, err := srv.beforeReturnRepo.ListDraftOrders(ctx)
// 	if err != nil {
// 		// หากเกิดข้อผิดพลาด อัปเดต Log ที่ Error
// 		logFinish("Failed", fmt.Errorf("❌ Failed to list draft orders : %v", err))
// 		srv.logger.Error("❌ Failed to list draft orders", zap.Error(err))
// 		return nil, err
// 	}

// 	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ

// 	return orders, nil
// }

// // Method สำหรับดึงรายการ Confirm Orders ทั้งหมด
// func (srv service) ListConfirmOrders(ctx context.Context) ([]response.ListDraftConfirmOrdersResponse, error) {
// 	// เริ่มต้น Logging ของ API Call
// 	logFinish := srv.logger.LogAPICall(ctx, "ListConfirmOrders")
// 	defer logFinish("Completed", nil)

// 	srv.logger.Info("🔎 Starting to list all confirm orders 🔎")

// 	// เรียก Repository เพื่อค้นหา Order ทั้งหมดที่ Status เป็น Confirm
// 	orders, err := srv.beforeReturnRepo.ListConfirmOrders(ctx)
// 	if err != nil {
// 		// หากเกิดข้อผิดพลาด อัปเดต Log ที่ Error
// 		logFinish("Failed", fmt.Errorf("❌ Failed to list confirm orders : %v", err))
// 		srv.logger.Error("❌ Failed to list confirm orders", zap.Error(err))
// 		return nil, err
// 	}

// 	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ

// 	return orders, nil
// }

// // Method สำหรับดึง Draft Confirm Order โดยใช้ OrderNo
// func (srv service) GetDraftConfirmOrderByOrderNo(ctx context.Context, orderNo string) (*response.DraftHeadResponse, error) {
// 	// เริ่มต้น Logging ของ API Call
// 	logFinish := srv.logger.LogAPICall(ctx, "GetDraftConfirmOrderByOrderNo", zap.String("OrderNo", orderNo))
// 	defer logFinish("Completed", nil)

// 	// Logging ว่าเริ่มการทำงาน
// 	srv.logger.Info("🔎 Starting to get draft order by order number 🔎", zap.String("OrderNo", orderNo))

// 	head, lines, err := srv.beforeReturnRepo.GetDraftConfirmOrderByOrderNo(ctx, orderNo)
// 	if err != nil {
// 		// อัปเดต Log ว่าไม่สามารถดึงข้อมูลได้
// 		logFinish("Failed", fmt.Errorf("❌ Failed to get draft order : %v", err))
// 		srv.logger.Error("❌ Failed to get draft order", zap.Error(err))
// 		return nil, err
// 	}

// 	head.OrderLines = lines

// 	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ

// 	return head, nil
// }

// // Method สำหรับดึง CodeR ทั้งหมด
// func (srv service) ListCodeR(ctx context.Context) ([]response.CodeRResponse, error) {
// 	// เริ่มต้น Logging ของ API Call
// 	logFinish := srv.logger.LogAPICall(ctx, "GetCodeR")
// 	defer logFinish("Completed", nil)

// 	// Logging ว่าเริ่มการทำงาน
// 	srv.logger.Info("🔎 Starting to get CodeR 🔎")

// 	// เรียก Repository เพื่อค้นหา CodeR ทั้งหมด
// 	codeR, err := srv.beforeReturnRepo.ListCodeR(ctx)
// 	if err != nil {
// 		// หากเกิดข้อผิดพลาด อัปเดต Log ที่ Error
// 		logFinish("Failed", fmt.Errorf("❌ Failed to get CodeR : %v", err))
// 		srv.logger.Error("❌ Failed to get CodeR", zap.Error(err))
// 		return nil, err
// 	}

// 	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ

// 	return codeR, nil
// }

// // Method สำหรับเพิ่ม CodeR
// func (srv service) AddCodeR(ctx context.Context, req request.CodeR) (*response.DraftLineResponse, error) {
// 	// เริ่มต้น Logging ของ API Call
// 	logFinish := srv.logger.LogAPICall(ctx, "AddCodeR")
// 	defer logFinish("Completed", nil)

// 	// Logging ว่าเริ่มการทำงาน
// 	srv.logger.Info("🔎 Starting to add CodeR 🔎")

// 	// ตรวจสอบว่า SKU มีอยู่แล้วหรือไม่
// 	existingLines, err := srv.beforeReturnRepo.GetBeforeReturnOrderLineByOrderNo(ctx, req.OrderNo)
// 	if err != nil {
// 		logFinish("Failed", fmt.Errorf("failed to check existing SKUs: %v", err))
// 		srv.logger.Error("❌ Failed to check existing SKUs", zap.Error(err))
// 		return nil, err
// 	}

// 	for _, line := range existingLines {
// 		if line.SKU == req.SKU {
// 			err := fmt.Errorf("SKU already exists for OrderNo: %s", req.OrderNo)

// 			srv.logger.Warn("⚠️ Duplicate SKU found", zap.String("OrderNo", req.OrderNo), zap.String("SKU", req.SKU))
// 			return nil, err
// 		}
// 	}

// 	// เรียกใช้ repository layer
// 	result, err := srv.beforeReturnRepo.AddCodeR(ctx, req)
// 	if err != nil {
// 		logFinish("Failed", fmt.Errorf("failed to add CodeR: %v", err))
// 		srv.logger.Error("❌ Failed to add CodeR", zap.Error(err))
// 		return nil, err
// 	}

// 	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ

// 	return result, nil
// }

// // Method สำหรับลบ CodeR
// func (srv service) DeleteCodeR(ctx context.Context, orderNo string, sku string) error {
// 	// เริ่มต้น Logging ของ API Call
// 	logFinish := srv.logger.LogAPICall(ctx, "DeleteCodeR", zap.String("OrderNo", orderNo), zap.String("SKU", sku))
// 	defer logFinish("Completed", nil)

// 	// เรียกใช้ repository layer
// 	if err := srv.beforeReturnRepo.DeleteCodeR(ctx, orderNo, sku); err != nil {
// 		// อัปเดต Log ว่าไม่สามารถลบ CodeR ได้ เนื่องจากเกิดข้อผิดพลาด
// 		logFinish("Failed", fmt.Errorf("failed to delete CodeR: %v", err))
// 		srv.logger.Error("❌ Failed to delete CodeR", zap.Error(err))
// 		return err
// 	}

// 	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ

// 	return nil
// }

// // Method สำหรับอัพเดท Draft Order
// func (srv service) UpdateDraftOrder(ctx context.Context, orderNo string, userID string) error {
// 	// เริ่มต้น Logging ของ API Call
// 	logFinish := srv.logger.LogAPICall(ctx, "UpdateDraftOrder", zap.String("OrderNo", orderNo), zap.String("UserID", userID))
// 	defer logFinish("Completed", nil)

// 	// Logging ว่าเริ่มการทำงาน
// 	srv.logger.Info("🔎 Starting draft order update process 🔎", zap.String("OrderNo", orderNo))

// 	// Update order status
// 	err := srv.beforeReturnRepo.UpdateOrderStatus(ctx, orderNo, 2, 3, userID) // StatusConfID = 2 (Confirm), StatusReturnID = 3 (Booking)
// 	if err != nil {

// 		srv.logger.Error("❌ Failed to update order status", zap.Error(err))
// 		return err
// 	}

// 	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ

// 	return nil
// }
