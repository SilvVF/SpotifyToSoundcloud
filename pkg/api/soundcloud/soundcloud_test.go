package soundcloud

import (
	"log"
	"path/filepath"
	"testing"

	"github.com/SilvVF/sptosc/pkg/api/internal/store"
)

func TestSoundCloud(t *testing.T) {

	api := New(t.Context())
	fp := filepath.Clean("C:\\Users\\DS\\dev\\SpotifyToSoundcloud\\apps\\desktop\\build\\bin\\temp")
	store.Initialize(filepath.Dir(fp))

	api.CheckAuth()

	c, err := api.AwaitClient()
	if err != nil {
		t.Error(err)
	}

	res, err := c.GetStreamUrls("https://api.soundcloud.com/tracks/soundcloud:tracks:1835023461/streams")
	if err != nil {
		t.Error("failed to get urls", err)
		t.FailNow()
	}

	log.Println(res)

	t.FailNow()
}
