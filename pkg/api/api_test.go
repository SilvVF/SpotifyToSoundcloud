package api

import (
	"testing"
	"time"
)

func TestAuthUrl(t *testing.T) {
	api := NewSoundCloud()
	if api.AuthUrl() == "" {
		t.Fail()
	}
	t.Log(api.AuthUrl())
}

func TestSoundCloud(t *testing.T) {

	api := NewSoundCloud()

	api.Register()
	Start()

	api.AuthUrl()

	timer := time.NewTimer(time.Second * 10)
	<-timer.C

	if api.token == "" {
		t.Fail()
	}
}
