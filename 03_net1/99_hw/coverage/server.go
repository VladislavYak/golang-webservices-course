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
	// yakovlev: add error handling
	mux := http.NewServeMux()
	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			MainPage(w, r, datapath)
		},
	)

	m := AuthMiddleware(mux)

	server := http.Server{
		Handler: m,
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

	// yakovlev: try to preinit p, pass it to parseParams (like unmarshal does)
	p, err := parseParams(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	res = QueryProcessing(p, res)

	// this is bad code :0
	res, err = Sorting(p, res)
	if err != nil {
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

	res = Offset(p, res)

	res, err = Limit(p, res)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		fmt.Println("err MainPage", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}

func QueryProcessing(p *params, rows []Row) []Row {
	if p.query == "" {
		return rows
	} else {
		res := []Row{}
		for _, row := range rows {
			if (strings.Contains(row.Name, p.query)) || (strings.Contains(row.About, p.query)) {
				res = append(res, row)
			}
		}
		return res
	}
}

// yakovlev: добавить нормальный еррор хенлдинг

// {"Id", "Age", "Name"}
func Sorting(p *params, rows []Row) ([]Row, error) {
	allowed_order_field := []string{"id", "age", "name"}
	allower_order_by := []string{"-1", "1", "0"}

	// fmt.Println("p.order_field", p.order_field)

	if !slices.Contains(allowed_order_field, strings.ToLower(p.order_field)) {
		return nil, ErrWrongOrderField
	}

	if !slices.Contains(allower_order_by, strings.ToLower(p.order_by)) {
		return nil, ErrWrongOrderBy
	}
	//  1 по возрастанию, 0 как встретилось, -1 по убыванию
	// OrderBy int

	if p.order_by == "0" {
		return rows, nil
	} else {
		// жесткий говнокод
		switch strings.ToLower(p.order_field) {
		case "id":
			sort.Slice(rows, func(i, j int) bool {
				if p.order_by == "-1" {
					return rows[i].Id < rows[j].Id
				} else {
					return rows[i].Id > rows[j].Id
				}
			})

		case "age":
			sort.Slice(rows, func(i, j int) bool {
				if p.order_by == "-1" {
					return rows[i].Age < rows[j].Age
				} else {
					return rows[i].Age > rows[j].Age
				}
			})
		case "name":
			sort.Slice(rows, func(i, j int) bool {
				if p.order_by == "-1" {
					return rows[i].Name < rows[j].Name
				} else {
					return rows[i].Name > rows[j].Name
				}
			})

		}

		return rows, nil
	}
}

func Offset(p *params, rows []Row) []Row {
	if p.offset == "" {
		return rows
	} else {
		offset, _ := strconv.Atoi(p.offset)
		// add error handling

		if len(rows)-1 > offset {
			return rows[offset:]
		} else {
			return []Row{}
		}
	}
}

func Limit(p *params, rows []Row) ([]Row, error) {
	if p.limit == "" {
		return rows, nil
	} else {
		limit, _ := strconv.Atoi(p.limit)

		if limit > len(rows) {
			return rows, nil
		} else if limit < 0 {
			return []Row{}, errors.New("invalid param")
		}
		// add error handling
		// validate bounds
		return rows[:limit], nil
	}
}

func main() {
	xml_path := "/Users/vi/personal_proj/golang_web_services_2024-04-26/03_net1/99_hw/coverage/dataset.xml"

	SearchServer(xml_path)
}
