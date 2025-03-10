package request

type Login struct {
	UserID   string `json:"userID" db:"UserID" example:"DC65060"`
	Password string `json:"password" db:"Password" example:"xxxxxxxx"`
}

type LoginWeb struct {
	UserName string `json:"userName" db:"UserName" example:"eknarin.ler"`
	Password string `json:"password" db:"Password" example:"EKna1234"` // change password lastest in 17 January 2025
}

type LoginLark struct {
	UserID   string `json:"userID" db:"UserID" example:"DC65060"`
	UserName string `json:"userName" db:"UserName" example:"eknarin.ler"`
}

type LoginJWT struct {
	UserID   string `json:"userID" db:"UserID" example:"DC53002"`
	UserName string `json:"userName" db:"UserName" example:"string"`
}

type AddUserRequest struct {
	UserID      string `json:"userID" binding:"required"`      // รหัสผู้ใช้
	RoleID      int    `json:"roleID" binding:"required"`      // รหัส Role
	WarehouseID string `json:"warehouseID" binding:"required"` // คลังสินค้า
}

type EditUserRequest struct {
	UserID      string  `json:"userID" binding:"required"` // รหัสผู้ใช้ที่ต้องการแก้ไข
	RoleID      *int    `json:"roleID,omitempty"`          // รหัส Role ใหม่
	WarehouseID *string `json:"warehouseID,omitempty"`     // คลังสินค้าใหม่
}
