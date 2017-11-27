package models

import (
	"time"
)

type Task struct {
	ID        int64     `xorm:"id bigint not null autoincr pk" json:"id" json:"id"`
	Name      string    `xorm:"varchar(100) not null" valid:"required,alphaSpaces" json:"name"`
	Step      int8      `xorm:"int not null" json:"step"`
	UserID    int64     `xorm:"user_id bigint not null" valid:"required" json:"user_id"`
	StartDate time.Time `xorm:"start_date" json:"start_date"`
	EndDate   time.Time `xorm:"end_date" json:"end_date"`
	Status    string    `xorm:"varchar(20) not null" valid:"required" json:"status"` //wait,running,failed,completed
	CreateAt  time.Time `xorm:"created" json:"-"`
	UpdateAt  time.Time `xorm:"updated" json:"-"`
}

//TableName name of table
func (task Task) TableName() string {
	return "tasks"
}
