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
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("When value does not exist, AuthAvatar.GetAvatarURL returns ErrNoAvatarURL")
	}
	// set value
	testUrl := "http://url-to-avatar"
	client.userData = map[string]interface{}{"avatar_url": testUrl}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("when value is set, AuthAvatar.GetAvatarURL should not return an error")
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
		t.Errorf("FileSystemAvatar.GetAvatarURL returns invalid url &s", url)
	}
}

