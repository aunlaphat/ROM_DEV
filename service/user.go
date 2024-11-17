package service

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

type UserService interface {
	Login(req request.LoginWeb) (response.Login, error)
	LoginLark(req request.LoginLark) (response.Login, error)
}

func (srv service) Login(req request.LoginWeb) (response.Login, error) {
	res := response.Login{}
	if req.UserName == "" || req.Password == "" {
		return res, errors.ValidationError("username or password must not be null")
	}

	hasher := md5.New()
	hasher.Write([]byte(req.Password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))
	res, err := srv.userRepo.GetUser(req.UserName, hashedPassword)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return res, errors.UnauthorizedError("username or password is not valid")
		default:
			srv.logger.Error(err)
			return res, errors.UnexpectedError()
		}
	}
	return res, nil
}

// Login: Lark
func (srv service) LoginLark(req request.LoginLark) (response.Login, error) {
	res := response.Login{}
	if req.UserName == "" || req.UserID == "" {
		return res, errors.ValidationError("username or password must not be null")
	}

	res, err := srv.userRepo.GetUserFromLark(req.UserName, req.UserID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return res, errors.UnauthorizedError("username or password is not valid")
		default:
			srv.logger.Error(err)
			return res, errors.UnexpectedError()
		}
	}
	return res, nil
}
