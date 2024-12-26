// ส่งกลับไปยังด้านหน้า
package response

import "time"

type Login struct {
	UserName   string `json:"userName,omitempty" db:"UserName" example:"userName"`
	UserID     string `json:"userID,omitempty" db:"UserID" example:"userID"`
	RoleID     int    `json:"roleID,omitempty" db:"RoleID" example:"1"`
	FullNameTH string `json:"fullNameTH,omitempty" db:"FullNameTH" example:"test1234"`
	NickName   string `json:"nickName,omitempty" db:"NickName" example:"test1234"`
	Platfrom   string `json:"platfrom" db:"Platfrom" example:"test1234"`
}

type UserInform struct {
	UserID     string `json:"userID,omitempty" db:"UserID" example:"DC64205"`
	UserName   string `json:"userName,omitempty" db:"Username" example:"aunlaphat.art"`
	NickName   string `json:"nickName,omitempty" db:"NickName" example:"fa"`
	FullNameTH string `json:"fullNameTH,omitempty" db:"FullNameTH" example:"อัญญ์ลภัส อาจสุริยงค์"`
	DepartmentNo string `json:"department,omitempty" db:"DepartmentNo" example:"G01"`
}

type OrderResponse struct {
	OrderNo         string     `json:"orderNo" db:"OrderNo" example:"OD0001"`
	BrandName       *string    `json:"brandName" db:"BrandName" example:"BEWELL"`
	CustName        *string    `json:"custName" db:"CustName" example:"Fa"`
	CustAddress     *string    `json:"custAddress" db:"CustAddress" example:"7/20"`
	CustDistrict    *string    `json:"custDistrict" db:"CustDistrict" example:"Bang-Kruay"`
	CustSubDistrict *string    `json:"custSubDistrict" db:"CustSubDistrict" example:"Bang-Kruay"`
	CustProvince    *string    `json:"custProvince" db:"CustProvince" example:"Nonthaburi"`
	CustPostCode    *string    `json:"custPostCode" db:"CustPostCode" example:"11130"`
	CustPhoneNum    *string    `json:"custPhoneNum" db:"CustPhoneNum" example:"0912345678"`
	CreateDate      *time.Time `json:"createDate" db:"CreateDate" example:"2024-11-22 09:45:33.260"`
	UserCreated     *string    `json:"userCreated" db:"UserCreated" example:"intern"`
	UpdateDate      *time.Time `json:"updateDate" db:"UpdateDate" example:"2024-11-30 09:45:33.260"`
	UserUpdated     *string    `json:"userUpdates" db:"UserUpdated" example:"intern"`

	OrderLines []OrderLineResponse `json:"orderLines"`
}

type OrderLineResponse struct {
	OrderNo  *string  `json:"orderNo" db:"OrderNo" example:"OD0001"`
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
	SaleReturn      *string    	  `json:"saleReturn" db:"SaleReturn"`          
	TrackingNo      *string    	  `json:"trackingNo" db:"TrackingNo"`         
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

	ReturnOrderLine		[]ReturnOrderLine		  `json:"ReturnOrderLine"`
}

type ReturnOrderLine struct {
	RecID      		int           `json:"recId" db:"RecID"`
	ReturnID        string        `json:"returnId" db:"ReturnID"`                    
	OrderNo         string     	  `json:"orderNo" db:"OrderNo"`                                  
	TrackingNo      string    	  `json:"trackingNo" db:"TrackingNo"`   
    SKU          	string    	  `json:"sku" db:"SKU"`         
    ReturnQTY    	int       	  `json:"returnQTY" db:"ReturnQTY"` 
    CheckQTY   	    *int       	  `json:"checkQTY" db:"CheckQTY"`
    Price        	float64    	  `json:"price" db:"Price"`     
	CreateBy        string     	  `json:"createBy" db:"CreateBy"`                      
	CreateDate      time.Time  	  `json:"createDate" db:"CreateDate"`   
    AlterSKU     	*string    	  `json:"alterSKU" db:"AlterSKU"`   
	UpdateBy        *string    	  `json:"updateBy" db:"UpdateBy"`              
	UpdateDate      *time.Time 	  `json:"updateDate" db:"UpdateDate"`      
}

type CancelReturnOrder struct {
	CancelID 		*int    	  `json:"platfId" db:"PlatfID"` 
	ReturnID		*string    	  `json:"returnId" db:"ReturnID"` 
	Remark			*string    	  `json:"remark" db:"Remark"` 
	CancelDate		*time.Time    `json:"cancelDate" db:"CancelDate"` 
	CancelBy		*string    	  `json:"cancelBy" db:"CancelBy"` 
}

type ReturnOrderHead struct {
	ReturnOrder		   ReturnOrder		  `json:"ReturnOrder"`
}