package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/juliotorresmoreno/pomodoro-server/models"
	"github.com/juliotorresmoreno/pomodoro-server/tasks"
	"github.com/juliotorresmoreno/pomodoro-server/util"
	"github.com/juliotorresmoreno/pomodoro-server/ws"
)

//Timer
type Timer struct {
	hub         *ws.Hub
	TaskManager tasks.TaskManager
}

//NewTimer
func NewTimer(hub *ws.Hub) Timer {
	return Timer{
		hub:         hub,
		TaskManager: tasks.NewTaskManager(),
	}
}

func (timer Timer) NewPomodoro(w http.ResponseWriter, r *http.Request, session models.Session) {
	data := util.GetPostParams(r)
	name := data.Get("name")
	task := models.Task{
		Name: name,
	}
	manager := timer.TaskManager
	_task, err := manager.NewTask(session.Username, task)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	_task.Start()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"task_id": _task.ID(),
	})
}
