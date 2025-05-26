package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/VladislavYak/taskbot/internal/db"
	"github.com/VladislavYak/taskbot/models"
)

func Tasks(w http.ResponseWriter, r *http.Request) {
	jsonResponse, err := json.Marshal(db.TASKS.GetAllTasks())

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func New(w http.ResponseWriter, r *http.Request) {
	TaskName := r.FormValue("taskname")
	UserID := r.FormValue("user_id")
	UserName := r.FormValue("username")

	taskID, err := db.TASKS.AddTask(&models.Task{TaskName: TaskName, ByUserID: UserID, ByUserName: UserName})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	message := fmt.Sprintf(`{"task_id": "%s"}`, strconv.Itoa(taskID))
	io.WriteString(w, message)

}

func Owner(w http.ResponseWriter, r *http.Request) {
	UserID := r.FormValue("user_id")

	if UserID == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	tasks := db.TASKS.GetAllTasksCreatedByUser(models.User{UserID: UserID})

	jsonResponse, err := json.Marshal(tasks)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func My(w http.ResponseWriter, r *http.Request) {
	UserID := r.FormValue("user_id")

	if UserID == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	tasks := db.TASKS.GetAllTasksAssignedToUser(models.User{UserID: UserID})

	jsonResponse, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func Assign(w http.ResponseWriter, r *http.Request) {
	UserID := r.FormValue("user_id")
	TaskID, _ := strconv.Atoi(r.FormValue("task_id"))

	err := db.TASKS.AssignTaskToUser(TaskID, UserID)
	if err != nil {

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func Unassign(w http.ResponseWriter, r *http.Request) {
	UserID := r.FormValue("user_id")
	TaskID, _ := strconv.Atoi(r.FormValue("task_id"))

	err := db.TASKS.UnassignTaskFromUser(TaskID, UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"Error": "cannot find user as an assignee for passed task id"}`)
	}

}

func Resolve(w http.ResponseWriter, r *http.Request) {
	TaskID, _ := strconv.Atoi(r.FormValue("task_id"))

	err := db.TASKS.DeleteTask(TaskID)
	if err != nil {

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
