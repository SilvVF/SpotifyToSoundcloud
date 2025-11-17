package soundcloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
)

func (c *SoundCloudClient) SearchTracks(q string) (*TracksPage, error) {
	res, err := c.http.Get(c.baseUrl + "tracks?q=" + q + "&limit=5")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Println(res.Status)
		b, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(string(b))
		}
	}

	var tp TracksPage

	err = json.NewDecoder(res.Body).Decode(&tp)
	if err != nil {
		return nil, err
	}

	return &tp, nil
}

type playlistPost struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Sharing     string `json:"sharing"`
	Tracks      []urn  `json:"tracks"`
}

type urn struct {
	Urn string `json:"urn"`
}

func (c *SoundCloudClient) CreatePlaylist(title, description, sharing string, ids []string) (*CreatedPlaylist, error) {
	tracks := make([]urn, len(ids))
	for i, id := range ids {
		tracks[i] = urn{id}
	}

	body, err := json.Marshal(playlistPost{
		Title:       title,
		Description: description,
		Sharing:     sharing,
		Tracks:      tracks,
	})

	if err != nil {
		return nil, err
	}

	res, err := c.http.Post(c.baseUrl+"playlists", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 201 {
		return nil, errors.New("failed to create playlist")
	}

	var cp CreatedPlaylist

	err = json.NewDecoder(res.Body).Decode(&cp)

	return &cp, nil
}
