package entity

import "time"

// type BeforeReturnOrder struct {
// 	RecID          int        `db:"RecID"`          // รหัสอ้างอิงอัตโนมัติ (PK - Auto Increment)
// 	OrderNo        string     `db:"OrderNo"`        // เลขที่ใบสั่งซื้อ
// 	SoNo           string     `db:"SoNo"`           // เลขที่ใบสั่งขาย
// 	SrNo           string     `db:"SrNo"`           // เลขที่ใบลดหนี้
// 	ChannelID      int        `db:"ChannelID"`      // รหัสช่องทางการขาย
// 	Reason         string     `db:"Reason"`         // เหตุผลในการคืนสินค้า
// 	CustomerID     string     `db:"CustomerID"`     // รหัสลูกค้า
// 	TrackingNo     string     `db:"TrackingNo"`     // เลขพัสดุ
// 	Logistic       string     `db:"Logistic"`       // ขนส่ง
// 	WarehouseID    int        `db:"WarehouseID"`    // รหัสคลังสินค้า
// 	SoStatus       *int       `db:"SoStatus"`       // สถานะใบสั่งขาย
// 	MkpStatus      *int       `db:"MkpStatus"`      // สถานะในตลาด
// 	ReturnDate     *time.Time `db:"ReturnDate"`     // วันที่คืนสินค้า
// 	StatusReturnID int        `db:"StatusReturnID"` // สถานะการคืนสินค้า
// 	StatusConfID   int        `db:"StatusConfID"`   // สถานะการยืนยัน
// 	ConfirmBy      *string    `db:"ConfirmBy"`      // ผู้ยืนยัน
// 	ConfirmDate    *time.Time `db:"ConfirmDate"`    // วันที่ยืนยัน
// 	CreateBy       string     `db:"CreateBy"`       // ผู้สร้างรายการ
// 	CreateDate     time.Time  `db:"CreateDate"`     // วันที่สร้างรายการ
// 	UpdateBy       *string    `db:"UpdateBy"`       // ผู้แก้ไขล่าสุด
// 	UpdateDate     *time.Time `db:"UpdateDate"`     // วันที่แก้ไขล่าสุด
// 	CancelID       *int       `db:"CancelID"`       // รหัสการยกเลิก
// }

// ในฟิลมีปรับ type ข้อมูลเพิ่ม //11/02
type BeforeReturnOrder struct {
	RecID      int     `db:"RecID"`      // รหัสอ้างอิงอัตโนมัติ (PK - Auto Increment)
	OrderNo    string  `db:"OrderNo"`    // เลขที่ใบสั่งซื้อ
	SoNo       string  `db:"SoNo"`       // เลขที่ใบสั่งขาย
	SrNo       *string `db:"SrNo"`       // เลขที่ใบลดหนี้ (Nullable)
	ChannelID  int     `db:"ChannelID"`  // รหัสช่องทางการขาย
	CustomerID string  `db:"CustomerID"` // รหัสลูกค้า
	Reason     string  `db:"Reason"`     // เหตุผลในการคืนสินค้า
	TrackingNo string  `db:"TrackingNo"` // เลขพัสดุ
	Logistic   string  `db:"Logistic"`   // ขนส่ง
	//Location       string     `db:"Location"`                               // ปลายทาง
	WarehouseID    int        `db:"WarehouseID"`                            // รหัสคลังสินค้า
	SoStatus       string     `db:"SoStatus" default:"open order"`          // สถานะใบสั่งขาย (Default: "open order")
	MkpStatus      string     `db:"MkpStatus" default:"complete"`           // สถานะของ Marketplace (Default: "complete")
	ReturnDate     time.Time  `db:"ReturnDate"`                             // วันที่คืนสินค้า
	StatusReturnID int        `db:"StatusReturnID"`                         // สถานะการคืนสินค้า
	StatusConfID   int        `db:"StatusConfID"`                           // สถานะการยืนยัน
	ConfirmBy      *string    `db:"ConfirmBy"`                              // ผู้ยืนยัน (Nullable)
	ConfirmDate    *time.Time `db:"ConfirmDate"`                            // วันที่ยืนยัน (Nullable)
	CreateBy       string     `db:"CreateBy"`                               // ผู้สร้างรายการ
	CreateDate     time.Time  `db:"CreateDate" default:"CURRENT_TIMESTAMP"` // วันที่สร้างรายการ (Default: Now)
	UpdateBy       *string    `db:"UpdateBy"`                               // ผู้แก้ไขล่าสุด (Nullable)
	UpdateDate     *time.Time `db:"UpdateDate"`                             // วันที่แก้ไขล่าสุด (Nullable)
	CancelID       *int       `db:"CancelID"`                               // รหัสการยกเลิก (Nullable)
	IsCNCreated    bool       `db:"IsCNCreated" default:"false"`            // สถานะการสร้าง CN (Default: false)
	IsEdited       bool       `db:"IsEdited" default:"false"`               // สถานะการแก้ไข (Default: false)
}

type BeforeReturnOrderLine struct {
	RecID      int        `db:"RecID"`      // รหัสอ้างอิงอัตโนมัติ (PK - Auto Increment)
	OrderNo    string     `db:"OrderNo"`    // เลขที่ใบสั่งซื้อ (FK -> BeforeReturnOrder)
	SKU        string     `db:"SKU"`        // รหัสสินค้า
	ItemName   string     `db:"ItemName"`   // ชื่อสินค้า
	QTY        int        `db:"QTY"`        // จำนวนสินค้าที่ซื้อ
	ReturnQTY  int        `db:"ReturnQTY"`  // จำนวนที่ต้องการคืน
	Price      float64    `db:"Price"`      // ราคาต่อหน่วย
	CreateBy   string     `db:"CreateBy"`   // ผู้สร้างรายการ
	CreateDate time.Time  `db:"CreateDate"` // วันที่สร้างรายการ
	AlterSKU   *string    `db:"AlterSKU"`   // รหัสสินค้าทดแทน (ถ้ามี)
	UpdateBy   *string    `db:"UpdateBy"`   // ผู้แก้ไขล่าสุด
	UpdateDate *time.Time `db:"UpdateDate"` // วันที่แก้ไขล่าสุด
	TrackingNo *string    `db:"TrackingNo"` // เลขพัสดุ (ถ้ามีกรณีส่งสินค้าคนละพัสดุ)
}

type Warehouse struct {
	WarehouseID   int    `db:"WarehouseID"`   // รหัสคลังสินค้า (PK - Auto Increment)
	WarehouseName string `db:"WarehouseName"` // ชื่อคลังสินค้า
	Location      string `db:"Location"`      // ที่ตั้งของคลังสินค้า
}

type CancelStatus struct {
	CancelID     int       `db:"CancelID"`     // รหัสการยกเลิก (PK - Auto Increment)
	RefID        string    `db:"RefID"`        // เลขที่ใบสั่งซื้อ (FK -> RecID(BeforeReturnOrder) || RuturnID(ReturnOrder))
	CancelStatus bool      `db:"CancelStatus"` // สถานะการยกเลิก
	Remark       string    `db:"Remark"`       // เหตุผลในการยกเลิก
	CancelBy     string    `db:"CancelBy"`     // ผู้ยกเลิก
	CancelDate   time.Time `db:"CancelDate"`   // วันที่ยกเลิก
}

type ROM_V_UserPermission struct {
	UserID       string `json:"userID,omitempty" db:"UserID"`           // รหัสผู้ใช้
	UserName     string `json:"userName,omitempty" db:"Username"`       // ชื่อผู้ใช้
	NickName     string `json:"nickName,omitempty" db:"NickName"`       // ชื่อเล่น
	FullNameTH   string `json:"fullNameTH,omitempty" db:"FullNameTH"`   // ชื่อเต็มภาษาไทย
	DepartmentNo string `json:"department,omitempty" db:"DepartmentNo"` // รหัสแผนก
	RoleID       int    `json:"roleID,omitempty" db:"RoleID"`           // รหัสบทบาท
	RoleName     string `json:"roleName,omitempty" db:"RoleName"`       // ชื่อบทบาท
	Description  string `json:"description,omitempty" db:"Description"` // คำอธิบาย
	Permission   string `json:"permission,omitempty" db:"Permission"`   // สิทธิ์การเข้าถึง
}


/********** Return Order ***************/

// ReturnOrder คือตารางสำหรับเก็บข้อมูลการคืนสินค้าที่ผ่านการตรวจสอบแล้ว
// เป็นขั้นตอนสุดท้ายของกระบวนการคืนสินค้า
type ReturnOrder struct {
	RecID         int        `db:"RecID"`         // รหัสอ้างอิงอัตโนมัติ (Auto Increment)
	ReturnID      string     `db:"ReturnID"`      // เลขที่ใบคืนสินค้า (PK - Generate จากระบบ)
	OrderNo       string     `db:"OrderNo"`       // เลขที่ใบสั่งซื้อ
	SoNo          string     `db:"SoNo"`          // เลขที่ใบกำกับภาษี
	SrNo          string     `db:"SrNo"`          // เลขที่ใบลดหนี้
	TrackingNo    string     `db:"TrackingNo"`    // เลขพัสดุ
	PlatfID       *int       `db:"PlatfID"`       // รหัสแพลตฟอร์ม (FK -> Platforms)
	ChannelID     *int       `db:"ChannelID"`     // รหัสช่องทางการขาย (FK -> Channel)
	OptStatusID   *int       `db:"OptStatusID"`   // สถานะการดำเนินการ
	AxStatusID    *int       `db:"AxStatusID"`    // สถานะในระบบ AX
	PlatfStatusID *int       `db:"PlatfStatusID"` // สถานะในแพลตฟอร์ม
	Reason        *string    `db:"Reason"`        // เหตุผลในการคืนสินค้า
	CreateBy      string     `db:"CreateBy"`      // ผู้สร้างรายการ
	CreateDate    time.Time  `db:"CreateDate"`    // วันที่สร้างรายการ
	UpdateBy      *string    `db:"UpdateBy"`      // ผู้แก้ไขล่าสุด
	UpdateDate    *time.Time `db:"UpdateDate"`    // วันที่แก้ไขล่าสุด
	StatusCheckID *int       `db:"StatusCheckID"` // สถานะการตรวจสอบ (FK -> StatusCheck)
	CheckBy       *string    `db:"CheckBy"`       // ผู้ตรวจสอบ
	CheckDate     *time.Time `db:"CheckDate"`     // วันที่ตรวจสอบ
	CancelID      *int       `db:"CancelID"`      // รหัสการยกเลิก
	Description   *string    `db:"Description"`   // รายละเอียดเพิ่มเติม
}

// ReturnOrderLine คือตารางสำหรับเก็บรายการสินค้าที่คืนและผ่านการตรวจสอบแล้ว
type ReturnOrderLine struct {
	RecID      int        `db:"RecID"`      // รหัสอ้างอิงอัตโนมัติ - (PK - Auto Increment)
	ReturnID   string     `db:"ReturnID"`   // เลขที่ใบคืนสินค้า (FK -> ReturnOrder)
	OrderNo    string     `db:"OrderNo"`    // เลขที่ใบสั่งซื้อ
	SKU        string     `db:"SKU"`        // รหัสสินค้า
	ItemName   string     `db:"ItemName"`   // ชื่อสินค้า
	QTY        int        `db:"QTY"`        // จำนวนสินค้าที่ซื้อ
	ReturnQTY  int        `db:"ReturnQTY"`  // จำนวนที่คืน
	ActualQTY  int        `db:"ActualQTY"`  // จำนวนที่ตรวจสอบแล้ว
	Price      float64    `db:"Price"`      // ราคาต่อหน่วย
	CreateBy   string     `db:"CreateBy"`   // ผู้สร้างรายการ
	CreateDate time.Time  `db:"CreateDate"` // วันที่สร้างรายการ
	UpdateBy   *string    `db:"UpdateBy"`   // ผู้แก้ไขล่าสุด
	UpdateDate *time.Time `db:"UpdateDate"` // วันที่แก้ไขล่าสุด
	TrackingNo string     `db:"TrackingNo"` // เลขพัสดุ
	AlterSKU   *string    `db:"AlterSKU"`   // รหัสสินค้าทดแทน (ถ้ามี)
}

