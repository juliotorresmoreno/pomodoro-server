package models

import (
	"time"
)

type Task struct {
	ID        int64     `xorm:"id bigint not null autoincr pk" json:"id"`
	Name      string    `xorm:"varchar(100) not null" valid:"required,alphaSpaces"`
	Step      int8      `xorm:"int not null" valid:"required"`
	UserID    int64     `xorm:"int not null" valid:"required"`
	StartDate time.Time `xorm:"start_date"`
	EndDate   time.Time `xorm:"end_date"`
	Status    string    `xorm:"varchar(20) not null" valid:"required,in(wait,running,failed,completed)"`
	CreateAt  time.Time `xorm:"created"`
	UpdateAt  time.Time `xorm:"updated"`
}

//TableName name of table
func (task Task) TableName() string {
	return "tasks"
}
