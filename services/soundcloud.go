/*
 * MumbleDJ
 * By Matthieu Grieger
 * services/soundcloud.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package services

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antonholmquist/jason"
	"layeh.com/gumble/gumble"
	"github.com/RichardNysater/mumbledj/bot"
	"github.com/RichardNysater/mumbledj/interfaces"
	"github.com/spf13/viper"
)

// SoundCloud is a wrapper around the SoundCloud API.
// https://developers.soundcloud.com/docs/api/reference
type SoundCloud struct {
	*GenericService
}

// NewSoundCloudService returns an initialized SoundCloud service object.
func NewSoundCloudService() *SoundCloud {
	return &SoundCloud{
		&GenericService{
			ReadableName: "SoundCloud",
			Format:       "bestaudio",
			TrackRegex: []*regexp.Regexp{
				regexp.MustCompile(`https?:\/\/(www\.)?soundcloud\.com\/([\w-]+)\/([\w-]+)(#t=\n\n?(:\n\n)*)?`),
			},
			PlaylistRegex: []*regexp.Regexp{
				regexp.MustCompile(`https?:\/\/(www\.)?soundcloud\.com\/([\w-]+)\/sets\/([\w-]+)`),
			},
		},
	}
}

// CheckAPIKey performs a test API call with the API key
// provided in the configuration file to determine if the
// service should be enabled.
func (sc *SoundCloud) CheckAPIKey() error {
	if viper.GetString("api_keys.soundcloud") == "" {
		return errors.New("No SoundCloud API key has been provided")
	}
	url := "http://api.soundcloud.com/tracks/13158665?client_id=%s"
	response, err := http.Get(fmt.Sprintf(url, viper.GetString("api_keys.soundcloud")))
	defer response.Body.Close()
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New(response.Status)
	}
	return nil
}

// GetTracks uses the passed URL to find and return
// tracks associated with the URL. An error is returned
// if any error occurs during the API call.
func (sc *SoundCloud) GetTracks(url string, submitter *gumble.User) ([]interfaces.Track, error) {
	var (
		apiURL string
		err    error
		resp   *http.Response
		v      *jason.Object
		track  bot.Track
		tracks []interfaces.Track
	)

	urlSplit := strings.Split(url, "#t=")

	apiURL = "http://api.soundcloud.com/resolve?url=%s&client_id=%s"

	if sc.isPlaylist(url) {
		// Submitter has added a playlist!
		resp, err = http.Get(fmt.Sprintf(apiURL, urlSplit[0], viper.GetString("api_keys.soundcloud")))
		defer resp.Body.Close()
		if err != nil {
			return nil, err
		}

		v, err = jason.NewObjectFromReader(resp.Body)
		if err != nil {
			return nil, err
		}

		title, _ := v.GetString("title")
		permalink, _ := v.GetString("permalink_url")
		playlist := &bot.Playlist{
			ID:        permalink,
			Title:     title,
			Submitter: submitter.Name,
			Service:   sc.ReadableName,
		}

		var scTracks []*jason.Object
		scTracks, err = v.GetObjectArray("tracks")
		if err != nil {
			return nil, err
		}

		dummyOffset, _ := time.ParseDuration("0s")
		for _, t := range scTracks {
			track, err = sc.getTrack(t, dummyOffset, submitter)
			if err != nil {
				// Skip this track.
				continue
			}
			track.Playlist = playlist
			tracks = append(tracks, track)
		}

		if len(tracks) == 0 {
			return nil, errors.New("Invalid playlist. No tracks were added")
		}
		return tracks, nil
	}

	// Submitter has added a track!

	offset := 0
	// Calculate track offset if needed
	if len(urlSplit) == 2 {
		timeSplit := strings.Split(urlSplit[1], ":")
		multiplier := 1
		for i := len(timeSplit) - 1; i >= 0; i-- {
			time, _ := strconv.Atoi(timeSplit[i])
			offset += time * multiplier
			multiplier *= 60
		}
	}
	playbackOffset, _ := time.ParseDuration(fmt.Sprintf("%ds", offset))

	resp, err = http.Get(fmt.Sprintf(apiURL, urlSplit[0], viper.GetString("api_keys.soundcloud")))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	v, err = jason.NewObjectFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	track, err = sc.getTrack(v, playbackOffset, submitter)
	if err != nil {
		return nil, err
	}

	tracks = append(tracks, track)
	return tracks, nil
}

func (sc *SoundCloud) getTrack(obj *jason.Object, offset time.Duration, submitter *gumble.User) (bot.Track, error) {
	title, _ := obj.GetString("title")
	idInt, _ := obj.GetInt64("id")
	id := strconv.FormatInt(idInt, 10)
	url, _ := obj.GetString("permalink_url")
	author, _ := obj.GetString("user", "username")
	authorURL, _ := obj.GetString("user", "permalink_url")
	durationMS, _ := obj.GetInt64("duration")
	duration, _ := time.ParseDuration(fmt.Sprintf("%dms", durationMS))
	thumbnail, err := obj.GetString("artwork_url")
	if err != nil {
		// Track has no artwork, using profile avatar instead.
		thumbnail, _ = obj.GetString("user", "avatar_url")
	}

	return bot.Track{
		ID:             id,
		URL:            url,
		Title:          title,
		Author:         author,
		AuthorURL:      authorURL,
		Submitter:      submitter.Name,
		Service:        sc.ReadableName,
		Filename:       id + ".track",
		ThumbnailURL:   thumbnail,
		Duration:       duration,
		PlaybackOffset: offset,
		Playlist:       nil,
	}, nil
}
