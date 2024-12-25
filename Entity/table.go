package entity

import "time"

// User schema
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

// BeforeReturnOrder represents the structure of the BeforeReturnOrder table
type BeforeReturnOrder struct {
	RecID          int        `db:"RecID"`
	OrderNo        string     `db:"OrderNo"`
	SaleOrder      string     `db:"SaleOrder"`
	SaleReturn     string     `db:"SaleReturn"`
	ChannelID      int        `db:"ChannelID"`
	ReturnType     string     `db:"ReturnType"`
	CustomerID     string     `db:"CustomerID"`
	TrackingNo     string     `db:"TrackingNo"`
	Logistic       string     `db:"Logistic"`
	WarehouseID    int        `db:"WarehouseID"`
	SoStatusID     *int       `db:"SoStatusID"`
	MkpStatusID    *int       `db:"MkpStatusID"`
	ReturnDate     *time.Time `db:"ReturnDate"`
	StatusReturnID int        `db:"StatusReturnID"`
	StatusConfID   int        `db:"StatusConfID"`
	ConfirmBy      *string    `db:"ConfirmBy"`
	CreateBy       string     `db:"CreateBy"`
	CreateDate     time.Time  `db:"CreateDate"`
	UpdateBy       *string    `db:"UpdateBy"`
	UpdateDate     *time.Time `db:"UpdateDate"`
	CancelID       *int       `db:"CancelID"`
}

// CancelReturnOrder represents the structure of the CancelReturnOrder table
type CancelReturnOrder struct {
	CancelID   int       `db:"CancelID"`   // รหัสการยกเลิก (Primary Key)
	OrderNo    string    `db:"OrderNo"`    // เลขที่ใบสั่งซื้อ (Foreign Key -> BeforeReturnOrder)
	Remark     string    `db:"Remark"`     // เหตุผลในการยกเลิก
	CancelBy   string    `db:"CancelBy"`   // ผู้ยกเลิก
	CancelDate time.Time `db:"CancelDate"` // วันที่ยกเลิก
}

// BeforeReturnOrderLine คือตารางสำหรับเก็บรายการสินค้าที่ต้องการคืน
// เป็นรายละเอียดของแต่ละรายการในใบคืนสินค้า
type BeforeReturnOrderLine struct {
	RecID      int        `db:"RecID"`      // รหัสอ้างอิงอัตโนมัติ
	OrderNo    string     `db:"OrderNo"`    // เลขที่ใบสั่งซื้อ (FK -> BeforeReturnOrder)
	SKU        string     `db:"SKU"`        // รหัสสินค้า
	QTY        *int       `db:"QTY"`        // จำนวนสินค้าที่ซื้อ
	ReturnQTY  *int       `db:"ReturnQTY"`  // จำนวนที่ต้องการคืน
	Price      *float64   `db:"Price"`      // ราคาต่อหน่วย
	CreateBy   string     `db:"CreateBy"`   // ผู้สร้างรายการ
	CreateDate time.Time  `db:"CreateDate"` // วันที่สร้างรายการ
	AlterSKU   *string    `db:"AlterSKU"`   // รหัสสินค้าทดแทน (ถ้ามี)
	UpdateBy   *string    `db:"UpdateBy"`   // ผู้แก้ไขล่าสุด
	UpdateDate *time.Time `db:"UpdateDate"` // วันที่แก้ไขล่าสุด
	TrackingNo string     `db:"TrackingNo"` // เลขพัสดุ
}

// ReturnOrder คือตารางสำหรับเก็บข้อมูลการคืนสินค้าที่ผ่านการตรวจสอบแล้ว
// เป็นขั้นตอนสุดท้ายของกระบวนการคืนสินค้า
type ReturnOrder struct {
	ReturnID      string     `db:"ReturnID"`      // เลขที่ใบคืนสินค้า (Primary Key)
	OrderNo       string     `db:"OrderNo"`       // เลขที่ใบสั่งซื้อ
	SaleOrder     string     `db:"SaleOrder"`     // เลขที่ใบกำกับภาษี
	SaleReturn    string     `db:"SaleReturn"`    // เลขที่ใบลดหนี้
	TrackingNo    string     `db:"TrackingNo"`    // เลขพัสดุ
	PlatfID       *int       `db:"PlatfID"`       // รหัสแพลตฟอร์ม (FK -> Platforms)
	ChannelID     *int       `db:"ChannelID"`     // รหัสช่องทางการขาย (FK -> Channel)
	OptStatusID   *int       `db:"OptStatusID"`   // สถานะการดำเนินการ
	AxStatusID    *int       `db:"AxStatusID"`    // สถานะในระบบ AX
	PlatfStatusID *int       `db:"PlatfStatusID"` // สถานะในแพลตฟอร์ม
	Remark        *string    `db:"Remark"`        // หมายเหตุ
	CreateBy      string     `db:"CreateBy"`      // ผู้สร้างรายการ
	CreateDate    time.Time  `db:"CreateDate"`    // วันที่สร้างรายการ
	UpdateBy      *string    `db:"UpdateBy"`      // ผู้แก้ไขล่าสุด
	UpdateDate    *time.Time `db:"UpdateDate"`    // วันที่แก้ไขล่าสุด
	CancelID      *int       `db:"CancelID"`      // รหัสการยกเลิก
	StatusCheckID *int       `db:"StatusCheckID"` // สถานะการตรวจสอบ (FK -> StatusCheck)
	CheckBy       *string    `db:"CheckBy"`       // ผู้ตรวจสอบ
	Description   *string    `db:"Description"`   // รายละเอียดเพิ่มเติม
}

// ReturnOrderLine คือตารางสำหรับเก็บรายการสินค้าที่คืนและผ่านการตรวจสอบแล้ว
type ReturnOrderLine struct {
	RecID      int        `db:"RecID"`      // รหัสอ้างอิงอัตโนมัติ
	ReturnID   string     `db:"ReturnID"`   // เลขที่ใบคืนสินค้า (FK -> ReturnOrder)
	OrderNo    string     `db:"OrderNo"`    // เลขที่ใบสั่งซื้อ
	TrackingNo string     `db:"TrackingNo"` // เลขพัสดุ
	SKU        string     `db:"SKU"`        // รหัสสินค้า
	ReturnQTY  int        `db:"ReturnQTY"`  // จำนวนที่คืน
	CheckQTY   *int       `db:"CheckQTY"`   // จำนวนที่ตรวจสอบแล้ว
	Price      float64    `db:"Price"`      // ราคาต่อหน่วย
	CreateBy   string     `db:"CreateBy"`   // ผู้สร้างรายการ
	CreateDate time.Time  `db:"CreateDate"` // วันที่สร้างรายการ
	AlterSKU   *string    `db:"AlterSKU"`   // รหัสสินค้าทดแทน
	UpdateBy   *string    `db:"UpdateBy"`   // ผู้แก้ไขล่าสุด
	UpdateDate *time.Time `db:"UpdateDate"` // วันที่แก้ไขล่าสุด
}

// สถานะต่างๆ ของการคืนสินค้า
const (
	// สถานะการคืนสินค้าเบื้องต้น
	StatusReturnPending  = 1 // รอดำเนินการ
	StatusReturnApproved = 2 // อนุมัติแล้ว
	StatusReturnRejected = 3 // ปฏิเสธการคืน

	// สถานะการยืนยัน
	StatusConfirmPending  = 1 // รอยืนยัน
	StatusConfirmApproved = 2 // ยืนยันแล้ว
	StatusConfirmRejected = 3 // ปฏิเสธการยืนยัน

	// สถานะการตรวจสอบ
	StatusCheckPending    = 1 // รอตรวจสอบ
	StatusCheckInProgress = 2 // กำลังตรวจสอบ
	StatusCheckCompleted  = 3 // ตรวจสอบเสร็จสิ้น
	StatusCheckFailed     = 4 // ตรวจสอบพบปัญหา
)
