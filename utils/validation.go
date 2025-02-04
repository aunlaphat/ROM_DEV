package utils

import (
	"errors"
	"fmt"
	"strings"

	req "boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
)

// 🛠️ ตรวจสอบว่าสถานะเป็น "ยกเลิก" หรือไม่
func IsStatusCanceled(statusConfID, statusReturnID *int) bool {
	return (statusConfID != nil && *statusConfID == 3) || (statusReturnID != nil && *statusReturnID == 2)
}

// 🛠️ ตรวจสอบว่าค่า string ไม่เป็นค่าว่าง
func validateRequiredString(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", field)
	}
	return nil
}

// 🛠️ ตรวจสอบว่า int ต้องมากกว่า 0 (รองรับ nil)
func validatePositiveInt(field string, value *int) error {
	if value == nil || *value <= 0 {
		return fmt.Errorf("invalid %s", field)
	}
	return nil
}

// 🛠️ ตรวจสอบสถานะของ BeforeReturnOrder
func ValidateOrderStatus(order *response.BeforeReturnOrderResponse, expectedStatusReturnID, expectedStatusConfID int) error {
	if order.StatusReturnID != nil && *order.StatusReturnID != expectedStatusReturnID {
		return fmt.Errorf("order is not in the expected return status")
	}
	if order.StatusConfID != nil && *order.StatusConfID != expectedStatusConfID {
		return fmt.Errorf("order is not in the expected confirm status")
	}
	return nil
}

// ✅ ValidateCreateSaleReturn - ตรวจสอบความถูกต้องของคำขอสร้าง Sale Return
func ValidateCreateSaleReturn(req req.CreateSaleReturnOrder) error {
	// 🔹 ตรวจสอบค่าที่ต้องไม่ว่างเปล่า
	requiredFields := []struct {
		name  string
		value string
	}{
		{"order number", req.OrderNo},
		{"SO number", req.SoNo},
		{"customer ID", req.CustomerID},
		{"reason", req.Reason},
		{"logistic", req.Logistic},
	}

	for _, field := range requiredFields {
		if err := validateRequiredString(field.name, field.value); err != nil {
			return err
		}
	}

	// 🔹 ตรวจสอบค่า int ที่ต้องมากกว่า 0
	requiredInts := []struct {
		name  string
		value *int
	}{
		{"channel ID", &req.ChannelID},
		{"warehouse ID", &req.WarehouseID},
	}

	for _, field := range requiredInts {
		if err := validatePositiveInt(field.name, field.value); err != nil {
			return err
		}
	}

	// 🔹 ต้องมีสินค้าขั้นต่ำ 1 รายการ
	if len(req.OrderLines) == 0 {
		return errors.New("at least one order line is required")
	}

	// 🔹 ตรวจสอบข้อมูลของแต่ละสินค้า
	for i, line := range req.OrderLines {
		if err := ValidateSaleReturnLine(line, i+1); err != nil {
			return err
		}
	}

	return nil
}

func ValidateSaleReturnLine(line req.CreateSaleReturnOrderLine, index int) error {
	// 🔹 ตรวจสอบค่า SKU ต้องไม่ว่างเปล่า
	if err := validateRequiredString(fmt.Sprintf("SKU for line %d", index), line.SKU); err != nil {
		return err
	}

	// 🔹 ตรวจสอบค่า Quantity ต้องมากกว่า 0
	if err := validatePositiveInt(fmt.Sprintf("quantity for line %d", index), &line.QTY); err != nil {
		return err
	}

	// 🔹 ตรวจสอบ ReturnQTY ว่าต้องไม่เป็นค่าลบ และต้องไม่มากกว่า QTY
	if line.ReturnQTY < 0 {
		return fmt.Errorf("return quantity cannot be negative for line %d", index)
	}
	if line.ReturnQTY > line.QTY {
		return fmt.Errorf("return quantity cannot be greater than quantity for line %d", index)
	}

	// 🔹 ตรวจสอบ Price ต้องไม่เป็นค่าลบ
	if line.Price < 0 {
		return fmt.Errorf("price cannot be negative for line %d", index)
	}

	return nil
}

// ✅ ValidateUpdateSaleReturn - ตรวจสอบความถูกต้องของคำขออัปเดต Sale Return
func ValidateUpdateSaleReturn(req req.UpdateSaleReturn) error {
	if err := validateRequiredString("order number", req.OrderNo); err != nil {
		return err
	}
	if err := validateRequiredString("SR number", req.SrNo); err != nil {
		return err
	}
	return nil
}

// ✅ ValidateCreateReturnOrder - ตรวจสอบความถูกต้องของการสร้าง Return Order
func ValidateCreateReturnOrder(req req.CreateReturnOrder) error {
	// 🔹 ตรวจสอบค่าพื้นฐาน
	if err := validateRequiredString("order number", req.OrderNo); err != nil {
		return err
	}
	if err := validateRequiredString("SO number", req.SoNo); err != nil {
		return err
	}
	if err := validatePositiveInt("channel ID", req.ChannelID); err != nil {
		return err
	}

	// 🔹 ตรวจสอบว่า order lines ต้องมีข้อมูล
	if len(req.ReturnOrderLine) == 0 {
		return fmt.Errorf("at least one order line is required")
	}

	// 🔹 ตรวจสอบค่าภายใน order lines
	for i, line := range req.ReturnOrderLine {
		if err := validateRequiredString(fmt.Sprintf("SKU for line %d", i+1), line.SKU); err != nil {
			return err
		}
		if err := validatePositiveInt(fmt.Sprintf("quantity for line %d", i+1), line.QTY); err != nil {
			return err
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
