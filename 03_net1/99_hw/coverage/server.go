package main

// yakovlev useful readings
// https://gowebexamples.com/forms/
// https://gowebexamples.com/sessions/

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
)

var ErrWrongOrderField = errors.New("found wrong order field")
var ErrWrongOrderBy = errors.New("found wrong order by")

// yakovlev: ошибки возвращать в джейсонах

// по сути, это мок внешней апи, которая отдавал бы данные
func SearchServer(datapath string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			MainPage(w, r, datapath)
		},
	)

	// yakovlev: temprorary
	// m := AuthMiddleware(mux)

	server := http.Server{
		// Handler: m,
		Handler: mux,
	}

	err := server.ListenAndServe()
	fmt.Printf("%v", err)
}

func AuthMiddleware(h http.Handler) http.Handler {
	VALID_TOKEN := "mytoken"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("AccessToken")

		if token != VALID_TOKEN {
			http.Error(w, "StatusUnauthorized", http.StatusUnauthorized)
		} else {
			h.ServeHTTP(w, r)
		}

	})
}

// order_by=-1&order_field=age&limit=1&offset=0&query=on
// тут писать SearchServer
// FindUsers отправляет запрос во внешнюю систему (на самом деле в searchServer, (по сути в Мок))
func MainPage(w http.ResponseWriter, r *http.Request, path string) {
	data, err := readXml(path)
	if err != nil {
		http.Error(w, "Bad request", http.StatusInternalServerError)
	}

	res := data.List

	p := parseParams(r)

	QueryProcessing(p, &res)

	if err := Sorting(p, &res); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if errors.Is(err, ErrWrongOrderBy) {
			io.WriteString(w, `{"Error": "OrderBy invalid"}`)
			return
		} else if errors.Is(err, ErrWrongOrderField) {
			io.WriteString(w, `{"Error": "OrderField invalid"}`)
			return
		} else {
			io.WriteString(w, `{"Error": "got unknown error"}`)
		}
	}

	if err := Offset(p, &res); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := Limit(p, &res); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		// yakovlev: add valid error handling here
		fmt.Println("err MainPage", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}

func QueryProcessing(p *params, rows *[]Row) {
	s := *rows
	if p.query == "" {
		*rows = s
	} else {
		s := *rows

		tmp := []Row{}
		for i := 0; i < len(s); i++ {
			if (strings.Contains(s[i].Name, p.query)) || (strings.Contains(s[i].About, p.query)) {
				tmp = append(tmp, s[i])
			}
		}
		*rows = tmp

	}
}

// yakovlev: добавить нормальный еррор хенлдинг

// {"Id", "Age", "Name"}
func Sorting(p *params, rows *[]Row) error {
	allowed_order_field := []string{"id", "age", "name"}
	allower_order_by := []string{"-1", "1", "0"}

	s := *rows

	if !slices.Contains(allowed_order_field, strings.ToLower(p.order_field)) {
		return ErrWrongOrderField
	}

	if !slices.Contains(allower_order_by, strings.ToLower(p.order_by)) {
		return ErrWrongOrderBy
	}

	if p.order_by == "0" {
		return nil
	} else {
		// жесткий говнокод
		switch strings.ToLower(p.order_field) {
		case "id":
			sort.Slice(s, func(i, j int) bool {
				if p.order_by == "-1" {
					return s[i].Id < s[j].Id
				} else {
					return s[i].Id > s[j].Id
				}
			})

		case "age":
			sort.Slice(s, func(i, j int) bool {
				if p.order_by == "-1" {
					return s[i].Age < s[j].Age
				} else {
					return s[i].Age > s[j].Age
				}
			})
		case "name":
			sort.Slice(s, func(i, j int) bool {
				if p.order_by == "-1" {
					return s[i].Name < s[j].Name
				} else {
					return s[i].Name > s[j].Name
				}
			})

		}

		*rows = s

		return nil
	}
}

func Offset(p *params, rows *[]Row) error {
	s := *rows

	if p.offset == "" {
		*rows = s
		return nil
	}

	offset, _ := strconv.Atoi(p.offset)
	// yakovlev :add error handling here

	if len(s)-1 > offset {
		*rows = s[offset:]
	} else {
		*rows = []Row{}
	}
	return nil

}

func Limit(p *params, rows *[]Row) error {
	s := *rows

	if p.limit == "" {
		*rows = s
		return nil

	}

	limit, _ := strconv.Atoi(p.limit)
	// yakovlev: add error handling herre

	if limit > len(s) {
		*rows = s
		return nil
	} else if limit < 0 {
		return errors.New("invalid param")
	} else {
		// yakovlev: validate bounds
		*rows = s[:limit]
		return nil
	}

}

func main() {
	xml_path := "/Users/vi/personal_proj/golang_web_services_2024-04-26/03_net1/99_hw/coverage/dataset.xml"

	SearchServer(xml_path)

	// 	r := []Row{Row{Id: 1, Name: "testeest"}, Row{Id: 2, Name: "vlad"}, Row{Id: 3, Name: "egor"}, Row{Id: 4, Name: "somte"}}
	// 	fmt.Println("r", r)

	// 	v := TestFunc(&r)
	// 	fmt.Println("v", v)

	// 	fmt.Println("r", r)
}
