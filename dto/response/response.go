// ส่งกลับไปยังด้านหน้า
package response

import "time"

/********** Login ***************/

/* type Login struct {
	UserID       string `json:"userID,omitempty" db:"UserID" example:"DC-XXXXX"`
	UserName     string `json:"userName,omitempty" db:"UserName" example:"userName"`
	RoleID       int    `json:"roleID,omitempty" db:"RoleID" example:"0"`
	FullNameTH   string `json:"fullNameTH,omitempty" db:"FullNameTH" example:"Firstname Lastname"`
	NickName     string `json:"nickName,omitempty" db:"NickName" example:"Nickname"`
	DepartmentNo string `json:"department,omitempty" db:"DepartmentNo" example:"G07"`
	Platform     string `json:"platform" db:"Platform" example:"Platform"`
} */

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
	TrackingNo string     `json:"trackingNo" db:"TrackingNo"`
	SKU        string     `json:"sku" db:"SKU"`
	ReturnQTY  int        `json:"returnQTY" db:"ReturnQTY"`
	ActualQTY  int        `json:"actualQTY" db:"ActualQTY"`
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
	SrNo          string    `json:"srNo" db:"SrNo" example:"SR0001"`
	TrackingNo    string    `json:"trackingNo" db:"TrackingNo" example:"12345678TH"`
	PlatfID       int       `json:"platfId" db:"PlatfID" example:"1"`
	ChannelID     int       `json:"channelId" db:"ChannelID" example:"2"`
	OptStatusID   int       `json:"optStatusId" db:"OptStatusID" example:"1"`
	AxStatusID    int       `json:"axStatusId" db:"AxStatusID" example:"1"`
	PlatfStatusID int       `json:"platfStatusId" db:"PlatfStatusID" example:"1"`
	Reason        string    `json:"reason" db:"Reason"`
	StatusCheckID int       `json:"statusCheckId" db:"StatusCheckID" example:"1"`
	Description   string    `json:"description" db:"Description" example:""`
	CreateBy      string    `json:"createBy" db:"CreateBy"`
	CreateDate    time.Time `json:"createDate" db:"CreateDate"` // MSSQL SYSDATETIME() function

	ReturnOrderLine []ReturnOrderLine `json:"ReturnOrderLine"`
}

type UpdateReturnOrder struct {
	OrderNo       string    `json:"-" db:"OrderNo"`
	SoNo          string    `json:"-" db:"SoNo"`
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
	UpdateDate    *time.Time `json:"updateDate" db:"UpdateDate"` // MSSQL SYSDATETIME() function
}

type DeleteReturnOrder struct {
	OrderNo string `db:"OrderNo"`
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
	TrackingNo string                    `json:"trackingNo" db:"TrackingNo"`
	CreateDate time.Time                 `json:"createDate" db:"CreateDate"`
	OrderLines []ImportOrderLineResponse `json:"orderLines"`
}

type ImportOrderLineResponse struct {
	SKU      string  `json:"sku" db:"SKU"`
	ItemName string  `json:"itemName" db:"ItemName"`
	QTY      int     `json:"qty" db:"QTY"`
	Price    float64 `json:"price" db:"Price"`
}

type ImageResponse struct {
	ImageID  int    `json:"imageID"`
	FilePath string `json:"filePath"`
}
