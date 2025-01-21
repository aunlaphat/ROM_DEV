package service

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/utils"
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
	SearchOrder(ctx context.Context, soNo, orderNo string) ([]response.SaleOrderResponse, error)
	CreateSaleReturn(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	UpdateSaleReturn(ctx context.Context, orderNo string, srNo string, updateBy string) error
	ConfirmSaleReturn(ctx context.Context, orderNo string, confirmBy string) error
	CancelSaleReturn(ctx context.Context, orderNo string, updateBy string, remark string) error

	ListDraftOrders(ctx context.Context) ([]response.ListDraftConfirmOrdersResponse, error)
	ListConfirmOrders(ctx context.Context) ([]response.ListDraftConfirmOrdersResponse, error)
	GetDraftConfirmOrderByOrderNo(ctx context.Context, orderNo string) (*response.DraftHeadResponse, error)
	GetCodeR(ctx context.Context) ([]response.CodeRResponse, error)
	AddCodeR(ctx context.Context, req request.CodeRRequest) error
	DeleteCodeR(ctx context.Context, sku string) error
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
	srv.logger.Info("✅ Successfully fetched return order lines")
	return lines, nil
}

// ************************ Create Sale Return ************************ //

func (srv service) SearchOrder(ctx context.Context, soNo, orderNo string) ([]response.SaleOrderResponse, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "SearchOrder", zap.String("SoNo", soNo), zap.String("OrderNo", orderNo))
	defer logFinish("Completed", nil)

	// Logging ว่าเริ่มการทำงาน
	srv.logger.Info("🔎 Starting to search sale order 🔎", zap.String("SoNo", soNo), zap.String("OrderNo", orderNo))

	// เรียก Repository เพื่อค้นหา Order ด้วย SoNo และ OrderNo
	order, err := srv.befRORepo.SearchOrder(ctx, soNo, orderNo)
	if err != nil {
		// หากเกิดข้อผิดพลาด อัปเดต Log ที่ Error
		logFinish("Failed", err)
		srv.logger.Error("❌ Failed to search sale orders", zap.Error(err))
		return nil, err
	}

	// กรณีไม่พบข้อมูล
	if order == nil {
		// หากเกิดข้อผิดพลาด อัปเดต Log ว่าไม่พบข้อมูล
		logFinish("Not Found", nil)
		srv.logger.Warn("⚠️ No sale order found", zap.String("SoNo", soNo), zap.String("OrderNo", orderNo))
		return nil, nil
	}

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	srv.logger.Info("✅ Successfully searched sale orders",
		zap.String("SoNo", soNo),
		zap.String("OrderNo", orderNo))

	// ส่งค่าผลลัพธ์กลับไป
	return []response.SaleOrderResponse{*order}, nil
}

func (srv service) CreateSaleReturn(ctx context.Context, req request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "CreateSaleReturn", zap.String("OrderNo", req.OrderNo))
	defer logFinish("Completed", nil) // สร้าง closure สำหรับบันทึกสถานะเมื่อฟังก์ชันจบ

	// Logging ว่าเริ่มการทำงาน
	srv.logger.Info("🔎 Starting sale return creation process 🔎", zap.String("OrderNo", req.OrderNo))

	// Validate request
	if err := utils.ValidateCreateSaleReturn(req); err != nil {
		logFinish("Failed", err)
		srv.logger.Error("❌ Validation failed", zap.Error(err))
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// ตรวจสอบว่า Order มีอยู่แล้วหรือไม่
	existingOrder, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		logFinish("Failed", err)
		srv.logger.Error("❌ Failed to fetch order", zap.Error(err))
		return nil, err
	}
	if existingOrder != nil {
		err := fmt.Errorf("order already exists: %s", req.OrderNo)
		logFinish("Failed", err)
		srv.logger.Warn("⚠️ Duplicate order found", zap.String("OrderNo", req.OrderNo))
		return nil, err
	}

	// สร้าง Sale Return Order
	createdOrder, err := srv.befRORepo.CreateSaleReturn(ctx, req)
	if err != nil {
		logFinish("Failed", err)
		srv.logger.Error("❌ Failed to create order", zap.Error(err))
		return nil, err
	}

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	srv.logger.Info("✅ Sale return created successfully", zap.String("OrderNo", req.OrderNo))
	return createdOrder, nil
}

func (srv service) UpdateSaleReturn(ctx context.Context, orderNo string, srNo string, updateBy string) error {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "UpdateSaleReturn", zap.String("OrderNo", orderNo), zap.String("SrNo", srNo), zap.String("UpdateBy", updateBy))
	defer logFinish("Completed", nil)

	// Logging ว่าเริ่มการทำงาน
	srv.logger.Info("🔎 Starting sale return update process 🔎",
		zap.String("OrderNo", orderNo),
		zap.String("SrNo", srNo),
		zap.String("UpdateBy", updateBy))

	// Validation ของ request
	if err := utils.ValidateUpdateSaleReturn(orderNo, srNo, updateBy); err != nil {
		// หากเกิดข้อผิดพลาด อัปเดต Log ว่าไม่สามารถอัพเดท order ได้
		logFinish("Failed", err)
		srv.logger.Error("❌ Invalid request", zap.Error(err))
		return err
	}

	// ตรวจสอบสถานะปัจจุบันของ order
	order, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo)
	if err != nil {
		// อัปเดต Log ว่าไม่สามารถดึงข้อมูล order ได้
		logFinish("Failed", err)
		srv.logger.Error("❌ Failed to get order", zap.Error(err))
		return err
	}
	if order == nil {
		// อัปเดต Log ว่าไม่พบ order
		logFinish("Not Found", nil)
		srv.logger.Warn("⚠️ Order not found", zap.String("OrderNo", orderNo))
		return fmt.Errorf("order not found")
	}

	// เพิ่มการตรวจสอบสถานะก่อนอัพเดท (ถ้าจำเป็น)
	if order.StatusConfID != nil && *order.StatusConfID == 3 { // ถ้าถูกยกเลิกแล้ว
		// อัปเดต Log ว่าไม่สามารถอัพเดท order ได้
		logFinish("Failed", fmt.Errorf("order is canceled"))
		srv.logger.Error("❌ Cannot update canceled order", zap.String("OrderNo", orderNo))
		return fmt.Errorf("order is canceled")
	}

	// เพิ่มการตรวจสอบสถานะเพิ่มเติม
	if order.StatusReturnID != nil && *order.StatusReturnID != 1 { // ถ้าไม่ใช่สถานะเริ่มต้น
		// อัปเดต Log ว่าไม่สามารถอัพเดท order ได้
		logFinish("Failed", fmt.Errorf("invalid status"))
		srv.logger.Error("❌ Cannot update SR number: invalid status", zap.String("OrderNo", orderNo))
		return fmt.Errorf("invalid status")
	}

	// อัพเดท SR number
	err = srv.befRORepo.UpdateSaleReturn(ctx, orderNo, srNo, updateBy)
	if err != nil {
		// อัปเดต Log ว่าไม่สามารถอัพเดท SR number ได้
		logFinish("Failed", err)
		srv.logger.Error("❌ Failed to update SR number", zap.Error(err))
		return err
	}

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	srv.logger.Info("✅ Successfully updated SR number",
		zap.String("OrderNo", orderNo),
		zap.String("SrNo", srNo),
		zap.String("UpdateBy", updateBy))

	return nil
}

func (srv service) ConfirmSaleReturn(ctx context.Context, orderNo string, confirmBy string) error {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "ConfirmSaleReturn", zap.String("OrderNo", orderNo), zap.String("ConfirmBy", confirmBy))
	defer logFinish("Completed", nil)

	// ตรวจสอบสถานะปัจจุบันของ order
	order, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo)
	if err != nil {
		// อัปเดต Log ว่าไม่สามารถยืนยัน order ได้ เนื่องจากเกิดข้อผิดพลาด
		logFinish("Failed", fmt.Errorf("failed to get order: %v", err))
		srv.logger.Error("❌ Failed to get order", zap.Error(err))
		return err
	}
	if order == nil {
		// อัปเดต Log ว่าไม่สามารถยืนยัน order ได้ เนื่องจากไม่พบ order
		err = fmt.Errorf("order not found: %s", orderNo)
		logFinish("Not Found", err)
		srv.logger.Warn("⚠️ Order not found", zap.String("OrderNo", orderNo))
		return err
	}

	// ตรวจสอบว่า order ถูก confirm ไปแล้วหรือไม่
	if order.StatusReturnID != nil && *order.StatusReturnID != 1 {
		// อัปเดต Log ว่าไม่สามารถยืนยัน order ได้ เนื่องจาก order ไม่ได้เริ่มต้น
		err = fmt.Errorf("order %s is not in pending status", orderNo)
		logFinish("Failed", err)
		srv.logger.Error("❌ Order is not in pending status", zap.String("OrderNo", orderNo))
		return err
	}
	if order.StatusConfID != nil && *order.StatusConfID == 1 {
		// อัปเดต Log ว่าไม่สามารถยืนยัน order ได้ เนื่องจาก order ถูกยืนยันข้อมูลไปแล้ว
		err = fmt.Errorf("order %s is already confirmed", orderNo)
		logFinish("Failed", err)
		srv.logger.Error("❌ Order is already confirmed", zap.String("OrderNo", orderNo))
		return err
	}

	// เรียกใช้ repository layer
	if err := srv.befRORepo.ConfirmSaleReturn(ctx, orderNo, confirmBy); err != nil {
		// อัปเดต Log ว่าไม่สามารถยืนยัน order ได้ เนื่องจากเกิดข้อผิดพลาด
		logFinish("Failed", fmt.Errorf("failed to confirm order: %v", err))
		srv.logger.Error("❌ Failed to confirm order", zap.Error(err))
		return err
	}

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	srv.logger.Info("✅ Successfully confirmed order", zap.String("OrderNo", orderNo))
	return nil
}

func (srv service) CancelSaleReturn(ctx context.Context, orderNo string, updateBy string, remark string) error {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "CancelSaleReturn", zap.String("OrderNo", orderNo), zap.String("UpdateBy", updateBy), zap.String("Remark", remark))
	defer logFinish("Completed", nil)

	// Input validation
	if orderNo == "" || updateBy == "" || remark == "" {
		// อัปเดต Log ว่าไม่สามารถยกเลิก order ได้ เนื่องจากข้อมูลไม่ครบ
		err := fmt.Errorf("orderNo, updateBy and remark are required")
		logFinish("Failed", err)
		srv.logger.Error("❌ Invalid input", zap.Error(err))
		return err
	}

	// ตรวจสอบสถานะปัจจุบันของ order
	order, err := srv.befRORepo.GetBeforeReturnOrderByOrderNo(ctx, orderNo)
	if err != nil {
		// อัปเดต Log ว่าไม่สามารถยกเลิก order ได้ เนื่องจากเกิดข้อผิดพลาด
		logFinish("Failed", fmt.Errorf("failed to get order: %v", err))
		srv.logger.Error("❌ Failed to get order", zap.Error(err))
		return err
	}
	if order == nil {
		// อัปเดต Log ว่าไม่สามารถยกเลิก order ได้ เนื่องจากไม่พบ order
		err = fmt.Errorf("order not found: %s", orderNo)
		logFinish("Not Found", err)
		srv.logger.Warn("⚠️ Order not found", zap.String("OrderNo", orderNo))
		return err
	}

	// ตรวจสอบว่าถูกยกเลิกไปแล้วหรือไม่
	if order.StatusConfID != nil && *order.StatusConfID == 3 {
		// อัปเดต Log ว่าไม่สามารถยกเลิก order ได้ เนื่องจาก order ถูกยกเลิกไปแล้ว
		err = fmt.Errorf("order %s is already canceled", orderNo)
		logFinish("Failed", err)
		srv.logger.Warn("⚠️ Order is already canceled", zap.String("OrderNo", orderNo))
		return err
	}
	if order.StatusReturnID != nil && *order.StatusReturnID == 2 {
		// อัปเดต Log ว่าไม่สามารถยกเลิก order ได้ เนื่องจาก order ถูกยกเลิกไปแล้ว
		err = fmt.Errorf("order %s is already canceled", orderNo)
		logFinish("Failed", err)
		srv.logger.Warn("⚠️ Order is already canceled", zap.String("OrderNo", orderNo))
		return err
	}

	// เรียกใช้ repository layer เพื่อยกเลิก order
	if err = srv.befRORepo.CancelSaleReturn(ctx, orderNo, updateBy, remark); err != nil {
		logFinish("Failed", fmt.Errorf("failed to cancel order: %v", err))
		srv.logger.Error("❌ Failed to cancel order", zap.Error(err))
		return err
	}

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	srv.logger.Info("✅ Successfully canceled order", zap.String("OrderNo", orderNo))
	return nil
}

func (srv service) ListDraftOrders(ctx context.Context) ([]response.ListDraftConfirmOrdersResponse, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "ListDraftOrders")
	defer logFinish("Completed", nil)

	srv.logger.Info("🏁 Starting to list all draft orders")

	// เรียก Repository เพื่อค้นหา Order ทั้งหมดที่ Status เป็น Draft
	orders, err := srv.befRORepo.ListDraftOrders(ctx)
	if err != nil {
		// หากเกิดข้อผิดพลาด อัปเดต Log ที่ Error
		logFinish("Failed", fmt.Errorf("❌ Failed to list draft orders : %v", err))
		srv.logger.Error("❌ Failed to list draft orders", zap.Error(err))
		return nil, err
	}

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	srv.logger.Info("✅ Successfully listed draft orders")
	return orders, nil
}

func (srv service) ListConfirmOrders(ctx context.Context) ([]response.ListDraftConfirmOrdersResponse, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "ListConfirmOrders")
	defer logFinish("Completed", nil)

	srv.logger.Info("🏁 Starting to list all confirm orders")

	// เรียก Repository เพื่อค้นหา Order ทั้งหมดที่ Status เป็น Confirm
	orders, err := srv.befRORepo.ListConfirmOrders(ctx)
	if err != nil {
		// หากเกิดข้อผิดพลาด อัปเดต Log ที่ Error
		logFinish("Failed", fmt.Errorf("❌ Failed to list confirm orders : %v", err))
		srv.logger.Error("❌ Failed to list confirm orders", zap.Error(err))
		return nil, err
	}

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	srv.logger.Info("✅ Successfully listed confirm orders")
	return orders, nil
}

func (srv service) GetDraftConfirmOrderByOrderNo(ctx context.Context, orderNo string) (*response.DraftHeadResponse, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "GetDraftOrderByOrderNo", zap.String("OrderNo", orderNo))
	defer logFinish("Completed", nil)

	srv.logger.Info("🏁 Starting to get draft order by order number", zap.String("OrderNo", orderNo))
	head, lines, err := srv.befRORepo.GetDraftConfirmOrderByOrderNo(ctx, orderNo)
	if err != nil {
		// อัปเดต Log ว่าไม่สามารถดึงข้อมูลได้
		logFinish("Failed", fmt.Errorf("❌ Failed to get draft order : %v", err))
		srv.logger.Error("❌ Failed to get draft order", zap.Error(err))
		return nil, err
	}

	head.OrderLines = lines

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	srv.logger.Info("✅ Successfully fetched draft order", zap.String("OrderNo", orderNo))
	return head, nil
}

func (srv service) GetCodeR(ctx context.Context) ([]response.CodeRResponse, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "GetCodeR")
	defer logFinish("Completed", nil)

	// เรียก Repository เพื่อค้นหา CodeR ทั้งหมด
	codeR, err := srv.befRORepo.GetCodeR(ctx)
	if err != nil {
		// หากเกิดข้อผิดพลาด อัปเดต Log ที่ Error
		logFinish("Failed", fmt.Errorf("❌ Failed to get CodeR : %v", err))
		srv.logger.Error("❌ Failed to get CodeR", zap.Error(err))
		return nil, err
	}

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	srv.logger.Info("✅ Successfully fetched CodeR")
	return codeR, nil
}

func (srv service) AddCodeR(ctx context.Context, req request.CodeRRequest) error {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "AddCodeR")
	defer logFinish("Completed", nil)

	// เรียกใช้ repository layer
	if err := srv.befRORepo.AddCodeR(ctx, req); err != nil {
		// อัปเดต Log ว่าไม่สามารถเพิ่ม CodeR ได้ เนื่องจากเกิดข้อผิดพลาด
		logFinish("Failed", fmt.Errorf("failed to add CodeR: %v", err))
		srv.logger.Error("❌ Failed to add CodeR", zap.Error(err))
		return err
	}

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	srv.logger.Info("✅ Successfully added CodeR")
	return nil
}

func (srv service) DeleteCodeR(ctx context.Context, sku string) error {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "DeleteCodeR", zap.String("SKU", sku))
	defer logFinish("Completed", nil)

	// เรียกใช้ repository layer
	if err := srv.befRORepo.DeleteCodeR(ctx, sku); err != nil {
		// อัปเดต Log ว่าไม่สามารถลบ CodeR ได้ เนื่องจากเกิดข้อผิดพลาด
		logFinish("Failed", fmt.Errorf("failed to delete CodeR: %v", err))
		srv.logger.Error("❌ Failed to delete CodeR", zap.Error(err))
		return err
	}

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	srv.logger.Info("✅ Successfully deleted CodeR", zap.String("SKU", sku))
	return nil
}
