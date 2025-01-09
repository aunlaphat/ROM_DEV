package service

// ใช้จัดการคำสั่งซื้อที่มีเข้ามา

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"context"
	"database/sql"

	"go.uber.org/zap"
)

type ReturnOrderService interface { // ตัวสื่อกลางในการรับส่งกับ API และประมวลผลข้อมูลที่รับมาจาก API
	AllGetReturnOrder(ctx context.Context) ([]response.ReturnOrder, error)
	GetReturnOrderByID(ctx context.Context, returnID string) (*response.ReturnOrder, error)
	GetAllReturnOrderLines(ctx context.Context) ([]response.ReturnOrderLine, error)
	GetReturnOrderLinesByReturnID(ctx context.Context, returnID string) ([]response.ReturnOrderLine, error)
	CreateReturnOrder(ctx context.Context, req request.CreateReturnOrder) error
	UpdateReturnOrder(ctx context.Context, req request.UpdateReturnOrder) error
	DeleteReturnOrder(ctx context.Context, returnID string) error
}

// AllGetReturnOrder - ดึงข้อมูล ReturnOrder ทั้งหมด
func (srv service) AllGetReturnOrder(ctx context.Context) ([]response.ReturnOrder, error) {
	allorder, err := srv.returnOrderRepo.AllGetReturnOrder(ctx)
	if err != nil {
		srv.logger.Error("Error fetching all return orders", zap.Error(err))
		return nil, errors.UnexpectedError()
	}

	return allorder, nil
}

// GetReturnOrderByID - ดึงข้อมูล ReturnOrder โดยใช้ ReturnID
func (srv service) GetReturnOrderByID(ctx context.Context, returnID string) (*response.ReturnOrder, error) {
	if returnID == "" {
		return nil, errors.ValidationError("ReturnID is required")
	}

	idorder, err := srv.returnOrderRepo.GetReturnOrderByID(ctx, returnID)
	if err != nil {
		if err == sql.ErrNoRows { 
			return nil, errors.NotFoundError("Return order not found")
		}
		srv.logger.Error("Error fetching ReturnOrder by ID", zap.Error(err))
		return nil, errors.UnexpectedError()
	}

	return idorder, nil
}

// GetAllReturnOrderLines - ดึงข้อมูล ReturnOrderLine ทั้งหมด
func (srv service) GetAllReturnOrderLines(ctx context.Context) ([]response.ReturnOrderLine, error) {
	lines, err := srv.returnOrderRepo.GetAllReturnOrderLines(ctx)
	if err != nil {
		srv.logger.Error("Error fetching all return order lines", zap.Error(err))
		return nil, errors.UnexpectedError()
	}

	return lines, nil
}

// GetReturnOrderLinesByReturnID - ดึงข้อมูล ReturnOrderLine โดยใช้ ReturnID
func (srv service) GetReturnOrderLinesByReturnID(ctx context.Context, returnID string) ([]response.ReturnOrderLine, error) {
	if returnID == "" {
		return nil, errors.ValidationError("ReturnID is required")
	}

	lines, err := srv.returnOrderRepo.GetReturnOrderLinesByReturnID(ctx, returnID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFoundError("No return order lines found for the given ReturnID")
		}
		srv.logger.Error("Error fetching return order lines by ReturnID", zap.Error(err))
		return nil, errors.UnexpectedError()
	}

	return lines, nil
}

// CreateReturnOrder - สร้าง ReturnOrder ใหม่
func (srv service) CreateReturnOrder(ctx context.Context, req request.CreateReturnOrder) error {
	if req.ReturnID == "" || req.OrderNo == "" {
		return errors.ValidationError("ReturnID or OrderNo are required")
	}

	err := srv.returnOrderRepo.CreateReturnOrder(ctx, req)
	if err != nil {
		srv.logger.Error("Error creating return order", zap.Error(err))
		return errors.UnexpectedError()
	}

	return nil
}

// UpdateReturnOrder - อัปเดตข้อมูล ReturnOrder
func (srv service) UpdateReturnOrder(ctx context.Context, req request.UpdateReturnOrder) error {
	if req.ReturnID == "" {
		return errors.ValidationError("ReturnID is required")
	}

	err := srv.returnOrderRepo.UpdateReturnOrder(ctx, req)
	if err != nil {
		srv.logger.Error("Error updating ReturnOrder", zap.Error(err))
		return errors.UnexpectedError()
	}

	return nil
}

// DeleteReturnOrder - ลบข้อมูล ReturnOrder
func (srv service) DeleteReturnOrder(ctx context.Context, returnID string) error {
	if returnID == "" {
		return errors.ValidationError("ReturnID is required")
	}

	err := srv.returnOrderRepo.DeleteReturnOrder(ctx, returnID)
	if err != nil {
		srv.logger.Error("Error deleting ReturnOrder", zap.Error(err))
		return errors.UnexpectedError()
	}

	return nil
}
