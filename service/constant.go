package service

import (
	"boilerplate-back-go-2411/dto/response"
	"context"

	"go.uber.org/zap"
)

type ConstantService interface {
	GetRoles(ctx context.Context) ([]response.RoleResponse, error)
	GetWarehouses(ctx context.Context) ([]response.WarehouseResponse, error)
}

func (srv service) GetRoles(ctx context.Context) ([]response.RoleResponse, error) {
	srv.logger.Info("🔍 [GetRoles] Fetching roles")

	roles, err := srv.constantRepo.GetRoles(ctx)
	if err != nil {
		srv.logger.Error("❌ Failed to fetch roles", zap.Error(err))
		return nil, err
	}

	srv.logger.Info("✅ Roles retrieved successfully", zap.Int("count", len(roles)))
	return roles, nil
}

func (srv service) GetWarehouses(ctx context.Context) ([]response.WarehouseResponse, error) {
	srv.logger.Info("🔍 [GetWarehouses] Fetching warehouses")

	warehouses, err := srv.constantRepo.GetWarehouses(ctx)
	if err != nil {
		srv.logger.Error("❌ Failed to fetch warehouses", zap.Error(err))
		return nil, err
	}

	srv.logger.Info("✅ Warehouses retrieved successfully", zap.Int("count", len(warehouses)))
	return warehouses, nil
}
