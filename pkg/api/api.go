package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/oauth2"
)

type ApiName string

const (
	SOUNDCLOUD ApiName = "soundcloud"
	SPOTIFY    ApiName = "spotify"
)

var fpath = ""

type store struct {
	Spotify struct {
		Token *oauth2.Token `json:"token"`
	} `json:"spotify"`
	SoundCloud struct {
		Token *oauth2.Token `json:"token"`
	} `json:"soundcloud"`
}

var m = sync.Mutex{}

func init() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fpath = filepath.Join(dir, "temp", "store.json")
	os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
	createStore()
}

func createStore() error {
	m.Lock()
	defer m.Unlock()
	f, err := os.Open(fpath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Create the file if it doesn't exist
			_, err := os.Create(fpath)
			return err // No need for f.Close() here since we are creating a new file.
		}
		return err // Return if there is another error opening the file
	}
	defer f.Close() // Ensure the file gets closed after reading

	b, err := io.ReadAll(f)
	if err != nil {
		log.Println("Error reading file:", err)
		b = []byte("{}") // Set to an empty JSON object as default.
	}

	var st store
	if err := json.Unmarshal(b, &st); err != nil {
		log.Println("Error unmarshaling JSON:", err)
		b, err = json.Marshal(store{}) // Default to an empty store
		if err != nil {
			return err // Return the error instead of writing on failure
		}
		// Write default store to the file
		if err := os.WriteFile(fpath, b, os.ModePerm); err != nil {
			return err // Return the error
		}
	}
	return nil
}

func writeStore(s store) error {
	m.Lock()
	defer m.Unlock()

	log.Println("writing store: ", s)
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(fpath, b, os.ModePerm)
}

func readStore() (*store, error) {
	m.Lock()
	defer m.Unlock()
	b, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	s := &store{}
	err = json.Unmarshal(b, s)
	return s, err
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
