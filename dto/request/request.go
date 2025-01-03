// ส่งข้อมูลเข้ามาด้านหลัง
package request

import "time"

/********** Login ***************/

type LoginWeb struct {
	UserName string `json:"userName" db:"userID" example:"eknarin"`
	Password string `json:"password," db:"password" example:"asdfhdskjf"`
}

type LoginLark struct {
	UserName string `json:"userName" db:"userName" example:"eknarin"`
	UserID   string `json:"userID" db:"userID" example:"DC99999"`
}

type Login struct {
	UserID string `json:"userID" db:"UserID" example:"DC53002"`
	Password string `json:"password," db:"Password" example:"string"`
}

type LoginJWT struct {
	UserID string `json:"userID" db:"UserID" example:"DC53002"`
	UserName string `json:"userName" db:"Username" example:"string"`
}

/********** Order ***************/

//ข้อมูลสำหรับคำสั่งซื้อ
//ข้อมูลลูกค้า แบรนด์/บริษัทในเครือที่สั่ง เชื่อมกับจน. รหัสสินค้าด้วย OrderNo ด้านล่าง
type CreateOrderRequest struct {
	OrderNo         string  `json:"orderNo" db:"OrderNo" example:"OD0001"`
	BrandName       *string `json:"brandName" db:"BrandName" example:"BEWELL"`
	CustName        *string `json:"custName" db:"CustName" example:"Fa"`
	CustAddress     *string `json:"custAddress" db:"CustAddress" example:"7/20"`
	CustDistrict    *string `json:"custDistrict" db:"CustDistrict" example:"Bang-Kruay"`
	CustSubDistrict *string `json:"custSubDistrict" db:"CustSubDistrict" example:"Bang-Kruay"`
	CustProvince    *string `json:"custProvince" db:"CustProvince" example:"Nonthaburi"`
	CustPostCode    *string `json:"custPostCode" db:"CustPostCode" example:"11130"`
	CustPhoneNum    *string `json:"custPhoneNum" db:"CustPhoneNum" example:"0912345678"`

	OrderLines []OrderLineRequest `json:"orderLines"`
}

type UpdateOrderRequest struct {
	OrderNo         string  `json:"-" db:"OrderNo"`
	CustName        *string `json:"custName" db:"CustName" example:"Fa"`
	CustAddress     *string `json:"custAddress" db:"CustAddress" example:"1/12"`
	CustDistrict    *string `json:"custDistrict" db:"CustDistrict" example:"Bang-Plad"`
	CustSubDistrict *string `json:"custSubDistrict" db:"CustSubDistrict" example:"Bang-Plad"`
	CustProvince    *string `json:"custProvince" db:"CustProvince" example:"Bangkok"`
	CustPostCode    *string `json:"custPostCode" db:"CustPostCode" example:"10600"`
	CustPhoneNum    *string `json:"custPhoneNum" db:"CustPhoneNum" example:"0921234567"`
}

//ข้อมูลของสินค้าแต่ละชิ้นที่ลูกค้าสั่ง เชื่อมกันกับด้านบนด้วย OrderNo
type OrderLineRequest struct {
	OrderNo  string   `json:"-" db:"OrderNo" `
	SKU      *string  `json:"sku" db:"SKU" example:"SKU12345"`
	ItemName *string  `json:"itemName" db:"ItemName" example:"เก้าอี้"`
	QTY      *int     `json:"qty" db:"QTY" example:"5"`
	Price    *float64 `json:"price" db:"Price" example:"5900.00"`
}

/********** Return Order ***************/

type ReturnOrder struct {
	ReturnID        string        `json:"returnId" db:"ReturnID"`                    
	OrderNo         string     	  `json:"orderNo" db:"OrderNo"`                         
	SaleOrder       string    	  `json:"saleOrder" db:"SaleOrder"`           
	SaleReturn      string    	  `json:"saleReturn" db:"SaleReturn"`          
	TrackingNo      string    	  `json:"trackingNo" db:"TrackingNo"`         
	PlatfID         *int       	  `json:"platfId" db:"PlatfID"`              
	ChannelID       *int       	  `json:"channelId" db:"ChannelID"`          
	OptStatusID     *int       	  `json:"optStatusId" db:"OptStatusID"`       
	AxStatusID      *int       	  `json:"axStatusId" db:"AxStatusID"`         
	PlatfStatusID   *int       	  `json:"platfStatusId" db:"PlatfStatusID"`    
	Remark          *string    	  `json:"remark" db:"Remark"`                  
	CreateBy        string     	  `json:"createBy" db:"CreateBy"`                      
	CreateDate      time.Time  	  `json:"createDate" db:"CreateDate"`                
	UpdateBy        *string    	  `json:"updateBy" db:"UpdateBy"`              
	UpdateDate      *time.Time 	  `json:"updateDate" db:"UpdateDate"`          
	CancelID        *int          `json:"cancelId" db:"CancelID"`             
	StatusCheckID   *int       	  `json:"statusCheckId" db:"StatusCheckID"`    
	CheckBy         *string    	  `json:"checkBy" db:"CheckBy"`                
	Description     *string    	  `json:"description" db:"Description"`        
}

type Platforms struct {
	PlatfID 		int    	  	  `json:"platfId" db:"PlatfID"` 
	PlatfName		string    	  `json:"platfName" db:"PlatfName"` 
}

type Channel struct {
	ChannelID 		int    	  	  `json:"channelId" db:"ChannelID"` 
	ChannelName		string    	  `json:"channelName" db:"ChannelName"` 
}

type CancelStatus struct {
	CancelID 		*int    	  `json:"platfId" db:"PlatfID"` 
	RefID		    *string    	  `json:"refId" db:"RefID"`
	CancelStatus	 bool         `db:"CancelStatus"` 
	Remark			*string    	  `json:"remark" db:"Remark"` 
	CancelDate		*time.Time    `json:"cancelDate" db:"CancelDate"` 
	CancelBy		*string    	  `json:"cancelBy" db:"CancelBy"` 
}

type StatusCheck struct {
	StatusCheckID	    string    `json:"statusCheckId" db:"StatusCheckID"` 
	StatusCheckName		string    `json:"statusCheckName" db:"StatusCheckName"` 
}

type ReturnOrderHead struct {
	ReturnOrder
	Platform           *Platforms         `json:"platform"`
	Channel            *Channel           `json:"channel"`
	CancelStatus       *CancelStatus      `json:"cancelStatus"`
	StatusCheck        *StatusCheck       `json:"statusCheck"`
}

type CreateReturnOrder struct {
	ReturnID        string          `json:"returnId" db:"ReturnID" example:"RID0001"`                    
	OrderNo         string          `json:"orderNo" db:"OrderNo" example:"ORD0001"`                         
	SaleOrder       string          `json:"saleOrder" db:"SaleOrder" example:"SO0001"`           
	SaleReturn      string          `json:"saleReturn" db:"SaleReturn" example:"SR0001"`          
	TrackingNo      string          `json:"trackingNo" db:"TrackingNo" example:"12345678TH"`         
	PlatfID         *int            `json:"platfId" db:"PlatfID" example:"1"`              
	ChannelID       *int            `json:"channelId" db:"ChannelID" example:"2"`          
	OptStatusID     *int            `json:"optStatusId" db:"OptStatusID" example:"1"`       
	AxStatusID      *int            `json:"axStatusId" db:"AxStatusID" example:"1"`         
	PlatfStatusID   *int            `json:"platfStatusId" db:"PlatfStatusID" example:"1"`    
	Remark          *string         `json:"remark" db:"Remark" example:""`
	CancelID        *int            `json:"cancelId" db:"CancelID" example:"1"`             
	StatusCheckID   *int            `json:"statusCheckId" db:"StatusCheckID" example:"1"`    
	CheckBy         *string         `json:"checkBy" db:"CheckBy" example:"dev03"`                
	Description     *string         `json:"description" db:"Description" example:""`

	ReturnOrderLine	[]ReturnOrderLine `json:"ReturnOrderLine"`
}

type UpdateReturnOrder struct {
	ReturnID        string          `json:"-" db:"ReturnID"`                    
	OrderNo         string          `json:"-" db:"OrderNo"`                         
	SaleOrder       string          `json:"-" db:"SaleOrder"`           
	SaleReturn      *string         `json:"saleReturn" db:"SaleReturn" example:"SR0001"`          
	TrackingNo      *string         `json:"trackingNo" db:"TrackingNo" example:"12345678TH"`         
	PlatfID         *int            `json:"platfId" db:"PlatfID" example:"1"`              
	ChannelID       *int            `json:"channelId" db:"ChannelID" example:"2"`          
	OptStatusID     *int            `json:"optStatusId" db:"OptStatusID" example:"1"`       
	AxStatusID      *int            `json:"axStatusId" db:"AxStatusID" example:"1"`         
	PlatfStatusID   *int            `json:"platfStatusId" db:"PlatfStatusID" example:"1"`    
	Remark          *string         `json:"remark" db:"Remark" example:""`
	CancelID        *int            `json:"cancelId" db:"CancelID" example:"1"`             
	StatusCheckID   *int            `json:"statusCheckId" db:"StatusCheckID" example:"1"`    
	CheckBy         *string         `json:"checkBy" db:"CheckBy" example:"dev03"`                
	Description     *string         `json:"description" db:"Description" example:""`
}

type ReturnOrderLine struct {       
	ReturnID        string          `json:"-" db:"ReturnID"`             
	OrderNo         string          `json:"-" db:"OrderNo"`                                  
	TrackingNo      string          `json:"-" db:"TrackingNo"`   
    SKU             string          `json:"sku" db:"SKU" example:"SKU12345"`         
    ReturnQTY       int             `json:"returnQTY" db:"ReturnQTY" example:"5"` 
    CheckQTY        *int            `json:"checkQTY" db:"CheckQTY" example:"5"`
    Price           float64         `json:"price" db:"Price" example:"199.99"`      
	AlterSKU        *string         `json:"-" db:"AlterSKU" `  
}


