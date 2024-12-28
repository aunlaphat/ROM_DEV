package service

// ใช้จัดการคำสั่งซื้อที่มีเข้ามา

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"database/sql"
	"fmt"
	"net/http"
)

type OrderService interface { // ตัวสื่อกลางในการรับส่งกับ api, ประมวลผลข้อมูลที่รับมาจาก api,
	GetOrderID(orderNo string) (response.OrderResponse, error)
	AllGetOrder() ([]response.OrderResponse, error)
	CreateOrder(req request.CreateOrderRequest) (response.OrderResponse, error)
	UpdateOrder(req request.UpdateOrderRequest) error
	DeleteOrder(orderNo string) error
}

// service เชื่อมกับ repo ต่อเพื่อดึงข้อมูลออกมา แต่ต้องมีการ validation ก่อนดึง

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

func (srv service) GetOrderID(orderNo string) (response.OrderResponse, error) {
	// ตรวจสอบ Input
	res := response.OrderResponse{}
	if orderNo == "" {
		return res, errors.ValidationError("order number is required")
	}

	// เรียก repository เพื่อดึงข้อมูล
	order, err := srv.orderRepo.GetOrderID(orderNo)
	if err != nil {
		switch e := err.(type) {
		case errors.AppError:
			if e.Code == http.StatusNotFound {
				return res, errors.NotFoundError(fmt.Sprintf("OrderNo '%s'", orderNo))
			}
			return res, errors.UnexpectedError()
		default:
			srv.logger.Error(err)
			return res, errors.UnexpectedError()
		}
	}

	// ตรวจสอบว่ามี OrderLines หรือไม่
	if len(order.OrderLines) == 0 {
		// บันทึกข้อมูลใน Log โดยไม่ใช้ Warnf
		srv.logger.Info(fmt.Sprintf("OrderNo '%s' has no order lines", orderNo))
	}

	return order, nil
}

func (srv service) CreateOrder(req request.CreateOrderRequest) (response.OrderResponse, error) {
	if req.OrderNo == "" {
		return response.OrderResponse{}, errors.ValidationError("order number is required")
	}

	// สร้างคำสั่งซื้อในฐานข้อมูล
	err := srv.orderRepo.CreateOrder(req)
	if err != nil {
		switch err.(type) {
		case errors.AppError:
			appErr := err.(errors.AppError)
			if appErr.Code == http.StatusConflict {
				return response.OrderResponse{}, errors.ValidationError(fmt.Sprintf("OrderNo '%s' already exists. Use a different OrderNo.", req.OrderNo))
			}
			return response.OrderResponse{}, errors.UnexpectedError()
		default:
			srv.logger.Error(err)
			return response.OrderResponse{}, errors.UnexpectedError()
		}
	}

	// ดึงข้อมูลคำสั่งซื้อที่เพิ่งสร้าง
	order, err := srv.GetOrderID(req.OrderNo)
	if err != nil {
		switch err.(type) {
		case errors.AppError:
			return response.OrderResponse{}, err
		default:
			srv.logger.Error(err)
			return response.OrderResponse{}, errors.UnexpectedError()
		}
	}

	return order, nil
}

func (srv service) UpdateOrder(req request.UpdateOrderRequest) error {
	if req.OrderNo == "" {
		return errors.ValidationError("order number is required")
	}

	// ตรวจสอบว่ามี OrderNo อยู่ในฐานข้อมูลหรือไม่
	orderExists, err := srv.orderRepo.CheckOrderExists(req.OrderNo)
	if err != nil {
		switch err.(type) {
		case errors.AppError:
			return err
		default:
			srv.logger.Error(err)
			return errors.UnexpectedError()
		}
	}

	if !orderExists {
		return errors.NotFoundError(fmt.Sprintf("OrderNo '%s'", req.OrderNo))
	}

	// ดำเนินการอัปเดตข้อมูล
	err = srv.orderRepo.UpdateOrder(req)
	if err != nil {
		switch err.(type) {
		case errors.AppError:
			return err
		default:
			srv.logger.Error(err)
			return errors.UnexpectedError()
		}
	}

	return nil
}

func (srv service) DeleteOrder(orderNo string) error {
	if orderNo == "" {
		return errors.ValidationError("order number is required")
	}

	// ตรวจสอบว่ามี OrderNo อยู่ในฐานข้อมูลหรือไม่
	orderExists, err := srv.orderRepo.CheckOrderExists(orderNo)
	if err != nil {
		switch err.(type) {
		case errors.AppError:
			return err
		default:
			srv.logger.Error(err)
			return errors.UnexpectedError()
		}
	}

	if !orderExists {
		return errors.NotFoundError(fmt.Sprintf("OrderNo '%s'", orderNo))
	}

	// ดำเนินการลบข้อมูล
	err = srv.orderRepo.DeleteOrder(orderNo)
	if err != nil {
		switch err.(type) {
		case errors.AppError:
			return err
		default:
			srv.logger.Error(err)
			return errors.UnexpectedError()
		}
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
