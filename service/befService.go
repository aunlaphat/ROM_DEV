package service

import (
	"boilerplate-backend-go/logs"
	"boilerplate-backend-go/repository"

	"github.com/jmoiron/sqlx"
)

type service struct {
	userRepo        repository.UserRepository //Repository of service
	returnOrderRepo repository.ReturnOrderRepository
	logger          logs.Logger //Logger of service
	constant        repository.Constants
}

type AllOfService struct {
	User        UserService
	ReturnOrder ReturnOrderService
	Constant    Constants
}

func NewService(db *sqlx.DB, logger logs.Logger) AllOfService {
	repo := repository.NewDB(db)
	srv := service{
		userRepo:        repo,
		returnOrderRepo: repo,
		logger:          logger,
		constant:        repo,
	}
	return AllOfService{
		User:        srv,
		ReturnOrder: srv,
		Constant:    srv,
	}
}
