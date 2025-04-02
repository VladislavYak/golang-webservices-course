package main

import (
	"net/http"
)

type params struct {
	query       string
	order_field string
	// yakovlev: whether it string or int
	order_by string
	offset   string
	limit    string
}

func parseParams(r *http.Request) *params {
	parsedUrl := r.URL.Query()

	orderField := parsedUrl.Get("order_field")
	if orderField == "" {
		orderField = "Name"
	}

	orderBy := parsedUrl.Get("order_by")
	if orderBy == "" {
		orderBy = "0"
	}
	p := params{
		order_by:    orderBy,
		order_field: orderField,
		limit:       parsedUrl.Get("limit"),
		offset:      parsedUrl.Get("offset"),
		query:       parsedUrl.Get("query"),
	}

	return &p
}
