package tasks

import (
	"fmt"
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
	running     *time.Timer
	taskRunning int64
}

type taskManager map[string]map[int64]Task

func NewTaskManager(hub *ws.Hub) TaskManager {
	return TaskManager{
		hub:         hub,
		taskManager: taskManager{},
	}
}

func (taskManager TaskManager) Load(session models.Session) ([]models.Task, error) {
	conn, err := db.NewConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return taskManager.load(conn, session)
}

func (taskManager TaskManager) load(conn *db.Connection, session models.Session) ([]models.Task, error) {
	_taskManager := taskManager.taskManager
	username := session.Username
	if _, ok := _taskManager[username]; !ok {
		fmt.Println(username, taskManager, taskManager.taskManager)
		taskManager.taskManager[username] = map[int64]Task{}
		conn.Where("user_id = ?", session.ID).Update(models.Task{
			Status: "wait",
		})
	}
	tasks := make([]models.Task, 0)
	conn.Where("user_id = ?", session.ID).Find(&tasks)
	for _, task := range tasks {
		_task := Task{
			task: task,
		}
		_taskManager[username][task.ID] = _task
	}
	return tasks, nil
}

func (taskManager TaskManager) Start(session models.Session, id int64, notification func(models.Session, []models.Task)) ([]models.Task, error) {
	conn, err := db.NewConnection()
	if err != nil {
		terror.Log(err)
		return nil, err
	}
	defer conn.Close()
	task := models.Task{}
	conn.Where("id = ? and user_id = ?", id, session.ID).
		Get(&task)
	if task.Step > 4 {
		task.Step = 0
	}
	task.StartDate = time.Now()
	task.Status = "running"
	conn.ShowSQL(true)
	conn.Where("id = ? and user_id = ?", id, session.ID).
		Update(task)
	conn.ShowSQL(false)
	task.ID = id
	taskManager.taskRunning = id

	taskManager.setTimeout(session, task, notification)
	if err != nil {
		terror.Log(err)
		return nil, err
	}
	return taskManager.load(conn, session)
}

func (taskManager TaskManager) setTimeout(session models.Session, task models.Task, notification func(models.Session, []models.Task)) {
	timer := time.AfterFunc(10*time.Second, func() {
		conn, err := db.NewConnection()
		if err != nil {
			terror.Log(err)
			return
		}
		defer conn.Close()
		conn.ShowSQL(true)
		task.Step = task.Step + 1
		_, err = conn.Where("id = ? and user_id = ?", task.ID, session.ID).
			Update(models.Task{
				Step: task.Step + 1,
			})
		conn.ShowSQL(false)
		if err != nil {
			task.Step = task.Step - 1
			terror.Log(err)
		}
		log.Info("Tiempo transcurrido")
		if task.Step < 3 && taskManager.taskRunning == task.ID {
			taskManager.sleep(session, task, notification)
		} else {
			taskManager.lastSleep(session, task, notification)
		}
		if notification != nil {
			tasks, _ := taskManager.load(conn, session)
			notification(session, tasks)
		}
	})
	taskManager.running = timer
}

func (taskManager TaskManager) sleep(session models.Session, task models.Task, notification func(models.Session, []models.Task)) {
	timer := time.AfterFunc(5*time.Second, func() {
		log.Info("Descanzo transcurrido")
		if taskManager.taskRunning == task.ID {
			taskManager.setTimeout(session, task, notification)
		}
	})
	if notification != nil {
		taskManager.running.Stop()
	}
	taskManager.running = timer
}

func (taskManager TaskManager) lastSleep(session models.Session, task models.Task, notification func(models.Session, []models.Task)) {
	timer := time.AfterFunc(5*time.Second, func() {
		conn, err := db.NewConnection()
		if err != nil {
			terror.Log(err)
			return
		}
		defer conn.Close()
		conn.Where("id = ? and user_id = ?", task.ID, session.ID).
			Cols("step", "end_date", "status").
			Update(models.Task{
				Step:    0,
				EndDate: time.Now(),
				Status:  "completed",
			})
		log.Info("Descanzo final transcurrido")
	})
	if taskManager.running != nil {
		taskManager.running.Stop()
	}
	taskManager.running = timer
}

func (taskManager TaskManager) Stop(session models.Session, id int64) ([]models.Task, error) {
	conn, err := db.NewConnection()
	if err != nil {
		terror.Log(err)
		return nil, err
	}
	defer conn.Close()
	task := models.Task{
		ID:     id,
		UserID: session.ID,
		Status: "wait",
	}
	conn.Where("id = ? and user_id = ?", id, session.ID).Update(task)
	if err != nil {
		terror.Log(err)
		return nil, err
	}
	if taskManager.running != nil {
		log.Info("Tiempo eliminado")
		taskManager.running.Stop()
		taskManager.taskRunning = 0
	}
	return taskManager.load(conn, session)
}

func (taskManager TaskManager) Delete(session models.Session, id int64) ([]models.Task, error) {
	conn, err := db.NewConnection()
	if err != nil {
		terror.Log(err)
		return nil, err
	}
	defer conn.Close()
	task := models.Task{
		ID:     id,
		UserID: session.ID,
	}
	conn.Delete(&task)
	log.Infof("delete %v", id)
	if err != nil {
		terror.Log(err)
		return nil, err
	}
	return taskManager.load(conn, session)
}

func (taskManager TaskManager) NewTask(session models.Session, task models.Task) (Task, error) {
	_taskManager := taskManager.taskManager
	username := session.Username
	if _, ok := _taskManager[username]; !ok {
		_taskManager[username] = map[int64]Task{}
	}
	_task := Task{
		task: task,
	}
	conn, err := db.NewConnection()
	if err != nil {
		terror.Log(err)
		return Task{}, err
	}
	defer conn.Close()
	task.UserID = session.ID
	task.Status = "wait"
	if _, err := conn.ValidateStruct(task); err != nil {
		terror.Log(err)
		return Task{}, err
	}
	_, err = conn.Insert(&task)
	_taskManager[username][task.ID] = _task
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
