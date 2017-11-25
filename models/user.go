package models

import "time"

//User model user
type User struct {
	ID       int64     `xorm:"id bigint not null autoincr pk" json:"id"`
	Name     string    `xorm:"varchar(100) not null" valid:"required,alphaSpaces"`
	LastName string    `xorm:"varchar(100) not null" valid:"required,alphanum"`
	Username string    `xorm:"varchar(100) not null" valid:"required,unique"`
	Password string    `xorm:"password" valid:"required,password"`
	Token    string    `xorm:"varchar(100)"`
	CreateAt time.Time `xorm:"created"`
	UpdateAt time.Time `xorm:"updated"`
}

//TableName name of table
func (user User) TableName() string {
	return "users"
}
