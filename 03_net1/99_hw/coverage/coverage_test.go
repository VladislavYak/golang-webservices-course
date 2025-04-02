package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// тут писать код тестов

const XML_PATH = "dataset.xml"

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
	resp, err := sc.FindUsers(sr)

	if err == nil && len(resp.Users) != 0 {
		t.Error("Can get some data from file, but its supposed to be invalid")
	}

}

func TestClientBadRequests(t *testing.T) {
	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, XML_PATH)
	})))

	params := []SearchRequest{
		{OrderField: "invalid_order_field"},
		{OrderBy: 100000},
	}
	sc := SearchClient{URL: ts.URL, AccessToken: "mytoken"}

	for _, param := range params {
		_, err := sc.FindUsers(param)

		if err == nil {
			t.Error("Supposed to get 400 error")
		}

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

func TestClientLimitOffset(t *testing.T) {
	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, XML_PATH)
	})))

	sr := SearchRequest{Limit: 15, Offset: 20}
	// i have 35 rows totaly

	sc := SearchClient{URL: ts.URL, AccessToken: "mytoken"}
	resp, _ := sc.FindUsers(sr)

	if len(resp.Users) != 15 {
		t.Error("Supposed to have exactly 15 rows")
	}
}

func TestLimit(t *testing.T) {
	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, XML_PATH)
	})))

	sc := SearchClient{URL: ts.URL, AccessToken: "mytoken"}

	testCases := []struct {
		CaseName string
		Params   *SearchRequest
		Expected *SearchResponse
		Error    bool
	}{
		{
			CaseName: "standart limit",
			Params:   &SearchRequest{Limit: 1},
			Expected: &SearchResponse{Users: []User{{ID: 0}}},
			Error:    false,
		},
		{
			CaseName: "negative limit",
			Params:   &SearchRequest{Limit: -1},
			Expected: nil,
			Error:    true,
		},

		{
			CaseName: "more then 25 limit",
			Params:   &SearchRequest{Limit: 100},
			Expected: &SearchResponse{Users: []User{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}}},
			Error:    false,
		},

		// {
		// 	CaseName: "no passed limit",
		// 	Params:   &SearchRequest{},
		// 	Expected: &SearchResponse{Users: []User{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}, {ID: 0}}},
		// },
	}

	for _, testCase := range testCases {
		resp, err := sc.FindUsers(*testCase.Params)

		if testCase.Error == false {
			if err != nil {
				t.Error("Erorr is not expected:", err, testCase.CaseName)
			}

			if len(resp.Users) != len(testCase.Expected.Users) {
				t.Error("Len of limit is not expected", len(resp.Users), len(testCase.Expected.Users), testCase.CaseName)
			}
		} else {
			if err == nil {
				t.Error("Erorr was expected:", err, testCase.CaseName)
			}

			if resp != nil {
				t.Error("Was expecting nil", err, testCase.CaseName)
			}
		}
	}
}

func TestQuery(t *testing.T) {
	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, XML_PATH)
	})))

	sr := SearchRequest{Query: "Twila", Limit: 10}

	sc := SearchClient{URL: ts.URL, AccessToken: "mytoken"}
	resp, err := sc.FindUsers(sr)

	if err == nil && len(resp.Users) != 1 && resp.Users[0].ID != 33 {
		t.Error("Expected to receive ID 33")
	}
}

func TestSorting(t *testing.T) {
	ts := httptest.NewServer(AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MainPage(w, r, XML_PATH)
	})))

	sc := SearchClient{URL: ts.URL, AccessToken: "mytoken"}

	testCases := []struct {
		CaseName string
		Params   *SearchRequest
		Expected *SearchResponse
		Error    bool
	}{
		{
			CaseName: "order by id desc",
			Params:   &SearchRequest{OrderField: "Id", OrderBy: -1, Limit: 2},
			Expected: &SearchResponse{Users: []User{{ID: 34}, {ID: 33}}},
			Error:    false,
		},
		{
			CaseName: "order by id asc",
			Params:   &SearchRequest{OrderField: "Id", OrderBy: 1, Limit: 2},
			Expected: &SearchResponse{Users: []User{{ID: 0}, {ID: 1}}},
			Error:    false,
		},
		{
			CaseName: "order by id no sortig",
			Params:   &SearchRequest{OrderField: "Id", OrderBy: 0, Limit: 2},
			Expected: &SearchResponse{Users: []User{{ID: 0}, {ID: 1}}},
			Error:    false,
		},
		{
			CaseName: "order by age desc",
			Params:   &SearchRequest{OrderField: "Age", OrderBy: -1, Limit: 2},
			Expected: &SearchResponse{Users: []User{{ID: 32}, {ID: 13}}},
			Error:    false,
		},
		{
			CaseName: "order by age asc",
			Params:   &SearchRequest{OrderField: "Age", OrderBy: 1, Limit: 2},
			Expected: &SearchResponse{Users: []User{{ID: 1}, {ID: 15}}},
			Error:    false,
		},
		{
			CaseName: "order by age no sorting",
			Params:   &SearchRequest{OrderField: "Age", OrderBy: 0, Limit: 2},
			Expected: &SearchResponse{Users: []User{{ID: 0}, {ID: 1}}},
			Error:    false,
		},
		{
			CaseName: "order by Name desc",
			Params:   &SearchRequest{OrderField: "Name", OrderBy: -1, Limit: 2},
			Expected: &SearchResponse{Users: []User{{ID: 13}, {ID: 33}}},
			Error:    false,
		},
		{
			CaseName: "order by Name asc",
			Params:   &SearchRequest{OrderField: "Name", OrderBy: 1, Limit: 2},
			Expected: &SearchResponse{Users: []User{{ID: 15}, {ID: 16}}},
			Error:    false,
		},
		{
			CaseName: "order by Name no sorting",
			Params:   &SearchRequest{OrderField: "Name", OrderBy: 0, Limit: 2},
			Expected: &SearchResponse{Users: []User{{ID: 0}, {ID: 1}}},
			Error:    false,
		},
	}

	for _, testCase := range testCases {
		resp, _ := sc.FindUsers(*testCase.Params)

		if len(resp.Users) != len(testCase.Expected.Users) {
			t.Error("expected len not equals to gotten")
		}

		for i := 0; i < len(resp.Users); i++ {
			if resp.Users[i].ID != testCase.Expected.Users[i].ID {
				t.Error("got unexpected order", resp.Users[i].ID, testCase.Expected.Users[i].ID)
			}
		}
	}
}
