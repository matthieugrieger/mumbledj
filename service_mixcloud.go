/*
 * MumbleDJ
 * By Matthieu Grieger
 * service_mixcloud.go
 * Copyright (c) 2016 Benjmain Klettbach (MIT License)
 */

package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/jsonq"
	"github.com/layeh/gumble/gumble"
)

// Regular expressions for mixcloud urls
var mixcloudSongPattern = `https?:\/\/(www\.)?mixcloud\.com\/([\w-]+)\/([\w-]+)(#t=\n\n?(:\n\n)*)?`

// Mixcloud implements the Service interface
type Mixcloud struct{}

// ------------------
//  MIXCLOUD SERVICE
// ------------------

// ServiceName is the human readable version of the service name
func (mc Mixcloud) ServiceName() string {
	return "Mixcloud"
}

// TrackName is the human readable version of the service name
func (mc Mixcloud) TrackName() string {
	return "Song"
}

// URLRegex checks to see if service will accept URL
func (mc Mixcloud) URLRegex(url string) bool {
	return RegexpFromURL(url, []string{mixcloudSongPattern}) != nil
}

// NewRequest creates the requested song and adds to the queue
func (mc Mixcloud) NewRequest(user *gumble.User, url string) ([]Song, error) {
	var apiResponse *jsonq.JsonQuery
	var songArray []Song
	var err error
	timesplit := strings.Split(url, "#t=")
	url = strings.Replace(timesplit[0], "www", "api", 1)
	if apiResponse, err = PerformGetRequest(url); err != nil {
		return nil, errors.New(INVALID_URL_MSG)
	}

	if strings.Contains(url, "playlists") {
		// PLAYLIST
		// Playlists from Mixcloud are not supported, because they do not provide an API for them.
		return nil, errors.New(fmt.Sprintf(NO_PLAYLISTS_SUPPORTED_MSG, mc.ServiceName()))
	} else {
		// SONG
		// Calculate offset
		offset := 0
		if len(timesplit) == 2 {
			timesplit = strings.Split(timesplit[1], ":")
			multiplier := 1
			for i := len(timesplit) - 1; i >= 0; i-- {
				time, _ := strconv.Atoi(timesplit[i])
				offset += time * multiplier
				multiplier *= 60
			}
		}

		// Add the track
		if song, err := mc.NewSong(user, apiResponse, offset); err == nil {
			return append(songArray, song), err
		}
		return nil, err
	}
}

// NewSong creates a track and adds to the queue
func (mc Mixcloud) NewSong(user *gumble.User, trackData *jsonq.JsonQuery, offset int) (Song, error) {
	title, _ := trackData.String("name")
	id, _ := trackData.String("slug")
	duration, _ := trackData.Int("audio_length")
	url, _ := trackData.String("url")
	thumbnail, err := trackData.String("pictures", "large")
	if err != nil {
		// Song has no artwork, using profile avatar instead
		userObj, _ := trackData.Object("user")
		thumbnail, _ = jsonq.NewQuery(userObj).String("pictures", "thumbnail")
	}

	song := &AudioTrack{
		id:        id,
		title:     title,
		url:       url,
		thumbnail: thumbnail,
		submitter: user,
		duration:  duration,
		offset:    offset,
		format:    "best",
		skippers:  make([]string, 0),
		dontSkip:  false,
		service:   mc,
	}
	return song, nil
}
