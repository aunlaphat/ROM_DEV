package entity

type User struct {
	UserID       string `db:"UserID" json:"userID"`
	UserName     string `db:"UserName" json:"userName"`
	RoleID       int    `db:"RoleID" json:"roleID"`
	PermissionID string `db:"PermissionID" json:"permissionID"`
	DeptNo       string `db:"DeptNo" json:"deptNo"`
	NickName     string `db:"NickName" json:"nickName"`
	FullNameTH   string `db:"FullNameTH" json:"fullNameTH"`
	FullNameEN   string `db:"FullNameEN" json:"fullNameEN"`
}

type DOM_V_User struct {
	UserID     string `json:"userID,omitempty" db:"UserID" example:"DC64205"`
	UserName   string `json:"userName,omitempty" db:"Username" example:"aunlaphat.art"`
	NickName   string `json:"nickName,omitempty" db:"NickName" example:"fa"`
	FullNameTH string `json:"fullNameTH,omitempty" db:"FullNameTH" example:"อัญญ์ลภัส อาจสุริยงค์"`
	DepartmentNo string `json:"department,omitempty" db:"DepartmentNo" example:"G01"`
}