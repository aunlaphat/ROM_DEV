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

// ✅ 1️⃣ AddUserRequest - ใช้สำหรับเพิ่มผู้ใช้ใหม่
type AddUserRequest struct {
	UserID      string `json:"userID" binding:"required"`    // รหัสผู้ใช้
	RoleID      int    `json:"roleID" binding:"required"`    // รหัส Role
	WarehouseID string `json:"warehouseID" binding:"required"` // คลังสินค้า
}

// ✅ 2️⃣ EditUserRequest - ใช้สำหรับแก้ไขข้อมูลผู้ใช้
type EditUserRequest struct {
	RoleID       int    `json:"roleID" binding:"required"`    // รหัส Role ใหม่
	RoleName     string `json:"roleName"`                     // ชื่อ Role ใหม่ (optional)
	OldRole      string `json:"oldRole"`                      // ชื่อ Role เดิม (optional)
	OldWarehouse string `json:"oldWarehouse"`                 // คลังสินค้าเดิม (optional)
	Warehouse    string `json:"warehouse" binding:"required"` // คลังสินค้าใหม่
}

// ✅ 3️⃣ ResetPasswordRequest - ใช้สำหรับเปลี่ยนรหัสผ่าน
type ResetPasswordRequest struct {
	UserID      string `json:"userID" binding:"required"`            // รหัสผู้ใช้
	NewPassword string `json:"newPassword" binding:"required,min=8"` // รหัสผ่านใหม่ (ต้องมีอย่างน้อย 8 ตัวอักษร)
}

// ✅ 4️⃣ DeleteUserRequest - ใช้สำหรับลบผู้ใช้ (Soft Delete)
type DeleteUserRequest struct {
	UserID string `json:"userID" binding:"required"` // รหัสผู้ใช้
}
