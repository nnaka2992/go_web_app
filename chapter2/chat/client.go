package main

import (
	"github.com/gorilla/websocket"
)

// Client indicate a user current in chat
type client struct {
	// socket is a websocket for this client
	socket *websocket.Conn
	// send is a sending message
	send chan *message
	// room is a chat room where user in
	room *room
	// userData store a data about user
	userData map[string]interface{}
}
/*
func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil{
			break
		}
	}
	c.socket.Close()
}
*/
