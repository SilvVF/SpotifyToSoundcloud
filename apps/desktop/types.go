package main

import "github.com/zmb3/spotify/v2"

type TracksWrapper struct {
	ForId  string        `json:"for_id"`
	Tracks []ScoredTrack `json:"tracks"`
}

type ScoredTrack struct {
	Track Track   `json:"track"`
	Score float64 `json:"score"`
}

type MatchProgress struct {
	ForId    string `json:"forId"`
	Total    int    `json:"total"`
	Progress int    `json:"progress"`
	Status   string `json:"status"`
}

type Img struct {
	H   int    `json:"h"`
	W   int    `json:"w"`
	Url string `json:"url"`
}

type Track struct {
	ID    string `json:"id"`
	Urn   string `json:"urn"`
	Title string `json:"title"`
	Imgs  []Img  `json:"imgs"`
}

type Playlist struct {
	ID          string `json:"id"`
	Urn         string `json:"urn"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Imgs        []Img  `json:"imgs"`
}

type PlaylistWithTracks struct {
	Playlist Playlist `json:"playlist"`
	Tracks   []Track  `json:"tracks"`
}

func toPlaylist(p spotify.SimplePlaylist) Playlist {
	return Playlist{
		ID:          p.ID.String(),
		Urn:         string(p.URI),
		Title:       p.Name,
		Description: p.Description,
		Imgs:        convertImgs(p.Images),
	}
}

func convertImgs(imgs []spotify.Image) []Img {
	arr := make([]Img, len(imgs))
	for i, img := range imgs {
		arr[i] = Img{
			H:   int(img.Height),
			W:   int(img.Width),
			Url: img.URL,
		}
	}
	return arr
}

func toPlaylistWithTracks(sp *spotify.FullPlaylist) PlaylistWithTracks {
	if sp == nil {
		return PlaylistWithTracks{}
	}

	convertTracks := func(tracks []spotify.PlaylistTrack) []Track {
		arr := make([]Track, len(tracks))
		for i, track := range tracks {
			arr[i] = Track{
				ID:    track.Track.ID.String(),
				Urn:   string(track.Track.URI),
				Title: track.Track.Name,
				Imgs:  convertImgs(track.Track.Album.Images),
			}
		}
		return arr
	}

	return PlaylistWithTracks{
		Playlist: Playlist{
			ID:          sp.ID.String(),
			Urn:         string(sp.URI),
			Title:       sp.Name,
			Description: sp.Description,
			Imgs:        convertImgs(sp.Images),
		},
		Tracks: convertTracks(sp.Tracks.Tracks),
	}
}
