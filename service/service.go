package service

import (
	"boilerplate-backend-go/logs"
	"boilerplate-backend-go/repository"

	"github.com/jmoiron/sqlx"
)

type service struct {
	userRepo         repository.UserRepository //Repository of service
	logger           logs.Logger               //Logger of service
	constant         repository.Constants
	beforeReturnRepo repository.BeforeReturnRepository
	returnOrderRepo  repository.ReturnOrderRepository
}
type AllOfService struct {
	User         UserService
	Constant     Constants
	BeforeReturn BeforeReturnService
	ReturnOrder  ReturnOrderService
}

func NewService(db *sqlx.DB, logger logs.Logger) AllOfService {
	repo := repository.NewDB(db)
	srv := service{
		userRepo:         repo,
		logger:           logger,
		constant:         repo,
		beforeReturnRepo: repo,
		returnOrderRepo:  repo,
	}
	return AllOfService{
		User:         srv,
		Constant:     srv,
		BeforeReturn: srv,
		ReturnOrder:  srv,
	}
}
