package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Tasks(w http.ResponseWriter, r *http.Request) {
	// gives list of tasks

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{tasks: test}`))
}

func New(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	fmt.Println("name", name)

	// yakovlev: name is a needed taskname

	w.WriteHeader(http.StatusOK)
	// w.Write([]byte(``))
}

func Owner(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{owner: test}`))
}

func My(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{my: test}`))
}
