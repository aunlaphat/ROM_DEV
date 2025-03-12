package service

import (
	"boilerplate-backend-go/logs"
	"boilerplate-backend-go/repository"

	"github.com/jmoiron/sqlx"
)

type service struct {
	logger           logs.Logger
	constantRepo     repository.ConstantRepository
	userRepo         repository.UserRepository
	orderRepo        repository.OrderRepository
	draftConfirmRepo repository.DraftConfirmRepository
}
type AllOfService struct {
	Constant     ConstantService
	Auth         AuthService
	User         UserService
	Order        OrderService
	DraftConfirm DraftConfirmService
}

func NewService(db *sqlx.DB, logger logs.Logger) AllOfService {
	repo := repository.NewDB(db)
	srv := service{
		logger:           logger,
		constantRepo:     repo,
		userRepo:         repo,
		orderRepo:        repo,
		draftConfirmRepo: repo,
	}
	return AllOfService{
		Constant:     srv,
		Auth:         srv,
		User:         srv,
		Order:        srv,
		DraftConfirm: srv,
	}
}
