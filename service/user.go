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
	GetUser(ctx context.Context, username, password string) (response.Login, error)
	GetUserFromLark(ctx context.Context, username, password string) (response.Login, error)
	GetUserWithPermission(ctx context.Context, username, password string) (response.UserPermission, error)
}

func (srv service) Login(req request.LoginWeb) (response.Login, error) {
	res := response.Login{}
	if req.UserName == "" || req.Password == "" {
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
			srv.logger.Warn("❌ No user found with provided credentials", zap.String("username", req.UserName))
			return res, errors.UnauthorizedError("username or password is not valid")
		default:
			srv.logger.Error("❌ Unexpected error occurred while getting user", zap.Error(err))
			return res, errors.UnexpectedError()
		}
	}
	return user, nil
}

// Login: Lark
func (srv service) LoginLark(req request.LoginLark) (response.Login, error) {
    srv.logger.Debug("🚀 Starting LoginLark", 
        zap.String("username", req.UserName),
        zap.String("userID", req.UserID))

    res := response.Login{}
    if req.UserName == "" || req.UserID == "" {
        srv.logger.Warn("❌ Invalid login attempt: empty username or userID",
            zap.String("username", req.UserName),
            zap.String("userID", req.UserID))
        return res, errors.ValidationError("username or userid must not be null")
    }

    ctx := context.Background()
    srv.logger.Debug("Attempting to get user from Lark",
        zap.String("username", req.UserName),
        zap.String("userID", req.UserID))

    user, err := srv.userRepo.GetUserFromLark(ctx, req.UserID, req.UserName)
    if err != nil {
        switch {
        case err == sql.ErrNoRows:
            srv.logger.Warn("❌ No user found with provided Lark credentials",
                zap.String("username", req.UserName),
                zap.String("userID", req.UserID))
            return res, errors.UnauthorizedError("user not found in system")
        default:
            srv.logger.Error("❌ Database error while getting user from Lark",
                zap.Error(err),
                zap.String("username", req.UserName),
                zap.String("userID", req.UserID))
            return res, errors.UnexpectedError()
        }
    }

    srv.logger.Info("✅ Successfully logged in via Lark",
        zap.String("username", user.UserName),
        zap.String("userID", user.UserID))
    return user, nil
}

func (srv service) GetUser(ctx context.Context, username, password string) (response.Login, error) {
	srv.logger.Debug("🚀 Starting GetUser", zap.String("username", username))

	user, err := srv.userRepo.GetUser(ctx, username, password)
	if err != nil {
		srv.logger.Error("❌ Failed to get user", zap.Error(err))
		return response.Login{}, fmt.Errorf("failed to get user: %w", err)
	}

	srv.logger.Debug("✅ Successfully retrieved user", zap.String("username", username))
	return user, nil
}

func (srv service) GetUserFromLark(ctx context.Context, userID, username string) (response.Login, error) {
	srv.logger.Debug("🚀 Starting GetUserFromLark", zap.String("username", username))

	user, err := srv.userRepo.GetUserFromLark(ctx, userID, username)
	if err != nil {
		srv.logger.Error("❌ Failed to get user from Lark", zap.Error(err))
		return response.Login{}, fmt.Errorf("failed to get user from Lark: %w", err)
	}

	srv.logger.Debug("✅ Successfully retrieved user from Lark", zap.String("username", username))
	return user, nil
}

func (srv service) GetUserWithPermission(ctx context.Context, username, password string) (response.UserPermission, error) {
	srv.logger.Debug("🚀 Starting GetUser", zap.String("username", username))

	user, err := srv.userRepo.GetUserWithPermission(ctx, username, password)
	if err != nil {
		srv.logger.Error("❌ Failed to get user", zap.Error(err))
		return response.UserPermission{}, fmt.Errorf("failed to get user: %w", err)
	}

	srv.logger.Debug("✅ Successfully retrieved user", zap.String("username", username))
	return user, nil
}
