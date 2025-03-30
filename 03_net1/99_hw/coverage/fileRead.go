package main

import (
	"encoding/xml"
	"os"
	"time"
)

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
	Name          string
}

func NewRow(r *Row) *Row {
	r.Name = r.FirstName + r.LastName
	return r
}
func readXml(path string) (Rows, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Rows{}, err
	}

	rows := new(Rows)
	if err := xml.Unmarshal(data, &rows); err != nil {
		return Rows{}, err
	}

	newRows := []Row{}
	for _, row := range rows.List {
		r := NewRow(&row)
		newRows = append(newRows, *r)
	}

	rows.List = newRows

	return *rows, nil

}
