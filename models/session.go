package models

//Session model user
type Session struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Username string `json:"username"`
}

//TableName name of table
func (session Session) TableName() string {
	return "user"
}
