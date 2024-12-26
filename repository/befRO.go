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

// ReturnOrderRepository interface กำหนด method สำหรับการทำงานกับฐานข้อมูล
type BefRORepository interface {
	// Create
	CreateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error
	CreateBeforeReturnOrderLine(ctx context.Context, orderNo string, lines []request.BeforeReturnOrderLine) error

	// Read
	ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error)
	GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error)
	ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error)
	ListBeforeReturnOrderLinesByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error)
	GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error)

	// Update
	UpdateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error
	UpdateBeforeReturnOrderLine(ctx context.Context, orderNo string, line request.BeforeReturnOrderLine) error

	// Transaction
	CreateReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error
	UpdateBeforeReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error

	//Cancle
}

// Implementation สำหรับ CreateBeforeReturnOrder
func (repo repositoryDB) CreateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error {
	log.Printf("🚀 Starting CreateBeforeReturnOrder for OrderNo: %s", order.OrderNo)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	queryOrder := `
        INSERT INTO BeforeReturnOrder (
            OrderNo, SaleOrder, SaleReturn, ChannelID, ReturnType, CustomerID, TrackingNo, Logistic, WarehouseID, SoStatusID, MkpStatusID, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CreateDate, CancelID
        ) VALUES (
            :OrderNo, :SaleOrder, :SaleReturn, :ChannelID, :ReturnType, :CustomerID, :TrackingNo, :Logistic, :WarehouseID, :SoStatusID, :MkpStatusID, :ReturnDate, :StatusReturnID, :StatusConfID, :ConfirmBy, :CreateBy, GETDATE(), :CancelID
        )
    `
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
		"CancelID":       order.CancelID,
	}

	_, err := repo.db.NamedExecContext(ctx, queryOrder, paramsOrder)
	if err != nil {
		log.Printf("❌ Failed to create BeforeReturnOrder: %v", err)
		return fmt.Errorf("failed to create BeforeReturnOrder: %w", err)
	}

	log.Printf("✅ Successfully created BeforeReturnOrder for OrderNo: %s", order.OrderNo)
	return nil
}

// Implementation สำหรับ CreateBeforeReturnOrderLine
func (repo repositoryDB) CreateBeforeReturnOrderLine(ctx context.Context, orderNo string, lines []request.BeforeReturnOrderLine) error {
	log.Printf("🚀 Starting CreateBeforeReturnOrderLine for OrderNo: %s", orderNo)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `
        INSERT INTO BeforeReturnOrderLine (
            OrderNo, SKU, QTY, ReturnQTY, Price, CreateBy, TrackingNo
        ) VALUES (
            :OrderNo, :SKU, :QTY, :ReturnQTY, :Price, :CreateBy, :TrackingNo
        )
    `
	for _, line := range lines {
		trackingNo := line.TrackingNo
		if trackingNo == "" {
			trackingNo = "N/A"
		}

		params := map[string]interface{}{
			"OrderNo":    orderNo,
			"SKU":        line.SKU,
			"QTY":        line.QTY,
			"ReturnQTY":  line.ReturnQTY,
			"Price":      line.Price,
			"CreateBy":   line.CreateBy,
			"TrackingNo": trackingNo,
		}
		_, err := repo.db.NamedExecContext(ctx, query, params)
		if err != nil {
			log.Printf("❌ Failed to create order line: %v", err)
			return fmt.Errorf("failed to create order line: %w", err)
		}
	}
	log.Printf("✅ Successfully created BeforeReturnOrderLine for OrderNo: %s", orderNo)
	return nil
}

// Implementation สำหรับ GetBeforeReturnOrderLineByOrderNo
func (repo repositoryDB) GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error) {
	log.Printf("🚀 Starting GetBeforeReturnOrderLineByOrderNo for OrderNo: %s", orderNo)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

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

	var lines []response.BeforeReturnOrderLineResponse
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		log.Printf("❌ Failed to prepare statement: %v", err)
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	err = nstmt.SelectContext(ctx, &lines, map[string]interface{}{"OrderNo": orderNo})
	if err != nil {
		log.Printf("❌ Failed to get order lines: %v", err)
		return nil, fmt.Errorf("failed to get order lines: %w", err)
	}

	log.Printf("✅ Successfully fetched BeforeReturnOrderLines for OrderNo: %s", orderNo)
	return lines, nil
}

// Implementation สำหรับ GetBeforeReturnOrderByOrderNo
func (repo repositoryDB) GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error) {
	log.Printf("🚀 Starting GetBeforeReturnOrderByOrderNo for OrderNo: %s", orderNo)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `
        SELECT OrderNo, SaleOrder, SaleReturn, ChannelID, ReturnType, CustomerID, TrackingNo, Logistic, WarehouseID, SoStatusID, MkpStatusID, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CreateDate, UpdateBy, UpdateDate, CancelID
        FROM BeforeReturnOrder WITH (NOLOCK)
        WHERE OrderNo = :OrderNo
    `
	order := new(response.BeforeReturnOrderResponse)
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		log.Printf("❌ Failed to prepare statement: %v", err)
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	log.Printf("📦 Executing query to fetch BeforeReturnOrder for OrderNo: %s", orderNo)

	err = nstmt.GetContext(ctx, order, map[string]interface{}{"OrderNo": orderNo})
	if err == sql.ErrNoRows {
		log.Printf("❗ No order found for OrderNo: %s", orderNo)
		return nil, nil
	}
	if err != nil {
		log.Printf("❌ Failed to fetch BeforeReturnOrder: %v", err)
		return nil, fmt.Errorf("failed to fetch BeforeReturnOrder: %w", err)
	}
	log.Printf("✅ Successfully fetched BeforeReturnOrder for OrderNo: %s", orderNo)

	lines, err := repo.ListBeforeReturnOrderLinesByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, err
	}
	order.BeforeReturnOrderLines = lines

	log.Printf("✅ Successfully fetched all lines for OrderNo: %s", orderNo)
	return order, nil
}

func (repo repositoryDB) ListBeforeReturnOrderLinesByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error) {
	log.Printf("🚀 Starting ListBeforeReturnOrderLinesByOrderNo for OrderNo: %s", orderNo)

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

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

	var lines []response.BeforeReturnOrderLineResponse
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		log.Printf("❌ Failed to prepare statement: %v", err)
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	log.Printf("📦 Executing query to fetch BeforeReturnOrderLines for OrderNo: %s", orderNo)
	err = nstmt.SelectContext(ctx, &lines, map[string]interface{}{"OrderNo": orderNo})
	if err != nil {
		log.Printf("❌ Failed to fetch BeforeReturnOrderLines: %v", err)
		return nil, fmt.Errorf("failed to get order lines: %w", err)
	}
	log.Printf("✅ Successfully fetched %d lines for OrderNo: %s", len(lines), orderNo)

	return lines, nil
}

func (repo repositoryDB) ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error) {
	log.Printf("🚀 Starting ListBeforeReturnOrderLines")

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

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
        ORDER BY RecID
    `

	var lines []response.BeforeReturnOrderLineResponse
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		log.Printf("❌ Failed to prepare statement: %v", err)
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	log.Printf("📦 Executing query to fetch BeforeReturnOrderLines")
	err = nstmt.SelectContext(ctx, &lines, map[string]interface{}{})
	if err != nil {
		log.Printf("❌ Failed to fetch BeforeReturnOrderLines: %v", err)
		return nil, fmt.Errorf("failed to get order lines: %w", err)
	}
	log.Printf("✅ Successfully fetched %d lines", len(lines))

	return lines, nil
}

func (repo repositoryDB) ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error) {
	log.Printf("🚀 Starting ListBeforeReturnOrders")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `
        SELECT OrderNo, SaleOrder, SaleReturn, ChannelID, ReturnType, CustomerID, TrackingNo, Logistic, WarehouseID, SoStatusID, MkpStatusID, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CreateDate, UpdateBy, UpdateDate, CancelID
        FROM BeforeReturnOrder WITH (NOLOCK)
        ORDER BY CreateDate ASC
    `
	var orders []response.BeforeReturnOrderResponse
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		log.Printf("❌ Failed to prepare statement: %v", err)
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	err = nstmt.SelectContext(ctx, &orders, map[string]interface{}{})
	if err != nil {
		log.Printf("❌ Failed to list orders: %v", err)
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	for i := range orders {
		lines, err := repo.ListBeforeReturnOrderLinesByOrderNo(ctx, orders[i].OrderNo)
		if err != nil {
			return nil, err
		}
		orders[i].BeforeReturnOrderLines = lines
	}

	log.Printf("✅ Successfully listed %d orders", len(orders))
	return orders, nil
}

// Implementation สำหรับ BeginTransaction CreateBeforeReturnOrder & CreateBeforeReturnOrderLine
func (repo repositoryDB) CreateReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error {
	log.Printf("🚀 Starting CreateReturnOrderWithTransaction for OrderNo: %s", order.OrderNo)
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("❌ Failed to start transaction: %v", err)
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	queryOrder := `
        INSERT INTO BeforeReturnOrder (
            OrderNo, SaleOrder, SaleReturn, ChannelID, ReturnType, CustomerID, TrackingNo, Logistic, WarehouseID, 
            SoStatusID, MkpStatusID, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CreateDate
        ) VALUES (
            :OrderNo, :SaleOrder, :SaleReturn, :ChannelID, :ReturnType, :CustomerID, :TrackingNo, :Logistic, :WarehouseID, 
            :SoStatusID, :MkpStatusID, :ReturnDate, :StatusReturnID, :StatusConfID, :ConfirmBy, :CreateBy, GETDATE()
        )
    `
	_, err = tx.NamedExecContext(ctx, queryOrder, map[string]interface{}{
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
	})
	if err != nil {
		tx.Rollback()
		log.Printf("❌ Failed to create BeforeReturnOrder: %v", err)
		return fmt.Errorf("failed to create BeforeReturnOrder: %w", err)
	}

	queryLine := `
        INSERT INTO BeforeReturnOrderLine (
            OrderNo, SKU, QTY, ReturnQTY, Price, CreateBy, TrackingNo
        ) VALUES (
            :OrderNo, :SKU, :QTY, :ReturnQTY, :Price, :CreateBy, :TrackingNo
        )
    `
	for _, line := range order.BeforeReturnOrderLines {
		// Ensure TrackingNo is not NULL
		trackingNo := line.TrackingNo
		if trackingNo == "" {
			trackingNo = "N/A" // Default value if TrackingNo is not provided
		}

		_, err = tx.NamedExecContext(ctx, queryLine, map[string]interface{}{
			"OrderNo":    order.OrderNo,
			"SKU":        line.SKU,
			"QTY":        line.QTY,
			"ReturnQTY":  line.ReturnQTY,
			"Price":      line.Price,
			"CreateBy":   "SYSTEM",
			"TrackingNo": trackingNo,
			// "CreateDate": line.CreateDate, // MSSQL GetDate() function
		})
		if err != nil {
			tx.Rollback()
			log.Printf("❌ Failed to create BeforeReturnOrderLine: %v", err)
			return fmt.Errorf("failed to create BeforeReturnOrderLine: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("❌ Failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("✅ Successfully created ReturnOrder with transaction for OrderNo: %s", order.OrderNo)
	return nil
}

func (repo repositoryDB) UpdateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error {
	log.Printf("🚀 Starting UpdateBeforeReturnOrder for OrderNo: %s", order.OrderNo)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `
        UPDATE BeforeReturnOrder 
        SET SaleOrder = COALESCE(:SaleOrder, SaleOrder),
            SaleReturn = COALESCE(:SaleReturn, SaleReturn),
            ChannelID = COALESCE(:ChannelID, ChannelID),
            ReturnType = COALESCE(:ReturnType, ReturnType),
            CustomerID = COALESCE(:CustomerID, CustomerID),
            TrackingNo = COALESCE(:TrackingNo, TrackingNo),
            Logistic = COALESCE(:Logistic, Logistic),
            WarehouseID = COALESCE(:WarehouseID, WarehouseID),
            SoStatusID = COALESCE(:SoStatusID, SoStatusID),
            MkpStatusID = COALESCE(:MkpStatusID, MkpStatusID),
            ReturnDate = COALESCE(:ReturnDate, ReturnDate),
            StatusReturnID = COALESCE(:StatusReturnID, StatusReturnID),
            StatusConfID = COALESCE(:StatusConfID, StatusConfID),
            ConfirmBy = COALESCE(:ConfirmBy, ConfirmBy),
            UpdateBy = COALESCE(:UpdateBy, UpdateBy),
            UpdateDate = GETDATE(),
            CancelID = COALESCE(:CancelID, CancelID)
        WHERE OrderNo = :OrderNo
    `
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

	_, err := repo.db.NamedExecContext(ctx, query, params)
	if err != nil {
		log.Printf("❌ Failed to update BeforeReturnOrder: %v", err)
		return fmt.Errorf("failed to update BeforeReturnOrder: %w", err)
	}

	log.Printf("✅ Successfully updated BeforeReturnOrder for OrderNo: %s", order.OrderNo)
	return nil
}

func (repo repositoryDB) UpdateBeforeReturnOrderLine(ctx context.Context, orderNo string, line request.BeforeReturnOrderLine) error {
	log.Printf("🚀 Starting UpdateBeforeReturnOrderLine for OrderNo: %s, SKU: %s", orderNo, line.SKU)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `
        UPDATE BeforeReturnOrderLine 
        SET QTY = COALESCE(:QTY, QTY),
            ReturnQTY = COALESCE(:ReturnQTY, ReturnQTY),
            Price = COALESCE(:Price, Price),
            UpdateBy = COALESCE(:UpdateBy, UpdateBy),
            UpdateDate = GETDATE(),
            TrackingNo = COALESCE(:TrackingNo, TrackingNo)
        WHERE OrderNo = :OrderNo
          AND SKU = :SKU
    `
	params := map[string]interface{}{
		"OrderNo":    orderNo,
		"SKU":        line.SKU,
		"QTY":        line.QTY,
		"ReturnQTY":  line.ReturnQTY,
		"Price":      line.Price,
		"UpdateBy":   line.UpdateBy,
		"TrackingNo": line.TrackingNo,
	}

	_, err := repo.db.NamedExecContext(ctx, query, params)
	if err != nil {
		log.Printf("❌ Failed to update BeforeReturnOrderLine: %v", err)
		return fmt.Errorf("failed to update BeforeReturnOrderLine: %w", err)
	}

	log.Printf("✅ Successfully updated BeforeReturnOrderLine for OrderNo: %s, SKU: %s", orderNo, line.SKU)
	return nil
}

// Implementation สำหรับ UpdateBeforeReturnOrderWithTransaction
func (repo repositoryDB) UpdateBeforeReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error {
	log.Printf("🚀 Starting UpdateBeforeReturnOrderWithTransaction for OrderNo: %s", order.OrderNo)

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("❌ Failed to start transaction: %v", err)
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Update BeforeReturnOrderLine first
	queryLine := `
        UPDATE BeforeReturnOrderLine 
        SET QTY = COALESCE(:QTY, QTY),
            ReturnQTY = COALESCE(:ReturnQTY, ReturnQTY),
            Price = COALESCE(:Price, Price),
            UpdateBy = COALESCE(:UpdateBy, UpdateBy),
            UpdateDate = GETDATE(),
            TrackingNo = COALESCE(:TrackingNo, TrackingNo)
        WHERE OrderNo = :OrderNo
          AND SKU = :SKU
    `

	for _, line := range order.BeforeReturnOrderLines {
		paramsLine := map[string]interface{}{
			"OrderNo":    line.OrderNo,
			"SKU":        line.SKU,
			"QTY":        line.QTY,
			"ReturnQTY":  line.ReturnQTY,
			"Price":      line.Price,
			"UpdateBy":   line.UpdateBy,
			"TrackingNo": line.TrackingNo,
		}

		log.Printf("🔍 Updating BeforeReturnOrderLine with OrderNo: %s, SKU: %s", line.OrderNo, line.SKU)
		result, err := tx.NamedExecContext(ctx, queryLine, paramsLine)
		if err != nil {
			log.Printf("❌ Failed to update BeforeReturnOrderLine: %v", err)
			return fmt.Errorf("failed to update BeforeReturnOrderLine: %w", err)
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			log.Printf("❗ No rows updated for OrderNo: %s, SKU: %s", line.OrderNo, line.SKU)
			return fmt.Errorf("no rows updated for OrderNo: %s, SKU: %s", line.OrderNo, line.SKU)
		}
	}

	// Update BeforeReturnOrder
	queryOrder := `
        UPDATE BeforeReturnOrder 
        SET SaleOrder = COALESCE(:SaleOrder, SaleOrder),
            SaleReturn = COALESCE(:SaleReturn, SaleReturn),
            ChannelID = COALESCE(:ChannelID, ChannelID),
            ReturnType = COALESCE(:ReturnType, ReturnType),
            CustomerID = COALESCE(:CustomerID, CustomerID),
            TrackingNo = COALESCE(:TrackingNo, TrackingNo),
            Logistic = COALESCE(:Logistic, Logistic),
            WarehouseID = COALESCE(:WarehouseID, WarehouseID),
            SoStatusID = COALESCE(:SoStatusID, SoStatusID),
            MkpStatusID = COALESCE(:MkpStatusID, MkpStatusID),
            ReturnDate = COALESCE(:ReturnDate, ReturnDate),
            StatusReturnID = COALESCE(:StatusReturnID, StatusReturnID),
            StatusConfID = COALESCE(:StatusConfID, StatusConfID),
            ConfirmBy = COALESCE(:ConfirmBy, ConfirmBy),
            UpdateBy = COALESCE(:UpdateBy, UpdateBy),
            UpdateDate = GETDATE(),
            CancelID = COALESCE(:CancelID, CancelID)
        WHERE OrderNo = :OrderNo
    `

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
		"UpdateBy":       order.UpdateBy,
		"CancelID":       order.CancelID,
	}

	_, err = tx.NamedExecContext(ctx, queryOrder, paramsOrder)
	if err != nil {
		log.Printf("❌ Failed to update BeforeReturnOrder: %v", err)
		return fmt.Errorf("failed to update BeforeReturnOrder: %w", err)
	}

	log.Printf("✅ Successfully updated BeforeReturnOrderWithTransaction for OrderNo: %s", order.OrderNo)
	return nil
}
