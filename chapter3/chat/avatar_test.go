package main
import (
	"os"
	"path/filepath"
	"fmt"
	"testing"
	"crypto/md5"
	"strings"
	"io"
	"io/ioutil"

	gomniauthtest "github.com/stretchr/gomniauth/test"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(testChatUser)
	if err != ErrNoAvatarURL {
		t.Error("When value does not exist, AuthAvatar.GetAvatarURL should return ErrNoAvatarURL")
	}
	testUrl := "http://url-to-avatar/"
	testUser = &gomniauthtest.TestUser{}
	testChatUser.User = testUser
	testUser.On("Avatar").Return(testUrl, nil)
	url, err = authAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("When value exists, AuthAvatar.GetAvatarURL should not return error")
	} else {
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarURL should return correct URL")
		}
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar

	m := md5.New()
	io.WriteString(m, strings.ToLower("MyEmailAddress@example.com"))
	client := new(client)
	client.userData = map[string]interface{}{
		"userid": fmt.Sprintf("%x", m.Sum(nil)),
	}
	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetaAvatarURL should not return an error")
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("GravatarAvatar.GetAvatarURL return incorrect value %s", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	// Generate test avatar image
	filename := filepath.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 077)
	defer func() { os.Remove(filename) } ()

	var fileSystemAvatar FileSystemAvatar
	client := new(client)
	client.userData = map[string]interface{} {"userid": "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURL should not return an error")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURL returns invalid url %s", url)
	}
}

func TestFileSystemAvatarPNG(t *testing.T) {
	// Generate test png avatar image
	filename := filepath.Join("avatars", "abc.png")
	ioutil.WriteFile(filename, []byte{}, 077)
	defer func() { os.Remove(filename) } ()

	var fileSystemAvatar FileSystemAvatar
	client := new(client)
	client.userData = map[string]interface{} {"userid": "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURL should not return an error")
	}
	if url != "/avatars/abc.png" {
		t.Errorf("FileSystemAvatar.GetAvatarURL returns invalid url %s", url)
	}
}


