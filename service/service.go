package service

import (
	"boilerplate-backend-go/logs"
	"boilerplate-backend-go/repository"

	"github.com/jmoiron/sqlx"
)

type service struct {
	userRepo  repository.UserRepository //Repository of service
	logger    logs.Logger               //Logger of service
	orderRepo repository.OrderRepository
}
type AllOfService struct {
	User  UserService
	Order OrderService
}

func NewService(db *sqlx.DB, logger logs.Logger) AllOfService {
	repo := repository.NewDB(db)
	srv := service{
		userRepo:  repo,
		logger:    logger,
		orderRepo: repo,
	}
	return AllOfService{
		User:  srv,
		Order: srv,
	}
}
