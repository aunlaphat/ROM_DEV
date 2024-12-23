package response

import "time"

// ReturnOrder represents the details of a return order
type ReturnOrder struct {
	ReturnID         string            `json:"returnID" db:"ReturnID"`
	OrderNo          string            `json:"orderNo" db:"OrderNo"`
	SaleOrder        string            `json:"SaleOrder" db:"SaleOrder"`
	SaleReturn       string            `json:"saleReturn" db:"SaleReturn"`
	TrackingNo       string            `json:"trackingNo" db:"TrackingNo"`
	PlafID           int               `json:"plafID" db:"PlafID"`
	ChannelID        int               `json:"channelID" db:"ChannelID"`
	OptStatusID      int               `json:"optStatusID" db:"OptStatusID"`
	AxStatusID       int               `json:"axStatusID" db:"AxStatusID"`
	PlatfStatusID    int               `json:"platfStatusID" db:"PlatfStatusID"`
	Remark           string            `json:"remark" db:"Remark"`
	CreateBy         string            `json:"createBy" db:"CreateBy"`
	CreateDate       *time.Time        `json:"createDate" db:"CreateDate"`
	UpdateBy         string            `json:"updateBy" db:"UpdateBy"`
	UpdateDate       *time.Time        `json:"updateDate" db:"UpdateDate"`
	CancelID         int               `json:"cancelID" db:"CancelID"`
	StatusCheckID    int               `json:"statusCheckID" db:"StatusCheckID"`
	CheckBy          string            `json:"checkBy" db:"CheckBy"`
	Description      string            `json:"description" db:"Description"`
	ReturnOrderLines []ReturnOrderLine `json:"returnOrderLines"`
}

// ReturnOrderLine represents the details of a return order line
type ReturnOrderLine struct {
	RecID      int        `json:"recID" db:"RecID"`
	ReturnID   string     `json:"returnID" db:"ReturnID"`
	OrderNo    string     `json:"orderNo" db:"OrderNo"`
	TrackingNo string     `json:"trackingNo" db:"TrackingNo"`
	SKU        string     `json:"sku" db:"SKU"`
	QTY        int        `json:"qty" db:"QTY"`
	ReturnQTY  int        `json:"returnQty" db:"ReturnQTY"`
	CheckQTY   int        `json:"checkQty" db:"CheckQTY"`
	Price      float64    `json:"price" db:"Price"`
	CreateBy   string     `json:"createBy" db:"CreateBy"`
	CreateDate *time.Time `json:"createDate" db:"CreateDate"`
	AlterSKU   string     `json:"alterSKU" db:"AlterSKU"`
	UpdateBy   string     `json:"updateBy" db:"UpdateBy"`
	UpdateDate *time.Time `json:"updateDate" db:"UpdateDate"`
}

// ReturnOrderList represents a list of return orders
type ReturnOrderList struct {
	Orders []ReturnOrder `json:"orders"`
	Total  int           `json:"total"`
}

// ReturnOrderResponse represents the response for a single return order
type ReturnOrderResponse struct {
	Success bool        `json:"success"`
	Data    ReturnOrder `json:"data"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// BeforeReturnOrderResponse represents the response after creating a return order
type BeforeReturnOrderResponse struct {
	Success   bool                   `json:"success"`
	OrderNo   string                 `json:"orderNo"`
	Message   string                 `json:"message,omitempty"`
	ErrorCode string                 `json:"errorCode,omitempty"`
	Data      *BeforeReturnOrderData `json:"data,omitempty"`
}

// BeforeReturnOrderData represents the details of a return order before it is processed
type BeforeReturnOrderData struct {
	OrderNo        string                      `json:"orderNo" db:"OrderNo"`
	SaleOrder      string                      `json:"saleOrder" db:"SaleOrder"`
	SaleReturn     string                      `json:"saleReturn" db:"SaleReturn"`
	ChannelID      int                         `json:"channelId" db:"ChannelID"`
	ReturnType     string                      `json:"returnType" db:"ReturnType"`
	CustomerID     string                      `json:"customerId" db:"CustomerID"`
	TrackingNo     string                      `json:"trackingNo" db:"TrackingNo"`
	Logistic       string                      `json:"logistic" db:"Logistic"`
	WarehouseID    int                         `json:"warehouseId" db:"WarehouseID"`
	StatusReturnID int                         `json:"statusReturnId" db:"StatusReturnID"`
	StatusConfID   int                         `json:"statusConfId" db:"StatusConfID"`
	CreateDate     time.Time                   `json:"createDate" db:"CreateDate"`
	CreateBy       string                      `json:"createBy" db:"CreateBy"`
	ReturnLines    []BeforeReturnOrderLineData `json:"returnLines"`
}

// BeforeReturnOrderLineData represents the details of a return order line before it is processed
type BeforeReturnOrderLineData struct {
	OrderNo    string    `json:"orderNo" db:"OrderNo"`
	SKU        string    `json:"sku" db:"SKU"`
	QTY        int       `json:"qty" db:"QTY"`
	ReturnQTY  int       `json:"returnQty" db:"ReturnQTY"`
	Price      float64   `json:"price" db:"Price"`
	TrackingNo string    `json:"trackingNo" db:"TrackingNo"`
	CreateDate time.Time `json:"createDate" db:"CreateDate"`
}

// BeforeReturnOrderListResponse represents the response for a paginated list of return orders
type BeforeReturnOrderListResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message,omitempty"`
	Page    int                     `json:"page"`
	Limit   int                     `json:"limit"`
	Total   int                     `json:"total"`
	Orders  []BeforeReturnOrderData `json:"orders"`
}

// ReturnPredictionResponse represents the response for return order prediction
type ReturnPredictionResponse struct {
	ReturnID      string    `json:"returnId"`
	PredictedDate time.Time `json:"predictedDate"`
	Status        string    `json:"status"`
	IsOverdue     bool      `json:"isOverdue"`
}

// WarehouseInspectionResponse represents the response for warehouse inspection
type WarehouseInspectionResponse struct {
	ReturnOrderNo  string            `json:"returnOrderNo" example:"RET20240101001"`
	InspectionDate time.Time         `json:"inspectionDate" example:"2024-01-01T10:00:00Z"`
	CheckedBy      string            `json:"checkedBy" example:"STAFF001"`
	Status         string            `json:"status" example:"Success" enums:"Success,Unsuccess"`
	CheckedItems   []CheckedItem     `json:"checkedItems"`
	Summary        InspectionSummary `json:"summary"`
}

// CheckedItem represents the details of a checked item during warehouse inspection
type CheckedItem struct {
	SKU           string `json:"sku"`
	ExpectedQty   int    `json:"expectedQty"`
	ReceivedQty   int    `json:"receivedQty"`
	CheckResult   string `json:"checkResult"`
	ItemCondition string `json:"itemCondition"`
}

// InspectionSummary represents the summary of a warehouse inspection
type InspectionSummary struct {
	TotalItems   int     `json:"totalItems"`
	CorrectItems int     `json:"correctItems"`
	WrongItems   int     `json:"wrongItems"`
	TotalAmount  float64 `json:"totalAmount"`
}

// FinancialProcessResponse represents the response for financial processing
type FinancialProcessResponse struct {
	ReturnOrderNo string    `json:"returnOrderNo"`
	DocumentNo    string    `json:"documentNo"`   // CN or DN number
	DocumentType  string    `json:"documentType"` // "CN" or "DN"
	Amount        float64   `json:"amount"`
	ProcessDate   time.Time `json:"processDate"`
	Status        string    `json:"status"`
}

// ReturnStatisticsResponse represents the response for return order statistics
type ReturnStatisticsResponse struct {
	Period       string                 `json:"period" example:"2024-01"`
	TotalReturns int                    `json:"totalReturns" example:"150"`
	SuccessRate  float64                `json:"successRate" example:"95.5"`
	Statistics   map[string]interface{} `json:"statistics" example:"{\"totalAmount\":15000.00,\"avgProcessingTime\":\"2d 4h\"}"`
}
