/*
 * MumbleDJ
 * By Matthieu Grieger
 * service_soundcloud.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
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

// Regular expressions for soundcloud urls
var soundcloudSongPattern = `https?:\/\/(www\.)?soundcloud\.com\/([\w-]+)\/([\w-]+)(#t=\n\n?(:\n\n)*)?`
var soundcloudPlaylistPattern = `https?:\/\/(www\.)?soundcloud\.com\/([\w-]+)\/sets\/([\w-]+)`

// SearchService name
var soundcloudSearchServiceName = "sc"

// SoundCloud implements the Service interface
type SoundCloud struct{}

// ------------------
// SOUNDCLOUD SERVICE
// ------------------

// ServiceName is the human readable version of the service name
func (sc SoundCloud) ServiceName() string {
	return "Soundcloud"
}

// TrackName is the human readable version of the service name
func (sc SoundCloud) TrackName() string {
	return "Song"
}

// URLRegex checks to see if service will accept URL
func (sc SoundCloud) URLRegex(url string) bool {
	return RegexpFromURL(url, []string{soundcloudSongPattern, soundcloudPlaylistPattern}) != nil
}

// SearchRegex checks to see if service will accept the searchString 
func (sc SoundCloud) SearchRegex(searchService string) bool {
	return searchService == soundcloudSearchServiceName
}

// NewRequest creates the requested song/playlist and adds to the queue
func (sc SoundCloud) NewRequest(user *gumble.User, url string) ([]Song, error) {
	var apiResponse *jsonq.JsonQuery
	var songArray []Song
	var err error
	timesplit := strings.Split(url, "#t=")
	url = fmt.Sprintf("http://api.soundcloud.com/resolve?url=%s&client_id=%s", timesplit[0], dj.conf.ServiceKeys.SoundCloud)
	if apiResponse, err = PerformGetRequest(url); err != nil {
		return nil, errors.New(fmt.Sprintf(INVALID_API_KEY, sc.ServiceName()))
	}

	tracks, err := apiResponse.ArrayOfObjects("tracks")
	if err == nil {
		// PLAYLIST
		// Create playlist
		title, _ := apiResponse.String("title")
		permalink, _ := apiResponse.String("permalink_url")
		playlist := &AudioPlaylist{
			id:    permalink,
			title: title,
		}

		if (dj.conf.General.MaxSongPerPlaylist > 0 && len(tracks) > dj.conf.General.MaxSongPerPlaylist){
		   tracks = tracks[:dj.conf.General.MaxSongPerPlaylist]
		}
		// Add all tracks
		for _, t := range tracks {
			if song, err := sc.NewSong(user, jsonq.NewQuery(t), 0, playlist); err == nil {
				songArray = append(songArray, song)
			}
		}
		return songArray, nil
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
		if song, err := sc.NewSong(user, apiResponse, offset, nil); err == nil {
			return append(songArray, song), err
		}
		return nil, err
	}
}

// SearchSong searches for a Song and returns the Songs URL
func (sc SoundCloud) SearchSong(searchString string) (string, error) {
	var returnString string
	url := fmt.Sprintf("https://api.soundcloud.com/tracks?q=%s&client_id=%s&limit=1", searchString, dj.conf.ServiceKeys.SoundCloud)

	if apiResponse, err := PerformGetRequest(url); err == nil {
		returnString, _ = apiResponse.String("json", "0", "permalink_url");
		return returnString, nil
	}
	return "", errors.New(fmt.Sprintf(INVALID_API_KEY, sc.ServiceName()))
}

// NewSong creates a track and adds to the queue
func (sc SoundCloud) NewSong(user *gumble.User, trackData *jsonq.JsonQuery, offset int, playlist Playlist) (Song, error) {
	title, _ := trackData.String("title")
	id, _ := trackData.Int("id")
	durationMS, _ := trackData.Int("duration")
	url, _ := trackData.String("permalink_url")
	thumbnail, err := trackData.String("artwork_url")
	if err != nil {
		// Song has no artwork, using profile avatar instead
		userObj, _ := trackData.Object("user")
		thumbnail, _ = jsonq.NewQuery(userObj).String("avatar_url")
	}

	song := &AudioTrack{
		id:        strconv.Itoa(id),
		title:     title,
		url:       url,
		thumbnail: thumbnail,
		submitter: user,
		duration:  durationMS / 1000,
		offset:    offset,
		format:    "mp3",
		playlist:  playlist,
		skippers:  make([]string, 0),
		dontSkip:  false,
		service:   sc,
	}
	return song, nil
}
