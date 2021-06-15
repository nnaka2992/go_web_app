package main

import (
	"time"
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

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
