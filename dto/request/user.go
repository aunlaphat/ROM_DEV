package request

type Login struct {
	UserID   string `json:"userID" db:"UserID"`
	Password string `json:"password" db:"Password"`
}

type LoginWeb struct {
	UserName string `json:"userName" db:"UserName"`
	Password string `json:"password" db:"Password"`
}

type LoginLark struct {
	UserID   string `json:"userID" db:"UserID"`
	UserName string `json:"userName" db:"UserName"`
}

type LoginJWT struct {
	UserID   string `json:"userID" db:"UserID"`
	UserName string `json:"userName" db:"UserName"`
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
