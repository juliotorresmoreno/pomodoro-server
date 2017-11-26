package models

//Session model user
type Session struct {
	ID       int64  `json:"-"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

//TableName name of table
func (session Session) TableName() string {
	return "user"
}
