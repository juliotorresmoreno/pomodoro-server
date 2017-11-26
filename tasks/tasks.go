package tasks

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/juliotorresmoreno/pomodoro-server/db"
	"github.com/juliotorresmoreno/pomodoro-server/models"
	"github.com/juliotorresmoreno/pomodoro-server/ws"
	"github.com/pingcap/tidb/terror"
)

type TaskManager struct {
	hub *ws.Hub
	taskManager
}

type taskManager map[string][]Task

func NewTaskManager(hub *ws.Hub) TaskManager {
	return TaskManager{hub: hub, taskManager: taskManager{}}
}

func (taskManager TaskManager) Load(session models.Session) ([]models.Task, error) {
	_taskManager := taskManager.taskManager
	username := session.Username
	if _, ok := _taskManager[username]; ok {
		return nil, nil
	}
	_taskManager[username] = make([]Task, 0)
	conn, err := db.NewConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	tasks := make([]models.Task, 0)
	conn.Where("user_id = ?", session.ID).Find(&tasks)
	for _, task := range tasks {
		_task := Task{
			task: task,
		}
		_taskManager[username] = append(_taskManager[username], _task)
	}
	return tasks, nil
}

func (taskManager TaskManager) NewTask(username string, task models.Task) (Task, error) {
	_taskManager := taskManager.taskManager
	if _, ok := _taskManager[username]; !ok {
		_taskManager[username] = make([]Task, 0)
	}
	_task := Task{
		task: task,
	}
	_taskManager[username] = append(_taskManager[username], _task)
	conn, err := db.NewConnection()
	if err != nil {
		terror.Log(err)
		return Task{}, err
	}
	defer conn.Close()
	task.ID, err = conn.Insert(task)
	log.Infof("TaskID: %v", task.ID)
	if err != nil {
		terror.Log(err)
		return Task{}, err
	}
	return _task, nil
}

type Task struct {
	task models.Task
}

func (task Task) Start() {
	conn, err := db.NewConnection()
	if err != nil {
		terror.Log(err)
		return
	}
	task.task.StartDate = time.Now()
	_, err = conn.Where("id = ?", task.task.ID).Update(task.task)
	if err != nil {
		terror.Log(err)
	}
}

func (task Task) ID() int64 {
	return task.task.ID
}
