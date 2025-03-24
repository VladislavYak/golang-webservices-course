package main

// yakovlev useful readings
// https://gowebexamples.com/forms/
// https://gowebexamples.com/sessions/

import (
	"fmt"
	"net/http"
	"sort"
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

	server.ListenAndServe()
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
	if p.order_by != "" {
		res = Sorting(p, res)
	}

	fmt.Fprintf(w, "%+v\n", res)

}

func QueryProcessing(p *params, rows []Row) []Row {

	res := []Row{}
	if p.query != "" {
		for _, row := range rows {
			if (strings.Contains(row.Name, p.query)) || (strings.Contains(row.About, p.query)) {
				res = append(res, row)
			}
		}
	} else {
		res = rows
	}
	return res
}

// {"Id", "Age", "Name"}
func Sorting(p *params, rows []Row) []Row {
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
	return rows
}

func main() {
	xml_path := "/Users/vi/personal_proj/golang_web_services_2024-04-26/03_net1/99_hw/coverage/dataset.xml"

	SearchServer(xml_path)
}
