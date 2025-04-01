package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// тут писать код тестов

// yakovlev: need to add valid tests

const XML_PATH = "/Users/vi/personal_proj/golang_web_services_2024-04-26/03_net1/99_hw/coverage/dataset.xml"

func TestSearchServer(t *testing.T) {
	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, XML_PATH)
	})))

	sr := SearchRequest{Limit: 5}

	sc := SearchClient{URL: ts.URL, AccessToken: "mytoken"}
	_, err := sc.FindUsers(sr)
	if err != nil {
		fmt.Println("err", err)
	}

}

func TestClientLimitLessZero(t *testing.T) {
	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, XML_PATH)
	})))

	sr := SearchRequest{Limit: -7}

	sc := SearchClient{URL: ts.URL, AccessToken: "mytoken"}
	resp, err := sc.FindUsers(sr)

	if err == nil && len(resp.Users) != 0 {
		t.Error("limit cannot be less then one")
	}
}

func TestClientLimitMore25(t *testing.T) {
	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, XML_PATH)
	})))

	sr := SearchRequest{Limit: 100}

	sc := SearchClient{URL: ts.URL, AccessToken: "mytoken"}
	res, _ := sc.FindUsers(sr)

	if len(res.Users) != 25 {
		t.Error("Cannot be more then 25")
	}

}

func TestClientOffsetLess0(t *testing.T) {
	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, XML_PATH)
	})))

	sr := SearchRequest{Offset: -20}

	sc := SearchClient{URL: ts.URL, AccessToken: "mytoken"}
	resp, err := sc.FindUsers(sr)

	if err == nil && len(resp.Users) != 0 {
		t.Error("Offset cannot be less then 0")
	}
}

func TestClientAuthError(t *testing.T) {
	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, XML_PATH)
	})))

	sr := SearchRequest{Offset: 10}

	sc := SearchClient{URL: ts.URL, AccessToken: "invalidToken"}
	resp, err := sc.FindUsers(sr)

	if err == nil && len(resp.Users) != 0 {
		t.Error("Passed invalidToken")
	}
}

func TestClientInvalidXml(t *testing.T) {
	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, "invalid_path___.xml")
	})))

	sr := SearchRequest{}

	sc := SearchClient{URL: ts.URL, AccessToken: "mytoken"}
	_, err := sc.FindUsers(sr)

	if err == nil {
		t.Error("Can get some data from file, but its supposed to be invalid")
	}

}

func TestClientBadRequest(t *testing.T) {
	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, XML_PATH)
	})))

	sr := SearchRequest{OrderField: "invalid_order_field"}

	sc := SearchClient{URL: ts.URL, AccessToken: "mytoken"}
	_, err := sc.FindUsers(sr)

	if err == nil {
		t.Error("Supposed to get 400 error")
	}

}

func TestClientBadRequestInvalidOrderBy(t *testing.T) {
	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, XML_PATH)
	})))

	sr := SearchRequest{OrderBy: 100000}

	sc := SearchClient{URL: ts.URL, AccessToken: "mytoken"}
	_, err := sc.FindUsers(sr)

	if err == nil {
		t.Error("Supposed to get 400 error")
	}

}

func TestClientTimeout(t *testing.T) {

	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, XML_PATH)
	})))
	defer ts.Close()

	// https://medium.com/@jac_ln/how-to-test-real-request-timeout-in-golang-with-httptest-dbc72ea23d1a
	http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = time.Microsecond

	sr := SearchRequest{OrderBy: 1}

	sc := SearchClient{URL: ts.URL, AccessToken: "mytoken"}
	_, err := sc.FindUsers(sr)
	http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = 0
	if err == nil {
		t.Error("Supposed to get 400 error")
	}

}

func TestClientLimit(t *testing.T) {
	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, XML_PATH)
	})))

	sr := SearchRequest{Limit: 15, Offset: 20}

	sc := SearchClient{URL: ts.URL, AccessToken: "mytoken"}
	_, err := sc.FindUsers(sr)
	if err == nil {

	}
}
