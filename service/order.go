package service

// ใช้จัดการคำสั่งซื้อที่มีเข้ามา

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"

	"database/sql"
	"fmt"
)

type OrderService interface {
	GetOrderID(orderNo string) (response.OrderResponse, error)
	AllGetOrder() ([]response.OrderResponse, error)
	CreateOrder(req request.CreateOrderRequest) (response.OrderResponse, error)
	UpdateOrder(req request.UpdateOrderRequest) error
	DeleteOrder(orderNo string) error
}

func (srv service) AllGetOrder() ([]response.OrderResponse, error) {
	allorder, err := srv.orderRepo.AllGetOrder()
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
	// ส่งออกข้อมูล
	return allorder, nil
}

func (srv service) CreateOrder(req request.CreateOrderRequest) (response.OrderResponse, error) {

	if req.OrderNo == "" {
		err := fmt.Errorf("order number is required")
		srv.logger.Error(err)
		return response.OrderResponse{}, err
	}

	// สร้างคำสั่งซื้อในฐานข้อมูล
	err := srv.orderRepo.CreateOrder(req)
	if err != nil {
		srv.logger.Error(err)
		return response.OrderResponse{}, fmt.Errorf("failed to create order: %w", err)
	}

	// ดึงข้อมูลคำสั่งซื้อที่เพิ่งสร้าง
	order, err := srv.GetOrderID(req.OrderNo)
	if err != nil {
		srv.logger.Error(err)
		return response.OrderResponse{}, fmt.Errorf("failed to fetch created order: %w", err)
	}

	return order, nil
}

func (srv service) GetOrderID(orderNo string) (response.OrderResponse, error) {
	orderID, err := srv.orderRepo.GetOrderID(orderNo)
	if err != nil {
		srv.logger.Error(err)
		if err.Error() == "sql: no rows in result set" {
			return response.OrderResponse{}, fmt.Errorf("order with OrderNo %s not found", orderNo)
		}
		return response.OrderResponse{}, fmt.Errorf("failed to fetch order: %w", err)
	}

	return orderID, nil
}

func (srv service) UpdateOrder(req request.UpdateOrderRequest) error {
	if req.OrderNo == "" {
		return fmt.Errorf("order number is required")
	}

	// ตรวจสอบว่ามี OrderNo อยู่ในฐานข้อมูลหรือไม่
	orderExists, err := srv.orderRepo.CheckOrderExists(req.OrderNo)
	if err != nil {
		srv.logger.Error(err)
		return fmt.Errorf("failed to check order existence: %w", err)
	}

	if !orderExists {
		return fmt.Errorf("no order data: %s", req.OrderNo)
	}

	// ดำเนินการอัปเดตข้อมูล
	err = srv.orderRepo.UpdateOrder(req)
	if err != nil {
		srv.logger.Error(err)
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

func (srv service) DeleteOrder(orderNo string) error {
	if orderNo == "" {
		return fmt.Errorf("order number is required")
	}

	// ตรวจสอบว่ามี OrderNo อยู่ในฐานข้อมูลหรือไม่
	orderExists, err := srv.orderRepo.CheckOrderExists(orderNo)
	if err != nil {
		srv.logger.Error(err)
		return fmt.Errorf("failed to check order existence: %w", err)
	}

	if !orderExists {
		return fmt.Errorf("no order data: %s", orderNo)
	}

	// ดำเนินการลบข้อมูล
	err = srv.orderRepo.DeleteOrder(orderNo)
	if err != nil {
		srv.logger.Error(err)
		return fmt.Errorf("failed to delete order: %w", err)
	}

	return nil
}


// func (srv service) GetOrder() ([]entity.Order, error) {
// 	// ดึงข้อมูลคำสั่งซื้อจาก repository
// 	orders, err := srv.orderRepo.GetOrder()
// 	if err != nil {
// 		switch err {
// 		case sql.ErrNoRows:
// 			srv.logger.Error(err)
// 			return nil, fmt.Errorf("no order data: %w", err)
// 		default:
// 			srv.logger.Error(err)
// 			return nil, fmt.Errorf("get order error: %w", err)
// 		}
// 	}

// 	/***	มีการประกาศไว้ใน repository แล้วในการเรียก GetOrderLinesByOrderNo รวมกับ GetOrder ประกาศใช้ไว้ในที่ใดที่นึงพอ	***/
// 	// ดึงข้อมูล OrderLine สำหรับแต่ละ OrderNo
// 	// for i, order := range orders {
// 	// 	// เรียกข้อมูล OrderLine ตาม OrderNo ของแต่ละคำสั่งซื้อ
// 	// 	orderLines, err := srv.orderRepo.GetOrderLinesByOrderNo(order.OrderNo)
// 	// 	if err != nil {
// 	// 		srv.logger.Error(err)
// 	// 		return nil, fmt.Errorf("get order lines error: %w", err)
// 	// 	}

// 	// 	// เก็บข้อมูล OrderLine ลงใน Order
// 	// 	orders[i].OrderLines = orderLines
// 	// }

// 	// เพิ่ม log เพื่อดูว่า getOrder ทำงานได้
// 	fmt.Println("Fetched Orders:", orders)

// 	// ส่งคืนข้อมูลทั้งหมด
// 	return orders, nil
// }
