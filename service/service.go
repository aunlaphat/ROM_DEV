package service

import (
	"boilerplate-backend-go/logs"
	"boilerplate-backend-go/repository"

	"github.com/jmoiron/sqlx"
)

type service struct {
	logger    logs.Logger
	userRepo  repository.UserRepository
	orderRepo repository.OrderRepository
}
type AllOfService struct {
	User  UserService
	Order OrderService
}

func NewService(db *sqlx.DB, logger logs.Logger) AllOfService {
	repo := repository.NewDB(db)
	srv := service{
		logger:    logger,
		userRepo:  repo,
		orderRepo: repo,
	}
	return AllOfService{
		User:  srv,
		Order: srv,
	}
}
