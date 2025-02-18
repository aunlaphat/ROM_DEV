package entity

import "time"

type ROM_V_User struct {
	UserID       string `db:"UserID"`       // รหัสพนักงานจาก ERP
	UserName     string `db:"UserName"`     // ชื่อผู้ใช้
	Password     string `db:"Password"`     // (MD5 Hash) ใช้สำหรับยืนยันตัวตน
	NickName     string `db:"NickName"`     // ชื่อเล่น
	FullNameTH   string `db:"FullNameTH"`   // ชื่อเต็มภาษาไทย
	DepartmentNo string `db:"DepartmentNo"` // รหัสแผนก
}

type ROM_V_UserDetail struct {
	UserID       string `db:"UserID"`       // รหัสผู้ใช้
	Password     string `db:"Password"`     // (MD5 Hash) ใช้สำหรับยืนยันตัวตน
	UserName     string `db:"UserName"`     // ชื่อผู้ใช้
	NickName     string `db:"NickName"`     // ชื่อเล่น
	FullNameTH   string `db:"FullNameTH"`   // ชื่อเต็มภาษาไทย
	DepartmentNo string `db:"DepartmentNo"` // รหัสแผนก
	RoleID       int    `db:"RoleID"`       // รหัสบทบาท
	RoleName     string `db:"RoleName"`     // ชื่อบทบาท
	Description  string `db:"Description"`  // รายละเอียดบทบาท
	IsActive     bool   `db:"IsActive"`     // สถานะการใช้งานบัญชี
}

type ROM_V_OrderHeadDetail struct {
	OrderNo     string    `db:"OrderNo"`     // เลขที่ใบสั่งซื้อ
	SoNo        string    `db:"SoNo"`        // เลขที่ใบสั่งขาย
	StatusMKP   string    `db:"StatusMKP"`   // สถานะในตลาด
	SalesStatus string    `db:"SalesStatus"` // สถานะการขาย
	CreateDate  time.Time `db:"CreateDate"`  // วันที่สร้างรายการ
	TrackingNo  string    `db:"TrackingNo"`  // เลขพัสดุ
}

type ROM_V_OrderLineDetail struct {
	OrderNo     string    `db:"OrderNo"`     // เลขที่ใบสั่งซื้อ
	SoNo        string    `db:"SoNo"`        // เลขที่ใบสั่งขาย
	StatusMKP   string    `db:"StatusMKP"`   // สถานะในตลาด
	SalesStatus string    `db:"SalesStatus"` // สถานะการขาย
	SKU         string    `db:"SKU"`         // รหัสสินค้า
	ItemName    string    `db:"ItemName"`    // ชื่อสินค้า
	QTY         int       `db:"QTY"`         // จำนวนสินค้าที่ซื้อ
	Price       float64   `db:"Price"`       // ราคาต่อหน่วย
	CreateDate  time.Time `db:"CreateDate"`  // วันที่สร้างรายการ
	TrackingNo  string    `db:"TrackingNo"`  // เลขพัสดุ
}
