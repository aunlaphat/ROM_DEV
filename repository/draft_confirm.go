package repository

import (
	"boilerplate-back-go-2411/dto/request"
	"boilerplate-back-go-2411/dto/response"
	"context"
	"fmt"
)

type DraftConfirmRepository interface {
	GetOrders(ctx context.Context, statusConfID int, startDate, endDate string) ([]response.OrderHeadResponse, error)
	GetOrderWithItems(ctx context.Context, orderNo string, statusConfID int) (*response.DraftConfirmResponse, error)
	ListCodeR(ctx context.Context) ([]response.ListCodeRResponse, error)
	AddItemToDraftOrder(ctx context.Context, req request.AddItem, userID string) ([]response.AddItemResponse, error)
	RemoveItemFromDraftOrder(ctx context.Context, orderNo string, sku string) (int64, error)
}

func (repo repositoryDB) GetOrders(ctx context.Context, statusConfID int, startDate, endDate string) ([]response.OrderHeadResponse, error) {
	var orders []response.OrderHeadResponse
	query := `
		SELECT OrderNo, SoNo, SrNo, CustomerID, TrackingNo, Logistic, ChannelID, CreateDate, WarehouseID
		FROM BeforeReturnOrder 
		WHERE StatusConfID = :statusConfID
		AND CreateDate BETWEEN :startDate AND :endDate
		ORDER BY CreateDate DESC
	`

	params := map[string]interface{}{
		"statusConfID": statusConfID,
		"startDate":    startDate,
		"endDate":      endDate,
	}

	rows, err := repo.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order response.OrderHeadResponse
		if err := rows.StructScan(&order); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (repo repositoryDB) GetOrderWithItems(ctx context.Context, orderNo string, statusConfID int) (*response.DraftConfirmResponse, error) {
	queryHead := `
		SELECT OrderNo, SoNo, SrNo
		FROM BeforeReturnOrder 
		WHERE OrderNo = :orderNo AND StatusConfID = :statusConfID
	`

	params := map[string]interface{}{
		"orderNo":      orderNo,
		"statusConfID": statusConfID,
	}

	rows, err := repo.db.NamedQueryContext(ctx, queryHead, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}
	defer rows.Close()

	var order response.DraftConfirmResponse
	if rows.Next() {
		if err := rows.StructScan(&order); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
	} else {
		return nil, fmt.Errorf("order not found")
	}

	queryLine := `
		SELECT OrderNo, SKU, ItemName, QTY, Price
		FROM BeforeReturnOrderLine 
		WHERE OrderNo = :orderNo
	`

	itemRows, err := repo.db.NamedQueryContext(ctx, queryLine, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order items: %w", err)
	}
	defer itemRows.Close()

	var items []response.DraftConfirmItem
	for itemRows.Next() {
		var item response.DraftConfirmItem
		if err := itemRows.StructScan(&item); err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, item)
	}

	order.Items = items

	return &order, nil
}

func (repo repositoryDB) ListCodeR(ctx context.Context) ([]response.ListCodeRResponse, error) {
	query := `
        SELECT SKU, NAMEALIAS
        FROM ROM_V_ProductAll
        WHERE SKU LIKE 'R%' -- ดึงเฉพาะ SKU ที่ขึ้นต้นด้วย 'R'
        ORDER BY NAMEALIAS ASC
    `

	var codeRList []response.ListCodeRResponse

	err := repo.db.SelectContext(ctx, &codeRList, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CodeR list: %w", err)
	}

	return codeRList, nil
}

func (repo repositoryDB) AddItemToDraftOrder(ctx context.Context, req request.AddItem, userID string) ([]response.AddItemResponse, error) {
	query := `
        INSERT INTO BeforeReturnOrderLine (OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, CreateBy, CreateDate)
        OUTPUT inserted.OrderNo, inserted.SKU, inserted.ItemName, inserted.QTY, inserted.ReturnQTY, inserted.Price, inserted.CreateBy, inserted.CreateDate
        VALUES (:OrderNo, :SKU, :ItemName, :QTY, :ReturnQTY, :Price, :CreateBy, GETDATE())
    `

	rows, err := repo.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"OrderNo":   req.OrderNo,
		"SKU":       req.SKU,
		"ItemName":  req.ItemName,
		"QTY":       req.QTY,
		"ReturnQTY": req.ReturnQTY,
		"Price":     req.Price,
		"CreateBy":  userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert item: %w", err)
	}
	defer rows.Close()

	var results []response.AddItemResponse
	for rows.Next() {
		var item response.AddItemResponse
		if err := rows.StructScan(&item); err != nil {
			return nil, fmt.Errorf("failed to scan inserted item: %w", err)
		}
		results = append(results, item)
	}

	return results, nil
}

func (repo repositoryDB) RemoveItemFromDraftOrder(ctx context.Context, orderNo string, sku string) (int64, error) {
	query := `
        DELETE FROM BeforeReturnOrderLine
        WHERE OrderNo = :OrderNo AND SKU = :SKU
    `

	result, err := repo.db.NamedExecContext(ctx, query, map[string]interface{}{
		"OrderNo": orderNo,
		"SKU":     sku,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to delete item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows: %w", err)
	}

	return rowsAffected, nil
}
