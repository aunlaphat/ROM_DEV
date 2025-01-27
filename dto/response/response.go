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

type UserInform struct {
	UserID       string `json:"userID,omitempty" db:"UserID" example:"DC64205"`
	UserName     string `json:"userName,omitempty" db:"Username" example:"aunlaphat.art"`
	NickName     string `json:"nickName,omitempty" db:"NickName" example:"fa"`
	FullNameTH   string `json:"fullNameTH,omitempty" db:"FullNameTH" example:"อัญญ์ลภัส อาจสุริยงค์"`
	DepartmentNo string `json:"department,omitempty" db:"DepartmentNo" example:"G01"`
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
	Remark        *string    `json:"remark" db:"Remark"`
	CreateBy      string     `json:"createBy" db:"CreateBy"`
	CreateDate    time.Time  `json:"createDate" db:"CreateDate"`
	UpdateBy      *string    `json:"updateBy" db:"UpdateBy"`
	UpdateDate    *time.Time `json:"updateDate" db:"UpdateDate"`
	CancelID      *int       `json:"cancelId" db:"CancelID"`
	StatusCheckID *int       `json:"statusCheckId" db:"StatusCheckID"`
	CheckBy       *string    `json:"checkBy" db:"CheckBy"`
	Description   *string    `json:"description" db:"Description"`

	ReturnOrderLine []ReturnOrderLine `json:"ReturnOrderLine"`
}

type ReturnOrderLine struct {
	OrderNo    string     `json:"orderNo" db:"OrderNo"`
	TrackingNo string     `json:"trackingNo" db:"TrackingNo"`
	SKU        string     `json:"sku" db:"SKU"`
	ReturnQTY  int        `json:"returnQTY" db:"ReturnQTY"`
	QTY        *int       `json:"qty" db:"QTY"`
	Price      float64    `json:"price" db:"Price"`
	CreateBy   string     `json:"createBy" db:"CreateBy"`
	CreateDate time.Time  `json:"createDate" db:"CreateDate"`
	AlterSKU   *string    `json:"alterSKU" db:"AlterSKU"`
	UpdateBy   *string    `json:"updateBy" db:"UpdateBy"`
	UpdateDate *time.Time `json:"updateDate" db:"UpdateDate"`
}

type CancelStatus struct {
	CancelID     *int       `json:"platfId" db:"PlatfID"`
	RefID        *string    `json:"refId" db:"RefID"` //fk of table beforeod and returnod with pk-RecID
	CancelStatus bool       `db:"CancelStatus"`
	Remark       *string    `json:"remark" db:"Remark"`
	CancelDate   *time.Time `json:"cancelDate" db:"CancelDate"`
	CancelBy     *string    `json:"cancelBy" db:"CancelBy"`
}

/********** OrderHead + Line data Project ***************/

type OrderDetail struct {
	// json => OrderHeadDetail[ OrderLineDetail[ {},{},..] ]
	OrderHeadDetail []OrderHeadDetail `json:"OrderHeadDetail"`
}

type OrderHeadDetail struct {
	OrderNo     string    `db:"OrderNo" json:"orderNo"`         // เลขที่ใบสั่งซื้อ
	SoNo        *string   `db:"SoNo" json:"soNo"`               // เลขที่ใบสั่งขาย
	StatusMKP   string    `db:"StatusMKP" json:"statusMKP"`     // สถานะในตลาด
	SalesStatus string    `db:"SalesStatus" json:"salesStatus"` // สถานะการขาย
	CreateDate  time.Time `db:"CreateDate" json:"-"`            // วันที่สร้างรายการ

	OrderLineDetail []OrderLineDetail `json:"OrderLineDetail"`
}

type OrderLineDetail struct {
	OrderNo     string    `db:"OrderNo" json:"-"`         // เลขที่ใบสั่งซื้อ
	SoNo        *string   `db:"SoNo" json:"-"`            // เลขที่ใบสั่งขาย
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
	OrderNo     string                  `json:"orderNo" db:"OrderNo"`
	SoNo        string                  `json:"soNo" db:"SoNo"`
	TrackingNo  string     `json:"trackingNo" db:"TrackingNo"`
	CreateDate  *time.Time              `json:"createDate" db:"CreateDate"`
	OrderLines  []ImportOrderLineResponse `json:"orderLines"`
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
