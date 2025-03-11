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
	returnOrderRepo  repository.ReturnOrderRepository
	importOrderRepo  repository.ImportOrderRepository
	beforeReturnRepo repository.BeforeReturnRepository
	constant         repository.Constants
}
type AllOfService struct {
	Auth         AuthService
	User         UserService
	Order        OrderService
	DraftConfirm DraftConfirmService
	ReturnOrder  ReturnOrderService
	ImportOrder  ImportOrderService
	BeforeReturn BeforeReturnService
	Constant     Constants
}

func NewService(db *sqlx.DB, logger logs.Logger) AllOfService {
	repo := repository.NewDB(db)
	srv := service{
		logger:           logger,
		constantRepo:     repo,
		userRepo:         repo,
		orderRepo:        repo,
		draftConfirmRepo: repo,
		returnOrderRepo:  repo,
		importOrderRepo:  repo,
		beforeReturnRepo: repo,
		constant:         repo,
	}
	return AllOfService{
		Auth:         srv,
		User:         srv,
		Order:        srv,
		DraftConfirm: srv,
		ReturnOrder:  srv,
		ImportOrder:  srv,
		BeforeReturn: srv,
		Constant:     srv,
	}
}
