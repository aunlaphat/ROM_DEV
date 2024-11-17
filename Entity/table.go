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
