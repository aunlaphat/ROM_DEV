package request

import "time"

type ReturnOrder struct {
	ReturnID         string            `json:"returnID"`
	OrderNo          string            `json:"orderNo"`
	SaleOrder        string            `json:"saleOrder"`
	SaleReturn       string            `json:"saleReturn"`
	TrackingNo       string            `json:"trackingNo"`
	PlafID           int               `json:"plafID"`
	ChannelID        int               `json:"channelID"`
	OptStatusID      int               `json:"optStatusID"`
	AxStatusID       int               `json:"axStatusID"`
	PlatfStatusID    int               `json:"platfStatusID"`
	Remark           string            `json:"remark"`
	CreateBy         string            `json:"createBy"`
	CreateDate       *time.Time        `json:"createDate"`
	UpdateBy         string            `json:"updateBy"`
	UpdateDate       *time.Time        `json:"updateDate"`
	CancelID         int               `json:"cancelID"`
	StatusCheckID    int               `json:"statusCheckID"`
	CheckBy          string            `json:"checkBy"`
	Description      string            `json:"description"`
	ReturnOrderLines []ReturnOrderLine `json:"returnOrderLines"`
}

type ReturnOrderLine struct {
	RecID      int        `json:"recID"`
	ReturnID   string     `json:"returnID"`
	OrderNo    string     `json:"orderNo"`
	TrackingNo string     `json:"trackingNo"`
	SKU        string     `json:"sku"`
	QTY        int        `json:"qty"`
	ReturnQTY  int        `json:"returnQty"`
	CheckQTY   int        `json:"checkQty"`
	Price      float64    `json:"price"`
	CreateBy   string     `json:"createBy"`
	CreateDate *time.Time `json:"createDate"`
	AlterSKU   string     `json:"alterSKU"`
	UpdateBy   string     `json:"updateBy"`
	UpdateDate *time.Time `json:"updateDate"`
}

// Module 1: Return Order Creation
// BeforeReturnOrderRequest - สำหรับสร้าง return order ใหม่
type BeforeReturnOrderRequest struct {
	OrderNo     string                         `json:"orderNo" validate:"required"`
	SaleOrder   string                         `json:"saleOrder" validate:"required"`
	SaleReturn  string                         `json:"saleReturn" validate:"required"`
	ChannelID   int                            `json:"channelID" validate:"required"`
	ReturnType  string                         `json:"returnType" validate:"required"`
	CustomerID  string                         `json:"customerID" validate:"required"`
	TrackingNo  string                         `json:"trackingNo" validate:"required"`
	Logistic    string                         `json:"logistic" validate:"required"`
	WarehouseID int                            `json:"warehouseID" validate:"required"`
	ReturnLines []BeforeReturnOrderLineRequest `json:"returnLines" validate:"required,min=1,dive"`
	CreateBy    string                         `json:"createBy" validate:"required"`
	ReturnDate  time.Time                      `json:"returnDate" validate:"required"`
}

// @Description Line item details for return order
type BeforeReturnOrderLineRequest struct {
	SKU        string  `json:"sku" validate:"required"`
	QTY        int     `json:"qty" validate:"required,min=1"`
	ReturnQTY  int     `json:"returnQty" validate:"required,min=1,lte=qty"`
	Price      float64 `json:"price" validate:"required,gt=0"`
	TrackingNo string  `json:"trackingNo" validate:"required"`
}

// Module 3: Warehouse Receiving & Inspection
// @Description Request model for warehouse inspection
type ReceivedReturnOrderRequest struct {
	ReturnOrderNo string               `json:"returnOrderNo" validate:"required" example:"RET20240101001" description:"Return order reference number"`
	ReceivedDate  time.Time            `json:"receivedDate" validate:"required" example:"2024-01-01T10:00:00Z" description:"Date items received at warehouse"`
	CheckedBy     string               `json:"checkedBy" validate:"required" example:"STAFF001" description:"Staff ID who performed inspection"`
	CheckedItems  []CheckedItemRequest `json:"checkedItems" description:"Inspection results for each item"`
}

type CheckedItemRequest struct {
	SKU           string `json:"sku" validate:"required"`
	CheckedQty    int    `json:"checkedQty" validate:"required"`
	CheckResult   string `json:"checkResult" validate:"required,oneof=Correct Incorrect Short Excess"`
	ItemCondition string `json:"itemCondition" validate:"required,oneof=Good Damaged"`
}

// Module 4: Financial Processing
type CreditNoteRequest struct {
	ReturnOrderNo string  `json:"returnOrderNo" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Reason        string  `json:"reason" validate:"required"`
}

type DebitNoteRequest struct {
	ReturnOrderNo string  `json:"returnOrderNo" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Reason        string  `json:"reason" validate:"required"`
}

// ค่า constant สำหรับ validation
const (
	StatusPending   = "Pending"
	StatusBooking   = "Booking"
	StatusWaiting   = "Waiting"
	StatusSuccess   = "Success"
	StatusUnsuccess = "Unsuccess"
	StatusClosed    = "Closed"
)

type UpdateStatusRequest struct {
	OrderNo  string `json:"orderNo" validate:"required"`
	StatusID int    `json:"statusId" validate:"required"`
	UpdateBy string `json:"updateBy" validate:"required"`
}

type ProcessReturnRequest struct {
	OrderNo     string `json:"orderNo" validate:"required"`
	ProcessType string `json:"processType" validate:"required,oneof=Approve Reject"`
	ProcessBy   string `json:"processBy" validate:"required"`
	Reason      string `json:"reason,omitempty"`
}

type CancelOrderRequest struct {
	CancelBy string `json:"cancelBy" validate:"required"`
}
