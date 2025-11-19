package main

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Settings struct {
	DebounceTime time.Duration
	Concurrency  int
}

func Load() (*Settings, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fpath := filepath.Join(dir, "temp", "settings.txt")
	f, err := os.Open(fpath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
			if err != nil {
				return nil, err
			}
			f, err = os.Create(fpath)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	settings := &Settings{}

	for scanner.Err() == nil {

		line := scanner.Text()

		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := parts[1]

		switch key {
		case "debounce":
			t, err := strconv.Atoi(value)
			if err != nil {
				continue
			}
			settings.DebounceTime = time.Duration(t) * time.Millisecond
		case "concurrency":
			v, err := strconv.Atoi(value)
			if err != nil {
				continue
			}
			settings.Concurrency = v
		}
	}

	return settings, nil
}
