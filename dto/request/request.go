package request

import "time"

type SearchOrder struct {
	SoNo    string `json:"soNo" db:"SoNo" form:"soNo"`
	OrderNo string `json:"orderNo" db:"OrderNo" form:"orderNo"`
}

type CreateBeforeReturnOrder struct {
	OrderNo     string                        `json:"orderNo" db:"OrderNo" binding:"required"`
	SoNo        string                        `json:"soNo" db:"SoNo" binding:"required"`
	ChannelID   int                           `json:"channelID" db:"ChannelID" binding:"required"`
	CustomerID  string                        `json:"customerID" db:"CustomerID" binding:"required"`
	Reason      string                        `json:"reason" db:"Reason" binding:"required"`
	SoStatus    string                        `json:"soStatus,omitempty" db:"SoStatus"`
	MkpStatus   string                        `json:"mkpStatus,omitempty" db:"MkpStatus"`
	WarehouseID int                           `json:"warehouseID" db:"WarehouseID" binding:"required"`
	ReturnDate  time.Time                     `json:"returnDate" db:"ReturnDate" binding:"required"`
	TrackingNo  string                        `json:"trackingNo" db:"TrackingNo" binding:"required"`
	Logistic    string                        `json:"logistic" db:"Logistic" binding:"required"`
	Items       []CreateBeforeReturnOrderItem `json:"items"`
}

type CreateBeforeReturnOrderItem struct {
	OrderNo    string  `json:"orderNo" db:"OrderNo" binding:"required"`
	SKU        string  `json:"sku" db:"SKU" binding:"required"`
	ItemName   string  `json:"itemName" db:"ItemName" binding:"required"`
	QTY        int     `json:"qty" db:"QTY" binding:"required"`
	ReturnQTY  int     `json:"returnQty" db:"ReturnQTY" binding:"required"`
	Price      float64 `json:"price" db:"Price" binding:"required"`
	CreateBy   string  `json:"createBy" db:"CreateBy" binding:"required"`
	TrackingNo *string `json:"trackingNo,omitempty" db:"TrackingNo"`
	AlterSKU   *string `json:"alterSKU,omitempty" db:"AlterSKU"`
}

type CancelOrder struct {
	RefID        string `json:"refID" db:"RefID"`
	SourceTable  string `json:"sourceTable" db:"SourceTable"`
	CancelReason string `json:"cancelReason" db:"CancelReason"`
}

type AddItem struct {
	OrderNo   string  `json:"orderNo" db:"OrderNo"`
	SKU       string  `json:"sku" db:"SKU"`
	ItemName  string  `json:"itemName" db:"ItemName"`
	QTY       int     `json:"qty" db:"QTY"`
	ReturnQTY int     `json:"returnQTY" db:"ReturnQTY"`
	Price     float64 `json:"price" db:"Price"`
}
