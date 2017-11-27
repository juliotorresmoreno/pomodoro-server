package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/pingcap/tidb/terror"

	"github.com/juliotorresmoreno/pomodoro-server/db"
	"github.com/juliotorresmoreno/pomodoro-server/models"
	"github.com/juliotorresmoreno/pomodoro-server/tasks"
	"github.com/juliotorresmoreno/pomodoro-server/util"
	"github.com/juliotorresmoreno/pomodoro-server/ws"
)

//Tasks
type Tasks struct {
	hub         *ws.Hub
	TaskManager tasks.TaskManager
}

//NewTasks
func NewTasks(hub *ws.Hub) Tasks {
	_tasks := Tasks{}
	manager := tasks.NewTaskManager(hub)
	_tasks.TaskManager = manager
	_tasks.hub = hub
	return _tasks
}

func (tasks Tasks) NewPomodoro(w http.ResponseWriter, r *http.Request, session models.Session) {
	data := util.GetPostParams(r)
	name := data.Get("name")
	task := models.Task{
		Name:   name,
		UserID: session.ID,
	}
	manager := tasks.TaskManager
	_, err := manager.NewTask(session, task)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

func (tasks Tasks) Notification(session models.Session, _tasks []models.Task) {
	log.Info("@tasks/setList")
	tasks.hub.SendJSON(session.Username, map[string]interface{}{
		"type": "@tasks/setList",
		"data": _tasks,
	})
}

func (tasks Tasks) Start(w http.ResponseWriter, r *http.Request, session models.Session) {
	data := util.GetPostParams(r)
	id, _ := strconv.Atoi(data.Get("id"))
	if id == 0 {
		log.Infof("value: %v", id)
		return
	}
	_tasks, err := tasks.TaskManager.Start(session, int64(id), tasks.Notification)
	tasks.sendList(w, session, _tasks, err)
}

func (tasks Tasks) Statistics(w http.ResponseWriter, r *http.Request, session models.Session) {
	conn, err := db.NewConnection()
	if err != nil {
		terror.Log(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	defer conn.Close()
	response := make([]map[string]interface{}, 0, 2)
	sql := `SELECT status, COUNT(*) cant FROM tasks WHERE user_id = ? and status = 'wait' GROUP BY status`
	result, err := conn.QueryString(sql, session.ID)
	if err != nil {
		terror.Log(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	value := 0
	name := "wait"
	if len(result) > 0 {
		name = result[0]["status"]
		value, _ = strconv.Atoi(result[0]["cant"])
	}
	response = append(response, map[string]interface{}{
		"name": result[0]["status"],
		"v":    value,
	})

	sql = `SELECT status, COUNT(*) cant FROM tasks WHERE user_id = ? and status = 'completed' GROUP BY status`
	result, err = conn.QueryString(sql, session.ID)
	if err != nil {
		terror.Log(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	value = 0
	name = "completed"
	if len(result) > 0 {
		name = result[0]["status"]
		value, _ = strconv.Atoi(result[0]["cant"])
	}
	response = append(response, map[string]interface{}{
		"name": name,
		"v":    value,
	})

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    response,
	})
}

func (tasks Tasks) Stop(w http.ResponseWriter, r *http.Request, session models.Session) {
	data := util.GetPostParams(r)
	id, _ := strconv.Atoi(data.Get("id"))
	_tasks, err := tasks.TaskManager.Stop(session, int64(id))
	tasks.sendList(w, session, _tasks, err)
}

func (tasks Tasks) List(w http.ResponseWriter, r *http.Request, session models.Session) {
	_tasks, err := tasks.TaskManager.Load(session)
	tasks.sendList(w, session, _tasks, err)
}

func (tasks Tasks) Delete(w http.ResponseWriter, r *http.Request, session models.Session) {
	data := mux.Vars(r)
	id, _ := strconv.Atoi(data["id"])
	_tasks, err := tasks.TaskManager.Delete(session, int64(id))
	tasks.sendList(w, session, _tasks, err)
}

func (tasks Tasks) sendList(w http.ResponseWriter, session models.Session, _tasks []models.Task, err error) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    _tasks,
	})
	tasks.hub.SendJSON(session.Username, map[string]interface{}{
		"type": "@tasks/setList",
		"data": _tasks,
	})
}
