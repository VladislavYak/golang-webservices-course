package db

import (
	"errors"
	"fmt"
	"sync"

	"github.com/VladislavYak/taskbot/models"
)

var TASKS = NewTasks()
var taskID = 1

type Tasks struct {
	TasksDB []models.Task
	*sync.Mutex
}

// pretty random setting mutex here, in production concurrency would be controlled by DB, so lazy do it here
func NewTasks() *Tasks {
	return &Tasks{
		[]models.Task{},
		&sync.Mutex{},
	}
}

func (t *Tasks) HasTask(task *models.Task) bool {
	for _, task_ := range t.TasksDB {
		if task_.TaskID == task.TaskID {
			return true
		}
	}
	return false
}

func (t *Tasks) AddTask(task *models.Task) (int, error) {

	oldTaskId := taskID
	if !t.HasTask(task) {
		t.Mutex.Lock()
		task.TaskID = oldTaskId
		taskID++

		fmt.Println("*task", *task)
		t.TasksDB = append(t.TasksDB, *task)
		t.Mutex.Unlock()
		return oldTaskId, nil
	} else {
		return 0, errors.New("the task already exists")
	}
}

func (t *Tasks) GetAllTasks() *[]models.Task {
	return &t.TasksDB
}

func (t *Tasks) GetAllTasksCreatedByUser(user models.User) *[]models.Task {
	out := []models.Task{}
	for _, task := range t.TasksDB {
		t.Mutex.Lock()
		if task.ByUserID == user.UserID {
			out = append(out, task)
		}
		t.Mutex.Unlock()
	}
	return &out

}

func (t *Tasks) GetAllTasksAssignedToUser(user models.User) *[]models.Task {
	out := []models.Task{}
	for _, task := range t.TasksDB {
		t.Mutex.Lock()
		if task.Asignee == user.UserID {
			out = append(out, task)
		}
		t.Mutex.Unlock()
	}
	return &out
}

func (t *Tasks) AssignTaskToUser(TaskID int, UserID string) error {
	for i, task := range t.TasksDB {
		t.Mutex.Lock()
		if task.TaskID == TaskID {

			t.TasksDB[i].Asignee = UserID
			return nil
		}
		t.Mutex.Unlock()
	}
	return errors.New("cannot find user as an assignee for passed task id")

}

func (t *Tasks) UnassignTaskFromUser(TaskID int, UserID string) error {
	for i, task := range t.TasksDB {
		t.Mutex.Lock()
		if (task.TaskID == TaskID) && (task.Asignee == UserID) {
			t.TasksDB[i].Asignee = ""
			return nil
		}
		t.Mutex.Unlock()
	}
	return errors.New("cannot find user as an assignee for passed task id")
}

func (t *Tasks) DeleteTask(TaskID int) error {
	index := -1
	for i, task := range t.TasksDB {
		t.Mutex.Lock()
		if task.TaskID == TaskID {
			index = i
		}
		t.Mutex.Unlock()
	}

	if index == -1 {
		return errors.New("cannot find that task")
	} else {
		t.Mutex.Lock()
		t.TasksDB = append(t.TasksDB[:index], t.TasksDB[index+1:]...)
		t.Mutex.Unlock()
		return nil
	}
}
