package service

import (
	"context"
	"fmt"
	"time"

	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"

	"go.uber.org/zap"
)

type UserService interface {
	GetUser(ctx context.Context, userID string) (response.UserResponse, error)
	GetUsers(ctx context.Context, isActive bool, limit, offset int) ([]response.UserResponse, error)
	AddUser(ctx context.Context, req request.AddUserRequest, adminID string, roleID int) (*response.AddUserResponse, error)
	EditUser(ctx context.Context, req request.EditUserRequest, adminID string, roleID int) (*response.EditUserResponse, error)
	DeleteUser(ctx context.Context, userID, adminID string, adminRoleID int) (*response.DeleteUserResponse, error)
}

func (srv service) GetUser(ctx context.Context, userID string) (response.UserResponse, error) {
	srv.logger.Info("🔍 [GetUser] Fetching user details", zap.String("userID", userID))

	user, err := srv.userRepo.GetUser(ctx, userID)
	if err != nil {
		srv.logger.Warn("❌ [GetUser] User not found",
			zap.String("userID", userID),
			zap.Error(err),
		)

		return response.UserResponse{}, fmt.Errorf("failed to fetch user details: %w", err)
	}

	return response.UserResponse{
		UserID:        user.UserID,
		UserName:      user.UserName,
		NickName:      user.NickName,
		FullNameTH:    user.FullNameTH,
		DepartmentNo:  user.DepartmentNo,
		RoleID:        user.RoleID,
		RoleName:      user.RoleName,
		WarehouseID:   user.WarehouseID,
		WarehouseName: user.WarehouseName,
		Description:   user.Description,
		IsActive:      user.IsActive,
	}, nil
}

func (srv service) GetUsers(ctx context.Context, isActive bool, limit, offset int) ([]response.UserResponse, error) {
	srv.logger.Info("📋 [GetUsers] Fetching user list",
		zap.Int("limit", limit),
		zap.Int("offset", offset),
		zap.Bool("isActive", isActive),
	)

	users, err := srv.userRepo.GetUsers(ctx, isActive, limit, offset)
	if err != nil {
		srv.logger.Warn("❌ [GetUsers] Failed to fetch users",
			zap.Bool("isActive", isActive),
			zap.Int("limit", limit),
			zap.Int("offset", offset),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	var userResponses []response.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, response.UserResponse{
			UserID:        user.UserID,
			UserName:      user.UserName,
			NickName:      user.NickName,
			FullNameTH:    user.FullNameTH,
			DepartmentNo:  user.DepartmentNo,
			RoleID:        user.RoleID,
			RoleName:      user.RoleName,
			WarehouseID:   user.WarehouseID,
			WarehouseName: user.WarehouseName,
			Description:   user.Description,
			IsActive:      user.IsActive,
		})
	}

	srv.logger.Info("✅ [GetUsers] Users retrieved successfully",
		zap.Int("totalUsers", len(userResponses)),
		zap.Bool("isActive", isActive),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	return userResponses, nil
}

func (srv service) AddUser(ctx context.Context, req request.AddUserRequest, adminID string, roleID int) (*response.AddUserResponse, error) {
	srv.logger.Info("➕ Adding new user", zap.String("userID", req.UserID), zap.String("adminID", adminID))

	if roleID != 1 {
		srv.logger.Warn("❌ Unauthorized attempt to add user", zap.String("adminID", adminID))
		return nil, errors.UnauthorizedError("you are not allowed to add a new user")
	}

	exists, err := srv.userRepo.CheckUserExistsInERP(ctx, req.UserID)
	if err != nil {
		srv.logger.Error("❌ Error checking user existence in ERP", zap.Error(err))
		return nil, err
	}
	if !exists {
		srv.logger.Warn("⚠️ User not found in ERP", zap.String("userID", req.UserID))
		return nil, errors.NotFoundError("user not found in ERP")
	}

	err = srv.userRepo.AddUser(ctx, req, adminID)
	if err != nil {
		srv.logger.Error("❌ Failed to add user", zap.Error(err))
		return nil, err
	}

	srv.logger.Info("✅ User added successfully", zap.String("userID", req.UserID))
	return &response.AddUserResponse{
		UserID:      req.UserID,
		RoleID:      req.RoleID,
		WarehouseID: req.WarehouseID,
		CreatedBy:   adminID,
	}, nil
}

func (srv service) EditUser(ctx context.Context, req request.EditUserRequest, adminID string, adminRoleID int) (*response.EditUserResponse, error) {
	srv.logger.Info("✏️ [EditUser] Editing user",
		zap.String("userID", req.UserID),
		zap.String("adminID", adminID),
	)

	if adminRoleID != 1 {
		srv.logger.Warn("❌ [EditUser] Unauthorized access",
			zap.String("adminID", adminID),
			zap.Int("adminRoleID", adminRoleID),
		)
		return nil, errors.UnauthorizedError("you are not allowed to edit this user")
	}

	exists, err := srv.userRepo.CheckUserExists(ctx, req.UserID)
	if err != nil {
		srv.logger.Error("❌ [EditUser] Error checking user existence",
			zap.String("userID", req.UserID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		srv.logger.Warn("⚠️ [EditUser] User not found",
			zap.String("userID", req.UserID),
		)
		return nil, errors.NotFoundError("user not found")
	}

	err = srv.userRepo.EditUser(ctx, req, adminID)
	if err != nil {
		srv.logger.Error("❌ [EditUser] Failed to update user",
			zap.String("userID", req.UserID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to edit user: %w", err)
	}

	updatedUser, err := srv.userRepo.GetUser(ctx, req.UserID)
	if err != nil {
		srv.logger.Error("❌ [EditUser] Failed to fetch updated user",
			zap.String("userID", req.UserID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to retrieve updated user: %w", err)
	}

	warehouseName, err := srv.constantRepo.GetWarehouseName(ctx, updatedUser.WarehouseID)
	if err != nil {
		srv.logger.Warn("⚠️ [EditUser] Warehouse not found",
			zap.Int("warehouseID", updatedUser.WarehouseID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to retrieve warehouse name: %w", err)
	}

	srv.logger.Info("✅ [EditUser] User edited successfully",
		zap.String("userID", updatedUser.UserID),
		zap.String("roleName", updatedUser.RoleName),
		zap.Int("warehouseID", updatedUser.WarehouseID),
		zap.String("warehouseName", warehouseName),
	)

	return &response.EditUserResponse{
		UserID:        updatedUser.UserID,
		RoleID:        &updatedUser.RoleID,
		RoleName:      updatedUser.RoleName,
		WarehouseID:   &updatedUser.WarehouseID,
		WarehouseName: warehouseName,
		UpdatedBy:     adminID,
		UpdatedAt:     time.Now(),
	}, nil
}

func (srv service) DeleteUser(ctx context.Context, userID, adminID string, adminRoleID int) (*response.DeleteUserResponse, error) {
	srv.logger.Info("🗑️ [DeleteUser] Deactivating user",
		zap.String("userID", userID),
		zap.String("adminID", adminID),
	)

	if adminRoleID != 1 {
		srv.logger.Warn("❌ [DeleteUser] Unauthorized access",
			zap.String("adminID", adminID),
			zap.Int("adminRoleID", adminRoleID),
		)
		return nil, errors.UnauthorizedError("you are not allowed to delete this user")
	}

	user, err := srv.userRepo.GetUser(ctx, userID)
	if err != nil {
		srv.logger.Warn("⚠️ [DeleteUser] User not found",
			zap.String("userID", userID),
		)
		return nil, errors.NotFoundError("user not found")
	}

	warehouseName, err := srv.constantRepo.GetWarehouseName(ctx, user.WarehouseID)
	if err != nil {
		srv.logger.Warn("⚠️ [DeleteUser] Warehouse not found",
			zap.Int("warehouseID", user.WarehouseID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to retrieve warehouse name: %w", err)
	}

	err = srv.userRepo.DeleteUser(ctx, userID, adminID)
	if err != nil {
		srv.logger.Error("❌ [DeleteUser] Failed to delete user",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	srv.logger.Info("✅ [DeleteUser] User deactivated successfully",
		zap.String("userID", userID),
		zap.String("adminID", adminID),
	)

	return &response.DeleteUserResponse{
		UserID:        user.UserID,
		UserName:      user.UserName,
		RoleID:        user.RoleID,
		RoleName:      user.RoleName,
		WarehouseID:   user.WarehouseID,
		WarehouseName: warehouseName,
		DeactivatedBy: adminID,
		DeactivatedAt: time.Now(),
		Message:       "User has been successfully deactivated",
	}, nil
}
