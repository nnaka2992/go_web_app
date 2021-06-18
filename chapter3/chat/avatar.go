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

type AuthAvatar struct{}
var UserAuthAvatar AuthAvatar
func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}
var UserGravater GravatarAvatar
func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct {}
var UserFileSystemAvatar FileSystemAvatar
func (_ FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "/avatars/" + useridStr + ".jpg", nil
		}
	}
	return "", ErrNoAvatarURL
}
