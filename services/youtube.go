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

// YouTube is a wrapper around the YouTube Data API.
// https://developers.google.com/youtube/v3/docs/
type YouTube struct {
	*GenericService
	dummyOffset         time.Duration
	playlistAPIURL      string
	playlistItemsAPIURL string
}

type ytPlaylist struct {
	id         string // id of url
	pageToken  string
	maxResults int // maxResults from YouTube playlist, max per request is 50
	// How many items we want from playlist in total i.e. 150. It is read from key queue.max_tracks_per_playlist
	maxItems int
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
		dummyOffset:         0,
		playlistAPIURL:      "https://www.googleapis.com/youtube/v3/playlists?part=snippet&id=%s&key=%s",
		playlistItemsAPIURL: "https://www.googleapis.com/youtube/v3/playlistItems?part=snippet,contentDetails&playlistId=%s&maxResults=%d&key=%s&pageToken=%s",
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
		id        string
		timestamp string
		err       error
		track     bot.Track
		tracks    []interfaces.Track
	)

	dummyOffset, _ := time.ParseDuration("0s")

	id, err = yt.getID(url)
	if err != nil {
		return nil, err
	}

	if yt.isPlaylist(url) {
		return yt.getPlaylist(id, submitter)
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

func (yt *YouTube) getPlaylist(id string, submitter *gumble.User) ([]interfaces.Track, error) {
	var (
		resp   *http.Response
		v      *jason.Object
		err    error
		tracks []interfaces.Track
	)
	resp, err = http.Get(fmt.Sprintf(yt.playlistAPIURL, id, viper.GetString("api_keys.youtube")))
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

	// helper struct
	ytp := &ytPlaylist{
		id:         id,
		pageToken:  "",
		maxResults: 50,
		maxItems:   maxItems,
	}

	// YouTube playlist searches return a max of 50 results per page
	if ytp.maxResults > ytp.maxItems {
		ytp.maxResults = ytp.maxItems
	}

	for len(tracks) < maxItems {
		tracks, _ = yt.getTrackFromPlaylist(tracks, submitter, playlist, ytp)
		if ytp.pageToken == "" {
			break
		}
	}

	if len(tracks) == 0 {
		return nil, errors.New("Invalid playlist. No tracks were added")
	}
	return tracks, nil
}

func (yt *YouTube) getTrackFromPlaylist(
	tracks []interfaces.Track,
	submitter *gumble.User,
	playlist interfaces.Playlist,
	ytp *ytPlaylist) ([]interfaces.Track, error) {

	var (
		resp *http.Response
		err  error
		v    *jason.Object
	)

	resp, err = http.Get(fmt.Sprintf(yt.playlistItemsAPIURL, ytp.id, ytp.maxResults, viper.GetString("api_keys.youtube"), ytp.pageToken))
	defer resp.Body.Close()
	if err != nil {
		return tracks, err
	}

	v, err = jason.NewObjectFromReader(resp.Body)
	if err != nil {
		return tracks, err
	}

	curTracks, _ := v.GetObjectArray("items")
	for _, track := range curTracks {
		videoID, _ := track.GetString("snippet", "resourceId", "videoId")

		// Unfortunately we have to execute another API call for each video as the YouTube API does not
		// return video durations from the playlistItems endpoint...
		newTrack, _ := yt.getTrack(videoID, submitter, yt.dummyOffset)
		newTrack.Playlist = playlist
		tracks = append(tracks, newTrack)

		if len(tracks) >= ytp.maxItems {
			return tracks, nil
		}
	}

	ytp.pageToken, _ = v.GetString("nextPageToken")
	return tracks, nil
}
