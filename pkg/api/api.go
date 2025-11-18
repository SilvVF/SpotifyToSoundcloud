package api

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/SilvVF/sptosc/pkg/api/internal/store"
	"github.com/joho/godotenv"
)

type ApiName string

const (
	SOUNDCLOUD ApiName = "soundcloud"
	SPOTIFY    ApiName = "spotify"
)

func init() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fpath := filepath.Join(dir, "temp", "api.env")
	f, err := os.Open(fpath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
			_, err := os.Create(fpath)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	defer f.Close()

	err = godotenv.Load(fpath)

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
