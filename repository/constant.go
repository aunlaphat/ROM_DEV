package repository

import (
	"boilerplate-back-go-2411/dto/response"
	"boilerplate-back-go-2411/errors"
	"context"
	"fmt"
)

type ConstantRepository interface {
	GetRoles(ctx context.Context) ([]response.RoleResponse, error)
	GetWarehouses(ctx context.Context) ([]response.WarehouseResponse, error)
	GetWarehouseName(ctx context.Context, warehouseID int) (string, error)
}

func (repo repositoryDB) GetRoles(ctx context.Context) ([]response.RoleResponse, error) {
	var roles []response.RoleResponse
	query := `SELECT RoleID, RoleName FROM Role ORDER BY RoleID`

	err := repo.db.SelectContext(ctx, &roles, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roles: %w", err)
	}

	if len(roles) == 0 {
		return []response.RoleResponse{}, nil
	}

	return roles, nil
}

func (repo repositoryDB) GetWarehouses(ctx context.Context) ([]response.WarehouseResponse, error) {
	var warehouses []response.WarehouseResponse
	query := `SELECT WarehouseID, WarehouseName FROM Warehouse ORDER BY WarehouseID`

	err := repo.db.SelectContext(ctx, &warehouses, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch warehouses: %w", err)
	}

	if len(warehouses) == 0 {
		return []response.WarehouseResponse{}, nil
	}

	return warehouses, nil
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
