package main

// yakovlev useful readings
// https://gowebexamples.com/forms/
// https://gowebexamples.com/sessions/

import (
	"fmt"
	"net/http"
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

	values := data.List

	// yakovlev: try to preinit p, pass it to parseParams (like unmarshal does)
	p, err := parseParams(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	res := QueryProcessing(&p, &values)

	fmt.Fprintf(w, "%+v\n", res)

}

func QueryProcessing(p *params, rows *[]Row) *[]Row {

	res := new([]Row)
	if p.query != "" {
		for _, row := range *rows {
			r := NewRow(&row)
			if (strings.Contains(r.Name, p.query)) || (strings.Contains(r.About, p.query)) {
				*res = append(*res, *r)
			}
		}
	} else {
		res = rows
	}
	return res
}

func main() {
	xml_path := "/Users/vi/personal_proj/golang_web_services_2024-04-26/03_net1/99_hw/coverage/dataset.xml"

	SearchServer(xml_path)
}
