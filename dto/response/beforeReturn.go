package response

import "time"

type BeforeReturnOrderResponse struct {
	OrderNo                string                          `json:"orderNo" db:"OrderNo"`
	SoNo                   string                          `json:"soNo" db:"SoNo"`
	SrNo                   string                          `json:"srNo" db:"SrNo"`
	ChannelID              int                             `json:"channelId" db:"ChannelID"`
	Reson                  string                          `json:"reason" db:"Reason"`
	CustomerID             string                          `json:"customerId" db:"CustomerID"`
	TrackingNo             string                          `json:"trackingNo" db:"TrackingNo"`
	Logistic               string                          `json:"logistic" db:"Logistic"`
	WarehouseID            int                             `json:"warehouseId" db:"WarehouseID"`
	SoStatusID             *int                            `json:"soStatusId" db:"SoStatusID"`
	MkpStatusID            *int                            `json:"mkpStatusId" db:"MkpStatusID"`
	ReturnDate             *time.Time                      `json:"returnDate" db:"ReturnDate"`
	StatusReturnID         *int                            `json:"statusReturnId" db:"StatusReturnID"`
	StatusConfID           *int                            `json:"statusConfId" db:"StatusConfID"`
	ConfirmBy              *string                         `json:"confirmBy" db:"ConfirmBy"`
	ConfirmDate            *time.Time                      `json:"confirmDate" db:"ConfirmDate"`
	CreateBy               string                          `json:"createBy" db:"CreateBy"`
	CreateDate             time.Time                       `json:"createDate" db:"CreateDate"`
	UpdateBy               *string                         `json:"updateBy" db:"UpdateBy"`
	UpdateDate             *time.Time                      `json:"updateDate" db:"UpdateDate"`
	CancelID               *int                            `json:"cancelId" db:"CancelID"`
	BeforeReturnOrderLines []BeforeReturnOrderLineResponse `json:"beforeReturnOrderLines"`
}

type BeforeReturnOrderLineResponse struct {
	OrderNo    string    `json:"orderNo" db:"OrderNo"`
	SKU        string    `json:"sku" db:"SKU"`
	ItemName   string    `json:"itemName" db:"ItemName"`
	QTY        int       `json:"qty" db:"QTY"`
	ReturnQTY  int       `json:"returnQty" db:"ReturnQTY"`
	Price      float64   `json:"price" db:"Price"`
	TrackingNo string    `json:"trackingNo" db:"TrackingNo"`
	CreateDate time.Time `json:"createDate" db:"CreateDate"`
}

type SaleOrderResponse struct {
	SoNo        string                  `json:"soNo" db:"SoNo"`
	OrderNo     string                  `json:"orderNo" db:"OrderNo"`
	StatusMKP   string                  `json:"statusMKP" db:"StatusMKP"`
	SalesStatus string                  `json:"salesStatus" db:"SalesStatus"`
	CreateDate  *time.Time              `json:"createDate" db:"CreateDate"`
	OrderLines  []SaleOrderLineResponse `json:"orderLines"`
}

type SaleOrderLineResponse struct {
	SoNo     string  `json:"soNo" db:"SoNo"`
	OrderNo  string  `json:"orderNo" db:"OrderNo"`
	SKU      string  `json:"sku" db:"SKU"`
	ItemName string  `json:"itemName" db:"ItemName"`
	QTY      int     `json:"qty" db:"QTY"`
	Price    float64 `json:"price" db:"Price"`
}

type UpdateSaleReturnResponse struct {
	OrderNo    string    `json:"orderNo" db:"OrderNo"`
	SrNo       string    `json:"srNo" db:"SrNo"`
	UpdateBy   string    `json:"updateBy" db:"UpdateBy"`
	UpdateDate time.Time `json:"updateDate" db:"UpdateDate"`
}

type ConfirmSaleReturnResponse struct {
	OrderNo     string    `json:"orderNo" db:"OrderNo"`
	ConfirmBy   string    `json:"confirmBy" db:"ConfirmBy"`
	ConfirmDate time.Time `json:"confirmDate" db:"ConfirmDate"`
}

type CancelSaleReturnResponse struct {
	RefID        string    `json:"refId" db:"RefID"`
	CancelStatus bool      `json:"cancelStatus" db:"CancelStatus"`
	Remark       string    `json:"remark" db:"Remark"`
	CancelBy     string    `json:"cancelBy" db:"CancelBy"`
	CancelDate   time.Time `json:"cancelDate" db:"CancelDate"`
}

type ListDraftConfirmOrdersResponse struct {
	OrderNo     string    `json:"orderNo" db:"OrderNo"`
	SoNo        string    `json:"soNo" db:"SoNo"`
	SrNo        string    `json:"srNo" db:"SrNo"`
	CustomerID  string    `json:"customerId" db:"CustomerID"`
	TrackingNo  string    `json:"trackingNo" db:"TrackingNo"`
	Logistic    string    `json:"logistic" db:"Logistic"`
	ChannelID   int       `json:"channelId" db:"ChannelID"`
	CreateDate  time.Time `json:"createDate" db:"CreateDate"`
	WarehouseID int       `json:"warehouseId" db:"WarehouseID"`
}

type DraftHeadResponse struct {
	OrderNo    string              `db:"OrderNo"`
	SoNo       string              `db:"SoNo"`
	SrNo       string              `db:"SrNo"`
	OrderLines []DraftLineResponse `db:"OrderLines"`
}

type DraftLineResponse struct {
	SKU      string  `db:"SKU"`
	ItemName string  `db:"ItemName"`
	QTY      int     `db:"QTY"`
	Price    float64 `db:"Price"`
}

type CodeRResponse struct {
	SKU       string `json:"sku" db:"SKU"`
	NameAlias string `json:"nameAlias" db:"NameAlias"`
}

// fa

type ConfirmReturnResponse struct {
	OrderNo     string    `json:"orderNo" db:"OrderNo"`
	ConfirmBy   string    `json:"confirmBy" db:"ConfirmBy"`
	ConfirmDate time.Time `json:"confirmDate" db:"ConfirmDate"`
}

type CancelReturnResponse struct {
	RefID        string    `json:"refId" db:"RefID"`
	CancelStatus bool      `json:"cancelStatus" db:"CancelStatus"`
	Remark       string    `json:"remark" db:"Remark"`
	CancelBy     string    `json:"cancelBy" db:"CancelBy"`
	CancelDate   time.Time `json:"cancelDate" db:"CancelDate"`
}
