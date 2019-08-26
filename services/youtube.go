/*
 * MumbleDJ
 * By Matthieu Grieger
 * services/youtube.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package services

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	neturl "net/url"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"

	"unicode"

	"github.com/antonholmquist/jason"
	duration "github.com/channelmeter/iso8601duration"
	"github.com/spf13/viper"
	"go.reik.pl/mumbledj/bot"
	"go.reik.pl/mumbledj/interfaces"
	"layeh.com/gumble/gumble"
)

var (
	// ErrNotFound is returned when search query return 0 results
	ErrNotFound = errors.New("Found nothing for given query")
	// ErrAPI is returned when API error occurs
	ErrAPI = errors.New("API error occured")
)

// YouTube is a wrapper around the YouTube Data API.
// https://developers.google.com/youtube/v3/docs/
type YouTube struct {
	*GenericService
	searchURL string
}

// NewYouTubeService returns an initialized YouTube service object.
func NewYouTubeService() *YouTube {

	return &YouTube{
		GenericService: &GenericService{
			ReadableName: "YouTube",
			Format:       "bestaudio",
			TrackRegex: []*regexp.Regexp{
				regexp.MustCompile(`https?:\/\/www.youtube.com\/watch\?v=(?P<id>[\w-]+)(?P<timestamp>\&t=\d*m?\d*s?)?`),
				regexp.MustCompile(`https?:\/\/youtube.com\/watch\?v=(?P<id>[\w-]+)(?P<timestamp>\&t=\d*m?\d*s?)?`),
				regexp.MustCompile(`https?:\/\/youtu.be\/(?P<id>[\w-]+)(\?t=(?P<timestamp>\d*m?\d*s?))?`),
				regexp.MustCompile(`https?:\/\/youtube.com\/v\/(?P<id>[\w-]+)(?P<timestamp>\?t=\d*m?\d*s?)?`),
				regexp.MustCompile(`https?:\/\/www.youtube.com\/v\/(?P<id>[\w-]+)(?P<timestamp>\?t=\d*m?\d*s?)?`),
			},
			PlaylistRegex: []*regexp.Regexp{
				regexp.MustCompile(`https?:\/\/www\.youtube\.com\/playlist\?list=(?P<id>[\w-]+)`),
			},
		},
		searchURL: "https://www.googleapis.com/youtube/v3/search?type=video&part=snippet&q=%s&key=%s",
	}
}

// CheckAPIKey performs a test API call with the API key
// provided in the configuration file to determine if the
// service should be enabled.
func (yt *YouTube) CheckAPIKey() error {
	var (
		response *http.Response
		v        *jason.Object
		err      error
	)

	if viper.GetString("api_keys.youtube") == "" {
		return errors.New("No YouTube API key has been provided")
	}
	url := "https://www.googleapis.com/youtube/v3/videos?part=snippet&id=KQY9zrjPBjo&key=%s"
	response, err = http.Get(fmt.Sprintf(url, viper.GetString("api_keys.youtube")))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if v, err = jason.NewObjectFromReader(response.Body); err != nil {
		return err
	}

	if v, err = v.GetObject("error"); err == nil {
		message, _ := v.GetString("message")
		code, _ := v.GetInt64("code")
		errArray, _ := v.GetObjectArray("errors")
		reason, _ := errArray[0].GetString("reason")

		return fmt.Errorf("%d: %s (reason: %s)", code, message, reason)
	}
	return nil
}

// GetTracks uses the passed URL to find and return
// tracks associated with the URL. An error is returned
// if any error occurs during the API call.
func (yt *YouTube) GetTracks(url string, submitter *gumble.User) ([]interfaces.Track, error) {
	var (
		playlistURL      string
		playlistItemsURL string
		id               string
		timestamp        string
		err              error
		resp             *http.Response
		v                *jason.Object
		track            bot.Track
		tracks           []interfaces.Track
	)

	dummyOffset, _ := time.ParseDuration("0s")

	playlistURL = "https://www.googleapis.com/youtube/v3/playlists?part=snippet&id=%s&key=%s"
	playlistItemsURL = "https://www.googleapis.com/youtube/v3/playlistItems?part=snippet,contentDetails&playlistId=%s&maxResults=%d&key=%s&pageToken=%s"
	id, err = yt.getID(url)
	if err != nil {
		return nil, err
	}

	if yt.isPlaylist(url) {
		resp, err = http.Get(fmt.Sprintf(playlistURL, id, viper.GetString("api_keys.youtube")))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		v, err = jason.NewObjectFromReader(resp.Body)
		if err != nil {
			return nil, err
		}

		items, _ := v.GetObjectArray("items")
		item := items[0]

		title, _ := item.GetString("snippet", "title")

		playlist := &bot.Playlist{
			ID:        id,
			Title:     title,
			Submitter: submitter.Name,
			Service:   yt.ReadableName,
		}

		maxItems := math.MaxInt32
		if viper.GetInt("queue.max_tracks_per_playlist") > 0 {
			maxItems = viper.GetInt("queue.max_tracks_per_playlist")
		}

		// YouTube playlist searches return a max of 50 results per page
		maxResults := 50
		if maxResults > maxItems {
			maxResults = maxItems
		}

		pageToken := ""
		for len(tracks) < maxItems {
			curResp, curErr := http.Get(fmt.Sprintf(playlistItemsURL, id, maxResults, viper.GetString("api_keys.youtube"), pageToken))
			if curErr != nil {
				// An error occurred, simply skip this track.
				continue
			}
			defer curResp.Body.Close()

			v, err = jason.NewObjectFromReader(curResp.Body)
			if err != nil {
				// An error occurred, simply skip this track.
				continue
			}

			curTracks, _ := v.GetObjectArray("items")
			for _, track := range curTracks {
				videoID, _ := track.GetString("snippet", "resourceId", "videoId")

				// Unfortunately we have to execute another API call for each video as the YouTube API does not
				// return video durations from the playlistItems endpoint...
				newTrack, _ := yt.getTrack(videoID, submitter, dummyOffset)
				newTrack.Playlist = playlist
				tracks = append(tracks, newTrack)

				if len(tracks) >= maxItems {
					break
				}
			}

			pageToken, _ = v.GetString("nextPageToken")
			if pageToken == "" {
				break
			}
		}

		if len(tracks) == 0 {
			return nil, errors.New("Invalid playlist. No tracks were added")
		}
		return tracks, nil
	}

	// Submitter added a track!

	// Set correct offset of YouTube video
	offset := dummyOffset
	u, _ := neturl.Parse(url)
	q := u.Query()
	timestamp = q.Get("t")

	if timestamp != "" {
		lastChar := len(timestamp) - 1
		if unicode.IsDigit(rune(timestamp[lastChar])) {
			timestamp += "s"
		}
		// time.ParseDuration returns offset 0 if err happen
		offset, _ = time.ParseDuration(timestamp)
	}

	track, err = yt.getTrack(id, submitter, offset)
	if err != nil {
		return nil, err
	}
	tracks = append(tracks, track)
	return tracks, nil
}

func (yt *YouTube) getTrack(id string, submitter *gumble.User, offset time.Duration) (bot.Track, error) {
	var (
		resp *http.Response
		err  error
		v    *jason.Object
	)

	videoURL := "https://www.googleapis.com/youtube/v3/videos?part=snippet,contentDetails&id=%s&key=%s"
	resp, err = http.Get(fmt.Sprintf(videoURL, id, viper.GetString("api_keys.youtube")))
	if err != nil {
		return bot.Track{}, err
	}
	defer resp.Body.Close()

	v, err = jason.NewObjectFromReader(resp.Body)
	if err != nil {
		return bot.Track{}, err
	}
	items, _ := v.GetObjectArray("items")
	if len(items) == 0 {
		return bot.Track{}, errors.New("This YouTube video is private")
	}
	item := items[0]
	title, _ := item.GetString("snippet", "title")
	thumbnail, _ := item.GetString("snippet", "thumbnails", "medium", "url")
	// download and convert thumbnail to base64
	thumbnailBinary, err := http.Get(thumbnail)
	var thumbnailBase64 string
	defer thumbnailBinary.Body.Close()
	if err != nil {
		logrus.WithField("url", thumbnail).Error("Unable to get thumbnail")
	} else {
		respThumbnailBase64Body, err := ioutil.ReadAll(thumbnailBinary.Body)
		if err != nil {
			logrus.WithField("url", thumbnail).Error("Unable to read response body")
		} else {
			thumbnailBase64 = base64.StdEncoding.EncodeToString(respThumbnailBase64Body)
		}
	}
	author, _ := item.GetString("snippet", "channelTitle")
	durationString, _ := item.GetString("contentDetails", "duration")
	durationConverted, _ := duration.FromString(durationString)
	duration := durationConverted.ToDuration()

	return bot.Track{
		ID:              id,
		URL:             "https://youtube.com/watch?v=" + id,
		Title:           title,
		Author:          author,
		Submitter:       submitter.Name,
		Service:         yt.ReadableName,
		Filename:        id + ".track",
		ThumbnailURL:    thumbnail,
		ThumbnailBase64: thumbnailBase64,
		Duration:        duration,
		PlaybackOffset:  offset,
		Playlist:        nil,
	}, nil
}

// SearchTrack uses YouTube Data API v3 to find results for given query and takes first video from
// results list and returns it. If nothing has been found or API call error occurred it should
// return empty bot.Track and error. `bot.Track, nil` otherwise.
func (yt *YouTube) SearchTrack(query string, submitter *gumble.User) (interfaces.Track, error) {
	resp, err := http.Get(fmt.Sprintf(yt.searchURL, neturl.QueryEscape(query), viper.GetString("api_keys.youtube")))
	if err != nil {
		return bot.Track{}, ErrAPI
	}
	defer resp.Body.Close()

	v, err := jason.NewObjectFromReader(resp.Body)
	if err != nil {
		return bot.Track{}, ErrAPI
	}

	items, err := v.GetObjectArray("items")
	if err != nil {
		return bot.Track{}, ErrAPI
	}

	if len(items) == 0 {
		return bot.Track{}, ErrNotFound
	}

	videoID, _ := items[0].GetString("id", "videoId")

	// TODO: replace last arg with dummyOffset
	return yt.getTrack(videoID, submitter, 0)
}
