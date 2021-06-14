package main
import (
	"testing"
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
	client.userData = map[string]interface{}{"avata_url": testUrl}
	if err != nil {
		t.Error("when value is set, AuthAvatar.GetAvatarURL should not return error")
	} else {
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarURL should return correct URL")
		}
	}
}
