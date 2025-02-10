package repository

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type OrderRepository interface {
	SearchOrder(ctx context.Context, req request.SearchOrder) (*response.SearchOrderResponse, error)
	CreateBeforeReturnOrder(ctx context.Context, req request.CreateBeforeReturnOrder, userID string) error
	GetBeforeReturnOrder(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error)
	GetBeforeReturnOrderItems(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderItem, error)
}

func (repo repositoryDB) SearchOrder(ctx context.Context, req request.SearchOrder) (*response.SearchOrderResponse, error) {
	if req.SoNo == "" && req.OrderNo == "" {
		return nil, fmt.Errorf("either SoNo or OrderNo must be provided")
	}

	queryConditions := []string{}

	if req.SoNo != "" {
		queryConditions = append(queryConditions, "SoNo = :SoNo")
	}
	if req.OrderNo != "" {
		queryConditions = append(queryConditions, "OrderNo = :OrderNo")
	}

	queryHead := `
        SELECT SoNo, OrderNo, StatusMKP, SalesStatus, CreateDate
        FROM ROM_V_OrderHeadDetail
    `
	if len(queryConditions) > 0 {
		queryHead += " WHERE " + strings.Join(queryConditions, " AND ")
	}

	stmt, err := repo.db.PrepareNamedContext(ctx, queryHead)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query head: %w", err)
	}
	defer stmt.Close()

	var order response.SearchOrderResponse
	err = stmt.GetContext(ctx, &order, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to fetch order head: %w", err)
	}

	queryLines := `
        SELECT SKU, ItemName, QTY, Price
        FROM ROM_V_OrderLineDetail
        WHERE SoNo = :SoNo
    `

	stmtLines, err := repo.db.PrepareNamedContext(ctx, queryLines)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query lines: %w", err)
	}
	defer stmtLines.Close()

	var items []response.SearchOrderItem
	err = stmtLines.SelectContext(ctx, &items, map[string]interface{}{"SoNo": order.SoNo})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order lines: %w", err)
	}

	order.Items = items
	return &order, nil
}

func (repo repositoryDB) CreateBeforeReturnOrder(ctx context.Context, req request.CreateBeforeReturnOrder, userID string) error {
	// 🔹 เริ่ม Transaction
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// ✅ ใช้ `defer` เพื่อจัดการ Rollback อัตโนมัติถ้ามี Error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 🔹 Insert `BeforeReturnOrder`
	queryHead := `
		INSERT INTO BeforeReturnOrder 
		(OrderNo, SoNo, ChannelID, CustomerID, Reason, SoStatus, MkpStatus, WarehouseID, ReturnDate, TrackingNo, Logistic, CreateBy, CreateDate)
		VALUES 
		(:OrderNo, :SoNo, :ChannelID, :CustomerID, :Reason, :SoStatus, :MkpStatus, :WarehouseID, :ReturnDate, :TrackingNo, :Logistic, :CreateBy, GETDATE())`

	_, err = tx.NamedExecContext(ctx, queryHead, map[string]interface{}{
		"OrderNo":     req.OrderNo,
		"SoNo":        req.SoNo,
		"ChannelID":   req.ChannelID,
		"CustomerID":  req.CustomerID,
		"Reason":      req.Reason,
		"SoStatus":    req.SoStatus,
		"MkpStatus":   req.MkpStatus,
		"WarehouseID": req.WarehouseID,
		"ReturnDate":  req.ReturnDate,
		"TrackingNo":  req.TrackingNo,
		"Logistic":    req.Logistic,
		"CreateBy":    userID,
	})
	if err != nil {
		return fmt.Errorf("failed to insert BeforeReturnOrder: %w", err)
	}

	// ✅ Batch Insert `BeforeReturnOrderLine` (ไม่ต้อง Loop ใช้ NamedExecContext)
	queryLines := `
		INSERT INTO BeforeReturnOrderLine 
		(OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, CreateBy, CreateDate, TrackingNo, AlterSKU)
		VALUES 
		(:OrderNo, :SKU, :ItemName, :QTY, :ReturnQTY, :Price, :CreateBy, GETDATE(), :TrackingNo, :AlterSKU)`

	// ✅ Map `req.Items` เพื่อให้ `OrderNo` ถูกต้อง
	for i := range req.Items {
		req.Items[i].OrderNo = req.OrderNo
		req.Items[i].CreateBy = userID
	}

	_, err = tx.NamedExecContext(ctx, queryLines, req.Items)
	if err != nil {
		return fmt.Errorf("failed to execute batch insert for BeforeReturnOrderLine: %w", err)
	}

	// ✅ Commit Transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (repo repositoryDB) GetBeforeReturnOrder(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error) {
	query := `
		SELECT OrderNo, SoNo, SrNo, ChannelID, CustomerID, Reason, TrackingNo, Logistic, WarehouseID, 
		       SoStatus, MkpStatus, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, ConfirmDate, 
		       CreateBy, CreateDate, UpdateBy, UpdateDate, CancelID, IsCNCreated, IsEdited
		FROM BeforeReturnOrder WHERE OrderNo = :OrderNo`

	// ✅ ใช้ NamedQueryContext() เพื่อรองรับ Named Parameters
	rows, err := repo.db.NamedQueryContext(ctx, query, map[string]interface{}{"OrderNo": orderNo})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 🛑 คืนค่า `nil` ถ้าไม่พบ OrderNo
		}
		return nil, fmt.Errorf("failed to retrieve order: %w", err)
	}
	defer rows.Close() // ✅ ปิด `rows` หลังใช้

	var order response.BeforeReturnOrderResponse
	if rows.Next() {
		if err := rows.StructScan(&order); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
	}

	return &order, nil
}

func (repo repositoryDB) GetBeforeReturnOrderItems(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderItem, error) {
	query := `
		SELECT OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, CreateBy, CreateDate, TrackingNo, AlterSKU
		FROM BeforeReturnOrderLine WHERE OrderNo = :OrderNo`

	rows, err := repo.db.NamedQueryContext(ctx, query, map[string]interface{}{"OrderNo": orderNo})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 🛑 คืนค่า `nil` ถ้าไม่พบรายการสินค้า
		}
		return nil, fmt.Errorf("failed to retrieve order items: %w", err)
	}
	defer rows.Close() // ✅ ปิด `rows` เมื่อใช้เสร็จ

	var items []response.BeforeReturnOrderItem
	for rows.Next() {
		var item response.BeforeReturnOrderItem
		if err := rows.StructScan(&item); err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}
