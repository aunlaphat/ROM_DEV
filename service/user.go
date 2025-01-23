package service

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"

	"go.uber.org/zap"
)

type UserService interface {
	Login(req request.LoginWeb) (response.Login, error)
	LoginLark(req request.LoginLark) (response.Login, error)
	GetUser(ctx context.Context, req request.LoginWeb) (response.Login, error)
	GetUserFromLark(ctx context.Context, username, password string) (response.Login, error)
	GetUserWithPermission(ctx context.Context, req request.LoginLark) (response.UserPermission, error)
}

func (srv service) Login(req request.LoginWeb) (response.Login, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(context.Background(), "Login", zap.String("username", req.UserName))
	defer logFinish("Completed", nil)

	res := response.Login{}
	if req.UserName == "" || req.Password == "" {
		logFinish("Failed", fmt.Errorf("username or password must not be null"))
		srv.logger.Warn("❌ Invalid login attempt: empty username or password", zap.String("username", req.UserName))
		return res, errors.ValidationError("username or password must not be null")
	}

	hasher := md5.New()
	hasher.Write([]byte(req.Password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))
	ctx := context.Background()
	user, err := srv.userRepo.GetUser(ctx, req.UserName, hashedPassword)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			logFinish("Not Found", nil)
			srv.logger.Warn("❌ No user found with provided credentials", zap.String("username", req.UserName))
			return res, errors.UnauthorizedError("username or password is not valid")
		default:
			logFinish("Failed", err)
			srv.logger.Error("❌ Unexpected error occurred while getting user", zap.Error(err))
			return res, errors.UnexpectedError()
		}
	}

	logFinish("Success", nil)
	srv.logger.Info("✅ Successfully logged in", zap.String("username", req.UserName))
	return user, nil
}

// Login: Lark
func (srv service) LoginLark(req request.LoginLark) (response.Login, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(context.Background(), "LoginLark", zap.String("username", req.UserName), zap.String("userID", req.UserID))
	defer logFinish("Completed", nil)

	srv.logger.Debug("🚀 Starting LoginLark", zap.String("username", req.UserName), zap.String("userID", req.UserID))

	res := response.Login{}
	if req.UserName == "" || req.UserID == "" {
		logFinish("Failed", fmt.Errorf("username or userid must not be null"))
		srv.logger.Warn("❌ Invalid login attempt: empty username or userID", zap.String("username", req.UserName), zap.String("userID", req.UserID))
		return res, errors.ValidationError("username or userid must not be null")
	}

	ctx := context.Background()
	srv.logger.Debug("Attempting to get user from Lark", zap.String("username", req.UserName), zap.String("userID", req.UserID))

	user, err := srv.userRepo.GetUserFromLark(ctx, req.UserID, req.UserName)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			logFinish("Not Found", nil)
			srv.logger.Warn("❌ No user found with provided Lark credentials", zap.String("username", req.UserName), zap.String("userID", req.UserID))
			return res, errors.UnauthorizedError("user not found in system")
		default:
			logFinish("Failed", err)
			srv.logger.Error("❌ Database error while getting user from Lark", zap.Error(err), zap.String("username", req.UserName), zap.String("userID", req.UserID))
			return res, errors.UnexpectedError()
		}
	}

	logFinish("Success", nil)
	srv.logger.Info("✅ Successfully logged in via Lark", zap.String("username", user.UserName), zap.String("userID", user.UserID))
	return user, nil
}

func (srv service) GetUser(ctx context.Context, req request.LoginWeb) (response.Login, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "GetUser", zap.String("username", req.UserName))
	defer logFinish("Completed", nil)

	srv.logger.Debug("🚀 Starting GetUser", zap.String("username", req.UserName))

	user, err := srv.userRepo.GetUser(ctx, req.UserName, req.Password)
	if err != nil {
		logFinish("Failed", err)
		srv.logger.Error("❌ Failed to get user", zap.Error(err))
		return response.Login{}, fmt.Errorf("failed to get user: %w", err)
	}

	logFinish("Success", nil)
	srv.logger.Debug("✅ Successfully retrieved user", zap.String("username", req.UserName))
	return user, nil
}

func (srv service) GetUserFromLark(ctx context.Context, userID, username string) (response.Login, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "GetUserFromLark", zap.String("username", username), zap.String("userID", userID))
	defer logFinish("Completed", nil)

	srv.logger.Debug("🚀 Starting GetUserFromLark", zap.String("username", username))

	user, err := srv.userRepo.GetUserFromLark(ctx, userID, username)
	if err != nil {
		logFinish("Failed", err)
		srv.logger.Error("❌ Failed to get user from Lark", zap.Error(err))
		return response.Login{}, fmt.Errorf("failed to get user from Lark: %w", err)
	}

	logFinish("Success", nil)
	srv.logger.Debug("✅ Successfully retrieved user from Lark", zap.String("username", username))
	return user, nil
}

func (srv service) GetUserWithPermission(ctx context.Context, req request.LoginLark) (response.UserPermission, error) {
	// เริ่มต้น Logging ของ API Call
	logFinish := srv.logger.LogAPICall(ctx, "GetUserWithPermission", zap.String("username", req.UserID))
	defer logFinish("Completed", nil)

	srv.logger.Debug("🚀 Starting GetUserWithPermission", zap.String("username", req.UserID))

	user, err := srv.userRepo.GetUserWithPermission(ctx, req.UserID, req.UserName)
	if err != nil {
		logFinish("Failed", err)
		srv.logger.Error("❌ Failed to get user with permission", zap.Error(err))
		return response.UserPermission{}, fmt.Errorf("failed to get user: %w", err)
	}

	logFinish("Success", nil)
	srv.logger.Debug("✅ Successfully retrieved user with permission", zap.String("username", req.UserID))
	return user, nil
}
