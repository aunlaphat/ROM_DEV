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
	GetUser(ctx context.Context, username string) (response.User, error)
	GetUserFromLark(ctx context.Context, userID, username string) (response.User, error)
	GetUserWithPermission(ctx context.Context, userID, username string) (response.UserPermission, error)
}

func (srv service) Login(ctx context.Context, req request.LoginWeb) (response.User, error) {
	logFinish := srv.logger.With(zap.String("username", req.UserName))
	logFinish.Info("🔑 Attempting login")

	if req.UserName == "" || req.Password == "" {
		logFinish.Warn("❌ Invalid login attempt: empty username or password")
		return response.User{}, errors.ValidationError("username or password must not be null")
	}

	hasher := md5.New()
	hasher.Write([]byte(req.Password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	user, err := srv.userRepo.GetUser(ctx, req.UserName)
	if err != nil {
		logFinish.Error("❌ User not found", zap.Error(err))
		return response.User{}, errors.UnauthorizedError("invalid username or password")
	}

	if hashedPassword != req.Password {
		logFinish.Warn("❌ Invalid password")
		return response.User{}, errors.UnauthorizedError("invalid username or password")
	}

	logFinish.Info("✅ Login successful")
	return user, nil
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
		logFinish.Warn("❌ User not found in Lark", zap.Error(err))
		return response.User{}, errors.UnauthorizedError("user not found")
	}

	logFinish.Info("✅ Lark login successful")
	return user, nil
}

func (srv service) GetUser(ctx context.Context, username string) (response.User, error) {
	logFinish := srv.logger.With(zap.String("username", username))
	logFinish.Info("🔍 Fetching user")

	user, err := srv.userRepo.GetUser(ctx, username)
	if err != nil {
		logFinish.Error("❌ Failed to fetch user", zap.Error(err))
		return response.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	logFinish.Info("✅ User fetched successfully")
	return user, nil
}

func (srv service) GetUserFromLark(ctx context.Context, userID, username string) (response.User, error) {
	logFinish := srv.logger.With(zap.String("username", username), zap.String("userID", userID))
	logFinish.Info("🔍 Fetching user from Lark")

	user, err := srv.userRepo.GetUserFromLark(ctx, userID, username)
	if err != nil {
		logFinish.Error("❌ Failed to fetch user from Lark", zap.Error(err))
		return response.User{}, fmt.Errorf("failed to get user from Lark: %w", err)
	}

	logFinish.Info("✅ User fetched successfully from Lark")
	return user, nil
}

func (srv service) GetUserWithPermission(ctx context.Context, userID, username string) (response.UserPermission, error) {
	logFinish := srv.logger.With(zap.String("userID", userID), zap.String("username", username))
	logFinish.Info("🔍 Fetching user permissions")

	userPermission, err := srv.userRepo.GetUserWithPermission(ctx, userID, username)
	if err != nil {
		logFinish.Error("❌ Failed to fetch user permissions", zap.Error(err))
		return response.UserPermission{}, fmt.Errorf("failed to get user permissions: %w", err)
	}

	logFinish.Info("✅ User permissions fetched successfully")
	return userPermission, nil
}
