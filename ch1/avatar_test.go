package main

import "testing"

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL if not present")
	}

	testUrl := "http://url-to-avatar"
	client.userData = map[string]interface{}{"avatar_url": testUrl}
	url, err = authAvatar.GetAvatarURL(client)

	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should not return an error if a value is present")
	}

	if url != testUrl {
		t.Error("AuthAvatar.GetAvatarURL should return the correct URL")
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := new(client)
	client.userData = map[string]interface{}{"userid": "b642b4217b34b1e8d3bd915fc65c4452"}
	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error")
	}
	if url != "www.gravatar.com/avatar/b642b4217b34b1e8d3bd915fc65c4452" {
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned %s", url)
	}
}
