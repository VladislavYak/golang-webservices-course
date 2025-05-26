package main

// сюда писать код

import (
	"context"
	"sync"

	"github.com/VladislavYak/taskbot/internal/client/telegram"
	router "github.com/VladislavYak/taskbot/internal/transport/http"
)

var (
	// @BotFather в телеграме даст вам это
	BotToken = "7930706080:AAGBqDBccmPjEGWl7-fIu4OIiJDPx7wYE-A"

	// урл выдаст вам игрок или хероку
	WebhookURL = "https://80c81c03134b28.lhr.life"
)

func startTaskBot(ctx context.Context) error {

	wg := &sync.WaitGroup{}
	wg.Add(2)
	// tg := telegram.Telegram{WebhookURL: WebhookURL, BotToken: BotToken, Url: "localhost:8081"}
	tg := telegram.NewTelegram(BotToken, WebhookURL)

	go router.Handlers(wg)
	go tg.Router(wg)

	wg.Wait()

	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}

	// time.Sleep(time.Second * 100)
}
