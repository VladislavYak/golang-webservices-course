package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/VladislavYak/taskbot/models"
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

func (t *Telegram) new(user *models.User, taskname string, chatID int64) {

	data := url.Values{}
	data.Set("taskname", taskname)
	data.Set("username", user.UserName)
	data.Set("user_id", user.UserID)

	resp, err := http.PostForm(t.Url+"/new", data)
	if err != nil {
		fmt.Println("err", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	taskID := &struct {
		TaskID string `json:"task_id"`
	}{}
	if err := json.Unmarshal(body, taskID); err != nil {

	}

	tgResponse := PrintNew(taskname, user, taskID.TaskID)
	msg := tgbotapi.NewMessage(chatID, tgResponse)
	t.bot.Send(msg)
}

func PrintNew(taskname string, user *models.User, taskID string) string {
	tgResponse := "Задача " + taskname + " создана, id=" + taskID
	return tgResponse
}
