/*
 * MumbleDJ
 * By Matthieu Grieger
 * services/youtube.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package services

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"regexp"

	"github.com/ChannelMeter/iso8601duration"
	"github.com/antonholmquist/jason"
	"github.com/layeh/gumble/gumble"
	"github.com/matthieugrieger/mumbledj/bot"
	"github.com/matthieugrieger/mumbledj/interfaces"
	"github.com/spf13/viper"
)

// YouTube is a wrapper around the YouTube Data API.
// https://developers.google.com/youtube/v3/docs/
type YouTube struct {
	*GenericService
}

// NewYouTubeService returns an initialized YouTube service object.
func NewYouTubeService() *YouTube {
	return &YouTube{
		&GenericService{
			ReadableName: "YouTube",
			Format:       "bestaudio",
			TrackRegex: []*regexp.Regexp{
				regexp.MustCompile(`https?:\/\/www.youtube.com\/watch\?v=(?P<id>[\w-]+)(?P<timestamp>\&t=\d*m?\d*s?)?`),
				regexp.MustCompile(`https?:\/\/youtube.com\/watch\?v=(?P<id>[\w-]+)(?P<timestamp>\&t=\d*m?\d*s?)?`),
				regexp.MustCompile(`https?:\/\/youtu.be\/(?P<id>[\w-]+)(?P<timestamp>\?t=\d*m?\d*s?)?`),
				regexp.MustCompile(`https?:\/\/youtube.com\/v\/(?P<id>[\w-]+)(?P<timestamp>\?t=\d*m?\d*s?)?`),
				regexp.MustCompile(`https?:\/\/www.youtube.com\/v\/(?P<id>[\w-]+)(?P<timestamp>\?t=\d*m?\d*s?)?`),
			},
			PlaylistRegex: []*regexp.Regexp{
				regexp.MustCompile(`https?:\/\/www\.youtube\.com\/playlist\?list=(?P<id>[\w-]+)`),
			},
		},
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
	defer response.Body.Close()
	if err != nil {
		return err
	}

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
		err              error
		resp             *http.Response
		v                *jason.Object
		track            bot.Track
		tracks           []interfaces.Track
	)

	playlistURL = "https://www.googleapis.com/youtube/v3/playlists?part=snippet&id=%s&key=%s"
	playlistItemsURL = "https://www.googleapis.com/youtube/v3/playlistItems?part=snippet,contentDetails&playlistId=%s&maxResults=%d&key=%s&pageToken=%s"
	id, err = yt.getID(url)
	if err != nil {
		return nil, err
	}

	if yt.isPlaylist(url) {
		resp, err = http.Get(fmt.Sprintf(playlistURL, id, viper.GetString("api_keys.youtube")))
		defer resp.Body.Close()
		if err != nil {
			return nil, err
		}

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
			defer curResp.Body.Close()
			if curErr != nil {
				// An error occurred, simply skip this track.
				continue
			}

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
				newTrack, _ := yt.getTrack(videoID, submitter)
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

	track, err = yt.getTrack(id, submitter)
	if err != nil {
		return nil, err
	}
	tracks = append(tracks, track)
	return tracks, nil
}

func (yt *YouTube) getTrack(id string, submitter *gumble.User) (bot.Track, error) {
	var (
		resp *http.Response
		err  error
		v    *jason.Object
	)

	videoURL := "https://www.googleapis.com/youtube/v3/videos?part=snippet,contentDetails&id=%s&key=%s"
	resp, err = http.Get(fmt.Sprintf(videoURL, id, viper.GetString("api_keys.youtube")))
	defer resp.Body.Close()
	if err != nil {
		return bot.Track{}, err
	}

	v, err = jason.NewObjectFromReader(resp.Body)
	if err != nil {
		return bot.Track{}, err
	}
	items, _ := v.GetObjectArray("items")
	item := items[0]
	title, _ := item.GetString("snippet", "title")
	thumbnail, _ := item.GetString("snippet", "thumbnails", "high", "url")
	author, _ := item.GetString("snippet", "channelTitle")
	durationString, _ := item.GetString("contentDetails", "duration")
	durationConverted, _ := duration.FromString(durationString)
	duration := durationConverted.ToDuration()

	return bot.Track{
		ID:           id,
		URL:          "https://youtube.com/watch?v=" + id,
		Title:        title,
		Author:       author,
		Submitter:    submitter.Name,
		Service:      yt.ReadableName,
		Filename:     id + ".track",
		ThumbnailURL: thumbnail,
		Duration:     duration,
		Playlist:     nil,
	}, nil
}
