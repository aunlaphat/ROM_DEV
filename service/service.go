package service

import (
	"boilerplate-backend-go/logs"
	"boilerplate-backend-go/repository"

	"github.com/jmoiron/sqlx"
)

type service struct {
	logger           logs.Logger
	userRepo         repository.UserRepository
	orderRepo        repository.OrderRepository
	draftConfirmRepo repository.DraftConfirmRepository
}
type AllOfService struct {
	Auth         AuthService
	User         UserService
	Order        OrderService
	DraftConfirm DraftConfirmService
}

func NewService(db *sqlx.DB, logger logs.Logger) AllOfService {
	repo := repository.NewDB(db)
	srv := service{
		logger:           logger,
		userRepo:         repo,
		orderRepo:        repo,
		draftConfirmRepo: repo,
	}
	return AllOfService{
		Auth:         srv,
		User:         srv,
		Order:        srv,
		DraftConfirm: srv,
	}
}
