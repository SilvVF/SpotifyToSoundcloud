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
			_, err := os.Create(fpath)
			return err
		}
		return err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		log.Println("Error reading file:", err)
		b = []byte("{}")
	}

	var st Store
	if err := json.Unmarshal(b, &st); err != nil {
		log.Println("Error unmarshaling JSON:", err)
		b, err = json.Marshal(Store{})
		if err != nil {
			return err
		}
		if err := os.WriteFile(fpath, b, os.ModePerm); err != nil {
			return err
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
