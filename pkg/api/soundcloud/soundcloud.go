package soundcloud

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/SilvVF/sptosc/pkg/api/internal/store"
	"golang.org/x/oauth2"
)

var redirectUri = os.Getenv("SOUNDCLOUD_REDIRECT")

var clientId = os.Getenv("SOUNDCLOUD_ID")
var clientSecret = os.Getenv("SOUNDCLOUD_SECRET")

var codeVerifier = os.Getenv("SOUNDCLOUD_VERIFIER")
var codeChallenge = os.Getenv("SOUNDCLOUD_CHALLENGE")

var ErrClientNotAuthenticated = errors.New("no authenticed client")

type SoundCloudApi struct {
	ctx    context.Context
	client *SoundCloudClient
	ch     chan *SoundCloudClient
	errCh  chan error
	state  string
	config oauth2.Config
}

type SoundCloudClient struct {
	http    *http.Client
	baseUrl string
	ctx     context.Context
	cancel  context.CancelFunc
}

func (s *SoundCloudApi) newClient(token *oauth2.Token, ctx context.Context) *SoundCloudClient {

	context, cancel := context.WithCancel(ctx)
	client := s.config.Client(context, token)

	return &SoundCloudClient{
		http:    client,
		baseUrl: "https://api.soundcloud.com/",
		ctx:     context,
		cancel:  cancel,
	}
}

func (s *SoundCloudApi) SearchTracks(query string) (*TracksPage, error) {

	if s.client == nil {
		return nil, ErrClientNotAuthenticated
	}

	return s.client.SearchTracks(query)
}

func New(ctx context.Context) *SoundCloudApi {
	s := &SoundCloudApi{
		ctx:   ctx,
		ch:    make(chan *SoundCloudClient),
		errCh: make(chan error),
		state: rand.Text(),
		config: oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Scopes:       []string{""},
			RedirectURL:  "http://localhost:8080/callback/soundcloud",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://secure.soundcloud.com/authorize" + "/oauth/authorize",
				TokenURL: "https://secure.soundcloud.com" + "/oauth/token",
			},
		},
	}
	return s
}

func (s *SoundCloudApi) CheckAuth() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)

		if store, err := store.Read(); err == nil {
			if t := store.SoundCloud.Token; t != nil {
				newToken, err := s.config.TokenSource(s.ctx, t).Token()
				if err != nil {
					return
				}
				s.ch <- s.newClient(newToken, s.ctx)
			}
		}
	}()
	return done
}

func (s *SoundCloudApi) Register() {
	http.HandleFunc("/callback/soundcloud", s.completeAuth)
}

func (s *SoundCloudApi) CreatePlaylist(title string, description string, sharing string, ids []string) (*CreatedPlaylist, error) {
	if s.client == nil {
		return nil, ErrClientNotAuthenticated
	}

	return s.client.CreatePlaylist(title, description, sharing, ids)
}

func (s *SoundCloudApi) AuthUrl() string {
	return fmt.Sprintf(
		"https://secure.soundcloud.com/authorize?client_id=%s&redirect_uri=%s&response_type=code&code_challenge=%s&code_challenge_method=S256&state=%s",
		clientId, redirectUri, codeChallenge, s.state,
	)
}

func (s *SoundCloudApi) AwaitClient() (*SoundCloudClient, error) {
	select {
	case client := <-s.ch:
		if s.client != nil {
			s.client.cancel()
		}
		s.client = client
		if data, err := store.Read(); err == nil {
			t, _ := client.Token()
			data.SoundCloud.Token = t
			store.Write(*data)
		}
		return client, nil
	case err := <-s.errCh:
		return nil, err
	}
}

// Token gets the client's current token.
func (c *SoundCloudClient) Token() (*oauth2.Token, error) {
	transport, ok := c.http.Transport.(*oauth2.Transport)
	if !ok {
		return nil, errors.New("souncloud: client not backed by oauth2 transport")
	}
	t, err := transport.Source.Token()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *SoundCloudApi) completeAuth(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	state := r.Form.Get("state")
	if state != s.state {
		http.Error(w, "State invalid", http.StatusBadRequest)
		return
	}
	code := r.Form.Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}
	token, err := s.config.Exchange(s.ctx, code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Received client token go back to ap to continue"))
	s.ch <- s.newClient(token, s.ctx)
}
