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

	// 🔹 Insert `BeforeReturnOrder`
	queryHead := `
		INSERT INTO BeforeReturnOrder 
		(OrderNo, SoNo, SoStatus, MkpStatus, WarehouseID, ReturnDate, TrackingNo, Logistic, CreateBy, CreateDate)
		VALUES 
		(:OrderNo, :SoNo, :SoStatus, :MkpStatus, :WarehouseID, :ReturnDate, :TrackingNo, :Logistic, :CreateBy, GETDATE())`

	_, err = tx.NamedExecContext(ctx, queryHead, map[string]interface{}{
		"OrderNo":     req.OrderNo,
		"SoNo":        req.SoNo,
		"SoStatus":    req.SoStatus,
		"MkpStatus":   req.MkpStatus,
		"WarehouseID": req.WarehouseID,
		"ReturnDate":  req.ReturnDate,
		"TrackingNo":  req.TrackingNo,
		"Logistic":    req.Logistic,
		"CreateBy":    userID,
	})
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert BeforeReturnOrder: %w", err)
	}

	// ✅ Batch Insert `BeforeReturnOrderLine`
	queryLines := `
		INSERT INTO BeforeReturnOrderLine (OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, CreateBy, CreateDate, TrackingNo, AlterSKU)
		VALUES `

	var placeholders []string
	var batchItems []map[string]interface{}

	for _, item := range req.Items {
		placeholders = append(placeholders, "(:OrderNo, :SKU, :ItemName, :QTY, :ReturnQTY, :Price, :CreateBy, GETDATE(), :TrackingNo, :AlterSKU)")
		batchItems = append(batchItems, map[string]interface{}{
			"OrderNo":    req.OrderNo,
			"SKU":        item.SKU,
			"ItemName":   item.ItemName,
			"QTY":        item.QTY,
			"ReturnQTY":  item.ReturnQTY,
			"Price":      item.Price,
			"CreateBy":   userID,
			"TrackingNo": item.TrackingNo,
			"AlterSKU":   item.AlterSKU,
		})
	}

	queryLines += strings.Join(placeholders, ",")

	_, err = tx.NamedExecContext(ctx, queryLines, batchItems)
	if err != nil {
		tx.Rollback()
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
		SELECT OrderNo, SoNo, SrNo, ChannelID, Reason, CustomerID, TrackingNo, Logistic, WarehouseID, 
		       SoStatus, MkpStatus, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, ConfirmDate, 
		       CreateBy, CreateDate, UpdateBy, UpdateDate, CancelID, IsCNCreated, IsEdited
		FROM BeforeReturnOrder WHERE OrderNo = :OrderNo`

	var order response.BeforeReturnOrderResponse
	err := repo.db.GetContext(ctx, &order, query, map[string]interface{}{"OrderNo": orderNo})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve order: %w", err)
	}

	return &order, nil
}

func (repo repositoryDB) GetBeforeReturnOrderItems(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderItem, error) {
	query := `
		SELECT OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, CreateBy, CreateDate, TrackingNo, AlterSKU
		FROM BeforeReturnOrderLine WHERE OrderNo = :OrderNo`

	var items []response.BeforeReturnOrderItem
	err := repo.db.SelectContext(ctx, &items, query, map[string]interface{}{"OrderNo": orderNo})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve order items: %w", err)
	}

	return items, nil
}
