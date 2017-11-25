package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/juliotorresmoreno/pomodoro-server/tasks"

	"github.com/juliotorresmoreno/pomodoro-server/models"
	"github.com/pingcap/tidb/terror"
	"golang.org/x/crypto/bcrypt"

	"github.com/juliotorresmoreno/pomodoro-server/db"
	"github.com/juliotorresmoreno/pomodoro-server/util"
	"github.com/juliotorresmoreno/pomodoro-server/ws"
)

//Auth
type Auth struct {
	hub         *ws.Hub
	taskManager tasks.TaskManager
}

//NewAuth s
func NewAuth(hub *ws.Hub, taskManager tasks.TaskManager) Auth {
	return Auth{
		hub:         hub,
		taskManager: taskManager,
	}
}

//Login login user
func (auth Auth) Login(w http.ResponseWriter, r *http.Request) {
	data := util.GetPostParams(r)
	username := data.Get("username")
	password := data.Get("password")
	conn, err := db.NewConnection()
	if err != nil {
		terror.Log(err)
	}
	defer conn.Close()
	user := models.User{}
	conn.Where("username = ?", username).Get(&user)
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err == nil {
		user.Token = util.GenerateRandomString(40)
		if _, err := conn.Where("id = ?", user.ID).Update(user); err != nil {
			terror.Log(err)
			w.WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "user or password not valid!",
			})
			return
		}
		auth.taskManager.Load(user.Username)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"token":   user.Token,
		})
		return
	}
	w.WriteHeader(http.StatusNotAcceptable)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"message": "user or password not valid!",
	})
}

//Session user
func (auth Auth) Session(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
		})
		return
	}
	conn, err := db.NewConnection()
	if err != nil {
		terror.Log(err)
	}
	defer conn.Close()
	user := models.User{}
	conn.Where("token = ?", token).Get(&user)

	if user.Token == "" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"session": map[string]interface{}{
			"name":     user.Name,
			"lastname": user.LastName,
			"username": user.Username,
			"token":    token,
		},
	})
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
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Database not found",
		})
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
	digest, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.Password = string(digest)

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
