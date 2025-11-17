package api

import (
	"log"
	"net/http"
	"os"

	"github.com/SilvVF/sptosc/pkg/api/internal/store"
	"github.com/joho/godotenv"
)

type ApiName string

const (
	SOUNDCLOUD ApiName = "soundcloud"
	SPOTIFY    ApiName = "spotify"
)

func init() {
	err := godotenv.Load("C:\\Users\\DS\\dev\\SpotifyToSoundcloud\\pkg\\api\\api.env")
	if err != nil {
		panic(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	store.Initialize(dir)
}

func Start() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Print(err)
		}
	}()
}
