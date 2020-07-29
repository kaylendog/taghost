package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/skyezerfox/taghost/routes"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type _routes []route

// NewRouter creates a router with the predefined API routes.
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range registeredRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var registeredRoutes = _routes{
	route{
		Name:        "GetAssets",
		Method:      "GET",
		Pattern:     "/assets/{id}",
		HandlerFunc: routes.GetAsset,
	},
}
