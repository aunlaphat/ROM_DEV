package repository

import (
	"boilerplate-back-go-2411/errors"
	"context"
	"fmt"
)

type ConstantRepository interface {
	GetWarehouseName(ctx context.Context, warehouseID int) (string, error)
}

func (repo repositoryDB) GetWarehouseName(ctx context.Context, warehouseID int) (string, error) {
	var warehouseName string
	query := `SELECT WarehouseName FROM Warehouse WHERE WarehouseID = :warehouseID`
	params := map[string]interface{}{"warehouseID": warehouseID}

	rows, err := repo.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return "", fmt.Errorf("query execution error: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&warehouseName)
		if err != nil {
			return "", fmt.Errorf("failed to scan warehouse name: %w", err)
		}
		return warehouseName, nil
	}

	return "", errors.NotFoundError("warehouse not found")
}
