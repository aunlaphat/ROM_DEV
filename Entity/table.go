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

// BeforeReturnOrder คือตารางสำหรับเก็บข้อมูลการขอคืนสินค้าก่อนการตรวจสอบ
// เป็นขั้นตอนแรกของกระบวนการคืนสินค้า
type BeforeReturnOrder struct {
	RecID          int        `db:"RecID"`          // รหัสอ้างอิงอัตโนมัติ
	OrderNo        string     `db:"OrderNo"`        // เลขที่ใบสั่งซื้อ
	SaleOrder      string     `db:"SaleOrder"`      // เลขที่ใบกำกับภาษี/ใบส่งของ
	SaleReturn     string     `db:"SaleReturn"`     // เลขที่ใบลดหนี้
	ChannelID      int        `db:"ChannelID"`      // รหัสช่องทางการขาย
	ReturnType     string     `db:"ReturnType"`     // ประเภทการคืน (เช่น คืนเงิน, เปลี่ยนสินค้า)
	CustomerID     string     `db:"CustomerID"`     // รหัสลูกค้า
	TrackingNo     string     `db:"TrackingNo"`     // เลขพัสดุ
	Logistic       string     `db:"Logistic"`       // ขนส่ง
	WarehouseID    int        `db:"WarehouseID"`    // รหัสคลังสินค้า
	SoStatusID     *int       `db:"SoStatusID"`     // สถานะใบกำกับภาษี
	MkpStatusID    *int       `db:"MkpStatusID"`    // สถานะในระบบ Marketplace
	ReturnDate     *time.Time `db:"ReturnDate"`     // วันที่คืนสินค้า
	StatusReturnID int        `db:"StatusReturnID"` // สถานะการคืน (FK -> StatusReturn)
	StatusConfID   int        `db:"StatusConfID"`   // สถานะการยืนยัน (FK -> StatusConfirm)
	ConfirmBy      *string    `db:"ConfirmBy"`      // ผู้ยืนยัน
	CreateBy       string     `db:"CreateBy"`       // ผู้สร้างรายการ
	CreateDate     time.Time  `db:"CreateDate"`     // วันที่สร้างรายการ
	UpdateBy       *string    `db:"UpdateBy"`       // ผู้แก้ไขล่าสุด
	UpdateDate     *time.Time `db:"UpdateDate"`     // วันที่แก้ไขล่าสุด
	CancelID       *int       `db:"CancelID"`       // รหัสการยกเลิก (ถ้ามี)
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
