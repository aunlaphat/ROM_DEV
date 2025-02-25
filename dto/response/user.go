package response

import "time"

type Login struct {
	UserID       string `json:"userID" db:"UserID"`
	UserName     string `json:"userName" db:"UserName"`
	RoleID       int    `json:"roleID" db:"RoleID"`
	RoleName     string `json:"roleName" db:"RoleName"`
	FullNameTH   string `json:"fullNameTH" db:"FullNameTH"`
	NickName     string `json:"nickName" db:"NickName"`
	DepartmentNo string `json:"departmentNo" db:"DepartmentNo"`
	Platform     string `json:"platform" db:"Platform"`
}

type UserResponse struct {
	UserID        string     `json:"userID"`                // รหัสพนักงาน
	UserName      string     `json:"userName"`              // ชื่อผู้ใช้
	NickName      string     `json:"nickName"`              // ชื่อเล่น
	FullNameTH    string     `json:"fullNameTH"`            // ชื่อเต็มภาษาไทย
	DepartmentNo  string     `json:"departmentNo"`          // รหัสแผนก
	RoleID        int        `json:"roleID"`                // รหัสบทบาท
	RoleName      string     `json:"roleName"`              // ชื่อบทบาท
	WarehouseID   int        `json:"warehouseID"`           // รหัสคลังสินค้า
	WarehouseName string     `json:"warehouseName"`         // ชื่อคลังสินค้า
	Description   string     `json:"description"`           // คำอธิบาย Role
	IsActive      bool       `json:"isActive"`              // สถานะบัญชี (Active/Inactive)
	LastLoginAt   *time.Time `json:"lastLoginAt,omitempty"` // เวลาล็อกอินล่าสุด (optional)
	CreatedAt     time.Time  `json:"createdAt"`             // เวลาสร้างบัญชี
	UpdatedAt     *time.Time `json:"updatedAt,omitempty"`   // เวลาล่าสุดที่มีการอัปเดต (optional)
}

type AddUserResponse struct {
	UserID      string `json:"userID"`      // รหัสพนักงานที่เพิ่ม
	RoleID      int    `json:"roleID"`      // รหัส Role ที่เพิ่มให้
	WarehouseID string `json:"warehouseID"` // รหัสคลังสินค้า
	CreatedBy   string `json:"createdBy"`   // ผู้สร้าง
}

type EditUserResponse struct {
	UserID        string    `json:"userID"`                // รหัสพนักงานที่แก้ไข
	RoleID        *int      `json:"roleID,omitempty"`      // Role ID เดิม
	RoleName      string    `json:"roleName"`              // ชื่อ Role เดิม
	WarehouseID   *int      `json:"warehouseID,omitempty"` // รหัสคลังสินค้าเดิม
	WarehouseName string    `json:"warehouseName"`         // ชื่อคลังสินค้าเดิม
	UpdatedBy     string    `json:"updatedBy"`             // ผู้ที่ทำการแก้ไข
	UpdatedAt     time.Time `json:"updatedAt"`             // เวลาที่อัปเดตล่าสุด
}

type DeleteUserResponse struct {
	UserID        string    `json:"userID"`        // รหัสพนักงานที่ถูกลบ
	UserName      string    `json:"userName"`      // ชื่อผู้ใช้ที่ถูกลบ
	RoleID        int       `json:"roleID"`        // Role ID ปัจจุบัน
	RoleName      string    `json:"roleName"`      // ชื่อ Role ปัจจุบัน
	WarehouseID   int       `json:"warehouseID"`   // รหัสคลังสินค้าปัจจุบัน
	WarehouseName string    `json:"warehouseName"` // ชื่อคลังสินค้าปัจจุบัน
	DeactivatedBy string    `json:"deactivatedBy"` // ผู้ที่ทำการลบ
	DeactivatedAt time.Time `json:"deactivatedAt"` // เวลาที่ลบล่าสุด
	Message       string    `json:"message"`       // ข้อความแจ้งเตือน
}
