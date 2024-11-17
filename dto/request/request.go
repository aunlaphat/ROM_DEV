package request

type LoginWeb struct {
	UserName string `json:"userName" db:"userID" example:"eknarin"`
	Password string `json:"password," db:"password" example:"asdfhdskjf"`
}
type LoginLark struct {
	UserName string `json:"userName" db:"userName" example:"eknarin"`
	UserID   string `json:"userID" db:"userID" example:"DC99999"`
}
