package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/jsonq"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumble_ffmpeg"
)

// Regular expressions for soundcloud urls
var soundcloudSongPattern = `https?:\/\/(www)?\.soundcloud\.com\/([\w-]+)\/([\w-]+)`
var soundcloudPlaylistPattern = `https?:\/\/(www)?\.soundcloud\.com\/([\w-]+)\/sets\/([\w-]+)`

// ------
// TYPES
// ------

// YouTube implements the Service interface
type SoundCloud struct{}

// YouTubeSong holds the metadata for a song extracted from a YouTube video.
type SoundCloudSong struct {
	submitter string
	title     string
	id        string
	offset    int
	filename  string
	duration  string
	thumbnail string
	skippers  []string
	playlist  Playlist
	dontSkip  bool
}

// YouTubePlaylist holds the metadata for a YouTube playlist.
type SoundCloudPlaylist struct {
	id    string
	title string
}

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
	url := fmt.Sprintf("http://api.soundcloud.com/resolve?url=%s&client_id=%s", url, os.Getenv("SOUNDCLOUD_API_KEY"))
	if apiResponse, err = PerformGetRequest(url); err != nil {
		return nil, errors.New(INVALID_API_KEY)
	}

	title, _ := apiResponse.String("title")
	tracks, err := apiResponse.ArrayOfObjects("tracks")

	if err == nil {
		if re.MatchString(url) {
			// PLAYLIST
			if dj.HasPermission(user.Name, dj.conf.Permissions.AdminAddPlaylists) {
				playlist, err := sc.NewPlaylist(user.Name, url)
				return playlist.Title(), err
			} else {
				return "", errors.New("NO_PLAYLIST_PERMISSION")
			}
		} else {

			// SONG
			song, err := sc.NewSong(user.Name, url, nil)
			return song.Title(), err
		}
	} else {
		return "", err
	}
}
