package api

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

const soundCloudRedirectURI = "http://localhost:8080/callback/soundcloud"
const soundCloudClientId = "9mHXJkFPzsdvV3kuCmYg150meFk1yLaG"
const soundCloudClientSecret = "9fDudeAnwj5z7uTXZCRbLeWW20GeGpQ7"

const codeVerifier = "4d07b770167cd8b6b8c2325c033ebddabe289a5e1209956cc4f4d5de"
const codeChallenge = "gv0322qzDOksupc253mV0Ub3xre7QDB1-UmG-EFzz9c"

type SoundCloudApi struct {
	token *oauth2.Token
	ch    chan *oauth2.Token
	errCh chan error
	state string
}

func NewSoundCloud() *SoundCloudApi {

	s := &SoundCloudApi{
		ch:    make(chan *oauth2.Token),
		errCh: make(chan error),
		state: rand.Text(),
	}

	s.CheckAuth()
	return s
}

func (s *SoundCloudApi) CheckAuth() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		if store, err := readStore(); err == nil {
			if t := store.SoundCloud.Token; t != nil {
				s.ch <- t
				log.Println("refreshed token from store")
			} else {
				log.Println(err)
			}
		} else {
			log.Println(err)
		}
		close(done)
	}()
	return done
}

func (s *SoundCloudApi) Register() {
	http.HandleFunc("/callback/soundcloud", s.completeAuth)
}

func (s *SoundCloudApi) AuthUrl() string {
	return fmt.Sprintf(
		"https://secure.soundcloud.com/authorize?client_id=%s&redirect_uri=%s&response_type=code&code_challenge=%s&code_challenge_method=S256&state=%s",
		soundCloudClientId, soundCloudRedirectURI, codeChallenge, s.state,
	)
}

func (s *SoundCloudApi) AwaitToken() (*oauth2.Token, error) {
	select {
	case token := <-s.ch:
		s.token = token
		log.Println("soundcloud token:", token)
		if store, err := readStore(); err == nil {
			store.SoundCloud.Token = token
			err = writeStore(*store)
			if err != nil {
				log.Println(err)
			}
		} else {
			log.Println(err)
		}
		log.Println("sending token")
		return token, nil
	case err := <-s.errCh:
		return nil, err
	}
}

func (s *SoundCloudApi) completeAuth(w http.ResponseWriter, r *http.Request) {
	// First, we need to get the value of the `code` query param
	//
	code := r.URL.Query().Get("code")

	if code == "" {
		s.errCh <- errors.New("failed to receive code")
		return
	}

	oauthURL := "https://secure.soundcloud.com/oauth/token"

	// Create URL-encoded form data
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("client_id", soundCloudClientId)
	form.Add("client_secret", soundCloudClientSecret)
	form.Add("redirect_uri", soundCloudRedirectURI)
	form.Add("code_verifier", codeVerifier)
	form.Add("code", code)

	res, err := http.Post(oauthURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Println("Error during HTTP POST:", err)
		return
	}
	defer res.Body.Close()

	// Parse the request body into the `OAuthAccessResponse` struct
	var t oauth2.Token
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.ch <- &t
}
