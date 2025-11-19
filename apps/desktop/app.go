package main

import (
	"context"
	"log"
	"sync"

	"github.com/SilvVF/sptosc/pkg/api"
	"github.com/SilvVF/sptosc/pkg/api/soundcloud"
	"github.com/SilvVF/sptosc/pkg/api/spotify"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	spotifyapi "github.com/zmb3/spotify/v2"
)

type CachedPlaylist struct {
	tracks   []spotifyapi.FullTrack
	playlist *spotifyapi.FullPlaylist
}

// App struct
type App struct {
	spotify         *spotify.SpotifyApi
	soundCloud      *soundcloud.SoundCloudApi
	ctx             context.Context
	jobsLock        sync.Mutex
	matchJobs       map[string]*Job
	cachedPlaylists map[string]CachedPlaylist
	cachedResults   map[string]map[string]TracksWrapper
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		jobsLock:        sync.Mutex{},
		cachedPlaylists: make(map[string]CachedPlaylist),
		cachedResults:   make(map[string]map[string]TracksWrapper),
		matchJobs:       make(map[string]*Job),
	}
}

type AuthEvent struct {
	Name string `json:"name"`
	Err  error  `json:"err"`
	Ok   bool   `json:"ok"`
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx

	a.spotify = spotify.New(ctx)
	a.soundCloud = soundcloud.New(ctx)

	a.spotify.Register()
	a.soundCloud.Register()

	api.Start()
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {
	a.RefreshAuthState()

	handleClientEvents(func() (any, error) {
		return a.soundCloud.AwaitClient()
	}, "soundcloud", a.ctx)
	handleClientEvents(func() (any, error) {
		return a.spotify.AwaitClient()
	}, "spotify", a.ctx)
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

func handleClientEvents(awaitClient func() (any, error), name string, ctx context.Context) {
	go func() {
		for {
			if ctx.Err() != nil {
				return
			}
			_, err := awaitClient()

			runtime.EventsEmit(ctx, "auth_event", AuthEvent{
				Err:  err,
				Name: name,
				Ok:   err == nil,
			})
		}
	}()
}

func (a *App) SoundCloudAuthUrl() string {
	return a.soundCloud.AuthUrl()
}

func (a *App) SpotifyAuthUrl() string {
	return a.spotify.AuthUrl()
}

func (a *App) SpotifyPlaylist(id string) (*PlaylistWithTracks, error) {

	if p, ok := a.cachedPlaylists[id]; ok && p.playlist != nil {
		mapped := toPlaylistWithTracks(p.playlist, p.tracks)
		return &mapped, nil
	}

	p, tracks, err := a.spotify.PlaylistItems(id)
	if err != nil {
		delete(a.cachedPlaylists, id)
	} else {
		a.cachedPlaylists[id] = CachedPlaylist{
			playlist: p,
			tracks:   tracks,
		}
	}
	mapped := toPlaylistWithTracks(p, tracks)

	return &mapped, err
}

func (a *App) SpotifyPlaylists() ([]Playlist, error) {
	res, err := a.spotify.UserPlaylists()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	mapped := make([]Playlist, len(res.Playlists))
	for i, p := range res.Playlists {
		mapped[i] = toPlaylist(p)
	}

	return mapped, nil
}

func (a *App) CreateSoundCloudPlaylist(title string, description string, sharing string, ids []string) (*soundcloud.CreatedPlaylist, error) {
	return a.soundCloud.CreatePlaylist(title, description, sharing, ids)
}

func (a *App) RefreshAuthState() {
	a.spotify.CheckAuth()
	a.soundCloud.CheckAuth()
}

func (a *App) GetStreams(urn string) (*soundcloud.AuthorizedStream, error) {
	return a.soundCloud.GetStreams(urn)
}
