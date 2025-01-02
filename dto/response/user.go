package response

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
