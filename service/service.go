package service

import (
	"boilerplate-backend-go/logs"
	"boilerplate-backend-go/repository"

	"github.com/jmoiron/sqlx"
)

type service struct {
	userRepo        repository.UserRepository //Repository of service
	befRORepo repository.BefRORepository
	logger          logs.Logger //Logger of service
	constant        repository.Constants
}

type AllOfService struct {
	User        UserService
	BefRO BefROService
	Constant    Constants
}

func NewService(db *sqlx.DB, logger logs.Logger) AllOfService {
	repo := repository.NewDB(db)
	srv := service{
		userRepo:        repo,
		befRORepo: repo,
		logger:          logger,
		constant:        repo,
	}
	return AllOfService{
		User:        srv,
		BefRO: srv,
		Constant:    srv,
	}
}
