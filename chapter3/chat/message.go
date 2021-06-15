package main

import (
	"time"
)

// message indicate one message
type message struct {
	Name string
	Message string
	When time.Time
	AvatarURL string
}
