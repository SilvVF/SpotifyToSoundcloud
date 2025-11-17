package store

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/oauth2"
)

var fpath = ""

type Store struct {
	Spotify struct {
		Token *oauth2.Token `json:"token"`
	} `json:"spotify"`
	SoundCloud struct {
		Token *oauth2.Token `json:"token"`
	} `json:"soundcloud"`
}

var m = sync.Mutex{}

func Initialize(dir string) {
	fpath = filepath.Join(dir, "temp", "store.json")
	os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
	Create()
}

func Create() error {
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

	var st Store
	if err := json.Unmarshal(b, &st); err != nil {
		log.Println("Error unmarshaling JSON:", err)
		b, err = json.Marshal(Store{}) // Default to an empty store
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

func Write(s Store) error {
	m.Lock()
	defer m.Unlock()

	log.Println("writing store: ", s)
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(fpath, b, os.ModePerm)
}

func Read() (*Store, error) {
	m.Lock()
	defer m.Unlock()
	b, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	s := &Store{}
	err = json.Unmarshal(b, s)
	return s, err
}
