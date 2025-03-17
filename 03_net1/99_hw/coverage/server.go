package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

// тут писать SearchServer

type Rows struct {
	Version string `xml:"version,attr"`
	List    []Row  `xml:"row"`
}

// ID      int    `xml:"id,attr"`
// Login   string `xml:"login"`
// Name    string `xml:"name"`
// Browser string `xml:"browser"`

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
	myFile, _ := readXml("/Users/vi/personal_proj/golang_web_services_2024-04-26/03_net1/99_hw/coverage/dataset.xml")

	// fmt.Println("myFile", myFile)

	for _, v := range myFile.List {
		fmt.Println("v.Registered", v.Registered)
	}
	fmt.Println("res", myFile.List[0].LastName)
}
