package main

// сюда писать код

import (
	"context"
	"net/http"

	handlers "github.com/VladislavYak/taskbot/handlers"
	"github.com/gorilla/mux"
)

var (
	// @BotFather в телеграме даст вам это
	BotToken = "XXX"

	// урл выдаст вам игрок или хероку
	WebhookURL = "https://525f2cb5.ngrok.io"
)

func startTaskBot(ctx context.Context) error {
	r := mux.NewRouter()

	r.HandleFunc("/tasks", handlers.Tasks).Methods("GET")
	r.HandleFunc("/new/{name}", handlers.New).Methods("POST")
	r.HandleFunc("/my", handlers.My).Methods("GET")
	r.HandleFunc("/owner", handlers.Owner).Methods("GET")

	http.ListenAndServe(":8080", r)

	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
