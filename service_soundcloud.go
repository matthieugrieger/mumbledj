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

	tracks, err := apiResponse.ArrayOfObjects("tracks")
	if err == nil {
		// PLAYLIST
		if dj.HasPermission(user.Name, dj.conf.Permissions.AdminAddPlaylists) {
			// Check duration of playlist
			// duration, _ := apiResponse.Int("duration")

			// Create playlist
			title, _ := apiResponse.String("title")
			permalink, _ := apiResponse.String("permalink_url")
			playlist := &SoundCloudPlaylist{
				id:    permalink,
				title: title,
			}

			// Add all tracks
			for _, t := range tracks {
				sc.NewSong(user.Name, jsonq.NewQuery(t), playlist)
			}
			return playlist.Title(), err
		} else {
			return "", errors.New("NO_PLAYLIST_PERMISSION")
		}
	} else {
		return sc.NewSong(user.Name, apiResponse, nil)
	}
}

// Creates a track and adds to the queue
func (sc SoundCloud) NewSong(user string, trackData *jsonq.JsonQuery, playlist SoundCloudPlaylist) (string, error) {
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

	song := &YoutubeDL{
		id:        id,
		title:     title,
		thumbnail: thumbnail,
		submitter: user.Name,
		duration:  duration,
		playlist:  playlist,
		skippers:  make([]string, 0),
		dontSkip:  false,
	}
	dj.queue.AddSong(song)
	return title, nil
}

// ----------------
// YOUTUBE PLAYLIST
// ----------------

// AddSkip adds a skip to the playlist's skippers slice.
func (p *SoundCloudPlaylist) AddSkip(username string) error {
	for _, user := range dj.playlistSkips[p.ID()] {
		if username == user {
			return errors.New("This user has already skipped the current song.")
		}
	}
	dj.playlistSkips[p.ID()] = append(dj.playlistSkips[p.ID()], username)
	return nil
}

// RemoveSkip removes a skip from the playlist's skippers slice. If username is not in the slice
// an error is returned.
func (p *YouTubePlaylist) RemoveSkip(username string) error {
	for i, user := range dj.playlistSkips[p.ID()] {
		if username == user {
			dj.playlistSkips[p.ID()] = append(dj.playlistSkips[p.ID()][:i], dj.playlistSkips[p.ID()][i+1:]...)
			return nil
		}
	}
	return errors.New("This user has not skipped the song.")
}

// DeleteSkippers removes the skippers entry in dj.playlistSkips.
func (p *YouTubePlaylist) DeleteSkippers() {
	delete(dj.playlistSkips, p.ID())
}

// SkipReached calculates the current skip ratio based on the number of users within MumbleDJ's
// channel and the number of usernames in the skippers slice. If the value is greater than or equal
// to the skip ratio defined in the config, the function returns true, and returns false otherwise.
func (p *YouTubePlaylist) SkipReached(channelUsers int) bool {
	if float32(len(dj.playlistSkips[p.ID()]))/float32(channelUsers) >= dj.conf.General.PlaylistSkipRatio {
		return true
	}
	return false
}

// ID returns the id of the YouTubePlaylist.
func (p *YouTubePlaylist) ID() string {
	return p.id
}

// Title returns the title of the YouTubePlaylist.
func (p *YouTubePlaylist) Title() string {
	return p.title
}
