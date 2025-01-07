package request

import "time"

type BeforeReturnOrder struct {
	OrderNo        string     `json:"orderNo" db:"OrderNo"`
	SoNo           string     `json:"soNo" db:"SoNo"`
	SrNo           string     `json:"srNo" db:"SrNo"`
	ChannelID      int        `json:"channelID" db:"ChannelID"`
	ReturnType     string     `json:"returnType" db:"ReturnType"`
	CustomerID     string     `json:"customerID" db:"CustomerID"`
	TrackingNo     string     `json:"trackingNo" db:"TrackingNo"`
	Logistic       string     `json:"logistic" db:"Logistic"`
	WarehouseID    int        `json:"warehouseID" db:"WarehouseID"`
	SoStatusID     *int       `json:"soStatusID" db:"SoStatusID"`
	MkpStatusID    *int       `json:"mkpStatusID" db:"MkpStatusID"`
	ReturnDate     *time.Time `json:"returnDate" db:"ReturnDate"`
	StatusReturnID int        `json:"statusReturnID" db:"StatusReturnID"`
	StatusConfID   int        `json:"statusConfID" db:"StatusConfID"`
	ConfirmBy      *string    `json:"confirmBy" db:"ConfirmBy"`
	CreateBy       string     `json:"createBy" db:"CreateBy"`
	// CreateDate             *time.Time              `json:"createDate" db:"CreateDate"` // MSSQL GetDate()
	UpdateBy *string `json:"updateBy" db:"UpdateBy"`
	//UpdateDate             *time.Time              `json:"updateDate" db:"UpdateDate"` // MSSQL GetDate()
	CancelID               *int                    `json:"cancelID" db:"CancelID"`
	BeforeReturnOrderLines []BeforeReturnOrderLine `json:"beforeReturnOrderLines"`
}

type BeforeReturnOrderLine struct {
	OrderNo    string  `json:"orderNo" db:"OrderNo"`
	SKU        string  `json:"sku" db:"SKU"`
	QTY        int     `json:"qty" db:"QTY"`
	ReturnQTY  int     `json:"returnQty" db:"ReturnQTY"`
	Price      float64 `json:"price" db:"Price"`
	CreateBy   string  `json:"createBy" db:"CreateBy"`
	TrackingNo string  `json:"trackingNo" db:"TrackingNo"`
	AlterSKU   *string `json:"alterSKU" db:"AlterSKU"`
	UpdateBy   *string `json:"updateBy" db:"UpdateBy"`
	//UpdateDate *time.Time `json:"updateDate" db:"UpdateDate"` // MSSQL GetDate()
	//CreateDate *time.Time `json:"createDate" db:"CreateDate"` // MSSQL GetDate()
}

type SearchOrderRequest struct {
	SoNo string `json:"soNo" db:"SoNo"`
}

type EditOrderRequest struct {
	OrderNo                string                  `json:"orderNo" db:"OrderNo"`
	SoNo                   string                  `json:"soNo" db:"SoNo"`
	SrNo                   string                  `json:"srNo" db:"SrNo"`
	ChannelID              int                     `json:"channelID" db:"ChannelID"`
	ReturnType             string                  `json:"returnType" db:"ReturnType"`
	CustomerID             string                  `json:"customerID" db:"CustomerID"`
	TrackingNo             string                  `json:"trackingNo" db:"TrackingNo"`
	Logistic               string                  `json:"logistic" db:"Logistic"`
	WarehouseID            int                     `json:"warehouseID" db:"WarehouseID"`
	SoStatusID             *int                    `json:"soStatusID" db:"SoStatusID"`
	MkpStatusID            *int                    `json:"mkpStatusID" db:"MkpStatusID"`
	ReturnDate             *time.Time              `json:"returnDate" db:"ReturnDate"`
	StatusReturnID         int                     `json:"statusReturnID" db:"StatusReturnID"`
	StatusConfID           int                     `json:"statusConfID" db:"StatusConfID"`
	ConfirmBy              *string                 `json:"confirmBy" db:"ConfirmBy"`
	UpdateBy               *string                 `json:"updateBy" db:"UpdateBy"`
	CancelID               *int                    `json:"cancelID" db:"CancelID"`
	BeforeReturnOrderLines []BeforeReturnOrderLine `json:"beforeReturnOrderLines"`
}

type ConfirmOrderRequest struct {
	OrderNo      string `json:"orderNo" db:"OrderNo"`
	StatusConfID int    `json:"statusConfID" db:"StatusConfID"`
	ConfirmBy    string `json:"confirmBy" db:"ConfirmBy"`
}
