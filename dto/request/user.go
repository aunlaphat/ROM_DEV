package request

type Login struct {
	UserID   string `json:"userID" db:"UserID" example:"DC65060"`
	Password string `json:"password" db:"Password" example:"xxxxxxxx"`
}

type LoginWeb struct {
	UserName string `json:"userName" db:"Username" example:"eknarin.ler"`
	Password string `json:"password" db:"Password" example:"EKna1234"` // change password lastest in 17 January 2025
}

type LoginLark struct {
	UserID   string `json:"userID" db:"userID" example:"DC65060"`
	UserName string `json:"userName" db:"userName" example:"eknarin.ler"`
}

type LoginJWT struct {
	UserID   string `json:"userID" db:"UserID" example:"DC53002"`
	UserName string `json:"userName" db:"Username" example:"string"`
}

type User struct {
	UserID       string `json:"userID" db:"UserID" example:"DC65060"`
	Password     string `json:"password" db:"Password" example:"xxxxxxxx"`
	UserName     string `json:"userName" db:"UserName" example:"userName"`
	NickName     string `json:"nickName" db:"NickName" example:"Nickname"`
	FullNameTH   string `json:"fullNameTH" db:"FullNameTH" example:"Firstname Lastname"`
	DepartmentNo string `json:"departmentNo" db:"DepartmentNo" example:"G07"`
	RoleID       int    `json:"roleID" db:"RoleID" example:"1"`
	WarehouseID  int    `json:"warehouseID" db:"WarehouseID" example:"1"`
}
