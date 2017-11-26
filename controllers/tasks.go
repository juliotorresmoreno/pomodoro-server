package controllers

import (
	"encoding/json"
	"net/http"

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

//NewTimer
func NewTasks(hub *ws.Hub) Tasks {
	return Tasks{
		hub:         hub,
		TaskManager: tasks.NewTaskManager(hub),
	}
}

func (tasks Tasks) NewPomodoro(w http.ResponseWriter, r *http.Request, session models.Session) {
	data := util.GetPostParams(r)
	name := data.Get("name")
	task := models.Task{
		Name:   name,
		UserID: session.ID,
	}
	manager := tasks.TaskManager
	_task, err := manager.NewTask(session.Username, task)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	_task.Start()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

func (tasks Tasks) List(w http.ResponseWriter, r *http.Request, session models.Session) {
	_tasks, err := tasks.TaskManager.Load(session)
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
}
