package main

import (
	"github.com/gorilla/websocket"
)

// client represents a single chat user, maintaining a websocket connection
type client struct {
	socket *websocket.Conn
	receive chan []byte
	room *room
}

// listens for incoming messages
func (c *client) read() {
	defer c.socket.Close() // ensure socket is closed when done
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

// listens on receive channel, writing each message to the websocket
func (c *client) write() {
	defer c.socket.Close() // ensure socket is closed when done
	for msg := range c.receive {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}