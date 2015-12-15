/*
 * MumbleDJ
 * By Matthieu Grieger
 * service.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package main

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/layeh/gumble/gumble"
)

// Service interface. Each service will implement these functions
type Service interface {
	ServiceName() string
	TrackName() string
	URLRegex(string) bool
	NewRequest(*gumble.User, string) ([]Song, error)
}

// Song interface. Each service will implement these
// functions in their Song types.
type Song interface {
	Download() error
	Play()
	Delete() error
	AddSkip(string) error
	RemoveSkip(string) error
	SkipReached(int) bool
	Submitter() string
	Title() string
	ID() string
	Filename() string
	Duration() time.Duration
	Thumbnail() string
	Playlist() Playlist
	DontSkip() bool
	SetDontSkip(bool)
}

// Playlist interface. Each service will implement these
// functions in their Playlist types.
type Playlist interface {
	AddSkip(string) error
	RemoveSkip(string) error
	DeleteSkippers()
	SkipReached(int) bool
	ID() string
	Title() string
}

var services []Service

// FindServiceAndAdd tries the given url with each service
// and adds the song/playlist with the correct service
func FindServiceAndAdd(user *gumble.User, url string) error {
	var urlService Service

	// Checks all services to see if any can take the URL
	for _, service := range services {
		if service.URLRegex(url) {
			urlService = service
		}
	}

	if urlService == nil {
		return errors.New(INVALID_URL_MSG)
	} else {
		var title string
		var songsAdded = 0
		var songArray []Song
		var err error

		// Get service to create songs
		if songArray, err = urlService.NewRequest(user, url); err != nil {
			return err
		}

		// Check Playlist Permission
		if len(songArray) > 1 && !dj.HasPermission(user.Name, dj.conf.Permissions.AdminAddPlaylists) {
			return errors.New(NO_PLAYLIST_PERMISSION_MSG)
		}

		// Loop through all songs and add to the queue
		oldLength := dj.queue.Len()
		for _, song := range songArray {
			// Check song is not too long
			if dj.conf.General.MaxSongDuration == 0 || int(song.Duration().Seconds()) <= dj.conf.General.MaxSongDuration {
				if !isNil(song.Playlist()) {
					title = song.Playlist().Title()
				} else {
					title = song.Title()
				}

				// Add song to queue
				dj.queue.AddSong(song)
				songsAdded++
			}
		}

		// Alert channel of added song/playlist
		if songsAdded == 0 {
			return errors.New(fmt.Sprintf(TRACK_TOO_LONG_MSG, urlService.ServiceName()))
		} else if songsAdded == 1 {
			dj.client.Self.Channel.Send(fmt.Sprintf(SONG_ADDED_HTML, user.Name, title), false)
		} else {
			dj.client.Self.Channel.Send(fmt.Sprintf(PLAYLIST_ADDED_HTML, user.Name, title), false)
		}

		// Starts playing the new song if nothing else is playing
		if oldLength == 0 && dj.queue.Len() != 0 && !dj.audioStream.IsPlaying() {
			if dj.conf.General.AutomaticShuffleOn {
				dj.queue.RandomNextSong(true)
			}
			if err := dj.queue.CurrentSong().Download(); err == nil {
				dj.queue.CurrentSong().Play()
			} else {
				dj.queue.CurrentSong().Delete()
				dj.queue.OnSongFinished()
				return errors.New(AUDIO_FAIL_MSG)
			}
		}
		return nil
	}
}

// FindServiceAndInsertNext tries the given url with each service
// and inserts the song/playlist with the correct service into the slot after the current one
func FindServiceAndInsertNext(user *gumble.User, url string) error {
	var urlService Service

	// Checks all services to see if any can take the URL
	for _, service := range services {
		if service.URLRegex(url) {
			urlService = service
		}
	}

	if urlService == nil {
		return errors.New(INVALID_URL_MSG)
	} else {
		var title string
		var songsAdded = 0
		var songArray []Song
		var err error

		// Get service to create songs
		if songArray, err = urlService.NewRequest(user, url); err != nil {
			return err
		}

		// Check Playlist Permission
		if len(songArray) > 1 && !dj.HasPermission(user.Name, dj.conf.Permissions.AdminAddPlaylists) {
			return errors.New(NO_PLAYLIST_PERMISSION_MSG)
		}

		// Loop through all songs and add to the queue
		i := 0
		for _, song := range songArray {
			i++
			// Check song is not too long
			if dj.conf.General.MaxSongDuration == 0 || int(song.Duration().Seconds()) <= dj.conf.General.MaxSongDuration {
				if !isNil(song.Playlist()) {
					title = song.Playlist().Title()
				} else {
					title = song.Title()
				}

				// Add song to queue
				dj.queue.InsertSong(song, i)
				songsAdded++
			}
		}

		// Alert channel of added song/playlist
		if songsAdded == 0 {
			return errors.New(fmt.Sprintf(TRACK_TOO_LONG_MSG, urlService.ServiceName()))
		} else if songsAdded == 1 {
			dj.client.Self.Channel.Send(fmt.Sprintf(NEXT_SONG_ADDED_HTML, user.Name, title), false)
		} else {
			dj.client.Self.Channel.Send(fmt.Sprintf(NEXT_PLAYLIST_ADDED_HTML, user.Name, title), false)
		}

		return nil
	}
}

// RegexpFromURL loops through an array of patterns to see if it matches the url
func RegexpFromURL(url string, patterns []string) *regexp.Regexp {
	for _, pattern := range patterns {
		if re, err := regexp.Compile(pattern); err == nil {
			if re.MatchString(url) {
				return re
			}
		}
	}
	return nil
}
