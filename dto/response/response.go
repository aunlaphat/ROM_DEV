package response

type Login struct {
	UserName   string `json:"userName,omitempty" db:"UserName" example:"userName"`
	UserID     string `json:"userID,omitempty" db:"UserID" example:"userID"`
	RoleID     int    `json:"roleID,omitempty" db:"RoleID" example:"1"`
	FullNameTH string `json:"fullNameTH,omitempty" db:"FullNameTH" example:"test1234"`
	NickName   string `json:"nickName,omitempty" db:"NickName" example:"test1234"`
	Platfrom   string `json:"platfrom" db:"Platfrom" example:"test1234"`
}
