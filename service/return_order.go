package service

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"boilerplate-backend-go/utils"
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

// ตัวสื่อกลางในการรับส่งกับ API และประมวลผลข้อมูลที่รับมาจาก API
type ReturnOrderService interface {
	GetAllReturnOrder(ctx context.Context) ([]response.ReturnOrder, error)
	GetReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.ReturnOrder, error)
	GetAllReturnOrderLines(ctx context.Context) ([]response.ReturnOrderLine, error)
	GetReturnOrderLinesByReturnID(ctx context.Context, orderNo string) ([]response.ReturnOrderLine, error)
	CreateReturnOrder(ctx context.Context, req request.CreateReturnOrder) (*response.CreateReturnOrder, error)
	UpdateReturnOrder(ctx context.Context, req request.UpdateReturnOrder, updateBy string) (*response.UpdateReturnOrder, error) 
	DeleteReturnOrder(ctx context.Context, orderNo string) error
}

func (srv service) GetAllReturnOrder(ctx context.Context) ([]response.ReturnOrder, error) {
	// Step 1: เรียก repository เพื่อดึงข้อมูล ReturnOrder ทั้งหมด
	allorder, err := srv.returnOrderRepo.GetAllReturnOrder(ctx)
	if err != nil {
		srv.logger.Error("Error fetching all return orders", zap.Error(err))
		// Step 2: หากเกิดข้อผิดพลาด ให้ส่ง Error กลับไปยัง API
		return nil, errors.UnexpectedError()
	}

	// Step 3: ส่งข้อมูล ReturnOrder ที่ได้กลับไปยัง API
	return allorder, nil
}

func (srv service) GetReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.ReturnOrder, error) {
	// Step 1: ตรวจสอบว่า OrderNo ไม่เป็นค่าว่าง
	if orderNo == "" {
		return nil, errors.ValidationError("OrderNo is required")
	}

	// Step 2: เรียก repository เพื่อดึงข้อมูล ReturnOrder โดยใช้ OrderNo
	idorder, err := srv.returnOrderRepo.GetReturnOrderByOrderNo(ctx, orderNo)
	if err != nil {
		if err == sql.ErrNoRows {
			// Step 3: หากไม่พบข้อมูล ReturnOrder ให้ส่ง Error กลับไปยัง API
			return nil, errors.NotFoundError("Return order not found")
		}
		srv.logger.Error("Error fetching ReturnOrder by ID", zap.Error(err))
		return nil, errors.UnexpectedError()
	}

	// Step 4: ส่งข้อมูล ReturnOrder ที่ได้กลับไปยัง API
	return idorder, nil
}

func (srv service) GetAllReturnOrderLines(ctx context.Context) ([]response.ReturnOrderLine, error) {
	lines, err := srv.returnOrderRepo.GetAllReturnOrderLines(ctx)
	if err != nil {
		srv.logger.Error("Error fetching all return order lines", zap.Error(err))
		return nil, errors.UnexpectedError()
	}

	return lines, nil
}

func (srv service) GetReturnOrderLinesByReturnID(ctx context.Context, orderNo string) ([]response.ReturnOrderLine, error) {
	if orderNo == "" {
		return nil, errors.ValidationError("OrderNo is required")
	}

	lines, err := srv.returnOrderRepo.GetReturnOrderLinesByReturnID(ctx, orderNo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFoundError("This Return Order Line not found")
		}
		srv.logger.Error("Error fetching return order lines by OrderNo", zap.Error(err))
		return nil, errors.UnexpectedError()
	}

	return lines, nil
}

func (srv service) CreateReturnOrder(ctx context.Context, req request.CreateReturnOrder) (*response.CreateReturnOrder, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "CreateReturnOrder", zap.String("OrderNo", req.OrderNo))
	defer logFinish("Completed", nil)

	srv.logger.Info("🔎 Starting return order creation process 🔎", zap.String("OrderNo", req.OrderNo))

	// Validate request
	if err := utils.ValidateCreateReturnOrder(req); err != nil {
		logFinish("Failed", err)
		srv.logger.Error("❌ Validation failed", zap.Error(err))
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// ตรวจสอบว่า OrderNo ซ้ำหรือไม่
	exists, err := srv.returnOrderRepo.CheckOrderNoExist(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("Failed to check OrderNo", zap.Error(err))
		return nil, errors.InternalError("Failed to check OrderNo")
	}
	if exists {
		srv.logger.Error("OrderNo already exists", zap.Error(err))
		return nil, errors.BadRequestError("OrderNo already exists")
	}

	err = srv.returnOrderRepo.CreateReturnOrder(ctx, req)
	if err != nil {
		srv.logger.Error("❌ Failed to create order with lines", zap.Error(err)) // Logging ว่าการสร้าง order ล้มเหลว
		return nil, err
	}

	// ดึงข้อมูล ReturnOrder ที่สร้างขึ้นมาใหม่ไปแสดง
	createdOrder, err := srv.returnOrderRepo.GetCreateReturnOrder(ctx, req.OrderNo)
	if err != nil {
		logFinish("Failed", err)
		srv.logger.Error("❌ Failed to fetch created order", zap.Error(err))
		return nil, err
	}

	// Logging สำเร็จ และอัปเดต Log ว่าสำเร็จ
	logFinish("Success", nil)
	return createdOrder, nil
}

func (srv service) UpdateReturnOrder(ctx context.Context, req request.UpdateReturnOrder, updateBy string) (*response.UpdateReturnOrder, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "UpdateReturnOrder", zap.String("UpdateBy", updateBy))
	defer logFinish("Completed", nil)
	srv.logger.Info("🔎 Starting return order update process 🔎")

	// Step 1: ตรวจสอบว่า OrderNo ไม่เป็นค่าว่าง
	if req.OrderNo == "" {
		return nil, errors.ValidationError("OrderNo is required")
	}

	exists, err := srv.returnOrderRepo.CheckOrderNoExist(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("Error checking OrderNo existence", zap.Error(err))
		return nil, errors.UnexpectedError()
	}
	if !exists {
		return nil, errors.NotFoundError("OrderNo not found")
	}

	// Step 2: เรียก repository เพื่ออัปเดต ReturnOrder
	err = srv.returnOrderRepo.UpdateReturnOrder(ctx, req)
	if err != nil {
		srv.logger.Error("Error updating ReturnOrder", zap.Error(err))
		return nil, errors.UnexpectedError()
	}

	updatedOrder, err := srv.returnOrderRepo.GetUpdateReturnOrder(ctx, req.OrderNo) // ดึงข้อมูล order ที่อัพเดทเสร็จแล้ว
	if err != nil {
		srv.logger.Error("❌ Failed to fetch updated order", zap.Error(err)) // Logging ว่าการดึงข้อมูล order ล้มเหลว
		return nil, err
	}

	logFinish("Success", nil)
	return updatedOrder, nil
}

func (srv service) DeleteReturnOrder(ctx context.Context, orderNo string) error {
	if orderNo == "" {
		return errors.ValidationError("OrderNo is required")
	}

	// Step 2: เรียก repository เพื่อลบ ReturnOrder
	err := srv.returnOrderRepo.DeleteReturnOrder(ctx, orderNo)
	if err != nil {
		srv.logger.Error("Error deleting ReturnOrder", zap.Error(err))
		return errors.UnexpectedError()
	}

	return nil
}
