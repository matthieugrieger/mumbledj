package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/jsonq"
	"github.com/layeh/gumble/gumble"
)

// Regular expressions for soundcloud urls
var soundcloudSongPattern = `https?:\/\/(www\.)?soundcloud\.com\/([\w-]+)\/([\w-]+)(#t=\n\n?(:\n\n)*)?`
var soundcloudPlaylistPattern = `https?:\/\/(www\.)?soundcloud\.com\/([\w-]+)\/sets\/([\w-]+)`

// SoundCloud implements the Service interface
type SoundCloud struct{}

// ------------------
// SOUNDCLOUD SERVICE
// ------------------

// URLRegex checks to see if service will accept URL
func (sc SoundCloud) URLRegex(url string) bool {
	return RegexpFromURL(url, []string{soundcloudSongPattern, soundcloudPlaylistPattern}) != nil
}

// NewRequest creates the requested song/playlist and adds to the queue
func (sc SoundCloud) NewRequest(user *gumble.User, url string) (string, error) {
	var apiResponse *jsonq.JsonQuery
	var err error
	timesplit := strings.Split(url, "#t=")
	url = fmt.Sprintf("http://api.soundcloud.com/resolve?url=%s&client_id=%s", timesplit[0], os.Getenv("SOUNDCLOUD_API_KEY"))
	if apiResponse, err = PerformGetRequest(url); err != nil {
		return "", errors.New(INVALID_API_KEY)
	}

	tracks, err := apiResponse.ArrayOfObjects("tracks")
	if err == nil {
		// PLAYLIST
		if dj.HasPermission(user.Name, dj.conf.Permissions.AdminAddPlaylists) {
			// Create playlist
			title, _ := apiResponse.String("title")
			permalink, _ := apiResponse.String("permalink_url")
			playlist := &YouTubePlaylist{
				id:    permalink,
				title: title,
			}

			// Add all tracks
			for _, t := range tracks {
				sc.NewSong(user, jsonq.NewQuery(t), 0, playlist)
			}
			return playlist.Title(), nil
		}
		return "", errors.New(NO_PLAYLIST_PERMISSION_MSG)
	} else {
		// SONG
		offset := 0
		if len(timesplit) == 2 {
			timesplit = strings.Split(timesplit[1], ":")
			multiplier := 1
			for i := len(timesplit) - 1; i >= 0; i-- {
				time, _ := strconv.Atoi(timesplit[i])
				offset += time * multiplier
				multiplier *= 60
			}
			fmt.Printf("Offset: " + strconv.Itoa(offset)) // DEBUG
		}
		return sc.NewSong(user, apiResponse, offset, nil)
	}
}

// NewSong creates a track and adds to the queue
func (sc SoundCloud) NewSong(user *gumble.User, trackData *jsonq.JsonQuery, offset int, playlist Playlist) (string, error) {
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

	// Check song is not longer than the MaxSongDuration
	if dj.conf.General.MaxSongDuration == 0 || (durationMS/1000) <= dj.conf.General.MaxSongDuration {
		timeDuration, _ := time.ParseDuration(strconv.Itoa(durationMS) + "ms")
		duration := timeDuration.String()

		song := &YouTubeSong{
			id:        strconv.Itoa(id),
			title:     title,
			url:       url,
			thumbnail: thumbnail,
			submitter: user,
			duration:  duration,
			offset:    offset,
			format:    "mp3",
			playlist:  playlist,
			skippers:  make([]string, 0),
			dontSkip:  false,
		}
		dj.queue.AddSong(song)
		return song.Title(), nil
	}
	return "", errors.New(VIDEO_TOO_LONG_MSG)
}
