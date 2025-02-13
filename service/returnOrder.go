package service

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"boilerplate-backend-go/utils"
	"context"
	"fmt"

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
}

func (srv service) GetAllReturnOrder(ctx context.Context) ([]response.ReturnOrder, error) {
	srv.logger.Info("[ Starting get return order process ]")

	allorder, err := srv.returnOrderRepo.GetAllReturnOrder(ctx)
	if err != nil {
		srv.logger.Error("[ Error fetching all return orders ]", zap.Error(err))
		return nil, fmt.Errorf("[ error fetching all return orders: %w ]", err)
	}

	// *️⃣ เช็คเมื่อไม่มีข้อมูลออเดอร์
	if len(allorder) == 0 {
		srv.logger.Info("[ No return orders found ]")
		return []response.ReturnOrder{}, nil
	}

	srv.logger.Info("[ Fetched all return orders ]", zap.Int("Total amount of data", len(allorder)))
	return allorder, nil
}

func (srv service) GetReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.ReturnOrder, error) {
	srv.logger.Info("[ Starting get return order process ]", zap.String("OrderNo", orderNo))

	// *️⃣ ตรวจสอบว่า OrderNo ไม่เป็นค่าว่าง
	if orderNo == "" {
		err := fmt.Errorf("[ orderNo is required ]")
		srv.logger.Error("[ Invalid input ]", zap.Error(err))
		return nil, err
	}

	// *️⃣ ตรวจสอบว่า OrderNo มีอยู่จริงใน ReturnOrderLine
	exists, err := srv.returnOrderRepo.CheckOrderNoExist(ctx, orderNo)
	if err != nil {
		srv.logger.Error("[ Error checking OrderNo existence ]", zap.Error(err))
		return nil, fmt.Errorf("[ error checking OrderNo existence: %w ]", err)
	}
	if !exists {
		err := fmt.Errorf("[ This OrderNo not found: %s ]", orderNo)
		srv.logger.Warn("[ OrderNo not found ]", zap.String("OrderNo", orderNo))
		return nil, err
	}

	idorder, err := srv.returnOrderRepo.GetReturnOrderByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("[ Error fetching ReturnOrder by OrderNo ]", zap.String("OrderNo", orderNo), zap.Error(err))
		return nil, fmt.Errorf("[ error fetching ReturnOrder by OrderNo %s: %w ]", orderNo, err)
	}

	// *️⃣ เช็คเมื่อไม่มีข้อมูลรายการส่งคืนในคำสั่งซื้อ
	if len(idorder.ReturnOrderLine) == 0 {
		srv.logger.Info("[ No lines found for this order ]")
		return idorder, nil
	}

	srv.logger.Info("[ Fetched return order by orderNo ]", zap.String("OrderNo", orderNo))
	return idorder, nil
}

func (srv service) GetAllReturnOrderLines(ctx context.Context) ([]response.ReturnOrderLine, error) {
	srv.logger.Info("[ Starting get return order lines process ]")

	lines, err := srv.returnOrderRepo.GetAllReturnOrderLines(ctx)
	if err != nil {
		srv.logger.Error("[ Error fetching all return order lines ]", zap.Error(err))
		return nil, fmt.Errorf("[ error fetching all return order lines: %w ]", err)
	}

	// *️⃣ เช็คเมื่อไม่มีข้อมูลรายการสั่งคืนในคำสั่งซื้อ
	if len(lines) == 0 {
		srv.logger.Info("[ No lines found ]")
		return lines, nil
	}

	srv.logger.Info("[ Fetched all return order lines ]", zap.Int("Total amount of data", len(lines)))
	return lines, nil
}

func (srv service) GetReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.ReturnOrderLine, error) {
	srv.logger.Info("[ Starting get return order line process ]", zap.String("OrderNo", orderNo))

	if orderNo == "" {
		err := fmt.Errorf("[ orderNo is required ]")
		srv.logger.Error("[ Invalid input ]", zap.Error(err))
		return nil, err
	}

	// *️⃣ ตรวจสอบว่า OrderNo มีอยู่ใน ReturnOrderLine
	exists, err := srv.returnOrderRepo.CheckOrderNoLineExist(ctx, orderNo)
	if err != nil {
		srv.logger.Error("[ Error checking OrderNo existence ]", zap.Error(err))
		return nil, fmt.Errorf("[ error checking OrderNo existence: %w ]", err)
	}
	if !exists {
		err := fmt.Errorf("[ This Return Order Line not found: %s ]", orderNo)
		srv.logger.Warn(err.Error())
		return nil, err
	}

	// *️⃣ ดึงข้อมูล ReturnOrderLines
	lines, err := srv.returnOrderRepo.GetReturnOrderLineByOrderNo(ctx, orderNo)
	if err != nil {
		srv.logger.Error("[ Error fetching return order lines by OrderNo ]", zap.Error(err))
		return nil, fmt.Errorf("[ error fetching return order lines by OrderNo: %w ]", err)
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
		return nil, errors.InternalError("[ failed to fetch Return Orders ]")
	}

	if len(orders) == 0 {
		srv.logger.Info("[ No order found ]")
		return []response.DraftTradeDetail{}, nil
	}

	srv.logger.Info("[ Successfully fetched Return Orders ]", zap.Int("StatusCheckID", statusCheckID), zap.Int("Total amount of data", len(orders)))
	return orders, nil
}

func (srv service) GetReturnOrdersByStatusAndDateRange(ctx context.Context, statusCheckID int, startDate, endDate string) ([]response.DraftTradeDetail, error) {
	srv.logger.Info("[ Starting get return order from range date process ]", zap.String("StartDate", startDate), zap.String("EndDate", endDate))

	orders, err := srv.returnOrderRepo.GetReturnOrdersByStatusAndDateRange(ctx, statusCheckID, startDate, endDate)
	if err != nil {
		srv.logger.Error("[ Failed to fetch Return Orders ]", zap.Int("StatusCheckID", statusCheckID), zap.Error(err))
		return nil, fmt.Errorf("[ failed to fetch Return Orders: %w ]", err)
	}

	// *️⃣ เช็คเมื่อไม่มีข้อมูลรายการสั่งคืนในคำสั่งซื้อ
	if len(orders) == 0 {
		srv.logger.Info("[ No order found within the specified date range ]")
		return []response.DraftTradeDetail{}, nil
	}

	srv.logger.Info("[ Successfully fetched Return Orders ]", zap.String("StartDate", startDate), zap.String("EndDate", endDate), zap.Int("Total amount of data", len(orders)))
	return orders, nil
}

func (srv service) CreateReturnOrder(ctx context.Context, req request.CreateReturnOrder) (*response.CreateReturnOrder, error) {
	srv.logger.Info("[ Starting return order creation process ]", zap.String("OrderNo", req.OrderNo), zap.String("CreateBy", req.CreateBy))

	// *️⃣ ตรวจสอบว่า ReturnOrderLine ต้องไม่เป็นค่าว่าง (มีอย่างน้อย 1 รายการ) จึงจะสามารถสร้างการคืนของได้
	if len(req.ReturnOrderLine) == 0 {
		err := fmt.Errorf("[ ReturnOrderLine cannot be empty ]")
		srv.logger.Error("[ Invalid input ]", zap.Error(err))
		return nil, err
	}

	// *️⃣ Validate request ที่ส่งมา
	if err := utils.ValidateCreateReturnOrder(req); err != nil {
		srv.logger.Error("[ Invalid return order request ]", zap.Error(err))
		return nil, fmt.Errorf("[ Validation failed: %w ]", err)
	}

	// *️⃣ ตรวจสอบว่า OrderNo ซ้ำหรือไม่
	exists, err := srv.returnOrderRepo.CheckOrderNoExist(ctx, req.OrderNo)
	if err != nil {

		srv.logger.Error("[ Failed to check OrderNo ]", zap.Error(err))
		return nil, fmt.Errorf("[ failed to check OrderNo: %w ]", err)
	}
	if exists {
		srv.logger.Error("[ OrderNo already exists ]", zap.Error(err))
		return nil, (fmt.Errorf("[ orderNo already exists: %s ]", req.OrderNo))
	}

	err = srv.returnOrderRepo.CreateReturnOrder(ctx, req)
	if err != nil {
		srv.logger.Error("[ Failed to create order with lines ]", zap.Error(err))
		return nil, fmt.Errorf("[ failed to create order with lines: %w ]", err)
	}

	// *️⃣ ดึงข้อมูล order ที่สร้างสำเร็จไปแสดง
	createdOrder, err := srv.returnOrderRepo.GetCreateReturnOrder(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("[ Failed to fetch created order ]", zap.Error(err))
		return nil, fmt.Errorf("[ failed to fetch created order: %w ]", err)
	}

	srv.logger.Info("[ Return order created successfully ]", zap.String("OrderNo", req.OrderNo), zap.String("CreateBy", req.CreateBy))
	return createdOrder, nil
}

func (srv service) UpdateReturnOrder(ctx context.Context, req request.UpdateReturnOrder) (*response.UpdateReturnOrder, error) {
	srv.logger.Info("[ Starting return order update process ]", zap.String("OrderNo", req.OrderNo), zap.String("UpdateBy", *req.UpdateBy))

	if req.OrderNo == "" {
		err := fmt.Errorf("[ OrderNo is required ]")
		srv.logger.Error("[ Invalid OrderNo ]", zap.Error(err))
		return nil, err
	}

	exists, err := srv.returnOrderRepo.CheckOrderNoExist(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("[ Error checking OrderNo existence ]", zap.Error(err))
		return nil, fmt.Errorf("[ error checking OrderNo existence: %w ]", err)
	}
	if !exists {
		srv.logger.Warn("[ OrderNo not found ]", zap.String("OrderNo", req.OrderNo))
		return nil, fmt.Errorf("[ OrderNo not found: %s ]", req.OrderNo)
	}

	err = srv.returnOrderRepo.UpdateReturnOrder(ctx, req)
	if err != nil {
		srv.logger.Error("[ Error updating ReturnOrder ]", zap.Error(err))
		return nil, fmt.Errorf("[ error updating ReturnOrder: %w ]", err)
	}

	// *️⃣ ดึงข้อมูล order ที่อัพเดทเสร็จแล้ว
	updatedOrder, err := srv.returnOrderRepo.GetUpdateReturnOrder(ctx, req.OrderNo)
	if err != nil {
		srv.logger.Error("[ Failed to fetch updated order ]", zap.Error(err))
		return nil, fmt.Errorf("[ failed to fetch updated order: %w ]", err)
	}

	srv.logger.Info("[ Return order updated successfully ]", zap.String("OrderNo", req.OrderNo), zap.String("UpdateBy", *req.UpdateBy))
	return updatedOrder, nil
}

func (srv service) DeleteReturnOrder(ctx context.Context, orderNo string) error {
	srv.logger.Info("[ Starting delete return order process ]", zap.String("OrderNo", orderNo))

	if orderNo == "" {
		err := fmt.Errorf("[ OrderNo are required ]")
		srv.logger.Error("[ Invalid OrderNo ]", zap.Error(err))
		return err
	}

	exists, err := srv.returnOrderRepo.CheckOrderNoExist(ctx, orderNo)
	if err != nil {
		srv.logger.Error("[ Error checking OrderNo existence ]", zap.Error(err))
		return fmt.Errorf("[ error checking OrderNo existence: %w ]", err)
	}
	if !exists {
		srv.logger.Warn("[ OrderNo not found ]", zap.String("OrderNo", orderNo))
		return fmt.Errorf("[ OrderNo not found: %s ]", orderNo)

	}

	err = srv.returnOrderRepo.DeleteReturnOrder(ctx, orderNo)
	if err != nil {
		srv.logger.Error("[ Error deleting ReturnOrder ]", zap.Error(err))
		return fmt.Errorf("[ error deleting ReturnOrder: %w ]", err)
	}

	srv.logger.Info("[ Return order deleted successfully ]", zap.String("OrderNo", orderNo))
	return nil
}
