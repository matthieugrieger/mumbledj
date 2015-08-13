package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jmoiron/jsonq"
	"github.com/layeh/gumble/gumble"
)

// Regular expressions for soundcloud urls
var soundcloudSongPattern = `https?:\/\/(www)?\.soundcloud\.com\/([\w-]+)\/([\w-]+)`
var soundcloudPlaylistPattern = `https?:\/\/(www)?\.soundcloud\.com\/([\w-]+)\/sets\/([\w-]+)`

// SoundCloud implements the Service interface
type SoundCloud struct{}

// ------------------
// SOUNDCLOUD SERVICE
// ------------------

// Name of the service
func (sc SoundCloud) ServiceName() string {
	return "SoundCloud"
}

// Checks to see if service will accept URL
func (sc SoundCloud) URLRegex(url string) bool {
	return RegexpFromURL(url, []string{soundcloudSongPattern, soundcloudPlaylistPattern}) != nil
}

// Creates the requested song/playlist and adds to the queue
func (sc SoundCloud) NewRequest(user *gumble.User, url string) (string, error) {
	var apiResponse *jsonq.JsonQuery
	var err error
	url = fmt.Sprintf("http://api.soundcloud.com/resolve?url=%s&client_id=%s", url, os.Getenv("SOUNDCLOUD_API_KEY"))
	if apiResponse, err = PerformGetRequest(url); err != nil {
		return "", errors.New(INVALID_API_KEY)
	}

	tracks, err := apiResponse.ArrayOfObjects("tracks")
	if err == nil {
		// PLAYLIST
		if dj.HasPermission(user.Name, dj.conf.Permissions.AdminAddPlaylists) {
			// Check duration of playlist
			duration, _ := apiResponse.Int("duration")

			// Create playlist
			title, _ := apiResponse.String("title")
			permalink, _ := apiResponse.String("permalink_url")
			playlist := &YouTubeDLPlaylist{
				id:    permalink,
				title: title,
			}

			// Add all tracks
			for _, t := range tracks {
				sc.NewSong(user, jsonq.NewQuery(t), playlist)
			}
			if err == nil {
				return playlist.Title(), nil
			} else {
				Verbose("soundcloud.NewRequest: " + err.Error())
				return "", err
			}
		} else {
			return "", errors.New("NO_PLAYLIST_PERMISSION")
		}
	} else {
		// SONG
		return sc.NewSong(user, apiResponse, nil)
	}
}

// Creates a track and adds to the queue
func (sc SoundCloud) NewSong(user *gumble.User, trackData *jsonq.JsonQuery, playlist Playlist) (string, error) {
	title, err := trackData.String("title")
	if err != nil {
		return "", err
	}
	id, err := trackData.String("id")
	if err != nil {
		return "", err
	}
	duration, err := trackData.Int("duration")
	if err != nil {
		return "", err
	}
	thumbnail, err := trackData.String("artwork_uri")
	if err != nil {
		return "", err
	}
	url, err := trackData.String("permalink_url")
	if err != nil {
		return "", err
	}

	song := &YouTubeDLSong{
		id:        id,
		title:     title,
		url:       url,
		thumbnail: thumbnail,
		submitter: user,
		duration:  duration,
		playlist:  playlist,
		skippers:  make([]string, 0),
		dontSkip:  false,
	}
	dj.queue.AddSong(song)
	return title, nil
}
