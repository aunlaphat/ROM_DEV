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

type ImportItem struct {
	OrderNo    string  `json:"orderNo" db:"OrderNo"`
	TrackingNo *string `json:"trackingNo" db:"TrackingNo"`
}

// ในฟิลมีปรับ type ข้อมูลเพิ่ม //11/02
type BeforeReturnOrderResponse struct {
	OrderNo        string                  `json:"orderNo" db:"OrderNo"`
	SoNo           string                  `json:"soNo" db:"SoNo"`
	SrNo           *string                 `json:"srNo" db:"SrNo"`
	ChannelID      int                     `json:"channelId" db:"ChannelID"`
	Reason         string                  `json:"reason" db:"Reason"`
	CustomerID     string                  `json:"customerId" db:"CustomerID"`
	TrackingNo     *string                 `json:"trackingNo" db:"TrackingNo"`
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
	OrderNo     string    `json:"orderNo" db:"OrderNo"`
	SKU         string    `json:"sku" db:"SKU"`
	ItemName    string    `json:"itemName" db:"ItemName"`
	QTY         int       `json:"qty" db:"QTY"`
	ReturnQTY   int       `json:"returnQty" db:"ReturnQTY"`
	Price       float64   `json:"price" db:"Price"`
	WarehouseID *int      `json:"warehouseID" db:"WarehouseID"`
	CreateBy    string    `json:"createBy" db:"CreateBy"`
	CreateDate  time.Time `json:"createDate" db:"CreateDate"`
	TrackingNo  *string   `json:"trackingNo,omitempty" db:"TrackingNo"`
	AlterSKU    *string   `json:"alterSKU,omitempty" db:"AlterSKU"`
}

/********** Return Order ***************/

type ReturnOrder struct {
	OrderNo       string     `json:"orderNo" db:"OrderNo"`
	SoNo          string     `json:"soNo" db:"SoNo"`
	SrNo          *string    `json:"srNo" db:"SrNo"`
	TrackingNo    *string    `json:"trackingNo" db:"TrackingNo"`
	PlatfID       *int       `json:"platfId" db:"PlatfID"`
	ChannelID     *int       `json:"channelId" db:"ChannelID"`
	OptStatusID   *int       `json:"optStatusId" db:"OptStatusID"`
	AxStatusID    *int       `json:"axStatusId" db:"AxStatusID"`
	PlatfStatusID *int       `json:"platfStatusId" db:"PlatfStatusID"`
	Reason        *string    `json:"reason" db:"Reason"`
	CreateBy      string     `json:"createBy" db:"CreateBy"`
	CreateDate    time.Time  `json:"createDate" db:"CreateDate"`
	UpdateBy      *string    `json:"updateBy" db:"UpdateBy"`
	UpdateDate    *time.Time `json:"updateDate" db:"UpdateDate"`
	CancelID      *int       `json:"cancelId" db:"CancelID"`
	StatusCheckID *int       `json:"statusCheckId" db:"StatusCheckID"`
	CheckBy       *string    `json:"checkBy" db:"CheckBy"`
	CheckDate     *time.Time `json:"checkDate" db:"CheckDate"`
	Description   *string    `json:"description" db:"Description"`

	ReturnOrderLine []ReturnOrderLine `json:"ReturnOrderLine"`
}

type ReturnOrderLine struct {
	OrderNo    string     `json:"orderNo" db:"OrderNo"`
	TrackingNo *string    `json:"trackingNo" db:"TrackingNo"`
	SKU        string     `json:"sku" db:"SKU"`
	ItemName   *string    `json:"itemName" db:"ItemName"`
	ReturnQTY  int        `json:"returnQTY" db:"ReturnQTY"`
	ActualQTY  *int       `json:"actualQTY" db:"ActualQTY"`
	QTY        int        `json:"qty" db:"QTY"`
	Price      float64    `json:"price" db:"Price"`
	CreateBy   string     `json:"createBy" db:"CreateBy"`
	CreateDate time.Time  `json:"createDate" db:"CreateDate"`
	AlterSKU   *string    `json:"alterSKU" db:"AlterSKU"`
	UpdateBy   *string    `json:"updateBy" db:"UpdateBy"`
	UpdateDate *time.Time `json:"updateDate" db:"UpdateDate"`
}

type CreateReturnOrder struct {
	OrderNo       string    `json:"orderNo" db:"OrderNo" example:"ORD0001"`
	SoNo          string    `json:"soNo" db:"SoNo" example:"SO0001"`
	SrNo          *string   `json:"srNo" db:"SrNo" example:"SR0001"`
	TrackingNo    *string   `json:"trackingNo" db:"TrackingNo" example:"12345678TH"`
	PlatfID       int       `json:"platfId" db:"PlatfID" example:"1"`
	ChannelID     int       `json:"channelId" db:"ChannelID" example:"2"`
	OptStatusID   int       `json:"optStatusId" db:"OptStatusID" example:"1"`
	AxStatusID    int       `json:"axStatusId" db:"AxStatusID" example:"1"`
	PlatfStatusID int       `json:"platfStatusId" db:"PlatfStatusID" example:"1"`
	Reason        string    `json:"reason" db:"Reason"`
	StatusCheckID int       `json:"statusCheckId" db:"StatusCheckID" example:"1"`
	Description   string    `json:"description" db:"Description" example:""`
	CreateBy      string    `json:"createBy" db:"CreateBy"`
	CreateDate    time.Time `json:"createDate" db:"CreateDate"`

	ReturnOrderLine []ReturnOrderLine `json:"ReturnOrderLine"`
}

type UpdateReturnOrder struct {
	OrderNo       string     `json:"-" db:"OrderNo"`
	SoNo          string     `json:"-" db:"SoNo"`
	SrNo          *string    `json:"srNo" db:"SrNo" example:"SR0001"`
	TrackingNo    *string    `json:"trackingNo" db:"TrackingNo" example:"12345678TH"`
	PlatfID       *int       `json:"platfId" db:"PlatfID" example:"1"`
	ChannelID     *int       `json:"channelId" db:"ChannelID" example:"2"`
	OptStatusID   *int       `json:"optStatusId" db:"OptStatusID" example:"1"`
	AxStatusID    *int       `json:"axStatusId" db:"AxStatusID" example:"1"`
	PlatfStatusID *int       `json:"platfStatusId" db:"PlatfStatusID" example:"1"`
	Reason        *string    `json:"reason" db:"Reason"`
	CancelID      *int       `json:"cancelId" db:"CancelID" example:"1"`
	StatusCheckID *int       `json:"statusCheckId" db:"StatusCheckID" example:"1"`
	CheckBy       *string    `json:"checkBy" db:"CheckBy" example:"dev03"`
	CheckDate     *time.Time `json:"checkDate" db:"CheckDate"`
	Description   *string    `json:"description" db:"Description" example:""`
	UpdateBy      *string    `json:"updateBy" db:"UpdateBy"`
	UpdateDate    *time.Time `json:"updateDate" db:"UpdateDate"`
}

type UpdateReturnOrderLine struct {
	OrderNo    string   `json:"-" db:"OrderNo"`
	TrackingNo *string  `json:"-" db:"TrackingNo"`
	SKU        string   `json:"sku" db:"SKU" example:"SKU12345"`
	QTY        *int     `json:"qty" db:"QTY" example:"5"`
	ReturnQTY  int      `json:"returnQTY" db:"ReturnQTY" example:"5"`
	ActualQTY  *int     `json:"actualQTY" db:"ActualQTY" example:"5"`
	Price      *float64 `json:"price" db:"Price" example:"199.99"`
	AlterSKU   *string  `json:"-" db:"AlterSKU" `
	UpdateBy   *string  `json:"updateBy" db:"UpdateBy"`
}

type DeleteReturnOrder struct {
	OrderNo string `db:"OrderNo"`
}

type DraftTradeDetail struct {
	OrderNo       string    `json:"orderNo" db:"OrderNo" example:"ORD0001"`
	SoNo          string    `json:"soNo" db:"SoNo" example:"SO0001"`
	CustomerID    *string   `json:"customerId" db:"CustomerID"`
	SrNo          *string   `json:"srNo" db:"SrNo" example:"SR0001"`
	TrackingNo    *string   `json:"trackingNo" db:"TrackingNo" example:"12345678TH"`
	Logistic      *string   `json:"logistic" db:"Logistic"`
	ChannelID     *int      `json:"channelId" db:"ChannelID" example:"2"`
	ChannelName   *string   `json:"channelName" db:"ChannelName"`
	CreateDate    time.Time `json:"createDate" db:"CreateDate"`
	WarehouseID   *int      `json:"warehouseId" db:"WarehouseID"`
	WarehouseName *string   `json:"warehouseName" db:"WarehouseName"`
	StatusCheckID int       `json:"statusCheckId" db:"StatusCheckID" example:"1"`
	// Reason        *string   `json:"reason" db:"Reason"`
	// CreateBy      string    `json:"createBy" db:"CreateBy"`
}

/********** OrderHead + Line data Project ***************/

type OrderDetail struct {
	// json => OrderHeadDetail[ OrderLineDetail[ {},{},..] ]
	OrderHeadDetail []OrderHeadDetail `json:"OrderHeadDetail"`
}

type OrderHeadDetail struct {
	OrderNo     string    `db:"OrderNo" json:"orderNo"`         // เลขที่ใบสั่งซื้อ
	SoNo        string    `db:"SoNo" json:"soNo"`               // เลขที่ใบสั่งขาย
	StatusMKP   string    `db:"StatusMKP" json:"statusMKP"`     // สถานะในตลาด
	SalesStatus string    `db:"SalesStatus" json:"salesStatus"` // สถานะการขาย
	CreateDate  time.Time `db:"CreateDate" json:"-"`            // วันที่สร้างรายการ

	OrderLineDetail []OrderLineDetail `json:"OrderLineDetail"`
}

type OrderLineDetail struct {
	OrderNo     string    `db:"OrderNo" json:"-"`         // เลขที่ใบสั่งซื้อ
	SoNo        string    `db:"SoNo" json:"-"`            // เลขที่ใบสั่งขาย
	StatusMKP   string    `db:"StatusMKP" json:"-"`       // สถานะ Market Place
	SalesStatus string    `db:"SalesStatus" json:"-"`     // สถานะการขาย
	SKU         string    `db:"SKU" json:"sku"`           // รหัสสินค้า
	ItemName    string    `db:"ItemName" json:"itemName"` // ชื่อสินค้า
	QTY         int       `db:"QTY" json:"qty"`           // จำนวนสินค้า
	Price       float64   `db:"Price" json:"price"`       // ราคาต่อหน่วย
	CreateDate  time.Time `db:"CreateDate" json:"-"`      // วันที่สร้างรายการ
}

/********** Import Order to Warehouse: Sale Return ***************/

type ImportOrderResponse struct {
	OrderNo    string                    `json:"orderNo" db:"OrderNo"`
	SoNo       string                    `json:"soNo" db:"SoNo"`
	TrackingNo *string                   `json:"trackingNo" db:"TrackingNo"`
	CreateDate time.Time                 `json:"createDate" db:"CreateDate"`
	OrderLines []ImportOrderLineResponse `json:"orderLines"`
}

type ImportOrderLineResponse struct {
	OrderNo    string  `json:"orderNo" db:"OrderNo"`
	TrackingNo *string `json:"trackingNo" db:"TrackingNo"`
	SKU        string  `json:"sku" db:"SKU"`
	ItemName   string  `json:"itemName" db:"ItemName"`
	QTY        int     `json:"qty" db:"QTY"`
	Price      float64 `json:"price" db:"Price"`
}

type ImageResponse struct {
	ImageID  int    `json:"imageID"`
	FilePath string `json:"filePath"`
}

/********** Trade Return (Offline) ***************/

type ConfirmToReturnOrder struct {
	OrderNo        string    `json:"orderNo" db:"OrderNo"`
	StatusReturnID string    `db:"StatusReturnID"`
	StatusCheckID  string    `db:"StatusCheckID"`
	UpdateBy       string    `json:"updateBy" db:"UpdateBy"`
	UpdateDate     time.Time `json:"updateDate" db:"UpdateDate"`
}

type ConfirmTradeReturnOrder struct {
	OrderNo    string    `json:"orderNo" db:"OrderNo"`
	UpdateBy   string    `json:"updateBy" db:"UpdateBy"`
	UpdateDate time.Time `json:"updateDate" db:"UpdateDate"`
}

type ConfirmReceipt struct {
	Identifier     string    `json:"identifier"`
	StatusReturnID string    `db:"StatusReturnID"`
	StatusCheckID  string    `db:"StatusCheckID"`
	UpdateBy       string    `json:"updateBy" db:"UpdateBy"`
	UpdateDate     time.Time `json:"updateDate" db:"UpdateDate"`
}

type ConfirmReturnResponse struct {
	OrderNo     string    `json:"orderNo" db:"OrderNo"`
	ConfirmBy   string    `json:"confirmBy" db:"ConfirmBy"`
	ConfirmDate time.Time `json:"confirmDate" db:"ConfirmDate"`
}

type ConfirmReturnOrderDetails struct {
	OrderNo       string  `db:"OrderNo"`
	SoNo          string  `db:"SoNo"`
	SrNo          *string `db:"SrNo"`
	ChannelID     int     `db:"ChannelID"`
	Reason        string  `db:"Reason"`
	TrackingNo    *string `db:"TrackingNo"`
	CreateBy      string  `db:"CreateBy"`
	CreateDate    string  `db:"CreateDate"`
	UpdateBy      string  `db:"UpdateBy"`
	UpdateDate    string  `db:"UpdateDate"`
	StatusCheckID int     `db:"StatusCheckID"`
	DeleteBy      string  `db:"DeleteBy"`
	DeleteDate    string  `db:"DeleteDate"`
	// ActualQTY     int    `db:"ActualQTY"`
	// Price         float64 `db:"Price"`
	// StatusDelete  bool `db:"StatusDelete"`
}

// type TradeReturnLine struct {
// 	TradeReturnLine []OrderLines `json:"tradeReturnLine"`
// }

// type OrderLines struct {
// 	SKU       string  `json:"sku" db:"SKU"`
// 	ItemName  string  `json:"itemName" db:"ItemName"`
// 	QTY       int     `json:"qty" db:"QTY"`
// 	ReturnQTY int     `json:"returnQty" db:"ReturnQTY"`
// 	Price     float64 `json:"price" db:"Price"`
// }

type ImportOrderSummary struct {
	OrderNo string `json:"orderNo"`
	SKU     string `json:"sku"`
	Photo   string `json:"photo"`
}

// type BeforeReturnOrderResponse struct {
// 	OrderNo                string                          `json:"orderNo" db:"OrderNo"`
// 	SoNo                   string                          `json:"soNo" db:"SoNo"`
// 	SrNo                   string                          `json:"srNo" db:"SrNo"`
// 	ChannelID              int                             `json:"channelId" db:"ChannelID"`
// 	Reason                 string                          `json:"reason" db:"Reason"`
// 	CustomerID             string                          `json:"customerId" db:"CustomerID"`
// 	TrackingNo             string                          `json:"trackingNo" db:"TrackingNo"`
// 	Logistic               string                          `json:"logistic" db:"Logistic"`
// 	WarehouseID            int                             `json:"warehouseId" db:"WarehouseID"`
// 	SoStatus               *int                            `json:"soStatusId" db:"SoStatus"`
// 	MkpStatus              *int                            `json:"mkpStatusId" db:"MkpStatus"`
// 	ReturnDate             *time.Time                      `json:"returnDate" db:"ReturnDate"`
// 	StatusReturnID         *int                            `json:"statusReturnId" db:"StatusReturnID"`
// 	StatusConfID           *int                            `json:"statusConfId" db:"StatusConfID"`
// 	ConfirmBy              *string                         `json:"confirmBy" db:"ConfirmBy"`
// 	ConfirmDate            *time.Time                      `json:"confirmDate" db:"ConfirmDate"`
// 	CreateBy               string                          `json:"createBy" db:"CreateBy"`
// 	CreateDate             time.Time                       `json:"createDate" db:"CreateDate"`
// 	UpdateBy               *string                         `json:"updateBy" db:"UpdateBy"`
// 	UpdateDate             *time.Time                      `json:"updateDate" db:"UpdateDate"`
// 	CancelID               *int                            `json:"cancelId" db:"CancelID"`
// 	BeforeReturnOrderLines []BeforeReturnOrderLineResponse `json:"beforeReturnOrderLines"`
// }

type CreateBeforeReturnOrderResponse struct {
	OrderNo     string     `json:"orderNo" db:"OrderNo"`
	SoNo        string     `json:"soNo" db:"SoNo"`
	SrNo        *string    `json:"srNo" db:"SrNo"`
	ChannelID   int        `json:"channelId" db:"ChannelID"`
	Reason      string     `json:"reason" db:"Reason"`
	CustomerID  string     `json:"customerId" db:"CustomerID"`
	TrackingNo  *string    `json:"trackingNo" db:"TrackingNo"`
	Logistic    string     `json:"logistic" db:"Logistic"`
	WarehouseID int        `json:"warehouseId" db:"WarehouseID"`
	SoStatus    *int       `json:"soStatusId" db:"SoStatus"`
	MkpStatus   *int       `json:"mkpStatusId" db:"MkpStatus"`
	ReturnDate  *time.Time `json:"returnDate" db:"ReturnDate"`
	// CreateBy               string                          `json:"createBy" db:"CreateBy"`
	CreateDate time.Time `json:"createDate" db:"CreateDate"`
	// UpdateBy               *string                         `json:"updateBy" db:"UpdateBy"`
	// UpdateDate             *time.Time                      `json:"updateDate" db:"UpdateDate"`
	// CancelID               *int                            `json:"cancelId" db:"CancelID"`
	BeforeReturnOrderLines []BeforeReturnOrderLineResponse `json:"beforeReturnOrderLines"`
}

type BeforeReturnOrderLineResponse struct {
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
	SrNo           *string   `json:"srNo" db:"SrNo"`
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
	SrNo        *string   `json:"srNo" db:"SrNo"`
	CustomerID  string    `json:"customerId" db:"CustomerID"`
	TrackingNo  *string   `json:"trackingNo" db:"TrackingNo"`
	Logistic    string    `json:"logistic" db:"Logistic"`
	ChannelID   int       `json:"channelId" db:"ChannelID"`
	CreateDate  time.Time `json:"createDate" db:"CreateDate"`
	WarehouseID int       `json:"warehouseId" db:"WarehouseID"`
}

type DraftConfirmResponse struct {
	OrderNo string             `json:"orderNo" db:"OrderNo"`
	SoNo    string             `json:"soNo" db:"SoNo"`
	SrNo    *string            `json:"srNo" db:"SrNo"`
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

type RoleResponse struct {
	RoleID   int    `json:"roleId" db:"RoleID"`
	RoleName string `json:"roleName" db:"RoleName"`
}

type WarehouseResponse struct {
	WarehouseID   int    `json:"warehouseId" db:"WarehouseID"`
	WarehouseName string `json:"warehouseName" db:"WarehouseName"`
}
