package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/SilvVF/sptosc/pkg/api"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	spotifyApi *api.SpotifyApi
	cursor     int // which our cursor is pointing at
	state      map[int]*PlatformState
}

func main() {

	spotifyApi := api.NewSpotify()
	spotifyApi.Start()

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	p := tea.NewProgram(initialModel(spotifyApi))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type PlatformState struct {
	authed  bool
	waiting bool
	err     error
	url     string
	name    string
}

func initialModel(spotifyApi *api.SpotifyApi) model {

	return model{
		spotifyApi: spotifyApi,
		cursor:     0,
		state: map[int]*PlatformState{
			0: {
				name:    "Spotify",
				url:     spotifyApi.AuthUrl(),
				waiting: true,
			},
			1: {
				name:    "Sound Cloud",
				url:     "",
				waiting: true,
			},
		},
	}
}

type statusMsg struct {
	err error
	id  int
}

func (m model) Init() tea.Cmd {

	return func() tea.Msg {

		_, err := m.spotifyApi.AwaitClient()

		return statusMsg{
			err: err,
			id:  0,
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case statusMsg:
		if msg.err != nil {
			m.state[msg.id].authed = true
			m.state[msg.id].waiting = false
		} else {
			m.state[msg.id].err = msg.err
			m.state[msg.id].authed = false
			m.state[msg.id].waiting = true
		}
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.state)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":

		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "Sign in to required platforms\n\n"

	// Iterate over our choices
	for i, platform := range m.state {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this platfrom authenticated?
		checked := " " // not selected
		url := ""
		if state, ok := m.state[i]; ok {
			if state.authed {
				checked = "x" // selected!
			}
			if state.waiting && state.url != "" {
				url = "\n" + state.url
			}
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s", cursor, checked, platform.name)

		sp := []string{}
		for _, c := range url {
			sp = append(sp, string(c))
		}
		for c := range slices.Chunk(sp, 10) {
			s += strings.Join(c, "") + "\n"
		}
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
