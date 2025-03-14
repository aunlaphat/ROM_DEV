package utils

import (
	req "boilerplate-back-go-2411/dto/request"
	"fmt"
	"strings"
)

func ValidateCreateReturnOrder(req req.CreateReturnOrder) error {
	var validate []string

	// ตรวจสอบข้อมูลออเดอร์
	if req.OrderNo == "" {
		validate = append(validate, "order number is required")
	}
	if req.SoNo == "" {
		validate = append(validate, "SO number is required")
	}

	if *req.ChannelID == 0 {
		validate = append(validate, "invalid channel ID")
	}

	// ตรวจสอบรายการคืน
	if len(req.ReturnOrderLine) == 0 {
		validate = append(validate, "at least one order line is required")
	}

	for i, line := range req.ReturnOrderLine {
		if line.SKU == "" {
			validate = append(validate, fmt.Sprintf("SKU is required for line %d", i+1))
		}
		if *line.QTY <= 0 {
			validate = append(validate, fmt.Sprintf("quantity must be greater than 0 for line %d", i+1))
		}
		// if line.ReturnQTY < 0 {
		// 	validate = append(validate, fmt.Sprintf("return quantity cannot be negative for line %d", i+1))
		// }
		// if line.ReturnQTY > *line.QTY {
		// 	validate = append(validate, fmt.Sprintf("return quantity cannot be greater than quantity for line %d", i+1))
		// }
		// if line.Price < 0 {
		// 	validate = append(validate, fmt.Sprintf("price cannot be negative for line %d", i+1))
		// }
	}

	// หากพบข้อผิดพลาด ให้ส่งคืนข้อผิดพลาดทั้งหมด
	if len(validate) > 0 {
		formattedErrors := make([]string, len(validate))

		for i, err := range validate {
			formattedErrors[i] = fmt.Sprintf("{%s}", err)
		}

		errorMsg := strings.Join(formattedErrors, ", ")
		return fmt.Errorf("%s", errorMsg)
	}

	return nil
}

func ValidateCreateTradeReturn(req req.BeforeReturnOrder) error {
	var validate []string

	// ตรวจสอบข้อมูลออเดอร์
	if req.OrderNo == "" {
		validate = append(validate, "order number is required")
	}
	if req.SoNo == "" {
		validate = append(validate, "SO number is required")
	}
	if req.CustomerID == "" {
		validate = append(validate, "customer ID is required")
	}

	if req.ChannelID == 0 {
		validate = append(validate, "invalid channel ID")
	}
	// if req.WarehouseID == 0 {
	// 	validate = append(validate, "invalid warehouse ID")
	// }

	// ตรวจสอบรายการคืน
	if len(req.BeforeReturnOrderLines) == 0 {
		validate = append(validate, "at least one order line is required")
	}

	for i, line := range req.BeforeReturnOrderLines {
		if line.SKU == "" {
			validate = append(validate, fmt.Sprintf("SKU is required for line %d", i+1))
		}
		if line.QTY <= 0 {
			validate = append(validate, fmt.Sprintf("quantity must be greater than 0 for line %d", i+1))
		}
		// if line.ReturnQTY < 0 {
		// 	validate = append(validate, fmt.Sprintf("return quantity cannot be negative for line %d", i+1))
		// }
		// if line.ReturnQTY > line.QTY {
		// 	validate = append(validate, fmt.Sprintf("return quantity cannot be greater than quantity for line %d", i+1))
		// }
		// if line.Price < 0 {
		// 	validate = append(validate, fmt.Sprintf("price cannot be negative for line %d", i+1))
		// }
		// if line.AlterSKU != nil && *line.AlterSKU == "" {
		// 	validate = append(validate, fmt.Sprintf("alter SKU cannot be empty if provided for line %d", i+1))
		// }
	}

	// หากพบข้อผิดพลาด ให้ส่งคืนข้อผิดพลาดทั้งหมด
	if len(validate) > 0 {
		formattedErrors := make([]string, len(validate))

		for i, err := range validate {
			formattedErrors[i] = fmt.Sprintf("{%s}", err)
		}

		errorMsg := strings.Join(formattedErrors, ", ")
		return fmt.Errorf("%s", errorMsg)
	}

	return nil
}

func ValidateCreateTradeReturnLine(lines []req.OrderLines) error {
	var validate []string

	for i, line := range lines {
		if line.SKU == "" {
			validate = append(validate, fmt.Sprintf("SKU is required for line %d", i+1))
		}
		if line.QTY <= 0 {
			validate = append(validate, fmt.Sprintf("quantity must be greater than 0 for line %d", i+1))
		}
		// if line.ReturnQTY < 0 {
		// 	validate = append(validate, fmt.Sprintf("return quantity cannot be negative for line %d", i+1))
		// }
		// if line.ReturnQTY > line.QTY {
		// 	validate = append(validate, fmt.Sprintf("return quantity cannot be greater than quantity for line %d", i+1))
		// }
		// if line.Price < 0 {
		// 	validate = append(validate, fmt.Sprintf("price cannot be negative for line %d", i+1))
		// }
		// if line.AlterSKU != nil && *line.AlterSKU == "" {
		// 	validate = append(validate, fmt.Sprintf("alter SKU cannot be empty if provided for line %d", i+1))
		// }
	}

	// หากพบข้อผิดพลาด ให้ส่งคืนข้อผิดพลาดทั้งหมด
	if len(validate) > 0 {
		formattedErrors := make([]string, len(validate))

		for i, err := range validate {
			formattedErrors[i] = fmt.Sprintf("{%s}", err)
		}

		errorMsg := strings.Join(formattedErrors, ", ")
		return fmt.Errorf("%s", errorMsg)
	}

	return nil
}
