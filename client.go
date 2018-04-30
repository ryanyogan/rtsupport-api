package main

import (
	"github.com/gorilla/websocket"
)

// Message type holds action and payload
type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

// FindHandler type takes a message, returns a handler and found
type FindHandler func(string) (Handler, bool)

// Client type hold the send message and socket connection
type Client struct {
	send        chan Message
	socket      *websocket.Conn
	findHandler FindHandler
}

func (client *Client) Read() {
	var message Message
	for {
		if err := client.socket.ReadJSON(&message); err != nil {
			break
		}

		if handler, found := client.findHandler(message.Name); found {
			handler(client, message.Data)
		}
	}
	client.socket.Close()
}

func (client *Client) Write() {
	for msg := range client.send {
		if err := client.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	client.socket.Close()
}

// NewClient returns a new client to send over a socket
func NewClient(socket *websocket.Conn, findHandler FindHandler) *Client {
	return &Client{
		send:        make(chan Message),
		socket:      socket,
		findHandler: findHandler,
	}
}
