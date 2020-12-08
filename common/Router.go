package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routesOpenPrefix {
		router.Methods(route.Method).
			PathPrefix(route.Pattern).
			Name(route.Name).
			Handler(loadTracing(route.HandlerFunc, route.Name))
	}
	// Add routes needing Authentication
	for _, route := range routesProtected {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(loadTracing(route.HandlerFunc, route.Name))
	}
	return router
}

func loadTracing(next http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(rw, req.WithContext(context.Background()))
	})
}