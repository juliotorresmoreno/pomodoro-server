package db

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"github.com/go-xorm/xorm"
	"github.com/juliotorresmoreno/pomodoro-server/models"
	"github.com/pingcap/tidb/terror"
	"gopkg.in/oleiade/reflections.v1"
)

func init() {
	govalidator.TagMap["unique"] = govalidator.Validator(func(str string) bool {
		return true
	})
	govalidator.TagMap["alphaSpaces"] = govalidator.Validator(func(str string) bool {
		patterm, _ := regexp.Compile("^([a-zA-Z]+( ){0,1}){1,}$")
		return patterm.MatchString(str)
	})
	govalidator.TagMap["username"] = govalidator.Validator(func(str string) bool {
		patterm, _ := regexp.Compile("^[a-zA-Z0-9_]{3,}$")
		return patterm.MatchString(str)
	})
	govalidator.TagMap["password"] = govalidator.Validator(func(str string) bool {
		return len(str) > 4
	})
}

type Connection struct {
	*xorm.Engine
}

//ValidateStruct d
func (conn Connection) ValidateStruct(model models.Model) (bool, error) {
	result, err := govalidator.ValidateStruct(model)
	if err != nil {
		return result, err
	}

	data := reflect.TypeOf(model)
	length := data.NumField()

	id, _ := reflections.GetField(model, "ID")
	for i := 0; i < length; i++ {
		field := data.Field(i)
		valid := field.Tag.Get("valid")
		if strings.Contains(valid, "unique") {
			value, _ := reflections.GetField(model, field.Name)
			count, _ := conn.Table(model.TableName()).
				Where("id != ?", id).
				And(field.Name+" = ?", value).
				Count()
			if count > 0 {
				return false, fmt.Errorf("%v: %v exists", field.Name, value)
			}
		}
	}
	return true, nil
}

//NewConnection new connection
func NewConnection() (*Connection, error) {
	//conn, err := xorm.NewEngine("sqlite3", "./data/database.db")
	conn, err := xorm.NewEngine("mysql", "root:paramore@tcp(127.0.0.1:3306)/pomodoro")
	if err != nil {
		terror.Log(err)
	}
	conn.ShowSQL(false)
	return &Connection{
		Engine: conn,
	}, err
}

func Sync() {
	conn, err := NewConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	if err = conn.Sync2(new(models.User)); err != nil {
		log.Fatal(err)
	}
	if err = conn.Sync2(new(models.Task)); err != nil {
		log.Fatal(err)
	}
}
