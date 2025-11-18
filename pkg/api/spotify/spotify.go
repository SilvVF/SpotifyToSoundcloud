package spotify

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/SilvVF/sptosc/pkg/api/internal/store"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/zmb3/spotify/v2"
)

// RedirectUri is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
var RedirectUri = os.Getenv("SPOTIFY_REDIRECT")

var ErrClientNotAuthenticated = errors.New("no authenticed client")

type SpotifyApi struct {
	ctx    context.Context
	client *spotify.Client
	auth   *spotifyauth.Authenticator
	ch     chan *spotify.Client
	errCh  chan error
	state  string
}

func New(ctx context.Context) *SpotifyApi {
	s := &SpotifyApi{
		ctx: ctx,
		auth: spotifyauth.New(
			spotifyauth.WithRedirectURL(RedirectUri),
			spotifyauth.WithScopes(
				spotifyauth.ScopePlaylistReadPrivate,
				spotifyauth.ScopePlaylistReadCollaborative,
				spotifyauth.ScopeUserLibraryRead,
				spotifyauth.ScopeUserFollowRead,
			),
		),
		ch:    make(chan *spotify.Client),
		state: rand.Text(),
	}

	return s
}

func (s *SpotifyApi) Register() {
	http.HandleFunc("/callback", s.completeAuth)
}

func (s *SpotifyApi) CheckAuth() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)

		if data, err := store.Read(); err == nil {
			if data.Spotify.Token != nil {
				t, err := s.auth.RefreshToken(s.ctx, data.Spotify.Token)
				if err != nil {
					return
				}
				client := spotify.New(s.auth.Client(s.ctx, t))
				s.ch <- client
			}
		}

	}()

	return done
}

func (s *SpotifyApi) AuthUrl() string {
	return s.auth.AuthURL(s.state)
}

func (s *SpotifyApi) AwaitClient() (*spotify.Client, error) {
	select {
	case client := <-s.ch:
		s.client = client
		if data, err := store.Read(); err == nil {
			t, _ := client.Token()
			data.Spotify.Token = t
			store.Write(*data)
		}
		return client, nil
	case err := <-s.errCh:
		s.client = nil
		return nil, err
	}
}

type PlaylistWithTracks struct {
	Name      string                  `json:"name"`
	Desc      string                  `json:"desc"`
	Imgs      []spotify.Image         `json:"imgs"`
	Followers spotify.Followers       `json:"followers"`
	Tracks    []spotify.PlaylistTrack `json:"tracks"`
}

func (s *SpotifyApi) PlaylistItems(id string) (*spotify.FullPlaylist, []spotify.FullTrack, error) {
	if s.client == nil {
		return nil, nil, ErrClientNotAuthenticated
	}

	p, err := s.client.GetPlaylist(s.ctx, spotify.ID(id), spotify.Limit(50), spotify.Offset(0))
	if err != nil {
		return nil, nil, err
	}

	tracks := []spotify.FullTrack{}
	for _, track := range p.Tracks.Tracks {
		tracks = append(tracks, track.Track)
	}

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	for offset := len(p.Tracks.Tracks); offset+1 < int(p.Tracks.Total); offset += 50 {
		wg.Go(func() {
			next, err := s.client.GetPlaylistItems(
				s.ctx,
				spotify.ID(id),
				spotify.Limit(50),
				spotify.Offset(offset),
			)

			if err != nil {
				return
			}

			log.Printf("%v\n", next)

			mutex.Lock()
			for _, item := range next.Items {
				tracks = append(tracks, *item.Track.Track)
			}
			mutex.Unlock()
		})
	}

	wg.Wait()

	return p, tracks, nil
}

func (s *SpotifyApi) UserPlaylists() (*spotify.SimplePlaylistPage, error) {

	if s.client == nil {
		return nil, ErrClientNotAuthenticated
	}

	p, err := s.client.CurrentUsersPlaylists(s.ctx, spotify.Limit(20), spotify.Offset(0))
	if err != nil {
		return nil, err
	}

	for len(p.Playlists) < int(p.Total) && p.Next != "" {
		next, err := s.client.CurrentUsersPlaylists(
			s.ctx,
			spotify.Limit(20),
			spotify.Offset(len(p.Playlists)),
		)
		if err != nil {
			break
		}

		p.Next = next.Next
		p.Limit = next.Limit
		p.Offset = next.Offset
		p.Total = next.Total
		p.Playlists = append(p.Playlists, next.Playlists...)
	}
	return p, err
}

func (s *SpotifyApi) completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := s.auth.Token(r.Context(), s.state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		s.errCh <- err
		return
	}
	if st := r.FormValue("state"); st != s.state {
		http.NotFound(w, r)
		s.errCh <- fmt.Errorf("state mismatch: %s != %s", st, s.state)
		return
	}

	w.Write([]byte("Received client token go back to ap to continue"))
	// use the token to get an authenticated client
	client := spotify.New(s.auth.Client(r.Context(), tok))
	s.ch <- client
}
