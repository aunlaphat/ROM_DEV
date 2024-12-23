package request

import "time"

type BeforeReturnOrder struct {
	OrderNo                string                  `json:"orderNo"`
	SaleOrder              string                  `json:"saleOrder"`
	SaleReturn             string                  `json:"saleReturn"`
	ChannelID              int                     `json:"channelID"`
	ReturnType             string                  `json:"returnType"`
	CustomerID             string                  `json:"customerID"`
	TrackingNo             string                  `json:"trackingNo"`
	Logistic               string                  `json:"logistic"`
	WarehouseID            int                     `json:"warehouseID"`
	SoStatusID             *int                    `json:"soStatusID"`
	MkpStatusID            *int                    `json:"mkpStatusID"`
	ReturnDate             *time.Time              `json:"returnDate"`
	StatusReturnID         int                     `json:"statusReturnID"`
	StatusConfID           int                     `json:"statusConfID"`
	ConfirmBy              *string                 `json:"confirmBy"`
	CreateBy               string                  `json:"createBy"`
	CreateDate             *time.Time              `json:"createDate"`
	UpdateBy               *string                 `json:"updateBy"`
	UpdateDate             *time.Time              `json:"updateDate"`
	CancelID               *int                    `json:"cancelID"`
	BeforeReturnOrderLines []BeforeReturnOrderLine `json:"returnOrderLines"`
}
type BeforeReturnOrderLine struct {
	OrderNo    string     `json:"orderNo"`
	SKU        string     `json:"sku"`
	QTY        int        `json:"qty"`
	ReturnQTY  int        `json:"returnQty"`
	Price      float64    `json:"price"`
	CreateBy   string     `json:"createBy"`
	CreateDate *time.Time `json:"createDate"`
	TrackingNo string     `json:"trackingNo"`
	AlterSKU   *string    `่หนื:"alterSKU"`
	UpdateBy   *string    `json:"updateBy"`
	UpdateDate *time.Time `json:"updateDate"`
}
