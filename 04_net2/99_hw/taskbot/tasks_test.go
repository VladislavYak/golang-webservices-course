package main

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"

	"fmt"

	tgbotapi "github.com/skinass/telegram-bot-api/v5"

	// "io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func init() {
	// upd global var for testing
	// we use patched version of gopkg.in/telegram-bot-api.v4 ( WebhookURL const -> var)
	WebhookURL = "http://127.0.0.1:8081"
	BotToken = "_golangcourse_test"
}

var (
	client = &http.Client{Timeout: time.Second}
)

// TDS is Telegram Dummy Server
type TDS struct {
	*sync.Mutex
	Answers map[int64]string
}

func NewTDS() *TDS {
	return &TDS{
		Mutex:   &sync.Mutex{},
		Answers: make(map[int64]string),
	}
}

func (srv *TDS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/getMe", func(w http.ResponseWriter, r *http.Request) {
		//nolint:errcheck
		w.Write([]byte(`{"ok":true,"result":{"id":` +
			strconv.Itoa(BotChatID) +
			`,"is_bot":true,"first_name":"game_test_bot","username":"game_test_bot"}}`))
	})
	mux.HandleFunc("/setWebhook", func(w http.ResponseWriter, r *http.Request) {
		//nolint:errcheck
		w.Write([]byte(`{"ok":true,"result":true,"description":"Webhook was set"}`))
	})
	mux.HandleFunc("/sendMessage", func(w http.ResponseWriter, r *http.Request) {
		//nolint:errcheck
		chatID, _ := strconv.ParseInt(r.FormValue("chat_id"), 10, 64)
		text := r.FormValue("text")
		srv.Lock()
		srv.Answers[chatID] = text
		srv.Unlock()

		//nolint:errcheck
		w.Write([]byte(`{"ok":true, "result":{"MessageID": 0}}`))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		panic(fmt.Errorf("unknown command %s", r.URL.Path))
	})

	handler := http.StripPrefix("/bot"+BotToken, mux)
	handler.ServeHTTP(w, r)
}

const (
	Ivanov     int64 = 256
	Petrov     int64 = 512
	Alexandrov int64 = 1024
	BotChatID        = 100500
)

var (
	users = map[int64]*tgbotapi.User{
		Ivanov: &tgbotapi.User{
			ID:           Ivanov,
			FirstName:    "Ivan",
			LastName:     "Ivanov",
			UserName:     "ivanov",
			LanguageCode: "ru",
			IsBot:        false,
		},
		Petrov: &tgbotapi.User{
			ID:           Petrov,
			FirstName:    "Petr",
			LastName:     "Pertov",
			UserName:     "ppetrov",
			LanguageCode: "ru",
			IsBot:        false,
		},
		Alexandrov: &tgbotapi.User{
			ID:           Alexandrov,
			FirstName:    "Alex",
			LastName:     "Alexandrov",
			UserName:     "aalexandrov",
			LanguageCode: "ru",
			IsBot:        false,
		},
	}

	updID uint64
	msgID uint64
)

func SendMsgToBot(userID int64, text string) error {
	// reqText := `{
	// 	"update_id":175894614,
	// 	"message":{
	// 		"message_id":29,
	// 		"from":{"id":133250764,"is_bot":false,"first_name":"Vasily Romanov","username":"rvasily","language_code":"ru"},
	// 		"chat":{"id":133250764,"first_name":"Vasily Romanov","username":"rvasily","type":"private"},
	// 		"date":1512168732,
	// 		"text":"THIS SEND FROM USER"
	// 	}
	// }`

	myUpdID := atomic.AddUint64(&updID, 1)

	// better have it per user, but lazy now
	myMsgID := atomic.AddUint64(&msgID, 1)

	user, ok := users[userID]
	if !ok {
		return fmt.Errorf("no user for %d", userID)
	}

	upd := &tgbotapi.Update{
		UpdateID: int(myUpdID),
		Message: &tgbotapi.Message{
			MessageID: int(myMsgID),
			From:      user,
			Chat: &tgbotapi.Chat{
				ID:        user.ID,
				FirstName: user.FirstName,
				UserName:  user.UserName,
				Type:      "private",
			},
			Text: text,
			Date: int(time.Now().Unix()),
			Entities: []tgbotapi.MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: len(strings.Split(text, " ")[0]),
				},
			},
		},
	}
	//nolint:errcheck
	reqData, _ := json.Marshal(upd)

	reqBody := bytes.NewBuffer(reqData)
	//nolint:errcheck
	req, _ := http.NewRequest(http.MethodPost, WebhookURL, reqBody)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	//nolint:govet
	defer resp.Body.Close()
	return err
}

func UpdateLastMessage(userID int64) error {

	//{
	//	"update_id": 443427246,
	//	"edited_message": {
	//	"message_id": 24,
	//		"from": {
	//		"id": 337749172,
	//			"is_bot": false,
	//			"first_name": "A",
	//			"username": "test",
	//			"language_code": "ru"
	//	},
	//	"chat": {
	//		"id": 337749172,
	//			"first_name": "A",
	//			"username": "test",
	//			"type": "private"
	//	},
	//	"date": 1699102952,
	//		"edit_date": 1699102963,
	//		"text": "12312323"
	//}
	//}

	myUpdID := atomic.AddUint64(&updID, 1)

	// better have it per user, but lazy now
	myMsgID := atomic.LoadUint64(&msgID)

	user, ok := users[userID]
	if !ok {
		return fmt.Errorf("no user for %d", userID)
	}

	upd := &tgbotapi.Update{
		UpdateID: int(myUpdID),
		EditedMessage: &tgbotapi.Message{
			MessageID: int(myMsgID),
			From:      user,
			Chat: &tgbotapi.Chat{
				ID:        user.ID,
				FirstName: user.FirstName,
				UserName:  user.UserName,
				Type:      "private",
			},
			Text:     "new_text",
			Date:     int(time.Now().Unix()),
			EditDate: int(time.Now().Unix()),
		},
	}
	//nolint:errcheck
	reqData, _ := json.Marshal(upd)

	reqBody := bytes.NewBuffer(reqData)
	//nolint:errcheck
	req, _ := http.NewRequest(http.MethodPost, WebhookURL, reqBody)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	//nolint:govet
	defer resp.Body.Close()
	return err
}

type testCase struct {
	user    int64
	command string
	answers map[int64]string
}

func TestTasks(t *testing.T) {

	tds := NewTDS()
	ts := httptest.NewServer(tds)
	tgbotapi.APIEndpoint = ts.URL + "/bot%s/%s"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := startTaskBot(ctx)
		if err != nil {
			//nolint:govet
			t.Fatalf("startTaskBot error: %s", err)
		}
	}()

	// give server time to start
	time.Sleep(100 * time.Millisecond)

	cases := []testCase{
		{
			// команда /tasks - выводит список всех активных задач
			Ivanov,
			"/tasks",
			map[int64]string{
				Ivanov: "Нет задач",
			},
		},
		{
			// команда /new - создаёт новую задачу, всё что после /new - идёт в название задачи
			Ivanov,
			"/new написать бота",
			map[int64]string{
				Ivanov: `Задача "написать бота" создана, id=1`,
			},
		},
		{
			Ivanov,
			"/tasks",
			map[int64]string{
				Ivanov: `1. написать бота by @ivanov
/assign_1`,
			},
		},
		{
			Ivanov,
			"/assign_1",
			map[int64]string{
				Ivanov: `Задача "написать бота" назначена на вас`,
			},
		},
		{
			// /assign_* - назначает задачу на себя
			Alexandrov,
			"/assign_1",
			map[int64]string{
				Alexandrov: `Задача "написать бота" назначена на вас`,
				Ivanov:     `Задача "написать бота" назначена на @aalexandrov`,
			},
		},
		{
			// в случае если задача была назначена на кого-то - он получает уведомление об этом
			// в данном случае она была назначена на Alexandrov, поэтому ему отправляется уведомление
			Petrov,
			"/assign_1",
			map[int64]string{
				Petrov:     `Задача "написать бота" назначена на вас`,
				Alexandrov: `Задача "написать бота" назначена на @ppetrov`,
			},
		},
		{
			// если задача назначена и на мне - показывается "на меня"
			Petrov,
			"/tasks",
			map[int64]string{
				Petrov: `1. написать бота by @ivanov
assignee: я
/unassign_1 /resolve_1`,
			},
		},
		{
			// если задача назначена и не на мне - показывается логин исполнителя
			// при
			Ivanov,
			"/tasks",
			map[int64]string{
				Ivanov: `1. написать бота by @ivanov
assignee: @ppetrov`,
			},
		},

		{
			// /unassign_ - снимает задачу с себя
			// нельзя снять задачу которая не на вас
			Alexandrov,
			"/unassign_1",
			map[int64]string{
				Alexandrov: `Задача не на вас`,
			},
		},

		{
			// /unassign_ - снимает задачу с себя
			// автору отправляется уведомление о том, что задача осталась без исполнителя
			Petrov,
			"/unassign_1",
			map[int64]string{
				Petrov: `Принято`,
				Ivanov: `Задача "написать бота" осталась без исполнителя`,
			},
		},

		{
			// повтор
			// в случае если задача была назначена на кого-то - автор получает уведомление об этом
			Petrov,
			"/assign_1",
			map[int64]string{
				Petrov: `Задача "написать бота" назначена на вас`,
				Ivanov: `Задача "написать бота" назначена на @ppetrov`,
			},
		},
		{
			// /resolve_* завершает задачу, удаляет её из хранилища
			// автору отправляется уведомление об этом
			Petrov,
			"/resolve_1",
			map[int64]string{
				Petrov: `Задача "написать бота" выполнена`,
				Ivanov: `Задача "написать бота" выполнена @ppetrov`,
			},
		},

		{
			Petrov,
			"/tasks",
			map[int64]string{
				Petrov: `Нет задач`,
			},
		},

		{
			// обратите внимание, id=2 - автоинкремент
			Petrov,
			"/new сделать ДЗ по курсу",
			map[int64]string{
				Petrov: `Задача "сделать ДЗ по курсу" создана, id=2`,
			},
		},
		{
			// обратите внимание, id=3 - автоинкремент
			Ivanov,
			"/new прийти на хакатон",
			map[int64]string{
				Ivanov: `Задача "прийти на хакатон" создана, id=3`,
			},
		},
		{
			Petrov,
			"/tasks",
			map[int64]string{
				Petrov: `2. сделать ДЗ по курсу by @ppetrov
/assign_2

3. прийти на хакатон by @ivanov
/assign_3`,
			},
		},
		{
			// повтор
			// в случае если задача была назначена на кого-то - автор получает уведомление об этом
			// если он автор задачи - ему не приходит дополнительного уведомления о том что она назначена на кого-то
			Petrov,
			"/assign_2",
			map[int64]string{
				Petrov: `Задача "сделать ДЗ по курсу" назначена на вас`,
			},
		},
		{
			Petrov,
			"/tasks",
			map[int64]string{
				Petrov: `2. сделать ДЗ по курсу by @ppetrov
assignee: я
/unassign_2 /resolve_2

3. прийти на хакатон by @ivanov
/assign_3`,
			},
		},
		{
			// /my показывает задачи которые назначены на меня
			// при этому тут нет метки assegnee
			Petrov,
			"/my",
			map[int64]string{
				Petrov: `2. сделать ДЗ по курсу by @ppetrov
/unassign_2 /resolve_2`,
			},
		},
		{
			// /owner - показывает задачи, которы я создал
			// при этому тут нет метки assegnee
			Ivanov,
			"/owner",
			map[int64]string{
				Ivanov: `3. прийти на хакатон by @ivanov
/assign_3`,
			},
		},
	}

	var unexpectedUpdate sync.Once

	for idx, item := range cases {

		tds.Lock()
		tds.Answers = make(map[int64]string)
		tds.Unlock()

		caseName := fmt.Sprintf("[case%d, %d: %s]", idx, item.user, item.command)
		err := SendMsgToBot(item.user, item.command)
		if err != nil {
			t.Fatalf("%s SendMsgToBot error: %s", caseName, err)
		}
		// give TDS time to process request
		time.Sleep(10 * time.Millisecond)

		// бот может получать разные апдейты. Например, если пользователь отредактирует сообщение, то вы получите update.Message == nil.
		// в этом тесте мы проверяем, что вы предусмотрели это, и добавили проверку на nil. Иначе, словите панику при попытке обратиться к update.Message.*
		unexpectedUpdate.Do(func() {
			err := UpdateLastMessage(item.user)
			if err != nil {
				t.Fatalf("%s UpdateMessage error: %s", caseName, err)
			}
		})

		tds.Lock()
		result := reflect.DeepEqual(tds.Answers, item.answers)
		if !result {
			t.Fatalf("%s bad results:\n\tWant: %v\n\tHave: %v", caseName, item.answers, tds.Answers)
		}
		tds.Unlock()

	}

}
