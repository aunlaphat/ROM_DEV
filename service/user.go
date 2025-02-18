package service

import (
	"context"
	"fmt"

	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"

	"go.uber.org/zap"
)

// ✅ **UserService Interface**
type UserService interface {
	GetUser(ctx context.Context, userID string) (response.UserResponse, error)
	GetUsers(ctx context.Context, isActive bool, limit, offset int) ([]response.UserResponse, error)
	AddUser(ctx context.Context, req request.AddUserRequest, adminID string, roleID int) (*response.AddUserResponse, error)
	EditUser(ctx context.Context, req request.EditUserRequest, adminID string, roleID int) (*response.EditUserResponse, error)
	DeleteUser(ctx context.Context, userID, adminID string, roleID int) error
	//ResetPassword(ctx context.Context, req request.ResetPasswordRequest, adminID string, roleID int) (*response.ResetPasswordResponse, error)
}

// ✅ 1️⃣ GetUser (ดึงข้อมูลผู้ใช้)
func (srv service) GetUser(ctx context.Context, userID string) (response.UserResponse, error) {
	srv.logger.Info("🔍 [GetUser] Fetching user details", zap.String("userID", userID))

	// 🟢 ดึงข้อมูลผู้ใช้จาก repository
	user, err := srv.userRepo.GetUser(ctx, userID)
	if err != nil {
		srv.logger.Warn("❌ [GetUser] User not found",
			zap.String("userID", userID),
			zap.Error(err),
		)

		// 🛑 แสดง error รายละเอียดของ SQL Query
		srv.logger.Debug("🛠 [SQL Debug] Failed Query Execution",
			zap.String("query", "SELECT UserID, UserName, NickName, FullNameTH, DepartmentNo, RoleID, RoleName, Description, IsActive FROM ROM_V_UserDetail WHERE UserID = :userID"),
			zap.String("param_userID", userID),
			zap.Error(err),
		)

		return response.UserResponse{}, fmt.Errorf("failed to fetch user details: %w", err)
	}

	return response.UserResponse{
		UserID:       user.UserID,
		UserName:     user.UserName,
		NickName:     user.NickName,
		FullNameTH:   user.FullNameTH,
		DepartmentNo: user.DepartmentNo,
		RoleID:       user.RoleID,
		RoleName:     user.RoleName,
		Description:  user.Description,
		IsActive:     user.IsActive,
	}, nil
}

// ✅ 2️⃣ GetUsers (ดึงรายชื่อผู้ใช้ทั้งหมด)
func (srv service) GetUsers(ctx context.Context, isActive bool, limit, offset int) ([]response.UserResponse, error) {
	srv.logger.Info("📋 [GetUsers] Fetching user list",
		zap.Int("limit", limit),
		zap.Int("offset", offset),
		zap.Bool("isActive", isActive),
	)

	// 🟢 ดึงข้อมูลผู้ใช้จาก repository
	users, err := srv.userRepo.GetUsers(ctx, isActive, limit, offset)
	if err != nil {
		srv.logger.Warn("❌ [GetUsers] Failed to fetch users",
			zap.Bool("isActive", isActive),
			zap.Int("limit", limit),
			zap.Int("offset", offset),
			zap.Error(err),
		)

		// 🛑 แสดง error รายละเอียดของ SQL Query
		srv.logger.Debug("🛠 [SQL Debug] Failed Query Execution",
			zap.String("query", "SELECT UserID, UserName, NickName, FullNameTH, DepartmentNo, RoleID, RoleName, Description, IsActive FROM ROM_V_UserDetail WHERE IsActive = :isActive ORDER BY UserID OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY"),
			zap.Bool("param_isActive", isActive),
			zap.Int("param_limit", limit),
			zap.Int("param_offset", offset),
			zap.Error(err),
		)

		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	var userResponses []response.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, response.UserResponse{
			UserID:       user.UserID,
			UserName:     user.UserName,
			NickName:     user.NickName,
			FullNameTH:   user.FullNameTH,
			DepartmentNo: user.DepartmentNo,
			RoleID:       user.RoleID,
			RoleName:     user.RoleName,
			Description:  user.Description,
			IsActive:     user.IsActive,
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

// ✅ **4️⃣ AddUser - เพิ่มผู้ใช้ใหม่**
func (srv service) AddUser(ctx context.Context, req request.AddUserRequest, adminID string, roleID int) (*response.AddUserResponse, error) {
	srv.logger.Info("➕ Adding new user", zap.String("userID", req.UserID), zap.String("adminID", adminID))

	// 🔹 SYSTEM_ADMIN เท่านั้นที่สามารถเพิ่มผู้ใช้ใหม่ได้
	if roleID != 1 {
		srv.logger.Warn("❌ Unauthorized attempt to add user", zap.String("adminID", adminID))
		return nil, errors.UnauthorizedError("you are not allowed to add a new user")
	}

	// 🔹 ตรวจสอบว่าผู้ใช้มีอยู่ใน ERP
	exists, err := srv.userRepo.CheckUserExistsInERP(ctx, req.UserID)
	if err != nil {
		srv.logger.Error("❌ Error checking user existence in ERP", zap.Error(err))
		return nil, err
	}
	if !exists {
		srv.logger.Warn("⚠️ User not found in ERP", zap.String("userID", req.UserID))
		return nil, errors.NotFoundError("user not found in ERP")
	}

	// 🔹 เพิ่ม User ในระบบ พร้อม Transaction
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

// ✅ **5️⃣ EditUser - แก้ไขข้อมูลผู้ใช้**
func (srv service) EditUser(ctx context.Context, req request.EditUserRequest, adminID string, adminRoleID int) (*response.EditUserResponse, error) {
	srv.logger.Info("✏️ [EditUser] Editing user",
		zap.String("userID", req.UserID),
		zap.String("adminID", adminID),
	)

	// 🟢 **ตรวจสอบสิทธิ์**
	if adminRoleID != 1 { // RoleID = 1 คือ ADMIN
		srv.logger.Warn("❌ [EditUser] Unauthorized access",
			zap.String("adminID", adminID),
			zap.Int("adminRoleID", adminRoleID),
		)
		return nil, errors.UnauthorizedError("you are not allowed to edit this user")
	}

	// 🟢 **ตรวจสอบว่าผู้ใช้มีอยู่จริงหรือไม่**
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

	// 🟢 **เรียกใช้ Repository เพื่ออัปเดตข้อมูล**
	err = srv.userRepo.EditUser(ctx, req, adminID)
	if err != nil {
		srv.logger.Error("❌ [EditUser] Failed to update user",
			zap.String("userID", req.UserID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to edit user: %w", err)
	}

	// 🟢 **Response**
	srv.logger.Info("✅ [EditUser] User edited successfully",
		zap.String("userID", req.UserID),
		zap.String("adminID", adminID),
	)

	return &response.EditUserResponse{
		UserID:         req.UserID,
		NewRoleID:      req.RoleID,
		NewWarehouseID: req.WarehouseID,
		UpdatedBy:      adminID,
	}, nil
}

// ✅ **6️⃣ DeleteUser - ปิดการใช้งานบัญชี (Soft Delete)**
func (srv service) DeleteUser(ctx context.Context, userID, adminID string, roleID int) error {
	srv.logger.Info("🗑️ Deleting user", zap.String("userID", userID), zap.String("adminID", adminID))

	// 🔹 SYSTEM_ADMIN เท่านั้นที่สามารถลบผู้ใช้ได้
	if roleID != 1 {
		srv.logger.Warn("❌ Unauthorized delete attempt", zap.String("adminID", adminID))
		return errors.UnauthorizedError("you are not allowed to delete this user")
	}

	// 🔹 ลบผู้ใช้แบบ Soft Delete
	err := srv.userRepo.DeleteUser(ctx, userID, adminID)
	if err != nil {
		srv.logger.Error("❌ Failed to delete user", zap.Error(err))
		return err
	}

	srv.logger.Info("✅ User deleted successfully", zap.String("userID", userID))
	return nil
}

// ✅ **7️⃣ GetCurrentPassword - ดึงรหัสผ่านปัจจุบันของ User**
func (srv service) GetCurrentPassword(ctx context.Context, userID string) (string, error) {
	srv.logger.Info("🔍 Fetching user password", zap.String("userID", userID))

	password, err := srv.userRepo.GetCurrentPassword(ctx, userID)
	if err != nil {
		srv.logger.Warn("❌ Failed to fetch password", zap.String("userID", userID))
		return "", errors.NotFoundError("password not found")
	}

	srv.logger.Info("✅ Password retrieved successfully", zap.String("userID", userID))
	return password, nil
}

/* // ✅ **6️⃣ ResetPassword - รีเซ็ตรหัสผ่าน**
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
} */
