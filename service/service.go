package service

import (
	"boilerplate-backend-go/logs"
	"boilerplate-backend-go/repository"

	"github.com/jmoiron/sqlx"
)

type service struct {
	userRepo repository.UserRepository //Repository of service
	logger   logs.Logger               //Logger of service
	constant repository.Constants
}
type AllOfService struct {
	User     UserService
	Constant Constants
}

func NewService(db *sqlx.DB, logger logs.Logger) AllOfService {
	repo := repository.NewDB(db)
	srv := service{
		userRepo: repo,
		logger:   logger,
		constant: repo,
	}
	return AllOfService{
		User:     srv,
		Constant: srv,
	}
}

type Login struct {
	UserID       string `json:"userID"`
	RoleID       int    `json:"roleID"`
	PermissionID string `json:"permissionID"`
	DeptNo       string `json:"deptNo"`
	NickName     string `json:"nickName"`
	FullNameTH   string `json:"fullNameTH"`
	FullNameEN   string `json:"fullNameEN"`
	Platfrom     string `json:"platfrom"`
}
