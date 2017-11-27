package models

//Session model user
type Session struct {
	ID       int64  `xorm:"id bigint not null autoincr pk" json:"-"`
	Name     string `xorm:"varchar(100) not null" json:"name"`
	LastName string `xorm:"varchar(100) not null" json:"lastname"`
	Username string `xorm:"varchar(100) not null" json:"username"`
	Token    string `xorm:"varchar(100)" json:"token"`
}

//TableName name of table
func (session Session) TableName() string {
	return "users"
}
