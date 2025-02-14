package response

import "time"

type SearchOrderResponse struct {
	SoNo        string            `db:"SoNo" json:"soNo"`
	OrderNo     string            `db:"OrderNo" json:"orderNo"`
	StatusMKP   string            `db:"StatusMKP" json:"statusMKP"`
	SalesStatus string            `db:"SalesStatus" json:"salesStatus"`
	CreateDate  time.Time         `db:"CreateDate" json:"createDate"`
	Items       []SearchOrderItem `json:"items"`
}

type SearchOrderItem struct {
	SKU      string  `db:"SKU" json:"sku"`
	ItemName string  `db:"ItemName" json:"itemName"`
	QTY      int     `db:"QTY" json:"qty"`
	Price    float64 `db:"Price" json:"price"`
}

type BeforeReturnOrderResponse struct {
	OrderNo        string                  `json:"orderNo" db:"OrderNo"`
	SoNo           string                  `json:"soNo" db:"SoNo"`
	SrNo           *string                 `json:"srNo" db:"SrNo"`
	ChannelID      int                     `json:"channelId" db:"ChannelID"`
	Reason         string                  `json:"reason" db:"Reason"`
	CustomerID     string                  `json:"customerId" db:"CustomerID"`
	TrackingNo     string                  `json:"trackingNo" db:"TrackingNo"`
	Logistic       string                  `json:"logistic" db:"Logistic"`
	WarehouseID    int                     `json:"warehouseId" db:"WarehouseID"`
	SoStatus       *string                 `json:"soStatus" db:"SoStatus"`
	MkpStatus      *string                 `json:"mkpStatus" db:"MkpStatus"`
	ReturnDate     *time.Time              `json:"returnDate" db:"ReturnDate"`
	StatusReturnID *int                    `json:"statusReturnId" db:"StatusReturnID"`
	StatusConfID   *int                    `json:"statusConfId" db:"StatusConfID"`
	ConfirmBy      *string                 `json:"confirmBy" db:"ConfirmBy"`
	ConfirmDate    *time.Time              `json:"confirmDate" db:"ConfirmDate"`
	CreateBy       string                  `json:"createBy" db:"CreateBy"`
	CreateDate     time.Time               `json:"createDate" db:"CreateDate"`
	UpdateBy       *string                 `json:"updateBy" db:"UpdateBy"`
	UpdateDate     *time.Time              `json:"updateDate" db:"UpdateDate"`
	CancelID       *int                    `json:"cancelId" db:"CancelID"`
	IsCNCreated    bool                    `json:"isCNCreated" db:"IsCNCreated"`
	IsEdited       bool                    `json:"isEdited" db:"IsEdited"`
	Items          []BeforeReturnOrderItem `json:"items"`
}

type BeforeReturnOrderItem struct {
	OrderNo    string    `json:"orderNo" db:"OrderNo"`
	SKU        string    `json:"sku" db:"SKU"`
	ItemName   string    `json:"itemName" db:"ItemName"`
	QTY        int       `json:"qty" db:"QTY"`
	ReturnQTY  int       `json:"returnQty" db:"ReturnQTY"`
	Price      float64   `json:"price" db:"Price"`
	CreateBy   string    `json:"createBy" db:"CreateBy"`
	CreateDate time.Time `json:"createDate" db:"CreateDate"`
	TrackingNo *string   `json:"trackingNo,omitempty" db:"TrackingNo"`
	AlterSKU   *string   `json:"alterSKU,omitempty" db:"AlterSKU"`
}

type UpdateSrNoResponse struct {
	OrderNo        string    `json:"orderNo" db:"OrderNo"`
	SrNo           string    `json:"srNo" db:"SrNo"`
	StatusReturnID *int      `json:"statusReturnID,omitempty" db:"StatusReturnID"`
	StatusConfID   *int      `json:"statusConfID,omitempty" db:"StatusConfID"`
	UpdateBy       string    `json:"updateBy" db:"UpdateBy"`
	UpdateDate     time.Time `json:"updateDate" db:"UpdateDate"`
}

type UpdateOrderStatusResponse struct {
	OrderNo        string    `json:"orderNo" db:"OrderNo"`
	StatusReturnID int       `json:"statusReturnID" db:"StatusReturnID"`
	StatusConfID   int       `json:"statusConfID" db:"StatusConfID"`
	ConfirmBy      string    `json:"confirmBy" db:"ConfirmBy"`
	ConfirmDate    time.Time `json:"confirmDate" db:"ConfirmDate"`
}

type CancelOrderResponse struct {
	RefID        string    `json:"refID" db:"RefID"`
	SourceTable  string    `json:"sourceTable" db:"SourceTable"`
	CancelReason string    `json:"cancelReason" db:"CancelReason"`
	CancelBy     string    `json:"cancelBy" db:"CancelBy"`
	CancelDate   time.Time `json:"cancelDate" db:"CancelDate"`
}

type OrderHeadResponse struct {
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

type DraftConfirmResponse struct {
	OrderNo string             `json:"orderNo" db:"OrderNo"`
	SoNo    string             `json:"soNo" db:"SoNo"`
	SrNo    string             `json:"srNo" db:"SrNo"`
	Items   []DraftConfirmItem `json:"items"`
}

type DraftConfirmItem struct {
	OrderNo  string  `json:"orderNo" db:"OrderNo"`
	SKU      string  `json:"sku" db:"SKU"`
	ItemName string  `json:"itemName" db:"ItemName"`
	QTY      int     `json:"qty" db:"QTY"`
	Price    float64 `json:"price" db:"Price"`
}

type ListCodeRResponse struct {
	SKU       string `json:"sku" db:"SKU"`
	NAMEALIAS string `json:"nameAlias" db:"NAMEALIAS"`
}

type AddItemResponse struct {
	OrderNo    string    `json:"orderNo" db:"OrderNo"`
	SKU        string    `json:"sku" db:"SKU"`
	ItemName   string    `json:"itemName" db:"ItemName"`
	QTY        int       `json:"qty" db:"QTY"`
	ReturnQTY  int       `json:"returnQty" db:"ReturnQTY"`
	Price      float64   `json:"price" db:"Price"`
	CreateBy   string    `json:"createBy" db:"CreateBy"`
	CreateDate time.Time `json:"createDate" db:"CreateDate"`
}
