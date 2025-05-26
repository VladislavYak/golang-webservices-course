package router

import (
	"net/http"
	"sync"

	"github.com/VladislavYak/taskbot/internal/transport/http/handlers"
	"github.com/gorilla/mux"
)

func Handlers(wg *sync.WaitGroup) {
	defer wg.Done()

	r := mux.NewRouter()

	r.HandleFunc("/tasks", handlers.Tasks).Methods("GET")
	r.HandleFunc("/new", handlers.New).Methods("POST")
	r.HandleFunc("/my", handlers.My).Methods("GET")
	r.HandleFunc("/owner", handlers.Owner).Methods("GET")
	r.HandleFunc("/assign", handlers.Assign).Methods("POST")
	r.HandleFunc("/unassign", handlers.Unassign).Methods("POST")
	r.HandleFunc("/resolve", handlers.Resolve).Methods("POST")

	go func() {
		http.ListenAndServe(":8081", r)
	}()
}
