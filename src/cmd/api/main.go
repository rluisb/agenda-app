package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Router interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Route struct {
	Path string
	Handler http.HandlerFunc
}

func NewRoute(path string, handler http.HandlerFunc) Route {
	return Route{path, handler}
}

type CustomRouter struct {
	BasePath string
	Routes map[string]Route
}

func NewCustomRouter(basePath string, routes map[string]Route) *CustomRouter {
	return &CustomRouter{basePath, routes}
}

func (c CustomRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := c.Routes[r.Method]
	log.Printf("ServeHTTP %s %s", r.Method, r.URL)
	route.Handler(w, r)
}



func main() {
	contactsMux := http.NewServeMux()
	contactsRouter := NewCustomRouter("/contacts", map[string]Route{
		http.MethodGet: NewRoute("", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("GET /contacts"))
		})),
		http.MethodPost: NewRoute("", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("POST /contacts"))
		})),
		http.MethodPut:  NewRoute("/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := strings.TrimPrefix(r.URL.Path, "/contacts/")
			if id == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("PUT /contacts/:id - id is required"))
				return
			}
			w.Write([]byte(fmt.Sprintf("PUT /contacts/%s", id)))
		})),
		http.MethodPatch:  NewRoute("/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := strings.TrimPrefix(r.URL.Path, "/contacts/")
			if id == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("PATCH /contacts/:id - id is required"))
				return
			}
			w.Write([]byte(fmt.Sprintf("PATCH /contacts/%s", id)))
		})),
		http.MethodDelete: NewRoute("/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := strings.TrimPrefix(r.URL.Path, "/contacts/")
			if id == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("DELETE /contacts/:id - id is required"))
				return
			}
			w.Write([]byte(fmt.Sprintf("DELETE /contacts/%s", id)))
		})),
	})

	contactsMux.Handle(contactsRouter.BasePath, contactsRouter)
	mainMux := http.NewServeMux()
	mainMux.Handle("/api/", http.StripPrefix("/api", contactsMux))
	http.ListenAndServe(":3000", mainMux)
}