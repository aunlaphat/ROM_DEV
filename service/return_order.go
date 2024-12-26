package service

// ใช้จัดการคำสั่งซื้อที่มีเข้ามา

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	//"boilerplate-backend-go/errors"
	"database/sql"
	"fmt"
	"log"
	// "os"
	// "time"
	//"net/http"
)

type ReturnOrderService interface { // ตัวสื่อกลางในการรับส่งกับ api, ประมวลผลข้อมูลที่รับมาจาก api,
	AllGetReturnOrder() ([]response.ReturnOrder, error)
	GetReturnOrderByID(returnID string) (*response.ReturnOrder, error)
	CreateReturnOrder(req request.CreateReturnOrder) error
	UpdateReturnOrder(req request.UpdateReturnOrder) error
	DeleteReturnOrder(returnID string) error
}

// service เชื่อมกับ repo ต่อเพื่อดึงข้อมูลออกมา แต่ต้องมีการ validation ก่อนดึง
func (srv service) AllGetReturnOrder() ([]response.ReturnOrder, error) {
	allorder, err := srv.returnOrderRepo.AllGetReturnOrder()
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

func (srv service) GetReturnOrderByID(returnID string) (*response.ReturnOrder, error) {
	idorder, err := srv.returnOrderRepo.GetReturnOrderByID(returnID)
	if err != nil {
		return nil, err
	}
	return idorder, nil
}

func (srv service) CreateReturnOrder(req request.CreateReturnOrder) error {

    // Validate input
    if req.ReturnID == "" || req.OrderNo == "" {
        return fmt.Errorf("ReturnID and OrderNo are required")
    }

    // Call the repository to create the return order
    err := srv.returnOrderRepo.CreateReturnOrder(req)
    if err != nil {
        log.Printf("Error: failed to create return order: %v\n", err)
        return fmt.Errorf("could not create return order: %w", err)
    }

    return nil
}

func (srv service) UpdateReturnOrder(req request.UpdateReturnOrder) error {
    if req.ReturnID == "" {
        return fmt.Errorf("ReturnID is required")
    }

    // ส่งไปยัง Repository Layer
    err := srv.returnOrderRepo.UpdateReturnOrder(req)
    if err != nil {
        return fmt.Errorf("failed to update return order: %w", err)
    }

    return nil
}

func (srv service) DeleteReturnOrder(returnID string) error {
	if returnID == "" {
		return fmt.Errorf("ReturnID is required")
	}

	// ส่งไปยัง Repository Layer
	err := srv.returnOrderRepo.DeleteReturnOrder(returnID)
	if err != nil {
		return fmt.Errorf("failed to delete return order: %w", err)
	}

	return nil
}