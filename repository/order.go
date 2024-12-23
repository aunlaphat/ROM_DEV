package repository

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"context"
	"database/sql"
	"fmt"
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

type ReturnOrderRepository interface {
	// Create
	CreateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrderRequest) error
	CreateBeforeReturnOrderLine(ctx context.Context, orderNo string, lines []request.BeforeReturnOrderLineRequest) error

	// Read
	GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderData, error)
	GetBeforeReturnOrderLines(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineData, error)
	ListBeforeReturnOrders(ctx context.Context, limit, offset int) ([]response.BeforeReturnOrderData, error)

	// Update
	UpdateBeforeReturnOrderStatus(ctx context.Context, orderNo string, statusID int, updateBy string) error

	// Delete (Soft Delete)
	CancelBeforeReturnOrder(ctx context.Context, orderNo string, cancelBy string) error

	// เพิ่ม method สำหรับดึงข้อมูล ReturnOrder
	GetOrderByReturnID(ctx context.Context, returnID string) (response.ReturnOrder, error)
}

// Implementation สำหรับ Create
func (repo repositoryDB) CreateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrderRequest) error {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	queryOrder := `
        INSERT INTO BeforeReturnOrder (
            OrderNo, CustomerID, TrackingNo, 
            StatusReturnID, CreateBy, CreateDate, ChannelID, SaleOrder, SaleReturn, ReturnType, Logistic, WarehouseID, StatusConfID
        ) VALUES (
            :OrderNo, :CustomerID, :TrackingNo,
            :StatusReturnID, :CreateBy, GETDATE(), :ChannelID, :SaleOrder, :SaleReturn, :ReturnType, :Logistic, :WarehouseID, :StatusConfID
        )
    `
	paramsOrder := map[string]interface{}{
		"OrderNo":        order.OrderNo,
		"CustomerID":     order.CustomerID,
		"TrackingNo":     order.TrackingNo,
		"StatusReturnID": StatusPending,
		"CreateBy":       order.CreateBy,
		"ChannelID":      order.ChannelID,
		"SaleOrder":      sql.NullString{String: order.SaleOrder, Valid: order.SaleOrder != ""},
		"SaleReturn":     sql.NullString{String: order.SaleReturn, Valid: order.SaleReturn != ""},
		"ReturnType":     order.ReturnType,
		"Logistic":       order.Logistic,
		"WarehouseID":    order.WarehouseID,
		"StatusConfID":   1,
	}

	_, err = tx.NamedExecContext(ctx, queryOrder, paramsOrder)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create BeforeReturnOrder: %w", err)
	}

	queryLine := `
        INSERT INTO BeforeReturnOrderLine (
            OrderNo, SKU, QTY, ReturnQTY, 
            Price, CreateBy, TrackingNo
        ) VALUES (
            :OrderNo, :SKU, :QTY, :ReturnQTY, 
            :Price, :CreateBy, :TrackingNo
        )
    `

	for _, line := range order.ReturnLines {
		paramsLine := map[string]interface{}{
			"OrderNo":    order.OrderNo,
			"SKU":        line.SKU,
			"QTY":        line.QTY,
			"ReturnQTY":  line.ReturnQTY,
			"Price":      line.Price,
			"CreateBy":   order.CreateBy,
			"TrackingNo": line.TrackingNo,
		}
		_, err = tx.NamedExecContext(ctx, queryLine, paramsLine)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create BeforeReturnOrderLine: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (repo repositoryDB) CreateBeforeReturnOrderLine(ctx context.Context, orderNo string, lines []request.BeforeReturnOrderLineRequest) error {
	query := `
        INSERT INTO BeforeReturnOrderLine (
            OrderNo, SKU, QTY, ReturnQTY,
            Price, CreateBy, TrackingNo
        ) VALUES (
            :OrderNo, :SKU, :QTY, :ReturnQTY,
            :Price, :CreateBy, :TrackingNo
        )
    `
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

func (repo repositoryDB) GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderData, error) {
	query := `
        SELECT OrderNo, CustomerID, TrackingNo, StatusReturnID, CreateBy, CreateDate, ChannelID, SaleOrder, SaleReturn, ReturnType, Logistic, WarehouseID, StatusConfID
        FROM BeforeReturnOrder WITH (NOLOCK)
        WHERE OrderNo = :OrderNo
    `
	order := new(response.BeforeReturnOrderData)
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	err = nstmt.GetContext(ctx, order, map[string]interface{}{"OrderNo": orderNo})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	lines, err := repo.GetBeforeReturnOrderLines(ctx, orderNo)
	if err != nil {
		return nil, err
	}
	order.ReturnLines = lines

	return order, nil
}

func (repo repositoryDB) GetBeforeReturnOrderLines(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineData, error) {
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
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	err = nstmt.SelectContext(ctx, &lines, map[string]interface{}{"OrderNo": orderNo})
	if err != nil {
		return nil, fmt.Errorf("failed to get order lines: %w", err)
	}

	return lines, nil
}

func (repo repositoryDB) ListBeforeReturnOrders(ctx context.Context, limit, offset int) ([]response.BeforeReturnOrderData, error) {
	query := `
        SELECT OrderNo, CustomerID, TrackingNo, StatusReturnID, CreateBy, CreateDate, ChannelID, SaleOrder, SaleReturn, ReturnType, Logistic, WarehouseID, StatusConfID
        FROM BeforeReturnOrder WITH (NOLOCK)
        ORDER BY CreateDate DESC
        OFFSET :Offset ROWS
        FETCH NEXT :Limit ROWS ONLY
    `
	var orders []response.BeforeReturnOrderData
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	err = nstmt.SelectContext(ctx, &orders, map[string]interface{}{
		"Offset": offset,
		"Limit":  limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	return orders, nil
}

func (repo repositoryDB) UpdateBeforeReturnOrderStatus(ctx context.Context, orderNo string, statusID int, updateBy string) error {
	query := `
        UPDATE BeforeReturnOrder 
        SET StatusReturnID = :StatusReturnID,
            UpdateBy = :UpdateBy,
            UpdateDate = GETDATE()
        WHERE OrderNo = :OrderNo
    `
	result, err := repo.db.NamedExecContext(ctx, query, map[string]interface{}{
		"StatusReturnID": statusID,
		"UpdateBy":       updateBy,
		"OrderNo":        orderNo,
	})
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("order not found: %s", orderNo)
	}
	return nil
}

func (repo repositoryDB) GetOrderByReturnID(ctx context.Context, returnID string) (response.ReturnOrder, error) {
	query := `
        SELECT 
            ReturnID, OrderNo, SaleOrder, SaleReturn, TrackingNo,
            PlatfID, ChannelID, OptStatusID, AxStatusID, PlatfStatusID,
            Remark, CreateBy, CreateDate, UpdateBy, UpdateDate,
            CancelID, StatusCheckID, CheckBy, Description
        FROM ReturnOrder WITH (NOLOCK)
        WHERE ReturnID = :ReturnID
    `
	var order response.ReturnOrder
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return response.ReturnOrder{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	err = nstmt.GetContext(ctx, &order, map[string]interface{}{"ReturnID": returnID})
	if err != nil {
		if err == sql.ErrNoRows {
			return response.ReturnOrder{}, fmt.Errorf("order not found: %s", returnID)
		}
		return response.ReturnOrder{}, fmt.Errorf("failed to get order: %w", err)
	}

	queryLines := `
        SELECT 
            RecID, ReturnID, OrderNo, TrackingNo, SKU,
            ReturnQTY, CheckQTY, Price, CreateBy, CreateDate,
            AlterSKU, UpdateBy, UpdateDate
        FROM ReturnOrderLine WITH (NOLOCK)
        WHERE ReturnID = :ReturnID
    `
	err = repo.db.SelectContext(ctx, &order.ReturnOrderLines, queryLines, map[string]interface{}{"ReturnID": returnID})
	if err != nil {
		return response.ReturnOrder{}, fmt.Errorf("failed to get order lines: %w", err)
	}

	return order, nil
}

func (repo repositoryDB) CancelBeforeReturnOrder(ctx context.Context, orderNo string, cancelBy string) error {
	query := `
        UPDATE BeforeReturnOrder 
        SET CancelID = 1,
            UpdateBy = :CancelBy,
            UpdateDate = GETDATE()
        WHERE OrderNo = :OrderNo
          AND CancelID IS NULL
    `
	result, err := repo.db.NamedExecContext(ctx, query, map[string]interface{}{
		"CancelBy": cancelBy,
		"OrderNo":  orderNo,
	})
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("order not found or already cancelled: %s", orderNo)
	}
	return nil
}
