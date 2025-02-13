package repository

import (
	entity "boilerplate-backend-go/Entity"
	"boilerplate-backend-go/dto/response"
	"context"
	"fmt"
)

type DraftConfirmRepository interface {
	GetOrders(ctx context.Context, statusConfID int, startDate, endDate string) ([]response.OrderHeadResponse, error)
	GetOrderWithItems(ctx context.Context, orderNo string) (*response.DraftConfirmResponse, error)
	AddItemToDraftOrder(ctx context.Context, orderNo string, item entity.BeforeReturnOrderLine) error
	RemoveItemFromDraftOrder(ctx context.Context, orderNo, sku string) error
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

	// ✅ ใช้ Named Parameters แบบ struct
	params := map[string]interface{}{
		"statusConfID": statusConfID,
		"startDate":    startDate,
		"endDate":      endDate,
	}

	// ✅ ใช้ `NamedQueryContext` สำหรับ Named Parameters
	rows, err := repo.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}
	defer rows.Close()

	// ✅ อ่านข้อมูลจาก `rows`
	for rows.Next() {
		var order response.OrderHeadResponse
		if err := rows.StructScan(&order); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (repo repositoryDB) GetOrderWithItems(ctx context.Context, orderNo string) (*response.DraftConfirmResponse, error) {
	// ✅ ดึงข้อมูลคำสั่งคืนสินค้า (HEAD)
	var order response.DraftConfirmResponse
	queryHead := `
		SELECT OrderNo, SoNo, SrNo FROM BeforeReturnOrder 
		WHERE OrderNo = :orderNo
	`
	err := repo.db.GetContext(ctx, &order, queryHead, map[string]interface{}{
		"orderNo": orderNo,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}

	// ✅ ดึงข้อมูลรายการสินค้า (LINE)
	var items []response.DraftConfirmItem
	queryLine := `
		SELECT OrderNo, SKU, ItemName, QTY, Price
		FROM BeforeReturnOrderLine 
		WHERE OrderNo = :orderNo
	`
	err = repo.db.SelectContext(ctx, &items, queryLine, map[string]interface{}{
		"orderNo": orderNo,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order items: %w", err)
	}

	// ✅ รวมข้อมูล HEAD + LINE
	order.Items = items

	return &order, nil
}

func (repo repositoryDB) AddItemToDraftOrder(ctx context.Context, orderNo string, item entity.BeforeReturnOrderLine) error {
	query := `
		INSERT INTO BeforeReturnOrderLine (OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, CreateBy, CreateDate)
		VALUES (:OrderNo, :SKU, :ItemName, :QTY, :ReturnQTY, :Price, :CreateBy, GETDATE())
	`
	_, err := repo.db.NamedExecContext(ctx, query, map[string]interface{}{
		"OrderNo":   orderNo,
		"SKU":       item.SKU,
		"ItemName":  item.ItemName,
		"QTY":       item.QTY,
		"ReturnQTY": item.ReturnQTY,
		"Price":     item.Price,
		"CreateBy":  item.CreateBy,
	})
	if err != nil {
		return fmt.Errorf("failed to add item: %w", err)
	}
	return nil
}

func (repo repositoryDB) RemoveItemFromDraftOrder(ctx context.Context, orderNo, sku string) error {
	query := `
		DELETE FROM BeforeReturnOrderLine
		WHERE OrderNo = :orderNo AND SKU = :sku
	`
	_, err := repo.db.NamedExecContext(ctx, query, map[string]interface{}{
		"orderNo": orderNo,
		"sku":     sku,
	})
	if err != nil {
		return fmt.Errorf("failed to remove item: %w", err)
	}
	return nil
}
