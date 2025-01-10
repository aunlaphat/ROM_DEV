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
	Login(req request.Login) (response.Login, error)
	LoginLark(req request.LoginLark) (response.Login, error)
	GetUser(ctx context.Context, username, password string) (response.Login, error)
	GetUserWithPermission(ctx context.Context, username, password string) (response.UserPermission, error)
}

func (srv service) Login(req request.Login) (response.Login, error) {
	res := response.Login{}
	if req.UserID == "" || req.Password == "" {
		return res, errors.ValidationError("userid or password must not be null")
	}

	hasher := md5.New()
	hasher.Write([]byte(req.Password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))
	ctx := context.Background()
	user, err := srv.userRepo.GetUser(ctx, req.UserID, hashedPassword)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			srv.logger.Warn("❌ No user found with provided credentials", zap.String("userid", req.UserID))
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
	res := response.Login{}
	if req.UserName == "" || req.UserID == "" {
		return res, errors.ValidationError("username or password must not be null")
	}

	user, err := srv.userRepo.GetUserFromLark(req.UserName, req.UserID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return res, errors.UnauthorizedError("username or password is not valid")
		default:
			srv.logger.Error(err)
			return res, errors.UnexpectedError()
		}
	}
	return user, nil
}

func (srv service) GetUser(ctx context.Context, username, password string) (response.Login, error) {
	srv.logger.Debug("🚀 Starting GetUser", zap.String("username", username))

	user, err := srv.userRepo.GetUser(ctx, username, password)
	if (err != nil) {
		srv.logger.Error("❌ Failed to get user", zap.Error(err))
		return response.Login{}, fmt.Errorf("failed to get user: %w", err)
	}

	srv.logger.Debug("✅ Successfully retrieved user", zap.String("username", username))
	return user, nil
}

func (srv service) GetUserWithPermission(ctx context.Context, username, password string) (response.UserPermission, error) {
	srv.logger.Debug("🚀 Starting GetUser", zap.String("username", username))

	user, err := srv.userRepo.GetUserWithPermission(ctx, username, password)
	if (err != nil) {
		srv.logger.Error("❌ Failed to get user", zap.Error(err))
		return response.UserPermission{}, fmt.Errorf("failed to get user: %w", err)
	}

	srv.logger.Debug("✅ Successfully retrieved user", zap.String("username", username))
	return user, nil
}
