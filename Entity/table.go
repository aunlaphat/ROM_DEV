package entity

import "time"

// BeforeReturnOrder represents the structure of the BeforeReturnOrder table
type BeforeReturnOrder struct {
	RecID          int        `db:"RecID"`          // รหัสอ้างอิงอัตโนมัติ (PK - Auto Increment)
	OrderNo        string     `db:"OrderNo"`        // เลขที่ใบสั่งซื้อ
	SoNo           string     `db:"SoNo"`           // เลขที่ใบสั่งขาย
	SrNo           string     `db:"SrNo"`           // เลขที่ใบลดหนี้
	ChannelID      int        `db:"ChannelID"`      // รหัสช่องทางการขาย
	Reason         string     `db:"Reason"`         // เหตุผลในการคืนสินค้า
	CustomerID     string     `db:"CustomerID"`     // รหัสลูกค้า
	TrackingNo     string     `db:"TrackingNo"`     // เลขพัสดุ
	Logistic       string     `db:"Logistic"`       // ขนส่ง
	WarehouseID    int        `db:"WarehouseID"`    // รหัสคลังสินค้า
	SoStatusID     *int       `db:"SoStatusID"`     // สถานะใบสั่งขาย
	MkpStatusID    *int       `db:"MkpStatusID"`    // สถานะในตลาด
	ReturnDate     *time.Time `db:"ReturnDate"`     // วันที่คืนสินค้า
	StatusReturnID int        `db:"StatusReturnID"` // สถานะการคืนสินค้า
	StatusConfID   int        `db:"StatusConfID"`   // สถานะการยืนยัน
	ConfirmBy      *string    `db:"ConfirmBy"`      // ผู้ยืนยัน
	ConfirmDate    *time.Time `db:"ConfirmDate"`    // วันที่ยืนยัน
	CreateBy       string     `db:"CreateBy"`       // ผู้สร้างรายการ
	CreateDate     time.Time  `db:"CreateDate"`     // วันที่สร้างรายการ
	UpdateBy       *string    `db:"UpdateBy"`       // ผู้แก้ไขล่าสุด
	UpdateDate     *time.Time `db:"UpdateDate"`     // วันที่แก้ไขล่าสุด
	CancelID       *int       `db:"CancelID"`       // รหัสการยกเลิก
}

// BeforeReturnOrderLine คือตารางสำหรับเก็บรายการสินค้าที่ต้องการคืน
// เป็นรายละเอียดของแต่ละรายการในใบคืนสินค้า
type BeforeReturnOrderLine struct {
	RecID      int        `db:"RecID"`      // รหัสอ้างอิงอัตโนมัติ (PK - Auto Increment)
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

// ReturnOrder คือตารางสำหรับเก็บข้อมูลการคืนสินค้าที่ผ่านการตรวจสอบแล้ว
// เป็นขั้นตอนสุดท้ายของกระบวนการคืนสินค้า
type ReturnOrder struct {
	ReturnID      string     `db:"ReturnID"`      // เลขที่ใบคืนสินค้า (PK - Generate จากระบบ)
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
	RecID      int        `db:"RecID"`      // รหัสอ้างอิงอัตโนมัติ - (PK - Auto Increment)
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

// View
type ROM_V_ProductAll struct {
	SKU       string  `db:"SKU" json:"sku"`             // รหัสสินค้า
	NameAlias string  `db:"NAMEALIAS" json:"nameAlias"` // ชื่อย่อของสินค้า
	Size      string  `db:"Size" json:"size"`           // ขนาดของสินค้า
	SizeID    string  `db:"SizeID" json:"sizeID"`       // รหัสขนาดของสินค้า
	Barcode   *string `db:"Barcode" json:"barcode"`     // บาร์โค้ดของสินค้า
	Type      *string `db:"Type" json:"type"`           // ประเภทของสินค้า
}

type ROM_V_OrderHeadDetail struct {
	OrderNo     string                  `db:"OrderNo" json:"orderNo"`         // เลขที่ใบสั่งซื้อ
	SoNo        string                  `db:"SoNo" json:"soNo"`               // เลขที่ใบสั่งขาย
	StatusMKP   string                  `db:"StatusMKP" json:"statusMKP"`     // สถานะในตลาด
	SalesStatus string                  `db:"SalesStatus" json:"salesStatus"` // สถานะการขาย
	CreateDate  time.Time               `db:"CreateDate" json:"createDate"`   // วันที่สร้างรายการ
	OrderLines  []ROM_V_OrderLineDetail `json:"orderLines"`                   // รายการสินค้า
}

type ROM_V_OrderLineDetail struct {
	OrderNo     string    `db:"OrderNo" json:"orderNo"`         // เลขที่ใบสั่งซื้อ
	SoNo        string    `db:"SoNo" json:"soNo"`               // เลขที่ใบสั่งขาย
	StatusMKP   string    `db:"StatusMKP" json:"statusMKP"`     // สถานะในตลาด
	SalesStatus string    `db:"SalesStatus" json:"salesStatus"` // สถานะการขาย
	SKU         string    `db:"SKU" json:"sku"`                 // รหัสสินค้า
	ItemName    string    `db:"ItemName" json:"itemName"`       // ชื่อสินค้า
	QTY         int       `db:"QTY" json:"qty"`                 // จำนวนสินค้า
	Price       float64   `db:"Price" json:"price"`             // ราคาต่อหน่วย
	CreateDate  time.Time `db:"CreateDate" json:"createDate"`   // วันที่สร้างรายการ
}

/* type ROM_V_User struct {
	UserID       string `json:"userID,omitempty" db:"UserID"`           // รหัสผู้ใช้
	UserName     string `json:"userName,omitempty" db:"Username"`       // ชื่อผู้ใช้
	NickName     string `json:"nickName,omitempty" db:"NickName"`       // ชื่อเล่น
	FullNameTH   string `json:"fullNameTH,omitempty" db:"FullNameTH"`   // ชื่อเต็มภาษาไทย
	DepartmentNo string `json:"department,omitempty" db:"DepartmentNo"` // รหัสแผนก
} */

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
