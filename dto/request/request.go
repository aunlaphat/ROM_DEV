// ส่งข้อมูลเข้ามาด้านหลัง
package request

// import "time"

type LoginWeb struct {
	UserName string `json:"userName" db:"userID" example:"eknarin"`
	Password string `json:"password," db:"password" example:"asdfhdskjf"`
}
type LoginLark struct {
	UserName string `json:"userName" db:"userName" example:"eknarin"`
	UserID   string `json:"userID" db:"userID" example:"DC99999"`
}

//ข้อมูลสำหรับคำสั่งซื้อ
//ข้อมูลลูกค้า แบรนด์/บริษัทในเครือที่สั่ง เชื่อมกับจน. รหัสสินค้าด้วย OrderNo ด้านล่าง
type GetOrderRequest struct {
	OrderNo         string  `json:"orderNo" db:"OrderNo" example:"AB0001"`
	BrandName       *string `json:"brandName" db:"BrandName" example:"BEWELL"`
	CustName        *string `json:"custName" db:"CustName" example:"Num"`
	CustAddress     *string `json:"custAddress" db:"CustAddress" example:"7/20 ซอย15/1"`
	CustDistrict    *string `json:"custDistrict" db:"CustDistrict" example:"บางกรวย"`
	CustSubDistrict *string `json:"custSubDistrict" db:"CustSubDistrict" example:"บางกรวย"`
	CustProvince    *string `json:"custProvince" db:"CustProvince" example:"นนทบุรี"`
	CustPostCode    *string `json:"custPostCode" db:"CustPostCode" example:"11130"`
	CustPhoneNum    *string `json:"custPhoneNum" db:"CustPhoneNum" example:"0921234567"`
	//CreateDate     	*time.Time  `json:"createDate" db:"CreateDate" example:"20/11/2567"`
	//UserCreated     *string 	`json:"userCreated" db:"UserCreated" example:"intern"`
	//UpdateDate      *time.Time  `json:"updateDate" db:"UpdateDate" example:"20/11/2568"`
	//UserUpdated     *string  	`json:"userUpdates" db:"UserUpdated" example:"intern"`

	OrderLines []OrderLineRequest `json:"orderLines"`
}

//ข้อมูลลูกค้า แบรนด์/บริษัทในเครือที่สั่ง เชื่อมกับจน. รหัสสินค้าด้วย OrderNo ด้านล่าง
type CreateOrderRequest struct {
	OrderNo         string  `json:"orderNo" db:"OrderNo" example:"AB0001"`
	BrandName       *string `json:"brandName" db:"BrandName" example:"BEWELL"`
	CustName        *string `json:"custName" db:"CustName" example:"Num"`
	CustAddress     *string `json:"custAddress" db:"CustAddress" example:"7/20 ซอย15/1"`
	CustDistrict    *string `json:"custDistrict" db:"CustDistrict" example:"บางกรวย"`
	CustSubDistrict *string `json:"custSubDistrict" db:"CustSubDistrict" example:"บางกรวย"`
	CustProvince    *string `json:"custProvince" db:"CustProvince" example:"นนทบุรี"`
	CustPostCode    *string `json:"custPostCode" db:"CustPostCode" example:"11130"`
	CustPhoneNum    *string `json:"custPhoneNum" db:"CustPhoneNum" example:"0921234567"`
	// CreateDate     	*time.Time  `json:"createDate" db:"CreateDate"`
	// UserCreated     *string 	`json:"userCreated" db:"UserCreated" example:"intern"`
	// UpdateDate      *time.Time  `json:"updateDate" db:"UpdateDate" example:"20/11/2568"`
	// UserUpdated     *string  	`json:"userUpdates" db:"UserUpdated" example:"intern"`

	OrderLines []OrderLineRequest `json:"orderLines"`
}

type UpdateOrderRequest struct {
	OrderNo         string  `json:"-" db:"OrderNo"`
	// BrandName       *string `json:"brandName" db:"BrandName" example:"BEWELL"`
	CustName        *string `json:"custName" db:"CustName" example:"Num"`
	CustAddress     *string `json:"custAddress" db:"CustAddress" example:"7/20 ซอย15/1"`
	CustDistrict    *string `json:"custDistrict" db:"CustDistrict" example:"บางกรวย"`
	CustSubDistrict *string `json:"custSubDistrict" db:"CustSubDistrict" example:"บางกรวย"`
	CustProvince    *string `json:"custProvince" db:"CustProvince" example:"นนทบุรี"`
	CustPostCode    *string `json:"custPostCode" db:"CustPostCode" example:"11130"`
	CustPhoneNum    *string `json:"custPhoneNum" db:"CustPhoneNum" example:"0921234567"`
	// CreateDate     	*time.Time  `json:"createDate" db:"CreateDate" example:"20/11/2567"`
	// UserCreated     *string 	`json:"userCreated" db:"UserCreated" example:"intern"`
	// UpdateDate      *time.Time  `json:"updateDate" db:"UpdateDate" example:"20/11/2568"`
	UserUpdated *string `json:"userUpdates" db:"UserUpdated" example:"intern"`
}

//ข้อมูลของสินค้าแต่ละชิ้นที่ลูกค้าสั่ง เชื่อมกันกับด้านบนด้วย OrderNo
type OrderLineRequest struct {
	OrderNo  string   `json:"orderNo" db:"OrderNo" example:"AB0001"`
	SKU      *string  `json:"sku" db:"SKU" Example:"SKU12345"`
	ItemName *string  `json:"itemName" db:"ItemName" example:"เก้าอี้"`
	QTY      *int     `json:"qty" db:"QTY" Example:"10"`
	Price    *float64 `json:"price" db:"Price" example:"199.05"`
}
