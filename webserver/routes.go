package main

import (
	"net/http"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
type routes []route

var webRoutes = routes{
	route{
		"Counter",
		"GET",
		"/ws/v1/add/{jobName}",
		setTaskReady,
	},
	route{
		"Counter",
		"GET",
		"/ws/v1/check/{jobName}",
		checkJobReady,
	},
}
