package response

type User struct {
	UserID       string `json:"userID" db:"UserID"`
	UserName     string `json:"userName" db:"UserName"`
	RoleID       int    `json:"roleID" db:"RoleID"`
	FullNameTH   string `json:"fullNameTH" db:"FullNameTH"`
	NickName     string `json:"nickName" db:"NickName"`
	DepartmentNo string `json:"departmentNo" db:"DepartmentNo"`
	Platform     string `json:"platform" db:"Platform"`
}

type UserRole struct {
	UserID       string `json:"userID" db:"UserID"`
	UserName     string `json:"userName" db:"UserName"`
	FullNameTH   string `json:"fullNameTH" db:"FullNameTH"`
	NickName     string `json:"nickName" db:"NickName"`
	DepartmentNo string `json:"departmentNo" db:"DepartmentNo"`
	RoleID       int    `json:"roleID" db:"RoleID"`
	RoleName     string `json:"roleName" db:"RoleName"`
	Description  string `json:"description" db:"Description"`
	Permission   string `json:"permission" db:"Permission"`
}
