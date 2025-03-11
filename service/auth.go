package service

import (
	"context"

	"boilerplate-back-go-2411/dto/request"
	"boilerplate-back-go-2411/dto/response"
	"boilerplate-back-go-2411/errors"
	"boilerplate-back-go-2411/utils"

	"go.uber.org/zap"
)

type AuthService interface {
	Login(ctx context.Context, req request.LoginWeb) (response.Login, error)
	LoginLark(ctx context.Context, req request.LoginLark) (response.Login, error)
}

func (srv service) Login(ctx context.Context, req request.LoginWeb) (response.Login, error) {
	srv.logger.Info("🔑 Attempting login", zap.String("username", req.UserName))

	if req.UserName == "" || req.Password == "" {
		srv.logger.Warn("❌ Invalid login attempt: empty username or password")
		return response.Login{}, errors.ValidationError("username or password must not be empty")
	}

	// ✅ ดึงข้อมูล User พร้อม Role และ Password
	user, err := srv.userRepo.Login(ctx, req.UserName)
	if err != nil {
		srv.logger.Warn("❌ User not found", zap.String("username", req.UserName), zap.Error(err))
		return response.Login{}, errors.UnauthorizedError("invalid username or password")
	}

	// ✅ ตรวจสอบรหัสผ่าน
	hashedPassword := utils.HashPassword(req.Password)
	if hashedPassword != user.Password {
		srv.logger.Warn("❌ Invalid password", zap.String("username", req.UserName))
		return response.Login{}, errors.UnauthorizedError("invalid username or password")
	}

	// ✅ สร้าง Response
	loginResponse := response.Login{
		UserID:       user.UserID,
		UserName:     user.UserName,
		FullNameTH:   user.FullNameTH,
		NickName:     user.NickName,
		RoleID:       user.RoleID,
		RoleName:     user.RoleName,
		DepartmentNo: user.DepartmentNo,
		Platform:     "web",
	}

	srv.logger.Info("✅ Login successful", zap.String("username", req.UserName))
	return loginResponse, nil
}

func (srv service) LoginLark(ctx context.Context, req request.LoginLark) (response.Login, error) {
	logFinish := srv.logger.With(zap.String("username", req.UserName), zap.String("userID", req.UserID))
	logFinish.Info("🔑 Attempting login via Lark")

	if req.UserName == "" || req.UserID == "" {
		logFinish.Warn("❌ Invalid login attempt: empty username or userID")
		return response.Login{}, errors.ValidationError("username or userID must not be null")
	}

	user, err := srv.userRepo.LoginLark(ctx, req.UserID, req.UserName)
	if err != nil {
		logFinish.Warn("⚠️ User not found in Lark", zap.String("username", req.UserName), zap.String("userID", req.UserID), zap.Error(err))
		return response.Login{}, errors.UnauthorizedError("user not found in system")
	}

	loginResponse := response.Login{
		UserID:       user.UserID,
		UserName:     user.UserName,
		FullNameTH:   user.FullNameTH,
		NickName:     user.NickName,
		RoleID:       user.RoleID,
		RoleName:     user.RoleName,
		DepartmentNo: user.DepartmentNo,
		Platform:     "lark",
	}

	logFinish.Info("✅ Lark login successful", zap.String("username", user.UserName))
	return loginResponse, nil
}
