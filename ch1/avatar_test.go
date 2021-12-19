package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	cu := new(chatUser)
	url, err := authAvatar.GetAvatarURL(cu)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL if not present")
	}

	testUrl := "http://url-to-avatar"
	cu.User.AvatarURL = testUrl
	url, err = authAvatar.GetAvatarURL(cu)

	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should not return an error if a value is present")
	}

	if url != testUrl {
		t.Error("AuthAvatar.GetAvatarURL should return the correct URL")
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	cu := new(chatUser)
	cu.uniqueID = "b642b4217b34b1e8d3bd915fc65c4452"
	url, err := gravatarAvatar.GetAvatarURL(cu)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error")
	}
	if url != "http://gravatar.com/avatar/b642b4217b34b1e8d3bd915fc65c4452" {
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned %s", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	filename := filepath.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer os.Remove(filename)
	var fileSystemAvatar FileSystemAvatar
	cu := new(chatUser)
	cu.uniqueID = "abc"
	url, err := fileSystemAvatar.GetAvatarURL(cu)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURL should not return an error")

	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURL wrongly returned %s", url)
	}
}
