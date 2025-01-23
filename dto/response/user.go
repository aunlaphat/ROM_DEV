package response

type Login struct {
	UserID       string `json:"userID,omitempty" db:"UserID" example:"DC-XXXXX"`
	UserName     string `json:"userName,omitempty" db:"UserName" example:"userName"`
	RoleID       int    `json:"roleID,omitempty" db:"RoleID" example:"0"`
	FullNameTH   string `json:"fullNameTH,omitempty" db:"FullNameTH" example:"Firstname Lastname"`
	NickName     string `json:"nickName,omitempty" db:"NickName" example:"Nickname"`
	DepartmentNo string `json:"department,omitempty" db:"DepartmentNo" example:"G07"`
	Platform     string `json:"platform" db:"Platform" example:"Platform"`
}

type UserPermission struct {
	UserID       string `json:"userID,omitempty" db:"UserID"`
	UserName     string `json:"userName,omitempty" db:"UserName"`
	NickName     string `json:"nickName,omitempty" db:"NickName"`
	FullNameTH   string `json:"fullNameTH,omitempty" db:"FullNameTH"`
	DepartmentNo string `json:"departmentNo,omitempty" db:"DepartmentNo"`
	RoleID       int    `json:"roleID,omitempty" db:"RoleID"`
	RoleName     string `json:"roleName,omitempty" db:"RoleName"`
	Description  string `json:"description,omitempty" db:"Description"`
	Permission   string `json:"permission,omitempty" db:"Permission"`
}
