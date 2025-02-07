package request

import "time"

type BeforeReturnOrder struct {
	//RecID		   int        `json:"recID" db:"RecID"` // (PK - Auto Increment)
	OrderNo        string     `json:"orderNo" db:"OrderNo"`
	SoNo           string     `json:"soNo" db:"SoNo"`
	SrNo           string     `json:"srNo" db:"SrNo"`
	ChannelID      int        `json:"channelID" db:"ChannelID"`
	Reason         string     `json:"reason" db:"Reason"`
	CustomerID     string     `json:"customerID" db:"CustomerID"`
	TrackingNo     string     `json:"trackingNo" db:"TrackingNo"`
	Logistic       string     `json:"logistic" db:"Logistic"`
	WarehouseID    int        `json:"warehouseID" db:"WarehouseID"`
	SoStatus       *int       `json:"soStatus" db:"SoStatus"`
	MkpStatus      *int       `json:"mkpStatus" db:"MkpStatus"`
	ReturnDate     *time.Time `json:"returnDate" db:"ReturnDate"`
	StatusReturnID *int       `json:"statusReturnID" db:"StatusReturnID"`
	StatusConfID   *int       `json:"statusConfID" db:"StatusConfID"`
	ConfirmBy      *string    `json:"confirmBy" db:"ConfirmBy"`
	//ConfirmDate  *time.Time `json:"confirmDate" db:"ConfirmDate"`
	CreateBy string `json:"createBy" db:"CreateBy"`
	//CreateDate  *time.Time `json:"createDate" db:"CreateDate"`
	UpdateBy *string `json:"updateBy" db:"UpdateBy"`
	//UpdateDate   *time.Time `json:"updateDate" db:"UpdateDate"`
	CancelID               *int                    `json:"cancelID" db:"CancelID"`
	BeforeReturnOrderLines []BeforeReturnOrderLine `json:"beforeReturnOrderLines"`
}

type BeforeReturnOrderLine struct {
	//RecID		  int        `json:"recID" db:"RecID"` // (PK - Auto Increment)
	OrderNo   string  `json:"orderNo" db:"OrderNo"`
	SKU       string  `json:"sku" db:"SKU"`
	ItemName  string  `json:"itemName" db:"ItemName"`
	QTY       int     `json:"qty" db:"QTY"`
	ReturnQTY int     `json:"returnQty" db:"ReturnQTY"`
	Price     float64 `json:"price" db:"Price"`
	CreateBy  string  `json:"createBy" db:"CreateBy"`
	//CreateDate *time.Time `json:"createDate" db:"CreateDate"`
	UpdateBy *string `json:"updateBy" db:"UpdateBy"`
	//UpdateDate *time.Time `json:"updateDate" db:"UpdateDate"`
	TrackingNo *string `json:"trackingNo" db:"TrackingNo"`
	AlterSKU   *string `json:"alterSKU" db:"AlterSKU"`
}

type UpdateSaleReturn struct {
	SrNo string `json:"srNo" validate:"required" example:"SR-TEST-123456"`
}

type CancelSaleReturn struct {
	//OrderNo      string `json:"orderNo" db:"OrderNo"`
	//CancelStatus bool   `json:"cancelStatus" db:"CancelStatus"`
	Remark string `json:"remark" db:"Remark"`
}

// Draft & Confirm MKP

type CodeR struct {
	OrderNo   string  `json:"orderNo" db:"OrderNo"`
	SKU       string  `json:"sku" db:"SKU"`
	ItemName  string  `json:"itemName" db:"ItemName"`
	QTY       int     `json:"qty" db:"QTY"`
	ReturnQTY int     `json:"returnQty" db:"ReturnQTY"`
	Price     float64 `json:"price" db:"Price"`
	CreateBy  string  `json:"createBy" db:"CreateBy"`
}

type CreateReturnOrder struct {
	OrderNo       string  `json:"orderNo" db:"OrderNo" example:"ORD0001"`
	SoNo          string  `json:"soNo" db:"SoNo" example:"SO0001"`
	SrNo          string  `json:"srNo" db:"SrNo" example:"SR0001"`
	TrackingNo    string  `json:"trackingNo" db:"TrackingNo" example:"12345678TH"`
	PlatfID       *int    `json:"platfId" db:"PlatfID" example:"1"`
	ChannelID     *int    `json:"channelId" db:"ChannelID" example:"2"`
	OptStatusID   *int    `json:"optStatusId" db:"OptStatusID" example:"1"`
	AxStatusID    *int    `json:"axStatusId" db:"AxStatusID" example:"1"`
	PlatfStatusID *int    `json:"platfStatusId" db:"PlatfStatusID" example:"1"`
	Reason        *string `json:"reason" db:"Reason"`
	StatusCheckID *int    `json:"statusCheckId" db:"StatusCheckID" example:"1"`
	Description   *string `json:"description" db:"Description" example:""`
	CreateBy      string  `json:"-" db:"CreateBy"`
	// CreateDate   *time.Time      `json:"createDate" db:"CreateDate"` // MSSQL GETDATE() function

	ReturnOrderLine []ReturnOrderLine `json:"ReturnOrderLine"`
}

type UpdateReturnOrder struct {
	OrderNo       string  `json:"-" db:"OrderNo"`
	SoNo          string  `json:"-" db:"SoNo"`
	SrNo          *string `json:"srNo" db:"SrNo" example:"SR0001"`
	TrackingNo    *string `json:"trackingNo" db:"TrackingNo" example:"12345678TH"`
	PlatfID       *int    `json:"platfId" db:"PlatfID" example:"1"`
	ChannelID     *int    `json:"channelId" db:"ChannelID" example:"2"`
	OptStatusID   *int    `json:"optStatusId" db:"OptStatusID" example:"1"`
	AxStatusID    *int    `json:"axStatusId" db:"AxStatusID" example:"1"`
	PlatfStatusID *int    `json:"platfStatusId" db:"PlatfStatusID" example:"1"`
	Reason        *string `json:"reason" db:"Reason" example:"CHANGE PRODUCTS"`
	CancelID      *int    `json:"cancelId" db:"CancelID" example:"1"`
	StatusCheckID *int    `json:"statusCheckId" db:"StatusCheckID" example:"1"`
	CheckBy       *string `json:"checkBy" db:"CheckBy" example:"dev03"`
	Description   *string `json:"description" db:"Description" example:""`
	UpdateBy      *string `json:"-" db:"UpdateBy"`
	// UpdateDate   *time.Time      `json:"updateDate" db:"UpdateDate"` // MSSQL GETDATE() function

}

type ReturnOrder struct {
	OrderNo     string `json:"orderNo" db:"OrderNo"`
	SoNo        string `json:"soNo" db:"SoNo"`
	SrNo        string `json:"srNo" db:"SrNo"`
	TrackingNo  string `json:"trackingNo" db:"TrackingNo"`
	PlatfID     int    `json:"platfID" db:"PlatfID"`
	ChannelID   int    `json:"channelID" db:"ChannelID"`
	OptStatusID int    `json:"optStatusID" db:"OptStatusID"`
	AxStatusID  int    `json:"axStatusID" db:"AxStatusID"`
	Reason      string `json:"reason" db:"Reason"`
	CreateBy    string `json:"createBy" db:"CreateBy"`
	//CreateDate       time.Time         `json:"createDate" db:"CreateDate"`
	UpdateBy string `json:"updateBy" db:"UpdateBy"`
	//UpdateDate       time.Time         `json:"updateDate" db:"UpdateDate"`
	CancelID         *int              `json:"cancelID" db:"CancelID"`
	ReturnOrderLines []ReturnOrderLine `json:"returnOrderLines"`
}

type ReturnOrderLine struct {
	OrderNo    string  `json:"-" db:"OrderNo"`
	TrackingNo string  `json:"-" db:"TrackingNo"`
	SKU        string  `json:"sku" db:"SKU" example:"SKU12345"`
	QTY        *int    `json:"qty" db:"QTY" example:"5"`
	ReturnQTY  int     `json:"returnQTY" db:"ReturnQTY" example:"5"`
	Price      float64 `json:"price" db:"Price" example:"199.99"`
	AlterSKU   *string `json:"-" db:"AlterSKU" `
}

/********** Trade Return (Offline) ***************/

type ConfirmTradeReturnRequest struct {
	Identifier  string                   `json:"-" `          // mean => OrderNo หรือ TrackingNo
	ImportLines []TradeReturnLineRequest `json:"importLines"` // รายการสินค้า
}

type TradeReturnLineRequest struct {
	SKU       string  `json:"sku" db:"SKU"`
	ItemName  string  `json:"itemName" db:"ItemName"`
	QTY       int     `json:"qty" db:"QTY"`
	ReturnQTY int     `json:"returnQty" db:"ReturnQTY"`
	Price     float64 `json:"price" db:"Price"`
	//TrackingNo string  `json:"trackingNo" db:"TrackingNo"`	// add form data BeforeReturnOrder
	CreateBy string `json:"-" db:"CreateBy" ` // from user login
	//CreateDate *time.Time `json:"createDate" db:"CreateDate"` // MSSQL GetDate()
	FilePath    string `json:"filePath" db:"FilePath"`       // เข้า Images
	ImageTypeID int    `json:"imageTypeID" db:"ImageTypeID"` // เข้า Images
}

type Image struct {
	ImageTypeID int    `json:"imageTypeID" db:"ImageTypeID"` // ID ของประเภทของรูปภาพ
	FilePath    string `json:"-" db:"FilePath"`              // เส้นทางของไฟล์รูปภาพ
}

type TradeReturnLine struct {
	TradeReturnLine []OrderLines `json:"tradeReturnLine"`
}

type OrderLines struct {
	SKU       string  `json:"sku" db:"SKU"`
	ItemName  string  `json:"itemName" db:"ItemName"`
	QTY       int     `json:"qty" db:"QTY"`
	ReturnQTY int     `json:"returnQty" db:"ReturnQTY"`
	Price     float64 `json:"price" db:"Price"`
	//TrackingNo string  `json:"trackingNo" db:"TrackingNo"`	// add form data BeforeReturnOrder
	CreateBy string `json:"-" db:"CreateBy"` // from user login
	//CreateDate *time.Time `json:"createDate" db:"CreateDate"` // MSSQL GetDate()
}

/********** Import Order to Warehouse: Sale Return ***************/

type Images struct {
	// ReturnID    string `json:"returnID" db:"ReturnID"`
	SKU         string `json:"sku" db:"SKU"`
	FilePath    string `json:"filePath" db:"FilePath"`
	ImageTypeID int    `json:"imageTypeID" db:"ImageTypeID"`
	CreateBy    string `json:"createBy" db:"CreateBy"`
	OrderNo     string `json:"orderNo" db:"OrderNo"`
}

type ConfirmToReturnRequest struct {
	OrderNo           string              `json:"-"`
	UpdateToReturn    []UpdateToReturn    `json:"updateToReturn"`    // เลข sr สุ่มจาก ax
	ImportLinesActual []ImportLinesActual `json:"importLinesActual"` // รายการสินค้าที่ผ่านการเช็คแล้วจากบัญชี
}

type UpdateToReturn struct {
	SrNo string `json:"srNo" db:"SrNo"`
}

type ImportLinesActual struct {
	SKU          string  `json:"sku" db:"SKU"`
	ActualQTY    int     `json:"actualQty" db:"ActualQTY"`
	Price        float64 `json:"price" db:"Price"`
	StatusDelete bool    `json:"statusDelete" db:"StatusDelete"`
}
