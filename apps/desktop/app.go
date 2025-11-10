package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SilvVF/sptosc/pkg/api"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	spotify    *api.SpotifyApi
	soundCloud *api.SoundCloudApi
	ctx        context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
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

	a.spotify = api.NewSpotify(ctx)
	a.soundCloud = api.NewSoundCloud()

	a.soundCloud.Register()
	a.spotify.Register()

	api.Start()
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {
	go func() {
		for {
			if ctx.Err() != nil {
				return
			}
			_, err := a.soundCloud.AwaitToken()

			runtime.LogDebugf(ctx, "received client err: %e", err)

			data, err := json.Marshal(AuthEvent{
				Err:  err,
				Name: "soundcloud",
				Ok:   err == nil,
			})
			if err != nil {
				runtime.LogError(ctx, err.Error())
			}
			runtime.EventsEmit(ctx, "auth_event", string(data))
		}
	}()
	go func() {
		for {
			if ctx.Err() != nil {
				return
			}
			_, err := a.spotify.AwaitClient()

			runtime.LogDebugf(ctx, "received client err: %e", err)

			data, err := json.Marshal(AuthEvent{
				Err:  err,
				Name: "spotify",
				Ok:   err == nil,
			})
			if err != nil {
				runtime.LogError(ctx, err.Error())
			}
			runtime.EventsEmit(ctx, "auth_event", string(data))
		}
	}()
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

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) SoundCloudAuthUrl() string {
	return a.soundCloud.AuthUrl()
}

func (a *App) SpotifyAuthUrl() string {
	return a.spotify.AuthUrl()
}

func (a *App) SpotifyPlaylists() ([]string, error) {
	runtime.LogDebug(a.ctx, "getting playlists")
	res, err := a.spotify.UserPlaylists(100, 0)
	if err != nil {
		runtime.LogDebug(a.ctx, err.Error())
		return make([]string, 0), err
	}
	runtime.LogDebugf(a.ctx, "%d", res.Total)
	strs := make([]string, len(res.Playlists))
	for i, p := range res.Playlists {
		strs[i] = p.Name
	}
	return strs, nil
}

func (a *App) RefreshAuthState() {
	a.spotify.CheckAuth(a.ctx)
	a.soundCloud.CheckAuth()
}
