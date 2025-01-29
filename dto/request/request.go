// ส่งข้อมูลเข้ามาด้านหลัง
package request

/*
type Login struct {
	UserID   string `json:"userID" db:"UserID" example:"DC65060"`
	Password string `json:"password" db:"Password" example:"xxxxxxxx"`
}

type LoginWeb struct {
	UserName string `json:"userName" db:"Username" example:"eknarin.ler"`
	Password string `json:"password" db:"Password" example:"xxxxxxxx"`
}

type LoginLark struct {
	UserID   string `json:"userID" db:"userID" example:"DC65060"`
	UserName string `json:"userName" db:"userName" example:"eknarin.ler"`
}

type LoginJWT struct {
	UserID   string `json:"userID" db:"UserID" example:"DC53002"`
	UserName string `json:"userName" db:"Username" example:"string"`
}

/********** Return Order (After) ***************/

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
	CancelID      *int    `json:"cancelId" db:"CancelID" example:"1"`
	StatusCheckID *int    `json:"statusCheckId" db:"StatusCheckID" example:"1"`
	CheckBy       *string `json:"checkBy" db:"CheckBy" example:"dev03"`
	Description   *string `json:"description" db:"Description" example:""`
	CreateBy      string  `json:"-" db:"CreateBy"`
	// CreateDate   *time.Time      `json:"createDate" db:"CreateDate"` // MSSQL SYSDATETIME() function

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
	UpdateBy      *string  `json:"-" db:"UpdateBy"` 
	// UpdateDate   *time.Time      `json:"updateDate" db:"UpdateDate"` // MSSQL SYSDATETIME() function

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
	Identifier  string                   `json:"-" `                                   // mean => OrderNo หรือ TrackingNo
	ImportLines []TradeReturnLineRequest `json:"importLines" validate:"required,dive"` // รายการสินค้า
}

type TradeReturnLineRequest struct {
	SKU       string   `json:"sku" db:"SKU" validate:"required"`
	ItemName  string   `json:"itemName" db:"ItemName" validate:"required"`
	QTY       int      `json:"qty" db:"QTY" validate:"required"`
	ReturnQTY int      `json:"returnQty" db:"ReturnQTY" validate:"required"`
	Price     float64  `json:"price" db:"Price" validate:"required"`
	//TrackingNo string  `json:"trackingNo" db:"TrackingNo"`	// add form data BeforeReturnOrder
	CreateBy string    `json:"-" db:"CreateBy" ` // from user login
	//CreateDate *time.Time `json:"createDate" db:"CreateDate"` // MSSQL GetDate()
	FilePath    string `json:"filePath" db:"FilePath" validate:"required"`
	ImageTypeID int    `json:"imageTypeID" db:"ImageTypeID" validate:"required"`
}

type Image struct {
	ImageTypeID int    `json:"imageTypeID" db:"ImageTypeID"` // ID ของประเภทของรูปภาพ
	FilePath    string `json:"-" db:"FilePath"`              // เส้นทางของไฟล์รูปภาพ
}

type TradeReturnLine struct {
	TradeReturnLine []TradeReturnLineRequest `json:"tradeReturnLine"`
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
