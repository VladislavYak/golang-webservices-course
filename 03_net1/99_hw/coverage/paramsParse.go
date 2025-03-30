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

	p := params{
		order_by:    parsedUrl.Get("order_by"),
		order_field: orderField,
		limit:       parsedUrl.Get("limit"),
		offset:      parsedUrl.Get("offset"),
		query:       parsedUrl.Get("query"),
	}
	return &p
}
