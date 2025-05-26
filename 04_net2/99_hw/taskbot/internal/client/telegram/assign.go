package telegram

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/VladislavYak/taskbot/models"
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

func (t *Telegram) assign(user models.User, taskID int) {
	data := url.Values{}
	data.Set("user_id", user.UserID)
	data.Set("task_id", string(taskID))

	resp, err := http.PostForm(t.Url+"/assign", data)
	if err != nil {
		fmt.Println("err", err)
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)

	tgResponse := PrintAssign()
	msg := tgbotapi.NewMessage(int64(user.ChatID), tgResponse)
	t.bot.Send(msg)
}

func PrintAssign(tasks *[]models.Task) string {
	return "Задача была назначена на вас"
}
