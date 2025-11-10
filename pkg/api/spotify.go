package api

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net/http"

	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/zmb3/spotify/v2"
)

// spotifyRedirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const spotifyRedirectURI = "http://localhost:8080/callback"
const spotifyClientId = "e6f8c62b14de4e75972c8a38e3628eaa"
const spotifyClientSecret = "bbf8669bc61c486ea09c8638adf20661"

type SpotifyApi struct {
	ctx    context.Context
	client *spotify.Client
	auth   *spotifyauth.Authenticator
	ch     chan *spotify.Client
	errCh  chan error
	state  string
}

func NewSpotify(ctx context.Context) *SpotifyApi {
	s := &SpotifyApi{
		ctx: ctx,
		auth: spotifyauth.New(
			spotifyauth.WithRedirectURL(spotifyRedirectURI),
			spotifyauth.WithScopes(
				spotifyauth.ScopePlaylistReadPrivate,
				spotifyauth.ScopePlaylistReadCollaborative,
				spotifyauth.ScopeUserLibraryRead,
				spotifyauth.ScopeUserFollowRead,
			),
			spotifyauth.WithClientID(spotifyClientId),
			spotifyauth.WithClientSecret(spotifyClientSecret),
		),
		ch:    make(chan *spotify.Client),
		state: rand.Text(),
	}

	s.CheckAuth(ctx)

	return s
}

func (s *SpotifyApi) Register() {
	http.HandleFunc("/callback", s.completeAuth)
}

func (s *SpotifyApi) CheckAuth(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		if store, err := readStore(); err == nil {
			if store.Spotify.Token != nil {
				if t, err := s.auth.RefreshToken(ctx, store.Spotify.Token); err == nil {
					client := spotify.New(s.auth.Client(ctx, t))
					s.ch <- client
					log.Println("refreshed token from store")
				} else {
					log.Println(err)
				}
			}
		} else {
			log.Println(err)
		}
		close(done)
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
		if store, err := readStore(); err == nil {
			t, _ := client.Token()
			store.Spotify.Token = t
			log.Println("spotify token: ", t)
			writeStore(*store)
		} else {
			log.Println(err)
		}
		log.Println("sending client")
		return client, nil
	case err := <-s.errCh:
		s.client = nil
		return nil, err
	}
}

func (s *SpotifyApi) UserPlaylists(limit, offset int) (*spotify.SimplePlaylistPage, error) {

	if s.client == nil {
		return nil, errors.New("no authenticed client")
	}

	p, err := s.client.CurrentUsersPlaylists(s.ctx, spotify.Limit(20), spotify.Offset(0))
	if err != nil {
		log.Println(err)
	}
	log.Printf("%v", p)
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
		s.errCh <- fmt.Errorf("state mismatch: %s != %s\n", st, s.state)
		return
	}

	// use the token to get an authenticated client
	client := spotify.New(s.auth.Client(r.Context(), tok))
	s.ch <- client
}
