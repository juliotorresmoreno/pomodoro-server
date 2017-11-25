package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/juliotorresmoreno/pomodoro-server/models"
	"github.com/pingcap/tidb/terror"

	"github.com/juliotorresmoreno/pomodoro-server/db"
	"github.com/juliotorresmoreno/pomodoro-server/util"
	"github.com/juliotorresmoreno/pomodoro-server/ws"
)

//Auth
type Auth struct {
	hub *ws.Hub
}

//Auth
func NewAuth(hub *ws.Hub) Auth {
	return Auth{
		hub: hub,
	}
}

//Login login user
func (auth Auth) Login(w http.ResponseWriter, r *http.Request) {
	data := util.GetPostParams(r)
	username := data.Get("username")
	password := data.Get("password")
	fmt.Fprintf(w, "username: %v, password: %v", username, password)
}

//Register dd
func (auth Auth) Register(w http.ResponseWriter, r *http.Request) {
	data := util.GetPostParams(r)
	name := data.Get("name")
	lastname := data.Get("lastname")
	username := data.Get("username")
	password := data.Get("password")

	conn, err := db.NewConnection()
	if err != nil {
		terror.Log(err)
		return
	}
	defer conn.Close()

	user := models.User{}
	user.Name = name
	user.LastName = lastname
	user.Username = username
	user.Password = password
	if _, err = conn.ValidateStruct(user); err != nil {
		terror.Log(err)
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	if _, err = conn.Insert(user); err != nil {
		terror.Log(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}
