package service

import (
	"context"

	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"boilerplate-backend-go/utils"

	"go.uber.org/zap"
)

type AuthService interface {
	Login(ctx context.Context, req request.LoginWeb) (response.Login, error)
	LoginLark(ctx context.Context, req request.LoginLark) (response.Login, error)
}

// ✅ 1️⃣ Login (เข้าสู่ระบบ)
func (srv service) Login(ctx context.Context, req request.LoginWeb) (response.Login, error) {
	srv.logger.Info("🔑 Attempting login", zap.String("username", req.UserName))

	if req.UserName == "" || req.Password == "" {
		srv.logger.Warn("❌ Invalid login attempt: empty username or password")
		return response.Login{}, errors.ValidationError("username or password must not be empty")
	}

	hashedPassword := utils.HashPassword(req.Password)

	user, err := srv.userRepo.GetUser(ctx, req.UserName)
	if err != nil {
		srv.logger.Warn("❌ User not found", zap.String("username", req.UserName))
		return response.Login{}, errors.UnauthorizedError("invalid username or password")
	}

	if hashedPassword != user.Password {
		srv.logger.Warn("❌ Invalid password", zap.String("username", req.UserName))
		return response.Login{}, errors.UnauthorizedError("invalid username or password")
	}

	userResponse := response.Login{
		UserID:       user.UserID,
		UserName:     user.UserName,
		RoleID:       user.RoleID,
		FullNameTH:   user.FullNameTH,
		NickName:     user.NickName,
		DepartmentNo: user.DepartmentNo,
		Platform:     "web",
	}

	srv.logger.Info("✅ Login successful", zap.String("username", req.UserName))
	return userResponse, nil
}

func (srv service) LoginLark(ctx context.Context, req request.LoginLark) (response.Login, error) {
	logFinish := srv.logger.With(zap.String("username", req.UserName), zap.String("userID", req.UserID))
	logFinish.Info("🔑 Attempting login via Lark")

	if req.UserName == "" || req.UserID == "" {
		logFinish.Warn("❌ Invalid login attempt: empty username or userID")
		return response.Login{}, errors.ValidationError("username or userID must not be null")
	}

	user, err := srv.userRepo.GetUserFromLark(ctx, req.UserID, req.UserName)
	if err != nil {
		logFinish.Warn("⚠️ User not found in Lark", zap.String("username", req.UserName), zap.String("userID", req.UserID), zap.Error(err))
		return response.Login{}, errors.UnauthorizedError("user not found in system")
	}

	user.Platform = "lark"

	logFinish.Info("✅ Lark login successful", zap.String("username", user.UserName))
	return user, nil
}
