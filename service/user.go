package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"

	"go.uber.org/zap"
)

type UserService interface {
	Login(ctx context.Context, req request.LoginWeb) (response.User, error)
	LoginLark(ctx context.Context, req request.LoginLark) (response.User, error)
	GetUser(ctx context.Context, username string) (response.UserRole, error)
}

func (srv service) Login(ctx context.Context, req request.LoginWeb) (response.User, error) {
	logFinish := srv.logger.With(zap.String("username", req.UserName))
	logFinish.Info("🔑 Attempting login")

	if req.UserName == "" || req.Password == "" {
		logFinish.Warn("❌ Invalid login attempt: empty username or password")
		return response.User{}, fmt.Errorf("username or password must not be null")
	}

	hasher := md5.New()
	hasher.Write([]byte(req.Password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	user, err := srv.userRepo.GetUser(ctx, req.UserName)
	if err != nil {
		logFinish.Warn("❌ User not found", zap.String("username", req.UserName))
		return response.User{}, fmt.Errorf("invalid username or password")
	}

	if hashedPassword != user.Password {
		logFinish.Warn("❌ Invalid password", zap.String("username", req.UserName))
		return response.User{}, fmt.Errorf("invalid username or password")
	}

	userResponse := response.User{
		UserID:       user.UserID,
		UserName:     user.UserName,
		RoleID:       user.RoleID,
		FullNameTH:   user.FullNameTH,
		NickName:     user.NickName,
		DepartmentNo: user.DepartmentNo,
		Platform:     "web",
	}

	logFinish.Info("✅ Login successful", zap.String("username", req.UserName))
	return userResponse, nil
}

func (srv service) LoginLark(ctx context.Context, req request.LoginLark) (response.User, error) {
	logFinish := srv.logger.With(zap.String("username", req.UserName), zap.String("userID", req.UserID))
	logFinish.Info("🔑 Attempting login via Lark")

	if req.UserName == "" || req.UserID == "" {
		logFinish.Warn("❌ Invalid login attempt: empty username or userID")
		return response.User{}, errors.ValidationError("username or userID must not be null")
	}

	user, err := srv.userRepo.GetUserFromLark(ctx, req.UserID, req.UserName)
	if err != nil {
		logFinish.Warn("⚠️ User not found in Lark", zap.String("username", req.UserName), zap.String("userID", req.UserID), zap.Error(err))
		return response.User{}, errors.UnauthorizedError("user not found in system")
	}

	user.Platform = "lark"

	logFinish.Info("✅ Lark login successful", zap.String("username", user.UserName))
	return user, nil
}

func (srv service) GetUser(ctx context.Context, username string) (response.UserRole, error) {
	logFinish := srv.logger.With(zap.String("username", username))
	logFinish.Info("🔍 Fetching user credentials")

	user, err := srv.userRepo.GetUser(ctx, username)
	if err != nil {
		if err.Error() == "user not found" {
			logFinish.Warn("❌ User not found", zap.String("username", username))
			return response.UserRole{}, fmt.Errorf("user not found")
		}
		logFinish.Error("❌ Failed to fetch user", zap.Error(err))
		return response.UserRole{}, fmt.Errorf("database error")
	}

	userResponse := response.UserRole{
		UserID:       user.UserID,
		UserName:     user.UserName,
		FullNameTH:   user.FullNameTH,
		NickName:     user.NickName,
		DepartmentNo: user.DepartmentNo,
		RoleID:       user.RoleID,
		RoleName:     user.RoleName,
		Description:  user.Description,
		Permission:   user.Permission,
	}

	logFinish.Info("✅ User credentials fetched successfully", zap.String("username", username))
	return userResponse, nil
}
