package main

// yakovlev useful readings
// https://gowebexamples.com/forms/
// https://gowebexamples.com/sessions/

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"time"
)

// order_by=-1&order_field=age&limit=1&offset=0&query=on
// тут писать SearchServer
// FindUsers отправляет запрос во внешнюю систему (на самом деле в searchServer, (по сути в Мок))
// func MainPage(w http.ResponseWriter, r *http.Request, data *Rows) {

// 	values := data.List

// 	// yakovlev: try to preinit p, pass it to parseParams (like unmarshal does)
// 	p := parseParams(r)

// 	fmt.Fprintln(w, "p:", p)

// 	fmt.Fprintln(w, "p.order_by:", p.order_by)

// 	res := []Row{}
// 	if p.query != "" {
// 		for _, row := range values {
// 			if (strings.Contains(row.Name(), p.query)) || (strings.Contains(row.About, p.query)) {
// 				res = append(res, row)
// 			}
// 		}
// 	} else {
// 		res = values
// 	}
// 	fmt.Fprintf(w, "%+v\n", res)

// }

// // по сути, это мок внешней апи, которая отдавал бы данные
// func SearchServer(datapath string) {
// 	// yakovlev: add error handling
// 	myData, _ := readXml(datapath)

// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/",
// 		func(w http.ResponseWriter, r *http.Request) {
// 			MainPage(w, r, &myData)
// 		},
// 	)

// 	server := http.Server{
// 		Handler: mux,
// 	}

// 	server.ListenAndServe()
// }

type Rows struct {
	Version string `xml:"version,attr"`
	List    []Row  `xml:"row"`
}

type customTime struct {
	time.Time
}

func (c *customTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const shortForm = "2006-01-02T15:04:05 -03:00" // yyyymmdd date format 2014-05-10T11:36:09 -03:00
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(shortForm, v)
	if err != nil {
		return err
	}

	*c = customTime{parse}
	return nil
}

type FullName struct {
	string
}

func (fn *FullName) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data struct {
		FirstName string `xml:"first_name"`
		LastName  string `xml:"last_name"`
	}

	d.DecodeElement(&data, &start)
	fmt.Println("----")
	fmt.Println("data.FirstName", data.FirstName)
	fmt.Println("----")
	fmt.Println("data.LastName", data.LastName)

	*fn = FullName{data.FirstName + data.LastName}
	return nil
}

type Row struct {
	Id            int        `xml:"id"`
	Guid          string     `xml:"guid"` /* uuid */
	IsActive      bool       `xml:"isActive"`
	Balance       string     `xml:"balance"`
	Picture       string     `xml:"picture"` /* url */
	Age           int        `xml:"age"`
	EyeColor      string     `xml:"eyeColor"`
	Name          FullName   `xml:"first_name"`
	Gender        string     `xml:"gender"` /* just 2 values */
	Company       string     `xml:"company"`
	Email         string     `xml:"email"`
	Phone         string     `xml:"phone"`
	Address       string     `xml:"address"`
	About         string     `xml:"about"`
	Registered    customTime `xml:"registered"` /* datetime */
	FavoriteFruit string     `xml:"favoriteFruit"`
}

func readXml(path string) (Rows, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Cannot read the file", err)
		return Rows{}, err
	}

	rows := new(Rows)
	if err := xml.Unmarshal(data, &rows); err != nil {
		fmt.Println("Cannot unmarshal", err)
		return Rows{}, err
	}

	return *rows, nil

}

// * Данные для работы лежит в файле `dataset.xml`
// * Параметр `query` ищет по полям `Name` и `About`
// * Параметр `order_field` работает по полям `Id`, `Age`, `Name`, если пустой - то возвращаем по `Name`, если что-то другое - SearchServer ругается ошибкой. `Name` - это first_name + last_name из xml.
// * Параметр `order_by` задает направление сортировки (по полю переданному в `order_field`) или ее отсутствие (OrderByAsIs)
// * Параметры `offset` и `limit` позволяют получать отсортированный список юзеров пачками с индекса `offset` не более `limit` штук.
// * Если `query` пустой, то делаем только сортировку, т.е. возвращаем все записи
// * Код нужно писать в файле coverage_test.go и server.go. Файл client.go трогать не надо
// * Как работать с XML смотрите в `3/6_xml/*`
// * Запускать как `go test -cover`
// * Построение покрытия: `go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html`.
// * В XML 2 поля с именем, наше поле Name это first_name + last_name из XML
// * <https://www.golangprograms.com/files-directories-examples.html> - в помощь для работы с файлами
// * проверка ошибок в функциях io.WriteString, io.ReadAll(и аналогичных им, что читают из Reader, пишут во Writer), а также json.Marshal/Unmarshal и xml.Unmarshal должны быть, но их можно не покрывать тестами. но если все же хочется чистых 100%, то вот подсказка: чтоб это было проще покрыть тестами можно сделать отдельную функцию для парсинга входящих параметров и отдельную для сериализации и отправки ответа, тогда довольно легко будет описать отдельный тест на негативные сценарии этих функций. что касается xml, то там можно обойтись без ReadAll если парсить через Decoder

// 1. query
// 2. offset & limit
// 3. order_by & order_field

type params struct {
	query       string
	order_field string
	// yakovlev: whether it string or int
	order_by string
	offset   string
	limit    string
}

func parseParams(r *http.Request) params {
	parsedUrl := r.URL.Query()

	p := params{
		order_by:    parsedUrl.Get("order_by"),
		order_field: parsedUrl.Get("order_field"),
		limit:       parsedUrl.Get("limit"),
		offset:      parsedUrl.Get("offset"),
		query:       parsedUrl.Get("query"),
	}
	return p
}

func main() {
	xml_path := "/Users/vi/personal_proj/golang_web_services_2024-04-26/03_net1/99_hw/coverage/dataset.xml"

	v, _ := readXml(xml_path)

	// fmt.Println("v", v)
	fmt.Printf("%+v\n", v)

	// SearchServer(xml_path)
}
