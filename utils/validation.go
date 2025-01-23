package utils

import (
	req "boilerplate-backend-go/dto/request"
	"fmt"
)

// เพิ่มฟังก์ชัน validate สำหรับ CreateSaleReturn
func ValidateCreateSaleReturn(req req.BeforeReturnOrder) error {
	// 1. ตรวจสอบข้อมูลพื้นฐาน
	if req.OrderNo == "" {
		return fmt.Errorf("order number is required")
	}
	if req.SoNo == "" {
		return fmt.Errorf("SO number is required")
	}
	if req.CustomerID == "" {
		return fmt.Errorf("customer ID is required")
	}

	// 2. ตรวจสอบค่าที่ต้องมากกว่า 0
	if req.ChannelID <= 0 {
		return fmt.Errorf("invalid channel ID")
	}
	if req.WarehouseID <= 0 {
		return fmt.Errorf("invalid warehouse ID")
	}

	// 3. ตรวจสอบ ReturnType
	/* validReturnTypes := map[string]bool{
		"NORMAL": true,
		"DAMAGE": true,
		// เพิ่ม type อื่นๆ ตามต้องการ
	}
	if !validReturnTypes[req.ReturnType] {
		return fmt.Errorf("invalid return type: %s", req.ReturnType)
	} */

	// 4. ตรวจสอบ order lines
	if len(req.BeforeReturnOrderLines) == 0 {
		return fmt.Errorf("at least one order line is required")
	}

	for i, line := range req.BeforeReturnOrderLines {
		if line.SKU == "" {
			return fmt.Errorf("SKU is required for line %d", i+1)
		}
		if line.QTY <= 0 {
			return fmt.Errorf("quantity must be greater than 0 for line %d", i+1)
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
		// ตรวจสอบ AlterSKU ถ้ามี
		if line.AlterSKU != nil && *line.AlterSKU == "" {
			return fmt.Errorf("alter SKU cannot be empty if provided for line %d", i+1)
		}
	}

	return nil
}

func ValidateUpdateSaleReturn(orderNo string, srNo string, updateBy string) error {
	if orderNo == "" {
		return fmt.Errorf("order number is required")
	}
	if srNo == "" {
		return fmt.Errorf("SR number is required")
	}
	if updateBy == "" {
		return fmt.Errorf("updater information is required")
	}

	return nil
}
