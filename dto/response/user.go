package response

import "time"

type Login struct {
	UserID       string `json:"userID" db:"UserID"`
	UserName     string `json:"userName" db:"UserName"`
	RoleID       int    `json:"roleID" db:"RoleID"`
	FullNameTH   string `json:"fullNameTH" db:"FullNameTH"`
	NickName     string `json:"nickName" db:"NickName"`
	DepartmentNo string `json:"departmentNo" db:"DepartmentNo"`
	Platform     string `json:"platform" db:"Platform"`
}

// ✅ 1️⃣ UserResponse - ใช้สำหรับการแสดงข้อมูลผู้ใช้
type UserResponse struct {
	UserID       string     `json:"userID"`                // รหัสพนักงาน
	UserName     string     `json:"userName"`              // ชื่อผู้ใช้
	NickName     string     `json:"nickName"`              // ชื่อเล่น
	FullNameTH   string     `json:"fullNameTH"`            // ชื่อเต็มภาษาไทย
	DepartmentNo string     `json:"department"`            // รหัสแผนก
	RoleID       int        `json:"roleID"`                // รหัสบทบาท
	RoleName     string     `json:"roleName"`              // ชื่อบทบาท
	Description  string     `json:"description"`           // คำอธิบาย Role
	Permission   string     `json:"permission"`            // สิทธิ์การเข้าถึง
	IsActive     bool       `json:"isActive"`              // สถานะบัญชี (Active/Inactive)
	LastLoginAt  *time.Time `json:"lastLoginAt,omitempty"` // เวลาล็อกอินล่าสุด (optional)
	CreatedAt    time.Time  `json:"createdAt"`             // เวลาสร้างบัญชี
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"`   // เวลาล่าสุดที่มีการอัปเดต (optional)
}

// ✅ 2️⃣ AddUserResponse - ใช้สำหรับแสดงข้อมูลหลังจากเพิ่มผู้ใช้สำเร็จ
type AddUserResponse struct {
	UserID    string `json:"userID"`    // รหัสพนักงานที่เพิ่ม
	RoleID    int    `json:"roleID"`    // รหัส Role ที่เพิ่มให้
	RoleName  string `json:"roleName"`  // ชื่อ Role
	Warehouse string `json:"warehouse"` // คลังสินค้า
	CreatedBy string `json:"createdBy"` // ผู้สร้าง
}

// ✅ 3️⃣ EditUserResponse - ใช้สำหรับแสดงข้อมูลหลังจากแก้ไขผู้ใช้สำเร็จ
type EditUserResponse struct {
	UserID      string `json:"userID"`      // รหัสพนักงานที่แก้ไข
	OldRoleID   int    `json:"oldRoleID"`   // Role ID เดิม
	NewRoleID   int    `json:"newRoleID"`   // Role ID ใหม่
	OldRoleName string `json:"oldRoleName"` // ชื่อ Role เดิม
	NewRoleName string `json:"newRoleName"` // ชื่อ Role ใหม่
	UpdatedBy   string `json:"updatedBy"`   // ผู้ที่ทำการแก้ไข
	UpdatedAt   string `json:"updatedAt"`   // เวลาที่อัปเดต
}

// ✅ 4️⃣ ResetPasswordResponse - ใช้สำหรับแสดงข้อมูลหลังจากเปลี่ยนรหัสผ่านสำเร็จ
type ResetPasswordResponse struct {
	UserID    string `json:"userID"`    // รหัสผู้ใช้ที่ถูกเปลี่ยนรหัสผ่าน
	UpdatedBy string `json:"updatedBy"` // ผู้ที่เปลี่ยนรหัสผ่าน
	Message   string `json:"message"`   // ข้อความแจ้งเตือน
}

// ✅ 5️⃣ DeleteUserResponse - ใช้สำหรับแสดงข้อมูลหลังจากลบผู้ใช้สำเร็จ
type DeleteUserResponse struct {
	UserID    string `json:"userID"`    // รหัสพนักงานที่ถูกลบ
	UpdatedBy string `json:"updatedBy"` // ผู้ที่ทำการลบ
	Message   string `json:"message"`   // ข้อความแจ้งเตือน
}
