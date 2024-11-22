//ดึงค่าออกมาโดยตรงจาก db
package entity

import "time"

type Order struct {
	OrderNo          string  	`json:"orderNo" db:"OrderNo" example:"AB0001"`
	BrandName       *string  	`json:"brandName" db:"BrandName" example:"BEWELL"`
	CustName        *string  	`json:"custName" db:"CustName" example:"Fa"`
	CustAddress     *string  	`json:"custAddress" db:"CustAddress" example:"7/20 ซอย15/1"`
	CustDistrict    *string  	`json:"custDistrict" db:"CustDistrict" example:"บางกรวย"`
	CustSubDistrict *string  	`json:"custSubDistrict" db:"CustSubDistrict" example:"บางกรวย"`
	CustProvince    *string  	`json:"custProvince" db:"CustProvince" example:"นนทบุรี"`
	CustPostCode    *string  	`json:"custPostCode" db:"CustPostCode" example:"11130"`
	CustPhoneNum    *string  	`json:"custPhoneNum" db:"CustPhoneNum" example:"0921234567"`
	CreateDate     	*time.Time  `json:"createDate" db:"CreateDate" example:"20/11/2567"`
	UserCreated     *string 	`json:"userCreated" db:"UserCreated" example:"intern"`
	UpdateDate      *time.Time  `json:"updateDate" db:"UpdateDate" example:"20/11/2568"`
	UserUpdated     *string  	`json:"userUpdates" db:"UserUpdated" example:"intern"`

	OrderLines []OrderLine `gorm:"foreignKey:OrderNo" json:"orderLine"`
}

type OrderLine struct {
	OrderNo  *string  `json:"orderNo" db:"OrderNo" example:"AB0001"`
	SKU      *string  `json:"sku" db:"SKU" Example:"SKU12345"`
	ItemName *string  `json:"itemName" db:"ItemName" example:"เก้าอี้"`
	QTY      *int     `json:"qty" db:"QTY" Example:"10"`
	Price    *float64 `json:"price" db:"Price" example:"199.05"`
}
