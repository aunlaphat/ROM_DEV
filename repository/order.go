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
