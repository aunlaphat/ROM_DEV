package response

import "time"

// BeforeReturnOrderLineResponse represents the response structure for return order line
type BeforeReturnOrderLineResponse struct {
	OrderNo    string    `json:"orderNo" db:"OrderNo"`
	SKU        string    `json:"sku" db:"SKU"`
	QTY        int       `json:"qty" db:"QTY"`
	ReturnQTY  int       `json:"returnQty" db:"ReturnQTY"`
	Price      float64   `json:"price" db:"Price"`
	TrackingNo string    `json:"trackingNo" db:"TrackingNo"`
	CreateDate time.Time `json:"createDate" db:"CreateDate"`
}

// BeforeReturnOrderResponse represents the response structure for a return order before processing
type BeforeReturnOrderResponse struct {
	OrderNo                string                          `json:"orderNo" db:"OrderNo"`
	SaleOrder              string                          `json:"saleOrder" db:"SaleOrder"`
	SaleReturn             string                          `json:"saleReturn" db:"SaleReturn"`
	ChannelID              int                             `json:"channelId" db:"ChannelID"`
	ReturnType             string                          `json:"returnType" db:"ReturnType"`
	CustomerID             string                          `json:"customerId" db:"CustomerID"`
	TrackingNo             string                          `json:"trackingNo" db:"TrackingNo"`
	Logistic               string                          `json:"logistic" db:"Logistic"`
	WarehouseID            int                             `json:"warehouseId" db:"WarehouseID"`
	SoStatusID             *int                            `json:"soStatusId" db:"SoStatusID"`
	MkpStatusID            *int                            `json:"mkpStatusId" db:"MkpStatusID"`
	ReturnDate             *time.Time                      `json:"returnDate" db:"ReturnDate"`
	StatusReturnID         int                             `json:"statusReturnId" db:"StatusReturnID"`
	StatusConfID           int                             `json:"statusConfId" db:"StatusConfID"`
	ConfirmBy              *string                         `json:"confirmBy" db:"ConfirmBy"`
	CreateBy               string                          `json:"createBy" db:"CreateBy"`
	CreateDate             time.Time                       `json:"createDate" db:"CreateDate"`
	UpdateBy               *string                         `json:"updateBy" db:"UpdateBy"`
	UpdateDate             *time.Time                      `json:"updateDate" db:"UpdateDate"`
	CancelID               *int                            `json:"cancelId" db:"CancelID"`
	BeforeReturnOrderLines []BeforeReturnOrderLineResponse `json:"beforeReturnOrderLines"`
}
