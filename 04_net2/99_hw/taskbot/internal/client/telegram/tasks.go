package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/VladislavYak/taskbot/models"
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

func (t *Telegram) tasks(chatID int64, user models.User) {
	resp, err := http.Get(t.Url + "/tasks")
	if err != nil {
		fmt.Println("err", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	tasks := &[]models.Task{}
	if err := json.Unmarshal(body, tasks); err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
	}

	fmt.Println("tasks (inside tasks)", tasks)
	fmt.Println("tasks (inside tasks)", tasks)

	tgResponse := PrintTasks(tasks, user.UserID)
	msg := tgbotapi.NewMessage(chatID, tgResponse)
	t.bot.Send(msg)
}

func PrintTasks(tasks *[]models.Task, userID string) string {

	if len(*tasks) == 0 {
		return "Нет задач"
	}

	out := ""
	for _, task := range *tasks {
		if task.Asignee == "" {
			out += strconv.Itoa(task.TaskID) + ". " + task.TaskName + " by " + task.ByUserName + "\n" + "/assign_" + strconv.Itoa(task.TaskID) + "\n"
		} else if task.ByUserID != userID {
			out += strconv.Itoa(task.TaskID) + ". " + task.TaskName + " by " + task.ByUserName + " assignee: " + task.Asignee + "\n"
		} else if task.ByUserID == userID {
			out += strconv.Itoa(task.TaskID) + ". " + task.TaskName + " by " + task.ByUserName + " assignee: я" + "\n"
		}
	}
	fmt.Println("tasks", tasks)
	return out
}
