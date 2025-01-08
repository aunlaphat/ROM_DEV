package service

import (
	"boilerplate-backend-go/logs"
	"boilerplate-backend-go/repository"

	"github.com/jmoiron/sqlx"
)

type service struct {
	userRepo  repository.UserRepository //Repository of service
	logger    logs.Logger               //Logger of service
	constant  repository.Constants
	befRORepo repository.BefRORepository
}
type AllOfService struct {
	User     UserService
	Constant Constants
	BefRO    BefROService
	// Login	    LoginService
}

func NewService(db *sqlx.DB, logger logs.Logger) AllOfService {
	repo := repository.NewDB(db)
	srv := service{
		userRepo:  repo,
		logger:    logger,
		constant:  repo,
		befRORepo: repo,
	}
	return AllOfService{
		User:     srv,
		Constant: srv,
		BefRO:    srv,
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
