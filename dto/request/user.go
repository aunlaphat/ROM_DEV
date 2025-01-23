package request

type Login struct {
	UserID   string `json:"userID" db:"UserID" example:"DC65060"`
	Password string `json:"password" db:"Password" example:"xxxxxxxx"`
}

type LoginWeb struct {
	UserName string `json:"userName" db:"Username" example:"eknarin.ler"`
	Password string `json:"password" db:"Password" example:"EKna1234"` // change password lastest in 17 January 2025
}

type LoginLark struct {
	UserID   string `json:"userID" db:"userID" example:"DC65060"`
	UserName string `json:"userName" db:"userName" example:"eknarin.ler"`
}

type LoginJWT struct {
	UserID   string `json:"userID" db:"UserID" example:"DC53002"`
	UserName string `json:"userName" db:"Username" example:"string"`
}
