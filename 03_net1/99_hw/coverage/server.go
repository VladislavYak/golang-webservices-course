package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"time"
)

// order_by=-1&order_field=age&limit=1&offset=0&query=on
// тут писать SearchServer
// FindUsers отправляет запрос во внешнюю систему (на самом деле в searchServer, (по сути в Мок))
func MainPage(w http.ResponseWriter, r *http.Request) {
	parsedUrl := r.URL.Query()
	orderBy := parsedUrl.Get("order_by")
	orderField := parsedUrl.Get("order_field")
	limit := parsedUrl.Get("limit")
	offset := parsedUrl.Get("offset")
	query := parsedUrl.Get("query")

	fmt.Fprintln(w, "hey")
	fmt.Fprintln(w, "orderBy:", orderBy)
	fmt.Fprintln(w, "orderField:", orderField)
	fmt.Fprintln(w, "limit:", limit)
	fmt.Fprintln(w, "offset:", offset)
	fmt.Fprintln(w, "query:", query)

	myFile, _ := readXml("/Users/vi/personal_proj/golang_web_services_2024-04-26/03_net1/99_hw/coverage/dataset.xml")

	v := myFile.List

	out, _ := json.Marshal(v[0])
	fmt.Fprintln(w, "random marshal:", string(out))
	fmt.Println("myFile", myFile)
}

// по сути, это мок внешней апи, которая отдавал бы данные
func SearchServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/",
		MainPage)

	server := http.Server{
		Handler: mux,
	}

	server.ListenAndServe()
}

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

type Row struct {
	Id            int        `xml:"id"`
	Guid          string     `xml:"guid"` /* uuid */
	IsActive      bool       `xml:"isActive"`
	Balance       string     `xml:"balance"`
	Picture       string     `xml:"picture"` /* url */
	Age           int        `xml:"age"`
	EyeColor      string     `xml:"eyeColor"`
	FirstName     string     `xml:"first_name"`
	LastName      string     `xml:"last_name"`
	Gender        string     `xml:"gender"` /* just 2 values */
	Company       string     `xml:"company"`
	Email         string     `xml:"email"`
	Phone         string     `xml:"phone"`
	Address       string     `xml:"address"`
	About         string     `xml:"about"`
	Registered    customTime `xml:"registered"` /* datetime */
	FavoriteFruit string     `xml:"favoriteFruit"`
}

func (r *Row) Name() string {
	return r.FirstName + " " + r.LastName
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

func main() {
	// myFile, _ := readXml("/Users/vi/personal_proj/golang_web_services_2024-04-26/03_net1/99_hw/coverage/dataset.xml")

	// fmt.Println("myFile", myFile)

	// for _, v := range myFile.List {
	// 	fmt.Println("v.Registered", v.Registered)
	// }
	// fmt.Println("res", myFile.List[0].LastName)
	SearchServer()
}
