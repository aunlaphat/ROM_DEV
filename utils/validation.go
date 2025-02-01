package utils

import (
	"errors"
	"fmt"
	"strings"

	req "boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
)

// IsStatusCanceled ตรวจสอบว่าสถานะเป็น "ยกเลิก" หรือ "ยืนยันแล้ว"
func IsStatusCanceled(statusConfID, statusReturnID *int) bool {
	if (statusConfID != nil && *statusConfID == 3) || (statusReturnID != nil && *statusReturnID == 2) {
		return true
	}
	return false
}

// ตรวจสอบว่าค่า string ไม่เป็นค่าว่าง
func validateRequiredString(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", field)
	}
	return nil
}

// ตรวจสอบว่า int มากกว่า 0
func validatePositiveInt(field string, value int) error {
	if value <= 0 {
		return fmt.Errorf("invalid %s", field)
	}
	return nil
}

func ValidateOrderStatus(order *response.BeforeReturnOrderResponse, expectedStatusReturnID, expectedStatusConfID int) error {
	if order.StatusReturnID != nil && *order.StatusReturnID != expectedStatusReturnID {
		return fmt.Errorf("order is not in the expected return status")
	}
	if order.StatusConfID != nil && *order.StatusConfID != expectedStatusConfID {
		return fmt.Errorf("order is not in the expected confirm status")
	}
	return nil
}

func ValidateCreateSaleReturn(req req.BeforeReturnOrder) error {
	if err := validateRequiredString("order number", req.OrderNo); err != nil {
		return err
	}
	if err := validateRequiredString("SO number", req.SoNo); err != nil {
		return err
	}
	if err := validateRequiredString("customer ID", req.CustomerID); err != nil {
		return err
	}
	if err := validatePositiveInt("channel ID", req.ChannelID); err != nil {
		return err
	}
	if err := validatePositiveInt("warehouse ID", req.WarehouseID); err != nil {
		return err
	}
	if len(req.BeforeReturnOrderLines) == 0 {
		return errors.New("at least one order line is required")
	}

	for i, line := range req.BeforeReturnOrderLines {
		if err := validateRequiredString(fmt.Sprintf("SKU for line %d", i+1), line.SKU); err != nil {
			return err
		}
		if err := validatePositiveInt(fmt.Sprintf("quantity for line %d", i+1), line.QTY); err != nil {
			return err
		}
		if line.ReturnQTY < 0 {
			return fmt.Errorf("return quantity cannot be negative for line %d", i+1)
		}
		if line.ReturnQTY > line.QTY {
			return fmt.Errorf("return quantity cannot be greater than quantity for line %d", i+1)
		}
		if line.Price < 0 {
			return fmt.Errorf("price cannot be negative for line %d", i+1)
		}
		if line.AlterSKU != nil && strings.TrimSpace(*line.AlterSKU) == "" {
			return fmt.Errorf("alter SKU cannot be empty for line %d", i+1)
		}
	}
	return nil
}

func ValidateUpdateSaleReturn(req req.UpdateSaleReturn) error {
	if err := validateRequiredString("order number", req.OrderNo); err != nil {
		return err
	}
	if err := validateRequiredString("SR number", req.SrNo); err != nil {
		return err
	}
	return nil
}

// ของ fa
func ValidateCreateReturnOrder(req req.CreateReturnOrder) error {
	// 1. ตรวจสอบข้อมูลพื้นฐาน
	if req.OrderNo == "" {
		return fmt.Errorf("order number is required")
	}
	if req.SoNo == "" {
		return fmt.Errorf("SO number is required")
	}

	// 2. ตรวจสอบค่าที่ต้องมากกว่า 0
	if *req.ChannelID <= 0 {
		return fmt.Errorf("invalid channel ID")
	}

	// 4. ตรวจสอบ order lines
	if len(req.ReturnOrderLine) == 0 {
		return fmt.Errorf("at least one order line is required")
	}

	for i, line := range req.ReturnOrderLine {
		if line.SKU == "" {
			return fmt.Errorf("SKU is required for line %d", i+1)
		}
		if *line.QTY <= 0 {
			return fmt.Errorf("quantity must be greater than 0 for line %d", i+1)
		}
		if line.ReturnQTY < 0 {
			return fmt.Errorf("return quantity cannot be negative for line %d", i+1)
		}
		if line.ReturnQTY > *line.QTY {
			return fmt.Errorf("return quantity cannot be greater than quantity for line %d", i+1)
		}
		if line.Price < 0 {
			return fmt.Errorf("price cannot be negative for line %d", i+1)
		}
	}

	return nil
}
