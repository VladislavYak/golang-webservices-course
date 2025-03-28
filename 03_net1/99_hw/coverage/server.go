package main

// yakovlev useful readings
// https://gowebexamples.com/forms/
// https://gowebexamples.com/sessions/

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// по сути, это мок внешней апи, которая отдавал бы данные
func SearchServer(datapath string) {
	// yakovlev: add error handling
	myData, _ := readXml(datapath)
	mux := http.NewServeMux()
	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			MainPage(w, r, &myData)
		},
	)

	server := http.Server{
		Handler: mux,
	}
	err := server.ListenAndServe()
	fmt.Printf("%v", err)
}

// order_by=-1&order_field=age&limit=1&offset=0&query=on
// тут писать SearchServer
// FindUsers отправляет запрос во внешнюю систему (на самом деле в searchServer, (по сути в Мок))
func MainPage(w http.ResponseWriter, r *http.Request, data *Rows) {
	res := data.List

	// yakovlev: try to preinit p, pass it to parseParams (like unmarshal does)
	p, err := parseParams(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	res = QueryProcessing(p, res)

	// this is bad code :0
	res, err = Sorting(p, res)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	res = Offset(p, res)

	res, err = Limit(p, res)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	fmt.Fprintf(w, "%+v\n", res)

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

// {"Id", "Age", "Name"}
func Sorting(p *params, rows []Row) ([]Row, error) {
	allowed := []string{"id", "age", "name"}

	if !slices.Contains(allowed, strings.ToLower(p.order_field)) {
		fmt.Println("i was here")
		return []Row{}, errors.New("invalid param")
	}
	if p.order_by == "" {
		return rows, nil
	} else {
		// жесткий говнокод
		switch p.order_field {
		case "Id":
			sort.Slice(rows, func(i, j int) bool {
				if p.order_by == "-1" {
					return rows[i].Id < rows[j].Id
				} else {
					return rows[i].Id > rows[j].Id
				}
			})

		case "Age":
			sort.Slice(rows, func(i, j int) bool {
				if p.order_by == "-1" {
					return rows[i].Age < rows[j].Age
				} else {
					return rows[i].Age > rows[j].Age
				}
			})
		case "Name":
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

// func main() {
// 	xml_path := "/Users/vi/personal_proj/golang_web_services_2024-04-26/03_net1/99_hw/coverage/dataset.xml"

// 	SearchServer(xml_path)
// }
