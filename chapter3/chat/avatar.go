package main
import (
	"errors"
)

// ErrNoAvatarURL is an error when avatar instance does not return url
var ErrNoAvatarURL = errors.New("chat: could not correct avater URL")
// Avater is a type representing profile pic
type Avatar interface {
	// GetAvatarURL returns avaters url
	// When any truble occurs return error
	// Especially, when could not get url, return ErrNoAvatarURL
	GetAvatarURL(c *client) (string, error)
}

