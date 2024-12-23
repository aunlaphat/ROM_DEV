package repository

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// เพิ่ม constant สำหรับ timeout
const (
	defaultTimeout = 10 * time.Second
	txTimeout      = 30 * time.Second
)

// เพิ่ม constants สำหรับ status
const (
	StatusPending    = 1
	StatusInProgress = 2
	StatusCompleted  = 3
	StatusCancelled  = 4
)

// ฟังก์ชันสำหรับ logging debug
func debugLog(format string, v ...interface{}) {
	log.Printf("🐞 DEBUG: "+format, v...)
}

// ReturnOrderRepository interface กำหนด method สำหรับการทำงานกับฐานข้อมูล
type ReturnOrderRepository interface {
	// Create
	CreateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error
	CreateBeforeReturnOrderLine(ctx context.Context, orderNo string, lines []request.BeforeReturnOrderLine) error

	// Read
	ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderData, error)
	GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderData, error)
	ListBeforeReturnOrderLines(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineData, error)
	GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderLineData, error)

	// Update
	UpdateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error
	UpdateBeforeReturnOrderLine(ctx context.Context, orderNo string, line request.BeforeReturnOrderLine) error

	// Transaction
	//BeginTransaction(ctx context.Context) (*sqlx.Tx, error)
}

// Implementation สำหรับ CreateBeforeReturnOrder
func (repo repositoryDB) CreateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error {
	// เพิ่ม context timeout
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// SQL query สำหรับสร้าง BeforeReturnOrder
	queryOrder := `
        INSERT INTO BeforeReturnOrder (
            OrderNo, SaleOrder, SaleReturn, ChannelID, ReturnType, CustomerID, TrackingNo, Logistic, WarehouseID, SoStatusID, MkpStatusID, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CreateDate, UpdateBy, UpdateDate, CancelID
        ) VALUES (
            :OrderNo, :SaleOrder, :SaleReturn, :ChannelID, :ReturnType, :CustomerID, :TrackingNo, :Logistic, :WarehouseID, :SoStatusID, :MkpStatusID, :ReturnDate, :StatusReturnID, :StatusConfID, :ConfirmBy, :CreateBy, GETDATE(), :UpdateBy, :UpdateDate, :CancelID
        )
    `
	// พารามิเตอร์สำหรับ SQL query
	paramsOrder := map[string]interface{}{
		"OrderNo":        order.OrderNo,
		"SaleOrder":      order.SaleOrder,
		"SaleReturn":     order.SaleReturn,
		"ChannelID":      order.ChannelID,
		"ReturnType":     order.ReturnType,
		"CustomerID":     order.CustomerID,
		"TrackingNo":     order.TrackingNo,
		"Logistic":       order.Logistic,
		"WarehouseID":    order.WarehouseID,
		"SoStatusID":     order.SoStatusID,
		"MkpStatusID":    order.MkpStatusID,
		"ReturnDate":     order.ReturnDate,
		"StatusReturnID": order.StatusReturnID,
		"StatusConfID":   order.StatusConfID,
		"ConfirmBy":      order.ConfirmBy,
		"CreateBy":       order.CreateBy,
		"UpdateBy":       order.UpdateBy,
		"UpdateDate":     order.UpdateDate,
		"CancelID":       order.CancelID,
	}

	// Execute SQL query สำหรับสร้าง BeforeReturnOrder
	_, err := repo.db.NamedExecContext(ctx, queryOrder, paramsOrder)
	if err != nil {
		return fmt.Errorf("failed to create BeforeReturnOrder: %w", err)
	}

	return nil
}

// Implementation สำหรับ CreateBeforeReturnOrderLine
func (repo repositoryDB) CreateBeforeReturnOrderLine(ctx context.Context, orderNo string, lines []request.BeforeReturnOrderLine) error {
	// เพิ่ม context timeout
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// SQL query สำหรับสร้าง BeforeReturnOrderLine
	query := `
        INSERT INTO BeforeReturnOrderLine (
            OrderNo, SKU, QTY, ReturnQTY, Price, CreateBy, TrackingNo
        ) VALUES (
            :OrderNo, :SKU, :QTY, :ReturnQTY, :Price, :CreateBy, :TrackingNo
        )
    `
	// Loop ผ่านแต่ละ line และ execute SQL query สำหรับสร้าง BeforeReturnOrderLine
	for _, line := range lines {
		params := map[string]interface{}{
			"OrderNo":    orderNo,
			"SKU":        line.SKU,
			"QTY":        line.QTY,
			"ReturnQTY":  line.ReturnQTY,
			"Price":      line.Price,
			"CreateBy":   "SYSTEM",
			"TrackingNo": line.TrackingNo,
		}
		_, err := repo.db.NamedExecContext(ctx, query, params)
		if err != nil {
			return fmt.Errorf("failed to create order line: %w", err)
		}
	}
	return nil
}

// Implementation สำหรับ GetBeforeReturnOrderLineByOrderNo
func (repo repositoryDB) GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderLineData, error) {
	// เพิ่ม context timeout
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// SQL query สำหรับดึงข้อมูล BeforeReturnOrderLine
	query := `
        SELECT 
            OrderNo,
            SKU,
            QTY,
            ReturnQTY,
            Price,
            TrackingNo,
            CreateDate
        FROM BeforeReturnOrderLine WITH (NOLOCK)
        WHERE OrderNo = :OrderNo
    `

	var line response.BeforeReturnOrderLineData
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	// Execute SQL query สำหรับดึงข้อมูล BeforeReturnOrderLine
	err = nstmt.GetContext(ctx, &line, map[string]interface{}{"OrderNo": orderNo})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order line: %w", err)
	}

	return &line, nil
}

// Implementation สำหรับ GetBeforeReturnOrderByOrderNo
func (repo repositoryDB) GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderData, error) {
	debugLog("🚀 Starting GetBeforeReturnOrderByOrderNo for OrderNo: %s", orderNo)
	// เพิ่ม context timeout
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// SQL query สำหรับดึงข้อมูล BeforeReturnOrder
	query := `
        SELECT OrderNo, SaleOrder, SaleReturn, ChannelID, ReturnType, CustomerID, TrackingNo, Logistic, WarehouseID, SoStatusID, MkpStatusID, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CreateDate, UpdateBy, UpdateDate, CancelID
        FROM BeforeReturnOrder WITH (NOLOCK)
        WHERE OrderNo = :OrderNo
    `
	order := new(response.BeforeReturnOrderData)
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		debugLog("❌ Failed to prepare statement: %v", err)
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	debugLog("📦 Executing query to fetch BeforeReturnOrder for OrderNo: %s", orderNo)

	// Execute SQL query สำหรับดึงข้อมูล BeforeReturnOrder
	err = nstmt.GetContext(ctx, order, map[string]interface{}{"OrderNo": orderNo})
	if err == sql.ErrNoRows {
		debugLog("❗ No order found for OrderNo: %s", orderNo)
		return nil, nil
	}
	if err != nil {
		debugLog("❌ Failed to fetch BeforeReturnOrder: %v", err)
		return nil, fmt.Errorf("failed to fetch BeforeReturnOrder: %w", err)
	}
	debugLog("✅ Successfully fetched BeforeReturnOrder for OrderNo: %s", orderNo)

	// ดึงข้อมูล BeforeReturnOrderLine ที่เกี่ยวข้อง
	lines, err := repo.ListBeforeReturnOrderLines(ctx, orderNo)
	if err != nil {
		return nil, err
	}
	order.ReturnLines = lines

	debugLog("✅ Successfully fetched all lines for OrderNo: %s", orderNo)
	return order, nil
}

func (repo repositoryDB) ListBeforeReturnOrderLines(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineData, error) {
	debugLog("🚀 Starting ListBeforeReturnOrderLines for OrderNo: %s", orderNo)

	// เพิ่ม context timeout
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// SQL query สำหรับดึงข้อมูล BeforeReturnOrderLine
	query := `
        SELECT 
            OrderNo,
            SKU,
            QTY,
            ReturnQTY,
            Price,
            TrackingNo,
            CreateDate
        FROM BeforeReturnOrderLine WITH (NOLOCK)
        WHERE OrderNo = :OrderNo
        ORDER BY RecID
    `

	var lines []response.BeforeReturnOrderLineData
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		debugLog("❌ Failed to prepare statement: %v", err)
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	debugLog("📦 Executing query to fetch BeforeReturnOrderLines for OrderNo: %s", orderNo)
	// Execute SQL query สำหรับดึงข้อมูล BeforeReturnOrderLine
	err = nstmt.SelectContext(ctx, &lines, map[string]interface{}{"OrderNo": orderNo})
	if err != nil {
		debugLog("❌ Failed to fetch BeforeReturnOrderLines: %v", err)
		return nil, fmt.Errorf("failed to get order lines: %w", err)
	}
	debugLog("✅ Successfully fetched %d lines for OrderNo: %s", len(lines), orderNo)

	return lines, nil
}

func (repo repositoryDB) ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderData, error) {
	// เพิ่ม context timeout
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// SQL query สำหรับดึงข้อมูล BeforeReturnOrder
	query := `
        SELECT OrderNo, SaleOrder, SaleReturn, ChannelID, ReturnType, CustomerID, TrackingNo, Logistic, WarehouseID, SoStatusID, MkpStatusID, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CreateDate, UpdateBy, UpdateDate, CancelID
        FROM BeforeReturnOrder WITH (NOLOCK)
        ORDER BY CreateDate DESC
    `
	var orders []response.BeforeReturnOrderData
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	// Execute SQL query สำหรับดึงข้อมูล BeforeReturnOrder
	err = nstmt.SelectContext(ctx, &orders, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	return orders, nil
}

// Implementation สำหรับ UpdateBeforeReturnOrder
func (repo repositoryDB) UpdateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error {
	// เพิ่ม context timeout
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// SQL query สำหรับอัปเดต BeforeReturnOrder
	query := `
        UPDATE BeforeReturnOrder 
        SET SaleOrder = :SaleOrder,
            SaleReturn = :SaleReturn,
            ChannelID = :ChannelID,
            ReturnType = :ReturnType,
            CustomerID = :CustomerID,
            TrackingNo = :TrackingNo,
            Logistic = :Logistic,
            WarehouseID = :WarehouseID,
            SoStatusID = :SoStatusID,
            MkpStatusID = :MkpStatusID,
            ReturnDate = :ReturnDate,
            StatusReturnID = :StatusReturnID,
            StatusConfID = :StatusConfID,
            ConfirmBy = :ConfirmBy,
            UpdateBy = :UpdateBy,
            UpdateDate = GETDATE(),
            CancelID = :CancelID
        WHERE OrderNo = :OrderNo
    `
	// พารามิเตอร์สำหรับ SQL query
	params := map[string]interface{}{
		"OrderNo":        order.OrderNo,
		"SaleOrder":      order.SaleOrder,
		"SaleReturn":     order.SaleReturn,
		"ChannelID":      order.ChannelID,
		"ReturnType":     order.ReturnType,
		"CustomerID":     order.CustomerID,
		"TrackingNo":     order.TrackingNo,
		"Logistic":       order.Logistic,
		"WarehouseID":    order.WarehouseID,
		"SoStatusID":     order.SoStatusID,
		"MkpStatusID":    order.MkpStatusID,
		"ReturnDate":     order.ReturnDate,
		"StatusReturnID": order.StatusReturnID,
		"StatusConfID":   order.StatusConfID,
		"ConfirmBy":      order.ConfirmBy,
		"UpdateBy":       order.UpdateBy,
		"CancelID":       order.CancelID,
	}

	// Execute SQL query สำหรับอัปเดต BeforeReturnOrder
	_, err := repo.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update BeforeReturnOrder: %w", err)
	}

	return nil
}

// Implementation สำหรับ UpdateBeforeReturnOrderLine
func (repo repositoryDB) UpdateBeforeReturnOrderLine(ctx context.Context, orderNo string, line request.BeforeReturnOrderLine) error {
	// เพิ่ม context timeout
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// SQL query สำหรับอัปเดต BeforeReturnOrderLine
	query := `
        UPDATE BeforeReturnOrderLine 
        SET SKU = :SKU,
            QTY = :QTY,
            ReturnQTY = :ReturnQTY,
            Price = :Price,
            UpdateBy = :UpdateBy,
            UpdateDate = GETDATE(),
            TrackingNo = :TrackingNo
        WHERE OrderNo = :OrderNo
          AND SKU = :SKU
    `
	// พารามิเตอร์สำหรับ SQL query
	params := map[string]interface{}{
		"OrderNo":    orderNo,
		"SKU":        line.SKU,
		"QTY":        line.QTY,
		"ReturnQTY":  line.ReturnQTY,
		"Price":      line.Price,
		"UpdateBy":   line.UpdateBy,
		"TrackingNo": line.TrackingNo,
	}

	// Execute SQL query สำหรับอัปเดต BeforeReturnOrderLine
	_, err := repo.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update BeforeReturnOrderLine: %w", err)
	}

	return nil
}

/* // Implementation สำหรับ BeginTransaction
func (repo repositoryDB) BeginTransaction(ctx context.Context) (*sqlx.Tx, error) {
	// เพิ่ม context timeout
	ctx, cancel := context.WithTimeout(ctx, txTimeout)
	defer cancel()

	// เริ่มต้น transaction
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return tx, nil
}
*/
