package service

import (
	"boilerplate-backend-go/logs"
	"boilerplate-backend-go/repository"

	"github.com/jmoiron/sqlx"
)

type service struct {
	userRepo         repository.UserRepository //Repository of service
	//usertestRepo     repository.UserTestRepository
	logger           logs.Logger //Logger of service
	constant         repository.Constants
	//befRORepo        repository.BefRORepository
	returnOrderRepo  repository.ReturnOrderRepository
	importOrderRepo  repository.ImportOrderRepository
	beforeReturnRepo repository.BeforeReturnRepository
}
type AllOfService struct {
	User         UserService
	//UserTest     UserTestService
	Constant     Constants
	//BefRO        BefROService
	ReturnOrder  ReturnOrderService
	ImportOrder  ImportOrderService
	BeforeReturn BeforeReturnService
	// Login	    LoginService
}

func NewService(db *sqlx.DB, logger logs.Logger) AllOfService {
	repo := repository.NewDB(db)
	srv := service{
		userRepo:         repo,
		//usertestRepo:     repo,
		logger:           logger,
		constant:         repo,
		//befRORepo:        repo,
		returnOrderRepo:  repo,
		importOrderRepo:  repo,
		beforeReturnRepo: repo,
	}
	return AllOfService{
		User:         srv,
		//UserTest:     srv,
		Constant:     srv,
		//BefRO:        srv,
		ReturnOrder:  srv,
		ImportOrder:  srv,
		BeforeReturn: srv,
	}
}

/* type Login struct {
	UserID       string `json:"userID"`
	RoleID       int    `json:"roleID"`
	PermissionID string `json:"permissionID"`
	DeptNo       string `json:"deptNo"`
	NickName     string `json:"nickName"`
	FullNameTH   string `json:"fullNameTH"`
	FullNameEN   string `json:"fullNameEN"`
	Platfrom     string `json:"platfrom"`
} */
