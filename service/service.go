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
	returnOrderRepo  repository.ReturnOrderRepository
	importOrderRepo  repository.ImportOrderRepository
	beforeReturnRepo repository.BeforeReturnRepository
	constant         repository.Constants
	constantRepo     repository.ConstantRepository
	draftConfirmRepo repository.DraftConfirmRepository
}
type AllOfService struct {
	Auth         AuthService
	User         UserService
	Order        OrderService
	ReturnOrder  ReturnOrderService
	ImportOrder  ImportOrderService
	BeforeReturn BeforeReturnService
	Constant     Constants
	DraftConfirm DraftConfirmService
}

func NewService(db *sqlx.DB, logger logs.Logger) AllOfService {
	repo := repository.NewDB(db)
	srv := service{
		logger:           logger,
		userRepo:         repo,
		orderRepo:        repo,
		returnOrderRepo:  repo,
		importOrderRepo:  repo,
		beforeReturnRepo: repo,
		constant:         repo,
		constantRepo:     repo,
		draftConfirmRepo: repo,
	}
	return AllOfService{
		Auth:         srv,
		User:         srv,
		Order:        srv,
		ReturnOrder:  srv,
		ImportOrder:  srv,
		BeforeReturn: srv,
		Constant:     srv,
		DraftConfirm: srv,
	}
}
