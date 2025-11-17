package soundcloud

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/SilvVF/sptosc/pkg/api/internal/store"
)

func TestSoundCloud(t *testing.T) {

	api := New(t.Context())

	api.Register()
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Print(err)
		}
	}()

	fp := filepath.Clean("C:\\Users\\DS\\dev\\SpotifyToSoundcloud\\apps\\desktop\\temp")
	store.Initialize(filepath.Dir(fp))

	f, _ := os.Create(filepath.Join(fp, "link.txt"))
	f.WriteString(api.AuthUrl())
	f.Close()

	api.CheckAuth()

	c, err := api.AwaitClient()
	if err != nil {
		t.Error(err)
	}

	res, err := c.SearchTracks("mirror")

	log.Printf("%v\n", res)

	t.Fail()
}
