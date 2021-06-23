package main
import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

type TryAvatars []Avatar

// ErrNoAvatarURL is an error when avatar instance does not return url
var ErrNoAvatarURL = errors.New("chat: could not correct avater URL")
// Avater is a type representing profile pic
type Avatar interface {
	// GetAvatarURL returns avaters url
	// When any truble occurs return error
	// Especially, when could not get url, return ErrNoAvatarURL
	GetAvatarURL(ChatUser) (string, error)
}

func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

type AuthAvatar struct{}
var UserAuthAvatar AuthAvatar
func (_ AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}
var UserGravater GravatarAvatar
func (_ GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

type FileSystemAvatar struct {}
var UserFileSystemAvatar FileSystemAvatar
func (_ FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("./avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := filepath.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}
