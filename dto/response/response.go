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
