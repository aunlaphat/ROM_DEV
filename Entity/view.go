package entity

import "time"

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
