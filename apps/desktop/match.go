package main

import (
	"context"
	"errors"
	"log"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/SilvVF/sptosc/pkg/api/soundcloud"
	"github.com/SilvVF/sptosc/pkg/levenshtein"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	spotifyapi "github.com/zmb3/spotify/v2"
)

const MatchEventName = "match_progress"
const MatchResultPrefix = "match_result_"

type Job struct {
	broadcaster *Broadcaster[map[string]TracksWrapper]
	ctx         context.Context
	cancel      context.CancelFunc
}

func (a *App) ScoredTrackBind() ScoredTrack {
	return ScoredTrack{}
}

func (a *App) TracksWrapperBind() TracksWrapper {
	return TracksWrapper{}
}

func (a *App) CancelMatching(id string) {

	a.jobsLock.Lock()
	defer a.jobsLock.Unlock()

	if job, ok := a.matchJobs[id]; ok {
		job.cancel()
		delete(a.matchJobs, id)
	}
}

func (a *App) GetMatches(id string) (map[string]TracksWrapper, error) {

	p, ok := a.cachedPlaylists[id]
	if !ok {
		log.Println("playlist not found in cache")
		return nil, errors.New("playlist no in cache")
	}

	if cached, ok := a.cachedResults[id]; ok {
		runtime.EventsEmit(
			a.ctx,
			MatchEventName,
			MatchProgress{
				ForId:    id,
				Total:    len(cached),
				Progress: len(cached),
				Status:   "done",
			},
		)

		return cached, nil
	}

	log.Println("finding matches for: ", id)

	a.jobsLock.Lock()
	// check for running job
	if job, ok := a.matchJobs[id]; ok && job.ctx.Err() == nil {
		// result broadcast
		ch := job.broadcaster.Listen()
		a.jobsLock.Unlock()

		res, ok := <-ch
		if !ok || res == nil {
			return nil, errors.New("job failed")
		}
		return res, nil
	}

	// setup new job
	broadcaster := NewBroadcaster[map[string]TracksWrapper]()
	ctx, cancel := context.WithCancel(a.ctx)
	defer cancel()
	job := &Job{
		broadcaster: broadcaster,
		ctx:         ctx,
		cancel:      cancel,
	}
	a.matchJobs[id] = job

	a.jobsLock.Unlock()

	res := findMatches(id, p.Tracks.Tracks, a.soundCloud, job.ctx)
	a.cachedResults[id] = res

	// push results and remove job
	a.jobsLock.Lock()
	defer a.jobsLock.Unlock()

	job.broadcaster.Broadcast(res)
	job.broadcaster.Close()
	delete(a.matchJobs, id)

	return res, nil
}

func createQuery(track spotifyapi.FullTrack) string {
	re := regexp.MustCompile(`$feat.*\..+$`)
	name := re.ReplaceAllString(track.Name, "")
	artists := make([]string, len(track.Artists))
	for i, artist := range track.Artists {
		artists[i] = artist.Name
	}
	names := strings.Join(artists, " ") + " " + strings.TrimSpace(name)
	query := strings.ReplaceAll(names, " &", "")
	return query
}

func findMatches(
	playlistId string,
	tracks []spotifyapi.PlaylistTrack,
	api *soundcloud.SoundCloudApi,
	ctx context.Context,
) map[string]TracksWrapper {

	matched := make(map[string]TracksWrapper)

	sem := make(chan struct{}, 4)
	m := sync.Mutex{}
	wg := sync.WaitGroup{}
	debounce := time.Millisecond * 300

	progress := atomic.Int32{}
	sendProgress := func() {
		runtime.EventsEmit(
			ctx,
			MatchEventName,
			MatchProgress{
				ForId:    playlistId,
				Total:    len(tracks),
				Progress: int(progress.Add(1)),
				Status:   "running",
			},
		)
	}

	for _, track := range tracks {
		wg.Go(func() {

			if ctx.Err() != nil {
				return
			}

			sem <- struct{}{}
			defer func() {
				wait, waitCancel := context.WithTimeout(ctx, debounce)
				defer waitCancel()

				sendProgress()

				<-wait.Done()
				<-sem
			}()

			query := createQuery(track.Track)
			log.Println("searching for: ", query)

			res, err := api.SearchTracks(query)

			if err != nil {
				log.Println("error making req: ", err.Error())
				return
			}

			if res == nil || len(*res) == 0 {
				log.Printf("No results found for query: %s\n", query)
				return
			}

			scores := scoreResults(track.Track, *res)

			sorted := make([]soundcloud.Track, len(*res))
			copy(sorted, *res)

			slices.SortFunc(sorted, func(t1, t2 soundcloud.Track) int {
				return int(scores[t2.ID] - scores[t1.ID])
			})
			mapped := make([]ScoredTrack, len(sorted))
			for i, track := range sorted {
				mapped[i] = ScoredTrack{
					Track: Track{
						ID:    strconv.Itoa(track.ID),
						Urn:   track.Urn,
						Title: track.Title,
						Imgs: []Img{
							{Url: track.ArtworkURL},
						},
					},
					Score: scores[track.ID],
				}
			}

			m.Lock()
			matched[track.Track.ID.String()] = TracksWrapper{
				ForId:  string(track.Track.ID),
				Tracks: mapped,
			}
			runtime.EventsEmit(ctx, MatchResultPrefix+playlistId, matched[track.Track.ID.String()])
			m.Unlock()
		})
	}

	wg.Wait()

	status := "done"
	if ctx.Err() != nil {
		status = "error"
	}
	runtime.EventsEmit(
		ctx,
		MatchEventName,
		MatchProgress{
			ForId:    playlistId,
			Total:    len(tracks),
			Progress: len(tracks),
			Status:   status,
		},
	)

	return matched
}

func scoreResults(track spotifyapi.FullTrack, results []soundcloud.Track) map[int]float64 {
	trackIdToScore := make(map[int]float64)

	targetArtists := make([]string, len(track.Artists))
	for i, artist := range track.Artists {
		targetArtists[i] = artist.Name
	}
	targetArtistString := strings.Join(targetArtists, " ")

	for _, res := range results {
		var totalScore float64
		var scoreCount float64

		artistScore := similarity(res.MetadataArtist, targetArtistString)
		totalScore += artistScore
		scoreCount += 1

		titleScore := similarity(track.Name, res.Title)
		totalScore += titleScore
		scoreCount += 1

		if res.Duration != 0 {
			durationScore := 1 - (math.Abs(float64(res.Duration-int(track.Duration))) * 2.0 / float64(track.Duration))
			totalScore += durationScore * 5
			scoreCount += 1
		}

		if scoreCount > 0 {
			trackIdToScore[res.ID] = totalScore / scoreCount
		}
	}

	return trackIdToScore
}

func similarity(s1, s2 string) float64 {

	r1 := []rune(strings.ToLower(s1))
	r2 := []rune(strings.ToLower(s2))

	distance := levenshtein.DistanceForStrings(r1, r2, levenshtein.DefaultOptions)
	maxLength := max(len(r1), len(r2))

	if maxLength == 0 {
		return 1.0
	}
	return 1 - float64(distance)/float64(maxLength)
}
