package tasks

import (
	"github.com/juliotorresmoreno/pomodoro-server/db"
	"github.com/juliotorresmoreno/pomodoro-server/models"
	"github.com/pingcap/tidb/terror"
)

type TaskManager map[string][]Task

func NewTaskManager() TaskManager {
	return TaskManager{}
}

func (taskManager TaskManager) Load(username string) error {
	if _, ok := taskManager[username]; ok {
		return nil
	}
	taskManager[username] = make([]Task, 0)
	conn, err := db.NewConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	tasks := make([]models.Task, 0)
	conn.Where("user_id = ?").Find(&tasks)
	for _, task := range tasks {
		_task := Task{
			task: task,
		}
		taskManager[username] = append(taskManager[username], _task)
	}
	return nil
}

func (taskManager TaskManager) NewTask(username string, task models.Task) (Task, error) {
	if _, ok := taskManager[username]; !ok {
		taskManager[username] = make([]Task, 0)
	}
	_task := Task{
		task: task,
	}
	taskManager[username] = append(taskManager[username], _task)
	conn, err := db.NewConnection()
	if err != nil {
		terror.Log(err)
		return Task{}, err
	}
	defer conn.Close()
	task.ID, err = conn.Insert(task)
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

}

func (task Task) ID() int64 {
	return task.task.ID
}
