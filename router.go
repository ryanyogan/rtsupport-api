package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Router type holds all the fields (rules) for a router
type Router struct {
	rules map[string]Handler
}

// Handler is an alias for anon function to handle client data
type Handler func(*Client, interface{})

// NewRouter returns a new instance of a Router Type
func NewRouter() *Router {
	return &Router{
		rules: make(map[string]Handler),
	}
}

// Handle is a basic function that maps a rule to a request
func (r *Router) Handle(msgName string, handler Handler) {
	r.rules[msgName] = handler
}

// FindHandler will match routing rules to the appropriate handler
func (r *Router) FindHandler(msgName string) (Handler, bool) {
	handler, found := r.rules[msgName]
	return handler, found
}

// ServeHTTP upgrades the http package with CORS = "*"
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	client := NewClient(socket, r.FindHandler)
	go client.Write()
	client.Read()
}
