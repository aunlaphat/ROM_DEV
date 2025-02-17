package service

import (
	"context"

	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"boilerplate-backend-go/utils"

	"go.uber.org/zap"
)

type UserService interface {
	GetUser(ctx context.Context, userID string) (response.UserResponse, error)
	GetUsers(ctx context.Context, isActive *bool, limit, offset int) ([]response.UserResponse, error)
	AddUser(ctx context.Context, req request.AddUserRequest, adminID string, roleID int) (*response.AddUserResponse, error)
	EditUser(ctx context.Context, userID string, req request.EditUserRequest, adminID string, roleID int) (*response.EditUserResponse, error)
	DeleteUser(ctx context.Context, userID, adminID string, roleID int) error
	ResetPassword(ctx context.Context, req request.ResetPasswordRequest, adminID string, roleID int) (*response.ResetPasswordResponse, error)
}

// ✅ 1️⃣ GetUser (ดึงข้อมูลผู้ใช้)
func (srv service) GetUser(ctx context.Context, userID string) (response.UserResponse, error) {
	srv.logger.Info("🔍 Fetching user details", zap.String("userID", userID))

	user, err := srv.userRepo.GetUser(ctx, userID)
	if err != nil {
		srv.logger.Warn("❌ User not found", zap.String("userID", userID))
		return response.UserResponse{}, errors.NotFoundError("user not found")
	}

	srv.logger.Info("✅ User details retrieved successfully", zap.String("userID", userID))
	return response.UserResponse(user), nil
}

// ✅ 2️⃣ GetUsers (ดึงรายชื่อผู้ใช้ทั้งหมด)
func (srv service) GetUsers(ctx context.Context, isActive *bool, limit, offset int) ([]response.UserResponse, error) {
	srv.logger.Info("📋 Fetching user list", zap.Int("limit", limit), zap.Int("offset", offset))

	users, err := srv.userRepo.GetUsers(ctx, isActive, limit, offset)
	if err != nil {
		srv.logger.Error("❌ Failed to fetch users", zap.Error(err))
		return nil, err
	}

	var userResponses []response.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, response.UserResponse(user))
	}

	srv.logger.Info("✅ Users retrieved successfully", zap.Int("totalUsers", len(userResponses)))
	return userResponses, nil
}

// ✅ 3️⃣ AddUser (เพิ่มผู้ใช้ใหม่)
func (srv service) AddUser(ctx context.Context, req request.AddUserRequest, adminID string, roleID int) (*response.AddUserResponse, error) {
	srv.logger.Info("➕ Adding new user", zap.String("userID", req.UserID), zap.String("adminID", adminID))

	// 🔹 SYSTEM_ADMIN เท่านั้นที่สามารถเพิ่มผู้ใช้ใหม่ได้
	if roleID != 5 {
		srv.logger.Warn("❌ Unauthorized attempt to add user", zap.String("adminID", adminID))
		return nil, errors.UnauthorizedError("you are not allowed to add a new user")
	}

	// 🔹 ตรวจสอบว่า User มีอยู่ใน ERP (ROM_V_User)
	exists, err := srv.userRepo.CheckUserExists(ctx, req.UserID)
	if err != nil {
		srv.logger.Error("❌ Error checking user existence", zap.Error(err))
		return nil, err
	}
	if !exists {
		srv.logger.Warn("⚠️ User not found in ERP", zap.String("userID", req.UserID))
		return nil, errors.NotFoundError("user not found in ERP")
	}

	// 🔹 เพิ่ม User ในระบบ (UserRole + UserStatus)
	err = srv.userRepo.AddUser(ctx, req.UserID, req.RoleID, adminID)
	if err != nil {
		srv.logger.Error("❌ Failed to add user", zap.Error(err))
		return nil, err
	}

	srv.logger.Info("✅ User added successfully", zap.String("userID", req.UserID))
	return &response.AddUserResponse{UserID: req.UserID, RoleID: req.RoleID, CreatedBy: adminID}, nil
}

// ✅ 4️⃣ EditUser (แก้ไขข้อมูลผู้ใช้)
func (srv service) EditUser(ctx context.Context, userID string, req request.EditUserRequest, adminID string, roleID int) (*response.EditUserResponse, error) {
	srv.logger.Info("✏️ Editing user", zap.String("userID", userID), zap.String("adminID", adminID))

	// 🔹 SYSTEM_ADMIN เท่านั้นที่สามารถแก้ไข Role ของ User อื่นได้
	if roleID != 5 {
		srv.logger.Warn("❌ Unauthorized access", zap.String("adminID", adminID))
		return nil, errors.UnauthorizedError("you are not allowed to edit this user")
	}

	// 🔹 อัปเดต Role ของ User
	err := srv.userRepo.EditUser(ctx, userID, req.RoleID, adminID)
	if err != nil {
		srv.logger.Error("❌ Failed to edit user", zap.Error(err))
		return nil, err
	}

	srv.logger.Info("✅ User edited successfully", zap.String("userID", userID))
	return &response.EditUserResponse{UserID: userID, NewRoleID: req.RoleID, UpdatedBy: adminID}, nil
}

// ✅ 5️⃣ DeleteUser (ลบผู้ใช้แบบ Soft Delete)
func (srv service) DeleteUser(ctx context.Context, userID, adminID string, roleID int) error {
	srv.logger.Info("🗑️ Deleting user", zap.String("userID", userID), zap.String("adminID", adminID))

	// 🔹 SYSTEM_ADMIN เท่านั้นที่สามารถลบผู้ใช้ได้
	if roleID != 5 {
		srv.logger.Warn("❌ Unauthorized delete attempt", zap.String("adminID", adminID))
		return errors.UnauthorizedError("you are not allowed to delete this user")
	}

	// 🔹 อัปเดต IsActive เป็น 0 (Soft Delete)
	err := srv.userRepo.DeleteUser(ctx, userID, adminID)
	if err != nil {
		srv.logger.Error("❌ Failed to delete user", zap.Error(err))
		return err
	}

	srv.logger.Info("✅ User deleted successfully", zap.String("userID", userID))
	return nil
}

// ✅ 6️⃣ ResetPassword (รีเซ็ตรหัสผ่าน)
func (srv service) ResetPassword(ctx context.Context, req request.ResetPasswordRequest, adminID string, roleID int) (*response.ResetPasswordResponse, error) {
	srv.logger.Info("🔄 Resetting password", zap.String("userID", req.UserID), zap.String("adminID", adminID))

	// 🔹 SYSTEM_ADMIN เท่านั้นที่สามารถ Reset Password ของ User อื่นได้
	if adminID != req.UserID && roleID != 5 {
		srv.logger.Warn("🚫 Unauthorized access", zap.String("adminID", adminID))
		return nil, errors.UnauthorizedError("you are not authorized to reset another user's password")
	}

	// 🔹 Hash รหัสผ่านใหม่
	hashedPassword := utils.HashPassword(req.NewPassword)

	// 🔹 อัปเดตรหัสผ่านในฐานข้อมูล
	err := srv.userRepo.UpdateUserPassword(ctx, req.UserID, hashedPassword, adminID)
	if err != nil {
		srv.logger.Error("❌ Failed to reset password", zap.Error(err))
		return nil, err
	}

	srv.logger.Info("✅ Password reset successfully", zap.String("userID", req.UserID))
	return &response.ResetPasswordResponse{UserID: req.UserID, UpdatedBy: adminID}, nil
}
