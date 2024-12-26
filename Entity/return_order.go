package entity

// import "time"

// type ReturnOrder struct {
// 	ReturnID        string        `json:"returnId" db:"ReturnID"`                    
// 	OrderNo         string     	  `json:"orderNo" db:"OrderNo"`                         
// 	SaleOrder       string    	  `json:"saleOrder" db:"SaleOrder"`           
// 	SaleReturn      string    	  `json:"saleReturn" db:"SaleReturn"`          
// 	TrackingNo      string    	  `json:"trackingNo" db:"TrackingNo"`         
// 	PlatfID         *int       	  `json:"platfId" db:"PlatfID"`              
// 	ChannelID       *int       	  `json:"channelId" db:"ChannelID"`          
// 	OptStatusID     *int       	  `json:"optStatusId" db:"OptStatusID"`       
// 	AxStatusID      *int       	  `json:"axStatusId" db:"AxStatusID"`         
// 	PlatfStatusID   *int       	  `json:"platfStatusId" db:"PlatfStatusID"`    
// 	Remark          *string    	  `json:"remark" db:"Remark"`                  
// 	CreateBy        string     	  `json:"createBy" db:"CreateBy"`                      
// 	CreateDate      time.Time  	  `json:"createDate" db:"CreateDate"`                
// 	UpdateBy        *string    	  `json:"updateBy" db:"UpdateBy"`              
// 	UpdateDate      *time.Time 	  `json:"updateDate" db:"UpdateDate"`          
// 	CancelID        *int          `json:"cancelId" db:"CancelID"`             
// 	StatusCheckID   *int       	  `json:"statusCheckId" db:"StatusCheckID"`    
// 	CheckBy         *string    	  `json:"checkBy" db:"CheckBy"`                
// 	Description     *string    	  `json:"description" db:"Description"`        
// }

// type ReturnOrderLine struct {
// 	ReturnID        string        `json:"returnId" db:"ReturnID"`                    
// 	OrderNo         string     	  `json:"orderNo" db:"OrderNo"`                                  
// 	TrackingNo      string    	  `json:"trackingNo" db:"TrackingNo"`   
//     SKU          	string    	  `json:"sku" db:"SKU"`         
//     ReturnQTY    	int       	  `json:"returnQTY" db:"ReturnQTY"` 
//     CheckQTY   	    *int       	  `json:"checkQTY" db:"CheckQTY"`
//     Price        	float64    	  `json:"price" db:"Price"`     
// 	CreateBy        string     	  `json:"createBy" db:"CreateBy"`                      
// 	CreateDate      time.Time  	  `json:"createDate" db:"CreateDate"`   
//     AlterSKU     	*string    	  `json:"alterSKU" db:"AlterSKU"`   
// 	UpdateBy        *string    	  `json:"updateBy" db:"UpdateBy"`              
// 	UpdateDate      *time.Time 	  `json:"updateDate" db:"UpdateDate"`      
// }

// type Platforms struct {
// 	PlatfID 		int    	  	  `json:"platfId" db:"PlatfID"` 
// 	PlatfName		string    	  `json:"platfName" db:"PlatfName"` 
// }

// type Channel struct {
// 	ChannelID 		int    	  	  `json:"channelId" db:"ChannelID"` 
// 	ChannelName		string    	  `json:"channelName" db:"ChannelName"` 
// }

// type CancelReturnOrder struct {
// 	CancelID 		*int    	  `json:"platfId" db:"PlatfID"` 
// 	ReturnID		*string    	  `json:"returnId" db:"ReturnID"` 
// 	Remark			*string    	  `json:"remark" db:"Remark"` 
// 	CancelDate		*time.Time    `json:"cancelDate" db:"CancelDate"` 
// 	CancelBy		*string    	  `json:"cancelBy" db:"CancelBy"` 
// }

// type StatusCheck struct {
// 	StatusCheckID	    string    `json:"statusCheckId" db:"StatusCheckID"` 
// 	StatusCheckName		string    `json:"statusCheckName" db:"StatusCheckName"` 
// }

// type BeforeReturnOrder struct {
//     RecID          int        `json:"recID" db:"RecID"`
//     OrderNo        string     `json:"orderNo" db:"OrderNo"`
//     SaleOrder      string     `json:"saleOrder" db:"SaleOrder"`
//     SaleReturn     string     `json:"saleReturn" db:"SaleReturn"`
//     ChannelID      int        `json:"channelID" db:"ChannelID"`
//     ReturnType     *string     `json:"returnType" db:"ReturnType"`
//     CustomerID     int        `json:"customerID" db:"CustomerID"`
//     TrackingNo     string     `json:"trackingNo" db:"TrackingNo"`
//     Logistic       string     `json:"logistic" db:"Logistic"`
//     WarehouseID    int        `json:"warehouseID" db:"WarehouseID"`
//     SoStatusID     *int        `json:"soStatusID" db:"SoStatusID"`
//     MkpStatusID    *int        `json:"mkpStatusID" db:"mkpStatusID"`
//     ReturnDate     *time.Time  `json:"returnDate" db:"ReturnDate"`
//     StatusReturnID int        `json:"statusReturnID" db:"StatusReturnID"`
//     StatusConfID   int        `json:"statusConfID" db:"StatusConfID"`
//     ConfirmBy      *string     `json:"confirmBy" db:"ConfirmBy"`
//     CreateBy       string     `json:"createBy" db:"CreateBy"`
//     CreateDate     time.Time  `json:"createDate" db:"CreateDate"`
//     UpdateBy       *string    `json:"updateBy" db:"UpdateBy"`
//     UpdateDate     *time.Time `json:"updateDate" db:"UpdateDate"`
//     CancelID       *int       `json:"cancelID" db:"CancelID"`
// }

// type BeforeReturnOrderLine struct {
//     RecID       int        `json:"recID" db:"RecID"`
//     OrderNo     string     `json:"orderNo" db:"OrderNo"`
//     SKU         string     `json:"sku" db:"SKU"`
//     QTY         *int        `json:"qty" db:"QTY"`
//     ReturnQTY   *int        `json:"returnQTY" db:"ReturnQTY"`
//     Price       *float64    `json:"price" db:"Price"`
//     CreateBy    string     `json:"createBy" db:"CreateBy"`
//     CreateDate  time.Time  `json:"createDate" db:"CreateDate"`
//     AlterSKU    *string    `json:"alterSKU" db:"AlterSKU"`
//     UpdateBy    *string    `json:"updateBy" db:"UpdateBy"`
//     UpdateDate  *time.Time `json:"updateDate" db:"UpdateDate"`
//     TrackingNo  string     `json:"trackingNo" db:"TrackingNo"`
// }