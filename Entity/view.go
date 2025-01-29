package entity

//for dropdown

type District struct {
	Code         int    `db:"Code" json:"code"`
	ProvinceCode int    `db:"ProvinceCode" json:"provinceCode"`
	NameEN       string `db:"NameEN" json:"nameEN"`
	NameTH       string `db:"NameTH" json:"nameTH"`
}

type SubDistrict struct {
	Code         int    `db:"Code" json:"code"`
	DistrictCode int    `db:"DistrictCode" json:"districtCode"`
	ZipCode      string `db:"ZipCode" json:"zipCode"`
	NameTH       string `db:"NameTH" json:"nameTH"`
	NameEN       string `db:"NameEN" json:"nameEN"`
}

type Province struct {
	Code   int    `db:"Code" json:"code"`
	NameTH string `db:"NameTH" json:"nameTH"`
	NameEN string `db:"NameEN" json:"nameEN"`
}

type PostCode struct {

}
