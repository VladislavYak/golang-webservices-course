package telegram

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/VladislavYak/taskbot/internal/db"
	models "github.com/VladislavYak/taskbot/models"
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

type Telegram struct {
	WebhookURL string
	BotToken   string
	Url        string
	// Wg         *sync.WaitGroup
	bot *tgbotapi.BotAPI
}

func NewTelegram(BotToken string, WebhookURL string) *Telegram {
	// defer wg.Done()
	bot, err := tgbotapi.NewBotAPI(BotToken)

	if err != nil {
		log.Fatalf("NewBotAPI failed: %s", err)
	}
	bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	wh, err := tgbotapi.NewWebhook(WebhookURL)
	if err != nil {
		log.Fatalf("NewWebhook failed: %s", err)
	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("SetWebhook failed: %s", err)
	}
	return &Telegram{WebhookURL: WebhookURL, BotToken: BotToken, Url: "http://localhost:8081", bot: bot}
}

func (t *Telegram) Router(wg *sync.WaitGroup) {
	defer wg.Done()

	updates := t.bot.ListenForWebhook("/")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	go func() {
		log.Fatalln("http err:", http.ListenAndServe(":"+port, nil))
	}()
	fmt.Println("start listen :" + port)

	for update := range updates {
		fmt.Println("update", update)
		myUser := update.SentFrom()
		chatID := update.FromChat().ChatConfig().ChatID

		user := models.User{
			ChatID:       int(chatID),
			UserID:       strconv.Itoa(int(myUser.ID)),
			UserName:     myUser.UserName,
			LanguageCode: myUser.LanguageCode,
		}
		db.USERS.AddUser(user)

		command := update.Message.Text

		if strings.HasPrefix(command, "/my") {
			t.my()
		} else if strings.HasPrefix(command, "/tasks") {
			t.tasks(chatID, user)
		} else if strings.HasPrefix(command, "/assign") {
			taskID := string(command[len(command)-1])

			taskIDconverted, err := strconv.Atoi(taskID)
			if err != nil {
				// ... handle error
				panic(err)
			}
			t.assign(user, taskIDconverted)
		} else if strings.HasPrefix(command, "/new") {
			taskName := strings.TrimSpace(strings.TrimPrefix(command, "/new"))

			t.new(&user, taskName, chatID)

		} else if strings.HasPrefix(command, "/owner") {
			// t.owner()
		}
	}
}
