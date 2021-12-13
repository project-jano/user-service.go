package model

import (
	"net/http"
)

type Route struct {
	Name          string
	Method        string
	Pattern       string
	HandlerFunc   http.HandlerFunc
	Authenticated bool
}

type Routes []Route
