package response

import "time"

type ReturnOrder struct {
	OrderNo          string            `json:"orderNo" db:"OrderNo"`
	SoNo             string            `json:"soNo" db:"SoNo"`
	SrNo             string            `json:"srNo" db:"SrNo"`
	TrackingNo       string            `json:"trackingNo" db:"TrackingNo"`
	PlatfID          int               `json:"platfID" db:"PlatfID"`
	ChannelID        int               `json:"channelID" db:"ChannelID"`
	OptStatusID      int               `json:"optStatusID" db:"OptStatusID"`
	AxStatusID       int               `json:"axStatusID" db:"AxStatusID"`
	Reason           string            `json:"reason" db:"Reason"`
	CreateBy         string            `json:"createBy" db:"CreateBy"`
	CreateDate       *time.Time        `json:"createDate" db:"CreateDate"`
	UpdateBy         string            `json:"updateBy" db:"UpdateBy"`
	UpdateDate       *time.Time        `json:"updateDate" db:"UpdateDate"`
	CancelID         *int              `json:"cancelID" db:"CancelID"`
	ReturnOrderLines []ReturnOrderLine `json:"returnOrderLines"`
}

type ReturnOrderLine struct {
	OrderNo    string     `json:"orderNo" db:"OrderNo"`
	SKU        string     `json:"sku" db:"SKU"`
	ItemName   string     `json:"itemName" db:"ItemName"`
	QTY        int        `json:"qty" db:"QTY"`
	ReturnQTY  int        `json:"returnQty" db:"ReturnQTY"`
	ActualQTY  int        `json:"actualQty" db:"ActualQTY"`
	Price      float64    `json:"price" db:"Price"`
	CreateBy   string     `json:"createBy" db:"CreateBy"`
	CreateDate *time.Time `json:"createDate" db:"CreateDate"`
	TrackingNo string     `json:"trackingNo" db:"TrackingNo"`
	AlterSKU   *string    `json:"alterSKU" db:"AlterSKU"`
	UpdateBy   *string    `json:"updateBy" db:"UpdateBy"`
	UpdateDate *time.Time `json:"updateDate" db:"UpdateDate"`
}
