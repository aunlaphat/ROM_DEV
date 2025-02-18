package entity

import "time"

type User struct {
	UserID       string `db:"UserID"`
	UserName     string `db:"UserName"`
	Password     string `db:"Password"`
	NickName     *string `db:"NickName"`
	FullNameTH   *string `db:"FullNameTH"`
	DepartmentNo string `db:"DepartmentNo"`
	RoleID       *int    `db:"RoleID"`
	RoleName     *string `db:"RoleName"`
	Description  *string `db:"Description"`
	// Permission   string `db:"Permission"`
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

type District struct {
	Code         int    `db:"Code" json:"code"`
	ProvinceCode int    `db:"ProvinceCode" json:"provinceCode"`
	NameEN       string `db:"NameEN" json:"nameEN"`
	NameTH       string `db:"NameTH" json:"nameTH"`
}

type SubDistrict struct {
	Code         int    `db:"Code" json:"code"`
	DistrictCode int    `db:"DistrictCode" json:"districtCode"`
	ZipCode      string `db:"ZipCode" json:"zipCode"`
	NameTH       string `db:"NameTH" json:"nameTH"`
	NameEN       string `db:"NameEN" json:"nameEN"`
}

type Province struct {
	Code   int    `db:"Code" json:"code"`
	NameTH string `db:"NameTH" json:"nameTH"`
	NameEN string `db:"NameEN" json:"nameEN"`
}

type PostCode struct {

}
