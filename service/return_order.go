package service

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"context"
	"database/sql"

	"go.uber.org/zap"
)

// ตัวสื่อกลางในการรับส่งกับ API และประมวลผลข้อมูลที่รับมาจาก API
type ReturnOrderService interface { 
	AllGetReturnOrder(ctx context.Context) ([]response.ReturnOrder, error)
	GetReturnOrderByID(ctx context.Context, returnID string) (*response.ReturnOrder, error)
	GetAllReturnOrderLines(ctx context.Context) ([]response.ReturnOrderLine, error)
	GetReturnOrderLinesByReturnID(ctx context.Context, returnID string) ([]response.ReturnOrderLine, error)
	CreateReturnOrder(ctx context.Context, req request.CreateReturnOrder) error
	UpdateReturnOrder(ctx context.Context, req request.UpdateReturnOrder) error
	DeleteReturnOrder(ctx context.Context, returnID string) error
}

func (srv service) AllGetReturnOrder(ctx context.Context) ([]response.ReturnOrder, error) {
	// Step 1: เรียก repository เพื่อดึงข้อมูล ReturnOrder ทั้งหมด
	allorder, err := srv.returnOrderRepo.AllGetReturnOrder(ctx)
	if err != nil {
		srv.logger.Error("Error fetching all return orders", zap.Error(err))
		// Step 2: หากเกิดข้อผิดพลาด ให้ส่ง Error กลับไปยัง API
		return nil, errors.UnexpectedError()
	}

	// Step 3: ส่งข้อมูล ReturnOrder ที่ได้กลับไปยัง API
	return allorder, nil
}

func (srv service) GetReturnOrderByID(ctx context.Context, returnID string) (*response.ReturnOrder, error) {
	// Step 1: ตรวจสอบว่า ReturnID ไม่เป็นค่าว่าง
	if returnID == "" {
		return nil, errors.ValidationError("ReturnID is required")
	}

	// Step 2: เรียก repository เพื่อดึงข้อมูล ReturnOrder โดยใช้ ReturnID
	idorder, err := srv.returnOrderRepo.GetReturnOrderByID(ctx, returnID)
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

func (srv service) GetReturnOrderLinesByReturnID(ctx context.Context, returnID string) ([]response.ReturnOrderLine, error) {
	if returnID == "" {
		return nil, errors.ValidationError("ReturnID is required")
	}

	lines, err := srv.returnOrderRepo.GetReturnOrderLinesByReturnID(ctx, returnID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFoundError("This Return Order Line not found")
		}
		srv.logger.Error("Error fetching return order lines by ReturnID", zap.Error(err))
		return nil, errors.UnexpectedError()
	}

	return lines, nil
}

func (srv service) CreateReturnOrder(ctx context.Context, req request.CreateReturnOrder) error {
	// Step 1: ตรวจสอบว่าฟิลด์ที่จำเป็นต้องไม่เป็นค่าว่าง
	if req.ReturnID == "" || req.OrderNo == "" {
		return errors.ValidationError("ReturnID or OrderNo are required")
	}

	// Step 2: เรียก repository เพื่อสร้าง ReturnOrder
	err := srv.returnOrderRepo.CreateReturnOrder(ctx, req)
	if err != nil {
		srv.logger.Error("Error creating return order", zap.Error(err))
		return errors.UnexpectedError()
	}

	// Step 3: หากสร้างสำเร็จ ให้ส่งข้อความยืนยันกลับไปยัง API
	return nil
}

func (srv service) UpdateReturnOrder(ctx context.Context, req request.UpdateReturnOrder) error {
	exists, err := srv.returnOrderRepo.CheckReturnIDExists(ctx, req.ReturnID)
	if err != nil {
        srv.logger.Error("Error checking ReturnID existence", zap.Error(err))
        return errors.UnexpectedError()
    }

    if !exists {
        return errors.NotFoundError("ReturnID not found")
    }

	// Step 1: ตรวจสอบว่า ReturnID ไม่เป็นค่าว่าง
	if req.ReturnID == "" {
		return errors.ValidationError("ReturnID is required")
	}

	// Step 2: เรียก repository เพื่ออัปเดต ReturnOrder
	err = srv.returnOrderRepo.UpdateReturnOrder(ctx, req)
	if err != nil {
		srv.logger.Error("Error updating ReturnOrder", zap.Error(err))
		return errors.UnexpectedError()
	}

	// Step 3: หากอัปเดตสำเร็จ ให้ส่งข้อความยืนยันกลับไปยัง API
	return nil
}

func (srv service) DeleteReturnOrder(ctx context.Context, returnID string) error {
	if returnID == "" {
		return errors.ValidationError("ReturnID is required")
	}

	// Step 2: เรียก repository เพื่อลบ ReturnOrder
	err := srv.returnOrderRepo.DeleteReturnOrder(ctx, returnID)
	if err != nil {
		srv.logger.Error("Error deleting ReturnOrder", zap.Error(err))
		return errors.UnexpectedError()
	}

	return nil
}
