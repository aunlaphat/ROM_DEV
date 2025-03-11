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
	UserID        string `db:"UserID"`        // รหัสผู้ใช้
	Password      string `db:"Password"`      // (MD5 Hash) ใช้สำหรับยืนยันตัวตน
	UserName      string `db:"UserName"`      // ชื่อผู้ใช้
	NickName      string `db:"NickName"`      // ชื่อเล่น
	FullNameTH    string `db:"FullNameTH"`    // ชื่อเต็มภาษาไทย
	DepartmentNo  string `db:"DepartmentNo"`  // รหัสแผนก
	RoleID        int    `db:"RoleID"`        // รหัสบทบาท
	RoleName      string `db:"RoleName"`      // ชื่อบทบาท
	WarehouseID   int    `db:"WarehouseID"`   // คลังสินค้า
	WarehouseName string `db:"WarehouseName"` // ชื่อคลังสินค้า
	Description   string `db:"Description"`   // รายละเอียดบทบาท
	IsActive      bool   `db:"IsActive"`      // สถานะการใช้งานบัญชี
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

/********** Constants for dropdown ***************/

type ROM_V_ProductAll struct {
	SKU       string  `db:"SKU" json:"sku"`             // รหัสสินค้า
	NameAlias string  `db:"NAMEALIAS" json:"nameAlias"` // ชื่อย่อของสินค้า
	Size      string  `db:"Size" json:"size"`           // ขนาดของสินค้า
	SizeID    string  `db:"SizeID" json:"sizeID"`       // รหัสขนาดของสินค้า
	Barcode   *string `db:"Barcode" json:"barcode"`     // บาร์โค้ดของสินค้า
	Type      *string `db:"Type" json:"type"`           // ประเภทของสินค้า
}

type InvoiceInformation struct {
	CustomerID   string  `db:"CustomerID" json:"customerID"`
	CustomerName *string `db:"CustomerName" json:"customerName"` // ใช้เป็น CustomerName and InvoiceName
	Address      *string `db:"Address" json:"address"`
	TaxID        string  `db:"TaxID" json:"taxID"`
}

type Province struct {
	ProvinceCode    int    `db:"ProvinceCode" json:"provinceCode"`
	ProvicesTH      string `db:"ProvicesTH" json:"provicesTH"`
}

type District struct {
	ProvinceCode    int    `db:"ProvinceCode" json:"provinceCode"`
	DistrictCode    int    `db:"DistrictCode" json:"districtCode"`
	DistrictTH      string `db:"DistrictTH" json:"districtTH"`
}

type SubDistrict struct {
	DistrictCode    int    `db:"DistrictCode" json:"districtCode"`
	SubdistrictCode int    `db:"SubdistrictCode" json:"subdistrictCode"`
	SubdistrictTH   string `db:"SubdistrictTH" json:"subdistrictTH"`
}

type PostalCode struct {
	SubdistrictCode int    `db:"SubdistrictCode" json:"subdistrictCode"`
	ZipCode         string `db:"ZipCode" json:"zipCode"`
}
