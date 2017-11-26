package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/Sirupsen/logrus"
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
	session, status, err := auth.login(conn, username, password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Database not found",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"session": session,
	})
}

func (auth Auth) login(conn *db.Connection, username, password string) (models.Session, int, error) {
	user := models.User{}
	conn.Where("username = ?", username).Get(&user)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err == nil {
		user.Token = util.GenerateRandomString(40)
		if _, err := conn.Where("id = ?", user.ID).Update(user); err != nil {
			terror.Log(err)
			return models.Session{}, http.StatusNotAcceptable, err
		}
		session := models.Session{
			ID:       user.ID,
			Name:     user.Name,
			LastName: user.LastName,
			Username: user.Username,
			Token:    user.Token,
		}
		auth.taskManager.Load(session)
		return session, http.StatusOK, nil
	}
	return models.Session{}, http.StatusNotAcceptable, errors.New("User or password not valid")
}

//Session user
func (auth Auth) Session(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		w.Header().Set("Content-Type", "application/json")
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"session": models.Session{
			ID:       user.ID,
			Name:     user.Name,
			LastName: user.LastName,
			Username: user.Username,
			Token:    user.Token,
		},
	})
}

//Register dd
func (auth Auth) Register(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Header)
	data := util.GetPostParams(r)
	name := data.Get("name")
	lastname := data.Get("lastname")
	username := data.Get("username")
	password := data.Get("password")

	conn, err := db.NewConnection()
	if err != nil {
		terror.Log(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
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

		w.Header().Set("Content-Type", "application/json")
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	session, status, err := auth.login(conn, username, password)
	if err != nil {
		terror.Log(err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"session": session,
	})
}
