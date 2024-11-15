package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// each room has a set of active clients and a channel for broadcasting messages
type room struct {
	clients map[*client]bool
	join chan *client
	leave chan *client
	forward chan []byte
}

// initializes a new room with the required channels and client map
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:	 make(chan *client),
		leave:	 make(chan *client),
		clients: make(map[*client]bool),
	}
}

// listen for join, leave, and message events
func (r *room) run() {
	for {
		select {
		// client has joined, add them to the map
		case client := <-r.join:
			r.clients[client] = true
		// client has left, remove them from the map
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)
		// broadcast message to the channel
		case msg := <- r.forward:
			for client := range r.clients {
				client.receive <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// configures websocket connection parameters
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: messageBufferSize}

// upgrades HTTP requests to websocket connections and registers new clients to the room
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// upgrade the HTTP connection to a websocket
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	// create a new client
	client := &client{
		socket: socket,
		receive: make(chan []byte, messageBufferSize),
		room: r,
	}

	// register the client
	r.join <- client
	
	// unregister the client when they leave
	defer func() { r.leave <- client} ()

	// start the client's write loop in a separate goroutine
	go client.write()

	// begin the read loop
	client.read()
}