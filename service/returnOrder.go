package service

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"boilerplate-backend-go/utils"
	"context"

	"go.uber.org/zap"
)

type ReturnOrderService interface {
	GetAllReturnOrder(ctx context.Context) ([]response.ReturnOrder, error)
	GetReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.ReturnOrder, error)
	GetAllReturnOrderLines(ctx context.Context) ([]response.ReturnOrderLine, error)
	GetReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.ReturnOrderLine, error)
	CreateReturnOrder(ctx context.Context, req request.CreateReturnOrder) (*response.CreateReturnOrder, error)
	UpdateReturnOrder(ctx context.Context, req request.UpdateReturnOrder) (*response.UpdateReturnOrder, error)
	DeleteReturnOrder(ctx context.Context, orderNo string) error

	GetReturnOrdersByStatus(ctx context.Context, statusCheckID int) ([]response.DraftTradeDetail, error)
	GetReturnOrdersByStatusAndDateRange(ctx context.Context, statusCheckID int, startDate, endDate string) ([]response.DraftTradeDetail, error)

	CheckOrderNoExist(ctx context.Context, orderNo string) error
}

// review all
func (srv service) GetAllReturnOrder(ctx context.Context) ([]response.ReturnOrder, error) {
	srv.logger.Info("[ Starting get return order process ]")

	allorder, err := srv.returnOrderRepo.GetAllReturnOrder(ctx)
	if err != nil {
		srv.logger.Error("[ Error fetching all return orders ]", zap.Error(err))
		return nil, errors.InternalError("[ Failed to fetch all return orders: %v ]", err)
	}

	srv.logger.Info("[ Fetched all return orders ]", zap.Int("Total amount of data", len(allorder)))
	return allorder, nil
}

func (srv service) GetReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.ReturnOrder, error) {
	srv.logger.Info("[ Starting get return order process ]", zap.String("OrderNo", orderNo))

	// *️⃣ ตรวจสอบ OrderNo
	err := srv.CheckOrderNoExist(ctx, orderNo)
	if err != nil {
		return nil, err
	}

	idorder, err := srv.returnOrderRepo.GetReturnOrderByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("[ Error fetching ReturnOrder by OrderNo ]", zap.String("OrderNo", orderNo), zap.Error(err))
		return nil, errors.InternalError("[ Error fetching ReturnOrder by OrderNo %s: %v ]", orderNo, err)
	}

	srv.logger.Info("[ Fetched return order by orderNo ]", zap.String("OrderNo", orderNo))
	return idorder, nil
}

func (srv service) GetAllReturnOrderLines(ctx context.Context) ([]response.ReturnOrderLine, error) {
	srv.logger.Info("[ Starting get return order lines process ]")

	lines, err := srv.returnOrderRepo.GetAllReturnOrderLines(ctx)
	if err != nil {
		srv.logger.Error("[ Error fetching all return order lines ]", zap.Error(err))
		return nil, errors.InternalError("[ Error fetching all return order lines: %v ]", err)
	}

	srv.logger.Info("[ Fetched all return order lines ]", zap.Int("Total amount of data", len(lines)))
	return lines, nil
}

func (srv service) GetReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.ReturnOrderLine, error) {
	srv.logger.Info("[ Starting get return order line process ]", zap.String("OrderNo", orderNo))

	// *️⃣ ตรวจสอบว่า OrderNo มีอยู่ใน ReturnOrderLine
	err := srv.CheckOrderNoExist(ctx, orderNo)
	if err != nil {
		return nil, err
	}

	// *️⃣ ดึงข้อมูล ReturnOrderLines
	lines, err := srv.returnOrderRepo.GetReturnOrderLineByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("[ Error fetching return order lines by OrderNo ]", zap.Error(err))
		return nil, errors.InternalError("[ Error fetching return order lines by OrderNo %s: %v ]", orderNo, err)
	}

	// *️⃣ เช็คเมื่อไม่มีข้อมูลรายการสั่งคืนในคำสั่งซื้อ
	if len(lines) == 0 {
		srv.logger.Info("[ No lines found for this order number ]")
		return lines, nil
	}

	srv.logger.Info("[ Fetched return order line by orderNo ]", zap.String("OrderNo", orderNo))
	return lines, nil
}

func (srv service) GetReturnOrdersByStatus(ctx context.Context, statusCheckID int) ([]response.DraftTradeDetail, error) {
	srv.logger.Info("[ Starting get return order process ]", zap.Int("StatusCheckID", statusCheckID))

	orders, err := srv.returnOrderRepo.GetReturnOrdersByStatus(ctx, statusCheckID)
	if err != nil {
		srv.logger.Error("[ Failed to fetch Return Orders ]", zap.Error(err))
		return nil, errors.InternalError("[ Failed to fetch Return Orders %v ]", err)
	}

	if len(orders) == 0 {
		srv.logger.Info("[ No order found ]")
		return orders, nil
	}

	srv.logger.Info("[ Successfully fetched Return Orders ]", zap.Int("StatusCheckID", statusCheckID), zap.Int("Total amount of data", len(orders)))
	return orders, nil
}

func (srv service) GetReturnOrdersByStatusAndDateRange(ctx context.Context, statusCheckID int, startDate, endDate string) ([]response.DraftTradeDetail, error) {
	srv.logger.Info("[ Starting get return order from range date process ]", zap.String("StartDate", startDate), zap.String("EndDate", endDate))

	orders, err := srv.returnOrderRepo.GetReturnOrdersByStatusAndDateRange(ctx, statusCheckID, startDate, endDate)
	if err != nil {
		srv.logger.Error("[ Failed to fetch Return Orders ]", zap.Int("StatusCheckID", statusCheckID), zap.Error(err))
		return nil, errors.InternalError("[ Failed to fetch Return Orders: %v ]", err)
	}

	// *️⃣ เช็คเมื่อไม่มีข้อมูลรายการสั่งคืนในคำสั่งซื้อ
	if len(orders) == 0 {
		srv.logger.Info("[ No order found within the specified date range ]")
		return orders, nil
	}

	srv.logger.Info("[ Successfully fetched Return Orders ]", zap.String("StartDate", startDate), zap.String("EndDate", endDate), zap.Int("Total amount of data", len(orders)))
	return orders, nil
}

func (srv service) CreateReturnOrder(ctx context.Context, req request.CreateReturnOrder) (*response.CreateReturnOrder, error) {
	srv.logger.Info("[ Starting return order creation process ]", zap.String("OrderNo", req.OrderNo), zap.String("CreateBy", req.CreateBy))

	// *️⃣ ตรวจสอบว่า ReturnOrderLine ต้องไม่เป็นค่าว่าง (มีอย่างน้อย 1 รายการ) จึงจะสามารถสร้างการคืนของได้
	if len(req.ReturnOrderLine) == 0 {
		srv.logger.Warn("[ ReturnOrderLine can't empty must be > 0 line ]")
		return nil, errors.ValidationError("[ ReturnOrderLine can't empty must be > 0 line ]")
	}

	// *️⃣ Validate request ที่ส่งมา
	if err := utils.ValidateCreateReturnOrder(req); err != nil {
		srv.logger.Warn("[ Validation failed ]", zap.Error(err))
		return nil, errors.ValidationError("[ Validation failed: %v ]", err)
	}

	// *️⃣ ตรวจสอบว่า OrderNo สร้างซ้ำหรือไม่
	exists, err := srv.returnOrderRepo.CheckOrderNoExist(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("[ Error checking OrderNo existence ]", zap.Error(err)) // db มีปัญหา
		return nil, errors.InternalError("[ Error checking OrderNo existence: %v ]", err)
	}
	if exists {
		srv.logger.Warn("[ OrderNo already exists ]", zap.String("OrderNo", req.OrderNo))
		return nil, errors.ConflictError("[ OrderNo %s already exists ]", req.OrderNo)
	}

	err = srv.returnOrderRepo.CreateReturnOrder(ctx, req)
	if err != nil {
		srv.logger.Error("[ Failed to create order with lines ]", zap.Error(err))
		return nil, errors.InternalError("[ Failed to create order with lines: %v ]", err)
	}

	// *️⃣ ดึงข้อมูล order ที่สร้างสำเร็จไปแสดง
	createdOrder, err := srv.returnOrderRepo.GetCreateReturnOrder(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("[ Failed to fetch created order ]", zap.Error(err))
		return nil, errors.InternalError("[ Failed to fetch created order: %v ]", err)
	}

	srv.logger.Info("[ Return order created successfully ]", zap.String("OrderNo", req.OrderNo), zap.String("CreateBy", req.CreateBy))
	return createdOrder, nil
}

func (srv service) UpdateReturnOrder(ctx context.Context, req request.UpdateReturnOrder) (*response.UpdateReturnOrder, error) {
	srv.logger.Info("[ Starting return order update process ]", zap.String("OrderNo", req.OrderNo), zap.String("UpdateBy", *req.UpdateBy))

	// *️⃣ ตรวจสอบ OrderNo
	err := srv.CheckOrderNoExist(ctx, req.OrderNo)
	if err != nil {
		return nil, err
	}

	err = srv.returnOrderRepo.UpdateReturnOrder(ctx, req)
	if err != nil {
		srv.logger.Error("[ Failed to update order with lines ]", zap.Error(err))
		return nil, errors.InternalError("[ Failed to update order with lines: %v ]", err)
	}

	// *️⃣ ดึงข้อมูล order ที่อัพเดทเสร็จแล้ว
	updatedOrder, err := srv.returnOrderRepo.GetUpdateReturnOrder(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("[ Failed to fetch updated order ]", zap.Error(err))
		return nil, errors.InternalError("[ Failed to fetch updated order: %v ]", err)
	}

	srv.logger.Info("[ Return order updated successfully ]", zap.String("OrderNo", req.OrderNo), zap.String("UpdateBy", *req.UpdateBy))
	return updatedOrder, nil
}

func (srv service) DeleteReturnOrder(ctx context.Context, orderNo string) error {
	srv.logger.Info("[ Starting delete return order process ]", zap.String("OrderNo", orderNo))

	// *️⃣ ตรวจสอบ OrderNo
	err := srv.CheckOrderNoExist(ctx, orderNo)
	if err != nil {
		return err
	}

	err = srv.returnOrderRepo.DeleteReturnOrder(ctx, orderNo)
	if err != nil {
		srv.logger.Error("[ Failed to delete order with lines ]", zap.Error(err))
		return errors.InternalError("Failed to delete order with lines: %v", err)
	}

	srv.logger.Info("[ Return order deleted successfully ]", zap.String("OrderNo", orderNo))
	return nil
}

// *️⃣ เช็คถ้า OrderNo มีในระบบ ทำงานต่อไปได้
func (srv service) CheckOrderNoExist(ctx context.Context, orderNo string) error {
	// *️⃣ ตรวจสอบ OrderNo ในฐานข้อมูล
	exists, err := srv.returnOrderRepo.CheckOrderNoExist(ctx, orderNo)
	if err != nil {
		srv.logger.Error("[ Error checking OrderNo existence ]", zap.Error(err)) // db มีปัญหา
		return errors.InternalError("[ Error checking OrderNo existence: %v ]", err)
	}

	// *️⃣ ตรวจเมื่อไม่พบ OrderNo
	if !exists {
		srv.logger.Warn("[ OrderNo not found ]", zap.String("OrderNo", orderNo))
		return errors.NotFoundError("[ This OrderNo not found: %s ]", orderNo)
	}

	return nil
}
