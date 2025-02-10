package request

import "time"

type SearchOrder struct {
	SoNo    string `json:"soNo" db:"SoNo" form:"soNo"`
	OrderNo string `json:"orderNo" db:"OrderNo" form:"orderNo"`
}

type CreateBeforeReturnOrder struct {
	OrderNo string `json:"orderNo" db:"OrderNo" binding:"required"`
	SoNo    string `json:"soNo" db:"SoNo" binding:"required"`
	//SrNo        *string                       `json:"srNo,omitempty" db:"SrNo"`
	SoStatus    string `json:"soStatus" db:"SoStatus"`
	MkpStatus   string `json:"mkpStatus" db:"MkpStatus"`
	WarehouseID int    `json:"warehouseID" db:"WarehouseID" binding:"required"`
	//Location    string                        `json:"location" db:"Location" binding:"required"`
	ReturnDate time.Time `json:"returnDate" db:"ReturnDate" binding:"required"`
	TrackingNo string    `json:"trackingNo" db:"TrackingNo" binding:"required"`
	Logistic   string    `json:"logistic" db:"Logistic" binding:"required"`
	//CreateBy    string                        `json:"createBy" db:"CreateBy" binding:"required"`
	Items []CreateBeforeReturnOrderItem `json:"items"`
}

type CreateBeforeReturnOrderItem struct {
	SKU        string  `json:"sku" db:"SKU" binding:"required"`
	ItemName   string  `json:"itemName" db:"ItemName" binding:"required"`
	QTY        int     `json:"qty" db:"QTY" binding:"required"`
	ReturnQTY  int     `json:"returnQty" db:"ReturnQTY" binding:"required"`
	Price      float64 `json:"price" db:"Price" binding:"required"`
	CreateBy   string  `json:"createBy" db:"CreateBy" binding:"required"`
	TrackingNo *string `json:"trackingNo,omitempty" db:"TrackingNo"`
	AlterSKU   *string `json:"alterSKU,omitempty" db:"AlterSKU"`
}

type UpdateSaleReturn struct {
	OrderNo string `json:"orderNo" db:"OrderNo" example:"ORD-TEST-123456"`
	SrNo    string `json:"srNo" db:"SrNo" example:"SR-TEST-123456"`
}

type CancelSaleReturn struct {
	OrderNo string `json:"orderNo" db:"OrderNo" example:"ORD-TEST-123456"`
	Remark  string `json:"remark" db:"Remark" example:"cancel order"`
}

// Draft & Confirm MKP

type AddCodeR struct {
	OrderNo   string  `json:"orderNo" db:"OrderNo"`
	SKU       string  `json:"sku" db:"SKU"`
	ItemName  string  `json:"itemName" db:"ItemName"`
	QTY       int     `json:"qty" db:"QTY"`
	ReturnQTY int     `json:"returnQty" db:"ReturnQTY"`
	Price     float64 `json:"price" db:"Price"`
}
