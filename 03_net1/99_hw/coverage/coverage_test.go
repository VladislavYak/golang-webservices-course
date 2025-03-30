package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// тут писать код тестов

const XML_PATH = "./dataset.xml"

func TestSearchServer(t *testing.T) {
	myData, _ := readXml(XML_PATH)

	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, &myData)
	})))

	sr := SearchRequest{Limit: 5}

	sc := SearchClient{URL: ts.URL, AccessToken: "mytoken"}
	_, err := sc.FindUsers(sr)
	if err != nil {
		fmt.Println("err", err)
	}

}

func TestSearchServerLimitLessZero(t *testing.T) {
	myData, _ := readXml(XML_PATH)

	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, &myData)
	})))

	sr := SearchRequest{Limit: -7}

	sc := SearchClient{URL: ts.URL, AccessToken: "mytoken"}
	_, err := sc.FindUsers(sr)

	if err == nil {
		t.Error("limit cannot be less then one")
	}
}

func TestSearchServerLimitMore25(t *testing.T) {
	myData, _ := readXml(XML_PATH)

	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, &myData)
	})))

	sr := SearchRequest{Limit: 100}

	sc := SearchClient{URL: ts.URL, AccessToken: "mytoken"}
	res, _ := sc.FindUsers(sr)

	if len(res.Users) != 25 {
		t.Error("Cannot be more then 25")
	}

}

func TestSearchServerOffsetLess0(t *testing.T) {
	myData, _ := readXml(XML_PATH)

	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, &myData)
	})))

	sr := SearchRequest{Offset: -20}

	sc := SearchClient{URL: ts.URL, AccessToken: "mytoken"}
	resp, err := sc.FindUsers(sr)

	if err == nil && len(resp.Users) != 0 {
		t.Error("Offset cannot be less then 0")
	}
}

func TestSearchServerAuthError(t *testing.T) {
	myData, _ := readXml(XML_PATH)

	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, &myData)
	})))

	sr := SearchRequest{Offset: 10}

	sc := SearchClient{URL: ts.URL, AccessToken: "invalidToken"}
	resp, err := sc.FindUsers(sr)

	if err == nil && len(resp.Users) != 0 {
		t.Error("Passed invalidToken")
	}
}
